package main

import (
	pb "golang-ride-sharing/shared/proto/driver"
	"golang-ride-sharing/shared/util"
	"math/rand"
	math "math/rand/v2"
	"sync"

	"github.com/mmcloughlin/geohash"
)


type driverInMap struct {
	Driver *pb.Driver
}

type DriverService struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewDriverService() *DriverService {
	return &DriverService{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *DriverService) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id: driverId,
		Geohash:  geohash,
		Location: &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:     "Charles Leclerc",
		PackageSlug:    packageSlug,
		ProfilePicture: util.GetRandomAvatar(rand.Intn(9) + 1),
		CarPlate:       GenerateRandomPlate(),
	}

	s.drivers = append(s.drivers, &driverInMap{Driver: driver})

	return driver, nil
}

func (s *DriverService) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

    for i, d := range s.drivers {
        if d.Driver.Id != driverId {
            s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
        }
    }
}