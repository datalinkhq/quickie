package main

import (
	json "encoding/json"
	fmt "fmt"
	http "net/http"
	os "os"
	time "time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	gin "github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

type dataValue struct {
	Table string `json:"table"`
	Value any    `json:"value"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	secret := os.Getenv("SECRET")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: client,
		Rate:        time.Second,
		Limit:       5,
	})
	limiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: requestsHandler,
		KeyFunc:      getClientIP,
	})
	router := gin.Default()
	router.POST("/set/:table", limiter, func(c *gin.Context) {
		_, clientInitErr := client.Ping(c).Result()

		if clientInitErr != nil {
			panic("Failed to initialize redis client! Is the server running?")
		}
		if c.GetHeader("Authorization") == secret {
			fmt.Println(c.GetHeader("Authorization"))
			fmt.Println(secret)
			value := c.Query("value")
			key := c.Query("key")
			table := c.Param("table")
			fmt.Println("Table: ", table)

			if value == "" || key == "" {
				data := generateResponseData(400, "error", "Invalid request")
				c.JSON(http.StatusOK, data)
				return

			} else {
				var dataValues []*dataValue
				dataValues = append(dataValues, &dataValue{Table: table, Value: value})

				s, _ := json.Marshal(dataValues)

				_, err := client.Set(c, fmt.Sprintf("%s_%s", table, key), s, 0).Result()

				if err != nil {
					c.JSON(http.StatusInternalServerError, generateResponseData(500, "Internal Server Error", "Saving data failed!", fmt.Sprint(err)))

					return
				}

				fetched, err := client.Get(c, fmt.Sprintf("%s_%s", table, key)).Result()
				trueValue := string(s)
				// Verify if the saved value matches the requested value to be saved
				if fetched != trueValue || err != nil {
					c.JSON(http.StatusInternalServerError, generateResponseData(500, "Internal Server Error", "Value matching failed!"))
					return
				} else if fetched == fmt.Sprint(s) && err == nil {
					data := generateResponseData(200, "success", fmt.Sprintf("Key: %s", key), fmt.Sprintf("Value: %s", fetched))
					c.JSON(http.StatusOK, data)
					return
				}

			}

		} else {
			data := generateResponseData(400, "Bad Request")
			c.JSON(http.StatusBadRequest, data)
			return
		}
	})

	router.GET("/get/:table", limiter, func(c *gin.Context) {
		if c.GetHeader("Authorization") == secret {
			// TODO: search for a key using a reference, like "findMany where condition" in prisma.
			key := c.Query("key")
			table := c.Param("table")
			if key == "" {
				c.JSON(http.StatusBadRequest, generateResponseData(400, "Bad Request"))
				return
			} else {
				value, err := client.Get(c, fmt.Sprintf("%s_%s", table, key)).Result()

				if err != nil {
					c.JSON(http.StatusInternalServerError, generateResponseData(500, "Internal Server Error"))
					return
				} else {
					c.JSON(http.StatusOK, generateResponseData(200, "success", value))
					return
				}
			}
		} else {
			data := generateResponseData(400, "Bad Request")
			c.JSON(http.StatusBadRequest, data)
			return
		}
	})

	router.Run()
}

func generateResponseData(code int, status string, other ...interface{}) gin.H {
	return gin.H{
		"code":    code,
		"status":  status,
		"details": other,
	}
}

func (s dataValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s dataValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func getClientIP(c *gin.Context) string {
	return c.ClientIP()
}

func requestsHandler(c *gin.Context, info ratelimit.Info) {
	resetIn := time.Until(info.ResetTime).String()
	c.JSON(http.StatusTooManyRequests, generateResponseData(429, "Too Many Requests", fmt.Sprintf("Try again in %s", resetIn)))
}
