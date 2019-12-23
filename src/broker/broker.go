package broker

import (
	"argovue/msg"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Broker defines broker data
type Broker struct {
	id             string
	Notifier       chan *msg.Msg
	newClients     chan chan *msg.Msg
	closingClients chan chan *msg.Msg
	cache          map[types.UID]interface{}
	clients        map[chan *msg.Msg]bool
}

func sendMsg(w http.ResponseWriter, filter string, flusher http.Flusher, m *msg.Msg) {
	mObj := m.Content.(v1.Object)
	if len(filter) > 0 && mObj.GetName() != filter {
		return
	}
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Broker can't encode message:%s", m)
		return
	}
	w.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonMsg)))
	flusher.Flush()
}

// New creates broker instance
func New(id string) (broker *Broker) {
	broker = &Broker{
		id:             id,
		Notifier:       make(chan *msg.Msg),
		newClients:     make(chan chan *msg.Msg),
		closingClients: make(chan chan *msg.Msg),
		cache:          make(map[types.UID]interface{}),
		clients:        make(map[chan *msg.Msg]bool),
	}
	go broker.listen()
	return
}

// Serve forwards events to HTTP client
func (broker *Broker) Serve(w http.ResponseWriter, filter string, flusher http.Flusher) {
	messageChan := make(chan *msg.Msg)
	broker.newClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		broker.closingClients <- messageChan
		log.Debugf("Broker: %s close connection", broker.id)
	}()

	defer func() {
		broker.closingClients <- messageChan
	}()

	for _, obj := range broker.cache {
		sendMsg(w, filter, flusher, msg.New("add", obj))
	}

	for {
		m, open := <-messageChan
		if !open {
			break
		}
		sendMsg(w, filter, flusher, m)
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
	log.Debugf("Broker: %s start", broker.id)
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			log.Debugf("Broker: %s add client, total:%d", broker.id, len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			log.Debugf("Broker: %s remove client, total:%d", broker.id, len(broker.clients))
		case msg, ok := <-broker.Notifier:
			if !ok {
				log.Debugf("Broker: %s stop", broker.id)
				return
			}
			broker.updateCache(msg)
			for clientMessageChan, _ := range broker.clients {
				select {
				case clientMessageChan <- msg:
				case <-time.After(patience):
					log.Print("Broker: %s client timeout", broker.id)
				}
			}
		}
	}
}
