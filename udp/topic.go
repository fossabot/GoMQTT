package udp

import "sync"

// This needs to be efficient, but I doubt it is.
// Raise my salary and maybe I'll fix it.
type topicNames struct {
	sync.RWMutex
	contents map[uint16]string
	next     uint16
}

// O(n)
func (repo *topicNames) containsTopic(topic string) bool {
	return repo.getId(topic) != 0
}

// O(1)
func (repo *topicNames) containsId(id uint16) bool {
	return repo.getTopic(id) != ""
}

// O(n)
func (repo *topicNames) getId(topic string) uint16 {
	defer repo.RUnlock()
	repo.RLock()
	var topicid uint16
	for id, topicVal := range repo.contents {
		if topicVal == topic {
			topicid = id
			break
		}
	}
	return topicid
}

// O(1)
func (repo *topicNames) getTopic(id uint16) string {
	defer repo.RUnlock()
	repo.RLock()
	topic := repo.contents[id]
	return topic
}

// O(1)
func (repo *topicNames) putTopic(topic string) uint16 {
	defer repo.Unlock()
	repo.Lock()
	repo.next++
	repo.contents[repo.next] = topic
	return repo.next
}
