package main

import "github.com/go-redis/redis/v8"

// Use this library for redis related operations as a wrapper
func InitializeRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
