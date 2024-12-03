package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	fiberHTML "github.com/gofiber/template/html/v2"
	"github.com/nats-io/nats.go"
)

var NATS_CONN *nats.Conn

const (
	APP_NAME         = "Streaming Coordinate with NATS.io"
	SUBJECT_LOCATION = "location"
)

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func init() {
	nc, err := nats.Connect("nats://127.0.0.1:4223")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	NATS_CONN = nc

	log.Println("Connected to NATS!")
}

func main() {
	waits := make(chan int)

	engine := fiberHTML.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName: APP_NAME,
		Views:   engine,
	})

	app.Get("/", HomePage)
	app.Get("/coord", Websocket, websocket.New(TrackCoordinate, websocket.Config{
		HandshakeTimeout: 100 * time.Second,
		ReadBufferSize:   1824,
		WriteBufferSize:  256,
	}))

	go Shutdown()

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Println("failed to start http server")
			CloseServices()
			os.Exit(1)
		}
	}()

	NATS_CONN.Subscribe(SUBJECT_LOCATION, func(msg *nats.Msg) {
		var coord Coordinate
		if err := json.Unmarshal(msg.Data, &coord); err != nil {
			log.Println("Failed to convert data to desired type from subject "+SUBJECT_LOCATION+" :", err)
			return
		}

		fmt.Println("Coordinate X: ", coord.X)
		fmt.Println("Coordinate Y: ", coord.Y)
		fmt.Println("+======================+")
	})

	<-waits
}

func HomePage(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Render("index", fiber.Map{
		"Title": APP_NAME,
	})
}

func Websocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func TrackCoordinate(conn *websocket.Conn) {
	var (
		msg []byte
		err error
	)
	for {
		if _, msg, err = conn.ReadMessage(); err != nil {
			log.Println("Read :", err)
			break
		}

		if err := NATS_CONN.Publish(SUBJECT_LOCATION, msg); err != nil {
			log.Println("Error when publishing coordinate:", err)
		}
	}
}
