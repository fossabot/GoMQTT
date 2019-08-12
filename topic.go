package main

import (
	"errors"
	"strings"
	"sync"
)

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

// Topic Names and Topic Filters
// The MQTT v3.1.1 spec clarifies a number of ambiguities with regard
// to the validity of Topic strings.
// - A Topic must be between 1 and 65535 bytes.
// - A Topic is case sensitive.
// - A Topic may contain whitespace.
// - A Topic containing a leading forward slash is different than a Topic without.
// - A Topic may be "/" (two levels, both empty string).
// - A Topic must be UTF-8 encoded.
// - A Topic may contain any number of levels.
// - A Topic may contain an empty level (two forward slashes in a row).
// - A TopicName may not contain a wildcard.
// - A TopicFilter may only have a # (multi-level) wildcard as the last level.
// - A TopicFilter may contain any number of + (single-level) wildcards.
// - A TopicFilter with a # will match the absence of a level
//     Example:  a subscription to "foo/#" will match messages published to "foo".

func ContainsWildcard(topic string) bool {
	if len(topic) == 1 && (topic == "+" || topic == "#") {
		return true
	}
	if len(topic) > 1 && (topic[len(topic)-2:] == "/#" || topic[len(topic)-2:] == "/+") {
		return true
	}
	return strings.Contains(topic, "/+/")
}

func ValidateTopicFilter(topic string) ([]string, error) {
	if len(topic) == 0 {
		return nil, errors.New("TopicFilter cannot be empty string")
	}

	levels := strings.Split(topic, "/")
	for i, level := range levels {
		if level == "#" && i != len(levels)-1 {
			return nil, errors.New("TopicFilter contains invalid wildcard")
		}
	}
	return levels, nil
}

func ValidateTopicName(topic string) ([]string, error) {
	if len(topic) == 0 {
		return nil, errors.New("TopicName cannot be empty string")
	}
	if len(topic) < 2 || topic[0] != '/' {
		return nil, errors.New("TopicName must start with a slash")
	}

	levels := strings.Split(topic, "/")
	for _, level := range levels {
		if level == "#" || level == "+" {
			return nil, errors.New("TopicName cannot contain wildcard")
		}
	}
	return levels, nil
}
