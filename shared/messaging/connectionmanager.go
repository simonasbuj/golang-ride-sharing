package messaging

import (
	"errors"
	"golang-ride-sharing/shared/contracts"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)


var (
	ErrConnectionNotFound = errors.New("connection not found")
)


type connWrapper struct {
	conn 	*websocket.Conn
	mutex	sync.Mutex
}

type ConnectionManager struct {
	connectons		map[string]*connWrapper
	mutex  			sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for now
	},
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connectons: make(map[string]*connWrapper),
	}
}

func (cm *ConnectionManager) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (cm *ConnectionManager) Add(id string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.connectons[id] = &connWrapper{
		conn: conn,
		mutex: sync.Mutex{},
	}

	log.Printf("added connection for user %s", id)
}

func (cm *ConnectionManager) Remove(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.connectons, id)
}

func (cm *ConnectionManager) Get(id string) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	wrapper, exists := cm.connectons[id]
	if !exists {
		return nil, false
	}

	return wrapper.conn, true
}

func (cm *ConnectionManager) SendMessage(id string, message contracts.WSMessage) error {
	cm.mutex.RLock()
	wrapper, exists := cm.connectons[id]
	cm.mutex.RUnlock()

	if !exists {
		return ErrConnectionNotFound
	}

	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()

	return wrapper.conn.WriteJSON(message)
}
