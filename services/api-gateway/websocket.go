package main

import (
	"encoding/json"
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

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request, rabbitmq *messaging.RabbitMQ) {
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

	// init queue consumers
	queues := []string{
		messaging.NotifyRiderNoDriversFoundQueue,
		messaging.NotifyDriverAssignedQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(rabbitmq, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("failed to start consumer for queue: %s, error: %v", q, err)
		}
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

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request, rabbitmq *messaging.RabbitMQ) {
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

	ctx := r.Context()

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

	// init queue consumers
	queues := []string{
		messaging.DriverCmdTripRequestQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(rabbitmq, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("failed to start consumer for queue: %s, error: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message, likely client disconnected: %v", err)
			break
		}
		log.Printf("received message: %s", message)

		type driverMessage struct {
			Type string 			`json:"type"`
			Data json.RawMessage	`json:"data"`
		}

		var driverMsg driverMessage
		if err := json.Unmarshal(message, &driverMsg); err != nil {
			log.Printf("ERROR: failed to unmarshal message: %v, error: %v", message, err)
			continue
		}

		// handle different message types
		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			// TODO: hanlde driver.cmd.location
			continue
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			// forward msg to RabbitMQ
			if err := rabbitmq.PublishMessage(ctx, driverMsg.Type, contracts.AmqpMessage{
				OwnerID: userID,
				Data: driverMsg.Data,
			}); err != nil {
				log.Printf("ERROR: failed to publish message: %v, error: %v", driverMsg.Data, err)
				continue
			}
		default:
			log.Printf("ERROR uknown messge type: %s", driverMsg.Type)
		}


	}
}