package main

import (
	"golang-ride-sharing/services/api-gateway/grpc_clients"
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/messaging"
	"log"
	"net/http"

	pb "golang-ride-sharing/shared/proto/driver"
)


var (
	connManager = messaging.NewConnectionManager()
)

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn ,err := connManager.Upgrade(w, r)
	if err != nil {
		log.Printf("Riders WebSocket upgrade failed %v", err)
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("no userID provided in request")
		return
	}

	connManager.Add(userID, conn)
	defer connManager.Remove(userID)

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
	conn ,err := connManager.Upgrade(w, r)
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

	connManager.Add(userID, conn)

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
		log.Printf("error in dirver-service.RegisterDriver: %v", err)
		http.Error(w, "failed to register driver", http.StatusInternalServerError)
		return
	}

	defer func() {
		connManager.Remove(userID)
		driverServiceClient.Client.UnregisterDriver(r.Context(), registerDriverRequest)
		driverServiceClient.Close()

		log.Printf("driver unregistered: %s", registerDriverRequest.DriverID)
	}()

	msg := contracts.WSMessage{
		Type: contracts.DriverCmdRegister, 
		Data: registerDriverResponse.Driver,
	}

	if err := connManager.SendMessage(userID, msg); err != nil {
		log.Printf("error while sending message: %v", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message, likely client disconnected: %v", err)
			break
		}

		log.Printf("received message: %s", message)
	}
}