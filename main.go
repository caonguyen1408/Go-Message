package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	Cronjob()
	// Routes
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", hola)
	e.GET("/getmessage", GetMessage)
	e.POST("/sendmessage", SendMessage)
	// Start server at localhost:1323
	e.Logger.Fatal(e.Start(":8080"))

}
func hola(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to APIs")
}
func Cronjob() {
	go func() {
		Worker()
	}()
}

//GET MESSAGE
func GetMessage(_ echo.Context) error {
	var message Message
	c := RedisConnect()
	defer c.Close()
	key, _ := redis.String(c.Do("GET", "message"))
	message.Message = key
	return echo.NewHTTPError(http.StatusOK, message)
}

//SEND MESSAGE
func SendMessage(c echo.Context) error {
	var message Message
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	CallQueue(message.Message)
	return echo.NewHTTPError(http.StatusOK, message)

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err.Error())
	}
}
