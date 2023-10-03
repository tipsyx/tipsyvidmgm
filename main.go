package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/streadway/amqp"

	"github.com/tipsyx/tipsyvidmgm/config"
	"github.com/tipsyx/tipsyvidmgm/internal/handlers"
	"github.com/tipsyx/tipsyvidmgm/internal/worker"
)

var ch *amqp.Channel
var db *gorm.DB

type Video struct {
	gorm.Model
	Filename     string
	Transcription string
}

func main() {
	r := gin.Default()

	var err error
	db, err = gorm.Open("mysql", config.DatabaseDSN)
	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	defer db.Close()
	db.AutoMigrate(&Video{})

	os.MkdirAll("./uploads", os.ModePerm)

	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v\n", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"videomgmex", // Exchange name
		"routekey",   // Routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello, RabbitMQ!"),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	msgs, err := ch.Consume(
		"queueA", // Queue name
		"",       // Consumer name
		true,
		false,
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	transcriptionWorker := worker.NewTranscriptionWorker(ch, db)
	go transcriptionWorker.Start()

	r.POST("/upload", func(c *gin.Context) {
		handlers.HandleUpload(c, ch, db)
	})

	r.GET("/playback/:id", func(c *gin.Context) {
		handlers.VideoPlayback(c, db)
	})

	r.GET("/listvideos", func(c *gin.Context) {
		handlers.ListUploadedVideos(c, db)
	})

	r.GET("/videodetails/:id", func(c *gin.Context) {
		handlers.ShowVideoDetails(c, db)
	})

	r.POST("/deletevideo", func(c *gin.Context) {
		handlers.DeleteVideoByID(c, db)
	})
	r.GET("/getvideo/:id", func(c *gin.Context){
		handlers.GetVideoByID(c, db)
	})
	fmt.Println("Server is listening on :8080")
	r.Run(":8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down...")
	db.Close()

	fmt.Println("Server has shut down.")
}
