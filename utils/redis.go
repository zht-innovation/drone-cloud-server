package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func GetRedisConn(ctx context.Context) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection build error: %v", err)
	}

	return &RedisClient{client: rdb}
}

func (r *RedisClient) PubChannel(ctx context.Context, channelName, msg string) error {
	err := r.client.Publish(ctx, channelName, msg).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) SubChannel(ctx context.Context, channelName string, msgChan chan<- string) error {
	pubsub := r.client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}

		msgChan <- msg.Payload
	}
}
