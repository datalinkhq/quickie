<div align="center"><img src="https://raw.githubusercontent.com/datalinkhq/datalink/main/assets/dark-wideshot.png#gh-dark-mode-only" width="50%" ></div>
<div align="center"><img src="https://raw.githubusercontent.com/datalinkhq/datalink/main/assets/light-wideshot.png#gh-light-mode-only" width="50%" ></div>

<h3 align="center">quickie - Yet another redis implementation</h1>

# What is this?
Quickie is a database caching solution to cache data using redis, written in Go. 

# Installation & Setup
Follow the steps below to install and setup quickie. 
## Prerequisites: 
- [Go](https://go.dev)
- [Gin](https://gin-gonic.com/)
- [Redis](https://redis.io)

Install the dependencies and set up the project by running:
```console
go get && go run src/impl.go
```

This will run a Gin webserver on port `8080`. 

Next, set up a environment variable known as `SECRET` with the secret key required to access the endpoints. 

# Usage

**NOTE: To access any endpoints, an Authorization header with a secret needs to be provided, for more details on configuring this, check [Installation & Setup](https://github.com/datalinkhq/quickie#installation-setup).**
<br>
<br>
This service exposes the following endpoints:

### `/set/:table` - Set a key, value pair in a table
#### Required Query Parameters: `value`, `key`

#### Example Request:
```
curl http://localhost:8080/set/datastore?key=name&value=jack
```

Where datastore is the table name — "name" and "jack" are the key and value respectively. 

### `/get/:table` - Set a key, value pair in a table
#### Required Query Parameters: `key`

#### Example Request:
```
curl http://localhost:8080/get/datastore?key=name
```

Where datastore is the table name and "name" is the key.
