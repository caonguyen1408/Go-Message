package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	mux "github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/streadway/amqp"
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
func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", "192.168.2.5:6379")
	failOnError(err, "Fail to connect db")
	return c
}
func Index(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	fmt.Fprintf(w, "<h1 style=\"font-family: Helvetica;\">Hello, welcome to blog service</h1>")
}

type Message struct {
	Message string `json:"message"`
}
type QueueMessage struct {
	Message string `json:"message"`
	TaskId  int64  `json:"task_id"`
}

func Worker() {
	conn, err := amqp.Dial("amqp://user:password@192.168.2.3:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Queue Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
			c1 := RedisConnect()
			defer c1.Close()
			setmes, err := c1.Do("SET", "message", d.Body)
			log.Printf("Queue Add to Database ")
			failOnError(err, "Fail to add message")
			log.Printf("%s", setmes)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func CallQueue(mes_string string) string {
	conn, err := amqp.Dial("amqp://user:password@192.168.2.3:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := mes_string
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Queue Sent %s", body)
	return string(body)
}

