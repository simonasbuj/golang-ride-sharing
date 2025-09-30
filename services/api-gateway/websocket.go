package main

import (
	"golang-ride-sharing/services/api-gateway/grpc_clients"
	"golang-ride-sharing/shared/contracts"
	"log"
	"net/http"

	pb "golang-ride-sharing/shared/proto/driver"

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

	driverServiceClient, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer driverServiceClient.Close()

	registerDriverRequest := &pb.RegisterDriverRequest{
		DriverID: userID,
		PackageSlug: packageSlug,
	}

	registerDriverResponse, err := driverServiceClient.Client.RegisterDriver(r.Context(), registerDriverRequest)
	if err != nil {
		log.Printf("error in trip-service.PreviewTrip: %v", err)
		http.Error(w, "failed to preview a trip", http.StatusInternalServerError)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: registerDriverResponse.Driver,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error while sending message: %v", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message, likely client disconnected: %v", err)
			driverServiceClient.Client.UnregisterDriver(r.Context(), registerDriverRequest)
			break
		}

		log.Printf("received message: %s", message)
	}
}