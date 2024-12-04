package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2"
	fiberHTML "github.com/gofiber/template/html/v2"

	"github.com/gofiber/contrib/websocket"
)

var NATS_CONN *nats.Conn

const APP_NAME = "Chat App with NATS.io"
const SUBJECT_CHAT = "chat"

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

	wsConf := websocket.Config{
		HandshakeTimeout: 100 * time.Second,
		ReadBufferSize:   1824,
		WriteBufferSize:  256,
	}

	app.Get("/", HomePage)
	app.Get("/chat/send", Websocket, websocket.New(ChatSend, wsConf))
	app.Get("/chat/receive", ChatReceive)

	go Shutdown()

	go func() {
		if err := app.Listen(":3001"); err != nil {
			log.Println("failed to start http server")
			CloseServices()
			os.Exit(1)
		}
	}()

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

func ChatSend(conn *websocket.Conn) {
	var (
		msg []byte
		err error
	)
	for {
		_, msg, err = conn.ReadMessage()

		if err != nil {
			log.Println("Read :", err)
			break
		}

		natsMsg, err := NATS_CONN.Request(SUBJECT_CHAT, msg, 10*time.Second)
		if err != nil {
			log.Println("failed to send request:", err)
		}

		fmt.Println("Reply:", string(natsMsg.Data))
	}
}

func ChatReceive(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		NATS_CONN.Subscribe(SUBJECT_CHAT, func(msgNats *nats.Msg) {
			log.Println("message:", string(msgNats.Data))

			fmt.Fprintf(w, "data: "+string(msgNats.Data)+"\n\n")
			if err := w.Flush(); err != nil {
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
			}

			msgNats.Respond([]byte("received"))
		})

		for {
		}
	}))

	return nil
}
