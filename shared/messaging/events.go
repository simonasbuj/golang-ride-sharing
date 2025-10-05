package messaging

import (
	pb "golang-ride-sharing/shared/proto/trip"
)

const (
	FindAvailableDriversQueue = "find_available_drivers"
	DriverCmdTripRequestQueue = "driver_cmd_trip_request"
	DriverTripResponseQueue = "driver_trip_response_queue"
	NotifyRiderNoDriversFoundQueue = "notify_rider_no_drivers_found"
)

type TripEventData struct {
	Trip *pb.Trip	`json:"trip"`
}
