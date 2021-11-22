package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Redis redis.Conn
	cc    string
)

func main() {
	// Check for Redis configuration and connection.
	var err error
	Redis, err = getRedis()
	if err != nil {
		log.Fatal(err)
	}

	// Set consistent somwehat random value for the life of this process.
	cc = gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: []string{"visa", "discover"}})

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Setup Prometheus metrics.
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Routes
	e.GET("/", hello)
	e.GET("/redis", redisRoute)
	e.GET("/healthz", healthz)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func redisRoute(c echo.Context) error {
	key := getKey(cc)

	// See if there's data at that key location
	val, err := Redis.Do("GET", key)
	if err != nil {
		log.Println("got a GET error from Redis:", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var finalVal string

	// If there's no data - fake some data and shove it in.
	if val == nil {
		log.Println("no data there - adding some fake data")
		j, err := getFakeData()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		Redis.Send("SET", key, string(j))
		// Set a random TTL - throw some randomness into here.
		randTTL := rand.Intn(120) + 10
		Redis.Send("EXPIRE", key, randTTL, "NX")
		rerr := Redis.Flush()
		if rerr != nil {
			log.Println("got a FLUSH error from Redis", rerr.Error())
			return c.String(http.StatusInternalServerError, rerr.Error())
		}
		finalVal = string(j)
	} else {
		finalVal = fmt.Sprintf("%s", val)
	}

	return c.String(http.StatusOK, finalVal)
}

func getRedis() (redis.Conn, error) {
	var r redis.Conn
	var err error

	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return r, errors.New("must set REDIS_URL")
	}
	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		return r, errors.New("must set REDIS_PASSWORD")
	}
	// Let's check for TLS - we dial differently if it's enabled.
	redisTLSEnabled, _ := os.LookupEnv("REDIS_CLIENT_USE_TLS")
	if redisTLSEnabled == "true" {
		// Yes - this is insecure - this is just a demo application.
		// Don't do this in real life.
		clientTLSConfig := &tls.Config{InsecureSkipVerify: true}
		r, err = redis.Dial("tcp", redisURL,
			redis.DialTLSConfig(clientTLSConfig),
			redis.DialUseTLS(true))
	} else {
		r, err = redis.Dial("tcp", redisURL)
	}

	if err != nil {
		return r, fmt.Errorf("redis.Dial problem: %w", err)
	}
	if redisPassword != "" {
		if _, err := r.Do("AUTH", redisPassword); err != nil {
			r.Close()
			return r, fmt.Errorf("redis AUTH problem: %w", err)
		}
	}
	return r, nil
}

func getKey(cc string) string {
	// Get the second we are running at.
	second := time.Now().Second()
	// combine that with the cc we set at startup to give each process 60 different keys.
	return fmt.Sprintf("%02d-%s", second, cc)
}

func getFakeData() ([]byte, error) {
	j, err := gofakeit.JSON(&gofakeit.JSONOptions{
		Type: "array",
		Fields: []gofakeit.Field{
			{Name: "id", Function: "autoincrement"},
			{Name: "first_name", Function: "firstname"},
			{Name: "last_name", Function: "lastname"},
			{Name: "address", Function: "address"},
			{Name: "animal", Function: "animal"},
			{Name: "browser", Function: "chromeuseragent"},
			{Name: "car", Function: "car"},
			{Name: "url", Function: "url"},
			{Name: "uuid", Function: "uuid"},
			{Name: "password", Function: "password", Params: map[string][]string{"special": {"false"}}},
		},
		RowCount: 30,
		Indent:   true,
	})
	return j, err
}

func healthz(c echo.Context) error {
	Redis.Send("PING")
	err := Redis.Flush()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Redis PING Error: %s", err))
	}
	return c.String(http.StatusOK, "OK")
}
