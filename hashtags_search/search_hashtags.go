package main

import (
	"encoding/json"
	"fmt"
	getdatabasecontroller "hashtags_search/get_database_controller"
	basecontroller "hashtags_search/modules/base_controller"
	"hashtags_search/modules/database"

	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
	Gcontrl  *getdatabasecontroller.GetDatabaseController
}

type Message struct {
	Query string `json:"query"`
}

func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	wsh.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := wsh.upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("An error occured: %s", err)
	}

	for {

		messageType, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message %s", err)
			break
		}

		// Разбор сообщения из JSON
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Printf("Error parsing message JSON: %s", err)
			continue
		}

		// Преобразуем UserId в целое число
		Hashtag := msg.Query

		jsonData, err := wsh.Gcontrl.GetHashtags(Hashtag)

		if err != nil {
			fmt.Printf("Error fetching users by ID: %s", err)
			return
		}
		// Process jsonData
		fmt.Println(jsonData)

		// Send the JSON response back to the client
		err = c.WriteMessage(messageType, jsonData)
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
			break
		}
	}

	defer c.Close()
}

func main() {

	DBPool := &database.DatabasePoolController{}
	DBPool.StartupEvent()

	BContrl := basecontroller.BaseController{Database: DBPool.GetDb()}
	GContrl := getdatabasecontroller.GetDatabaseController{BaseController: &BContrl}

	wsHandler := WebSocketHandler{
		upgrader: websocket.Upgrader{},
		Gcontrl:  &GContrl,
	}

	http.Handle("/", &wsHandler)
	fmt.Println("Запуск сервера на http://localhost:8380")

	server := &http.Server{Addr: "localhost:8380"}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибкка при запуске сервера: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	DBPool.ShutdownEvent()

	// Корректное завершение работы сервера
	if err := server.Close(); err != nil {
		fmt.Printf("Ошибка при завершении работы сервера: %v", err)
	}

	fmt.Println("Сервер был корректно завершен.")
}
