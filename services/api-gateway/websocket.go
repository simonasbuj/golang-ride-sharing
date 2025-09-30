package main

import (
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/util"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn ,err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Riders WebSocket upgrade failed %v", err)
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("no userID provided in request")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		log.Printf("received message: %s", message)
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn ,err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Drivers WebSocket upgrade failed %v", err)
	}
	defer conn.Close()
	
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("no userID provided in request")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("no packageSlug provided in request")
		return
	}

	type Driver struct {
		ID 				string 	`jsong:"id"`
		Name 			string	`jsong:"name"`
		ProfilePicture	string	`jsong:"profilePicture"`
		CarPlate		string	`jsong:"carPlate"`
		PackageSlug		string	`jsong:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID: userID,
			Name: "sim sim",
			ProfilePicture: util.GetRandomAvatar(3),
			CarPlate: "LLP831",
			PackageSlug: packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error while sending message: %v", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		log.Printf("received message: %s", message)
	}
}