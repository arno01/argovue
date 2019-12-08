package broker

import (
	"encoding/json"
	"fmt"
	"kubevue/msg"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Broker defines broker data
type Broker struct {
	Notifier       chan *msg.Msg
	newClients     chan chan *msg.Msg
	closingClients chan chan *msg.Msg
	cache          map[types.UID]interface{}
	clients        map[chan *msg.Msg]bool
}

func sendMsg(rw http.ResponseWriter, flusher http.Flusher, m *msg.Msg) {
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Can't encode message:%s", m)
		return
	}
	rw.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonMsg)))
	flusher.Flush()
}

// New creates broker instance
func New() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan *msg.Msg, 1),
		newClients:     make(chan chan *msg.Msg),
		closingClients: make(chan chan *msg.Msg),
		cache:          make(map[types.UID]interface{}),
		clients:        make(map[chan *msg.Msg]bool),
	}
	go broker.listen()
	return
}

// Serve forwards events to HTTP client
func (broker *Broker) Serve(rw http.ResponseWriter, flusher http.Flusher) {
	messageChan := make(chan *msg.Msg)
	broker.newClients <- messageChan

	notify := rw.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		broker.closingClients <- messageChan
		log.Debugf("Closing connection")
	}()

	defer func() {
		broker.closingClients <- messageChan
	}()

	for _, obj := range broker.cache {
		sendMsg(rw, flusher, msg.New("add", obj))
	}

	for {
		m, open := <-messageChan
		if !open {
			break
		}
		sendMsg(rw, flusher, m)
	}

}

func (broker *Broker) updateCache(m *msg.Msg) {
	mObj := m.Content.(v1.Object)
	if m.Action == "delete" {
		delete(broker.cache, mObj.GetUID())
	} else {
		broker.cache[mObj.GetUID()] = mObj
	}
}

const patience time.Duration = time.Second * 1

func (broker *Broker) listen() {
	log.Debugf("Starting message broker")
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			log.Debugf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			log.Debugf("Removed client. %d registered clients", len(broker.clients))
		case msg := <-broker.Notifier:
			broker.updateCache(msg)
			for clientMessageChan, _ := range broker.clients {
				select {
				case clientMessageChan <- msg:
				case <-time.After(patience):
					log.Print("Skipping client.")
				}
			}
		}
	}
}
