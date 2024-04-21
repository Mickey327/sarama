package sarama

import (
	"sync"
	"sync/atomic"
)

type UsedTopicRegistry interface {
	Add(topic string)
	Has(topic string) bool
	All() []string
	Len() int
}

type usedTopicsRegistry struct {
	mu     *sync.RWMutex
	topics map[string]struct{}
	length int32
}

func newTopicRegistry() UsedTopicRegistry {
	return &usedTopicsRegistry{
		topics: make(map[string]struct{}),
		mu:     &sync.RWMutex{},
	}
}

func (r *usedTopicsRegistry) Len() int {
	return int(atomic.LoadInt32(&r.length))
}

func (r *usedTopicsRegistry) Add(topic string) {
	r.mu.Lock()
	r.topics[topic] = struct{}{}
	atomic.AddInt32(&r.length, 1)
	r.mu.Unlock()
}

func (r *usedTopicsRegistry) Has(topic string) bool {
	r.mu.RLock()
	_, ok := r.topics[topic]
	r.mu.RUnlock()
	return ok
}

func (r *usedTopicsRegistry) All() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	topics := make([]string, 0, len(r.topics))
	for topic := range r.topics {
		topics = append(topics, topic)
	}
	return topics
}

var (
	GlobalUsedTopicRegistry = newTopicRegistry()
)
