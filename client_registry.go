package sarama

import (
	"expvar"
	"fmt"
	"sync"

	"github.com/IBM/sarama/mickey/boolswitch"
)

func newClientRegistry() *clientRegistry {
	return &clientRegistry{
		clients: make(map[string]*client),
		mu:      &sync.RWMutex{},
	}
}

type clientRegistry struct {
	mu      *sync.RWMutex
	clients map[string]*client
}

func (r *clientRegistry) Add(c *client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%p", c)
	r.clients[key] = c
}

func (r *clientRegistry) dumpExpvar() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topicsByClient := make(map[string]interface{})

	for key, client := range r.clients {
		if topics := client.dumpDebugVars(); topics != nil {
			topicsByClient[key] = topics
		}
	}
	return map[string]interface{}{
		"metadata_topics_by_client": topicsByClient,
	}
}

func (r *clientRegistry) resetMetadata() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		client.resetMetadata()
	}
}

var (
	globalClientsRegistry      = newClientRegistry()
	FilterMetadataTopicsSwitch boolswitch.Switch
)

func init() {
	FilterMetadataTopicsSwitch = boolswitch.NewDisabled()
	FilterMetadataTopicsSwitch.AddEnableCallback(globalClientsRegistry.resetMetadata)
	FilterMetadataTopicsSwitch.AddDisableCallback(globalClientsRegistry.resetMetadata)

	expvar.Publish("mickey_kafka", expvar.Func(func() interface{} {
		return globalClientsRegistry.dumpExpvar()
	}))
}
