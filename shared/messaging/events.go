package messaging

import (
	pb "golang-ride-sharing/shared/proto/trip"
	pbd "golang-ride-sharing/shared/proto/driver"
)

const (
	FindAvailableDriversQueue 		= "find_available_drivers"
	DriverCmdTripRequestQueue 		= "driver_cmd_trip_request"
	DriverTripResponseQueue 		= "driver_trip_response_queue"
	NotifyRiderNoDriversFoundQueue 	= "notify_rider_no_drivers_found"
	NotifyDriverAssignedQueue 		= "notify_driver_assigned"
)

type TripEventData struct {
	Trip *pb.Trip	`json:"trip"`
}

type DriverTripResponseData struct {
	Driver 		*pbd.Driver	`json:"driver"`
	TripID 		string		`json:"tripID"`
	RiderID 	string		`json:"riderID"`
}
