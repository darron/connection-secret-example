package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Redis redis.Conn
	ctx   context.Context
)

func main() {
	// Check for Redis configuration and connection.
	var err error
	Redis, err = redisCheck()
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/redis", redisRoute)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func redisRoute(c echo.Context) error {
	// Connect to Redis
	rdb, err := redisCheck()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Get the second we are running at.
	t := time.Now()
	second := t.Second()
	key := fmt.Sprintf("%d", second)

	// See if there's data at that key.
	val, err := rdb.Do("GET", key)
	if err != nil {
		log.Println("got a GET error from Redis:", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var finalVal string

	// If there's no data - fake some data and shove it in.
	if val == nil {
		log.Println("no value there - adding some")
		finalVal = "Some faked value"
		err := rdb.Send("SET", key, finalVal)
		if err != nil {
			log.Println("got a SET error from Redis", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		rdb.Flush()
	} else {
		finalVal = fmt.Sprintf("%s", val)
	}

	return c.String(http.StatusOK, finalVal)
}

func redisCheck() (redis.Conn, error) {
	var r redis.Conn
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Println("REDIS_URL problem")
		return r, errors.New("must set REDIS_URL")
	}
	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		log.Println("REDIS_PASSWORD problem")
		return r, errors.New("must set REDIS_PASSWORD")
	}
	r, err := redis.Dial("tcp", redisURL)
	if err != nil {
		log.Println("redis.Dial problem")
		return r, err
	}
	if redisPassword != "" {
		if _, err := r.Do("AUTH", redisPassword); err != nil {
			log.Println("redis AUTH problem")
			r.Close()
			return r, err
		}
	}
	return r, nil
}
