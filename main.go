package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Redis *redis.Client
	ctx   context.Context
)

func main() {
	// Check for Redis configuration and connection.
	err := redisCheck()
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
	// Get the second we are running at.
	t := time.Now()
	second := t.Second()
	key := fmt.Sprintf("%d", second)

	// See if there's data at that key.
	val, err := Redis.Get(ctx, key).Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// If not - fake some data and shove it in.
	if val == "" {
		val = "Some faked value"
		err := Redis.Set(ctx, key, val, 0).Err()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusOK, val)
}

func redisCheck() error {
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return errors.New("must set REDIS_URL")
	}
	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		return errors.New("must set REDIS_PASSWORD")
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	})
	return nil
}
