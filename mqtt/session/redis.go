package session

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

type redisProvider struct {
	redisClient *redis.Client
	sessionName string
}

// NewRedisProvider creates a Redis-based session storage provider using given
// connection parameters and session name
func NewRedisProvider(session string, options *redis.Options) *SessionsProvider {

	provider := SessionsProvider(&redisProvider{
		redisClient: redis.NewClient(options),
		sessionName: session,
	})
	return &provider
}

func (r *redisProvider) New(id string) (*Session, error) {
	val, _ := jsoniter.Marshal(&Session{id: id})

	err := r.redisClient.HSet(r.sessionName, id, val).Err()
	if err != nil {
		return nil, err
	}

	result, err := r.redisClient.HGet(r.sessionName, id).Bytes()
	if err != nil {
		return nil, err
	}

	sess := Session{}
	err = jsoniter.Unmarshal(result, &sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func (r *redisProvider) Get(id string) (*Session, error) {

	result, err := r.redisClient.HGet(r.sessionName, id).Bytes()
	if err != nil {
		return nil, err
	}

	sess := Session{}
	err = jsoniter.Unmarshal(result, &sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func (r *redisProvider) Delete(id string) {
	r.redisClient.HDel(r.sessionName, id)
}

func (r *redisProvider) Save(id string) error {
	return nil
}

func (r *redisProvider) Count() int {
	return int(r.redisClient.HLen(r.sessionName).Val())
}

func (r *redisProvider) Close() error {
	return r.redisClient.Del(r.sessionName).Err()
}
