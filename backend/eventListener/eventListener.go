package eventListener

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

// ======= PUBLIC ========
type channelType struct {
	id      string
	pubSub  *redis.PubSub
	clients map[string]func(data string)
	lock    sync.RWMutex
}

func Subscribe(channel string, key string, fn func(data string)) {
	rdb := GetPubSub()

	sub, ok := getSubscription(channel)
	if !ok {
		sub = new(channelType)
		sub.clients = make(map[string]func(data string))
		setSubscription(channel, sub)
		sub.pubSub = rdb.Subscribe(context.Background(), channel)
		go listen(sub)
	}
	sub.clients[key] = fn
}

func Unsubscribe(channel string, key string) {
	sub, ok := getSubscription(channel)
	if !ok {
		return
	}
	sub.lock.Lock()
	delete(sub.clients, key)

	if len(sub.clients) == 0 {
		sub.close()
		deleteSubscription(channel)
	}
	sub.lock.Unlock()
}

// ==== PRIVATE ======

var subsLock = sync.RWMutex{}
var subs = make(map[string]*channelType)

func (s *channelType) close() {
	s.pubSub.Close()

	subsLock.Lock()
	defer subsLock.Unlock()
	delete(subs, s.id)
}

func getSubscription(channel string) (*channelType, bool) {
	subsLock.RLock()
	defer subsLock.RUnlock()

	sub, ok := subs[channel]
	return sub, ok

}

func setSubscription(channel string, sub *channelType) {
	subsLock.Lock()
	defer subsLock.Unlock()

	subs[channel] = sub
}

func deleteSubscription(channel string) {
	subsLock.Lock()
	defer subsLock.Unlock()
	delete(subs, channel)
}

func listen(sub *channelType) {
	for msg := range sub.pubSub.Channel() {
		sub.lock.RLock()
		for _, handler := range sub.clients {
			handler(msg.Payload)
		}
		sub.lock.RUnlock()
	}
}
