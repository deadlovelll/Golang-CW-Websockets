package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	getdatabasecontroller "places_search/get_database_controller"
	basecontroller "places_search/modules/base_controller"
	"places_search/modules/database"
	"syscall"

	"strings"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
	GContrl  *getdatabasecontroller.GetDatabaseController
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
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message %s", err)
			break
		}

		// Разбор сообщения из JSON
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error parsing message JSON: %s", err)
			continue
		}

		// Преобразуем UserId в целое число
		PlaceName := msg.Query

		if !strings.HasPrefix(PlaceName, "#") {

			// Если преобразование успешно, то это целое число
			jsonData, err := wsh.GContrl.GetPlaceByName(PlaceName)
			if err != nil {
				log.Printf("Error fetching users by ID: %s", err)
				return
			}
			// Process jsonData
			fmt.Println(jsonData)

			// Send the JSON response back to the client
			err = c.WriteMessage(messageType, jsonData)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}

			log.Printf("Получено сообщение %s", message)

		} else {

			// Если преобразование успешно, то это целое число
			jsonData, err := wsh.GContrl.GetPlaceWithHashtag(PlaceName)
			if err != nil {
				log.Printf("Error fetching users by ID: %s", err)
				return
			}
			// Process jsonData
			fmt.Println(jsonData)

			// Send the JSON response back to the client
			err = c.WriteMessage(messageType, jsonData)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}

			log.Printf("Получено сообщение %s", message)
		}
	}

	defer c.Close()
}

func main() {

	// Инициализация пула базы данных
	DBPool := &database.DatabasePoolController{}
	DBPool.StartupEvent()

	BContrl := basecontroller.BaseController{Database: DBPool.GetDb()}
	GContrl := getdatabasecontroller.GetDatabaseController{BaseController: &BContrl}

	// Настройка WebSocket обработчика
	wsHandler := &WebSocketHandler{
		upgrader: websocket.Upgrader{},
		GContrl:  &GContrl,
	}

	http.Handle("/", wsHandler)
	fmt.Println("Запуск сервера на http://localhost:8285")

	// Запуск HTTP-сервера в отдельной горутине
	server := &http.Server{Addr: "localhost:8285"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Захват сигналов завершения работы (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала завершения работы
	<-quit

	// Вызов функции завершения работы
	DBPool.ShutdownEvent()

	// Корректное завершение работы сервера
	if err := server.Close(); err != nil {
		fmt.Printf("Ошибка при завершении работы сервера: %v", err)
	}

	fmt.Println("Сервер был корректно завершен.")
}
