package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pipe-scaling-analyzer/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client 			*redis.Client
	ttl    			time.Duration
	logger 			*zap.SugaredLogger
	
}

func NewRedisCache (addr, password string, logger *zap.SugaredLogger) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     	addr,
		Password: 	password,
		DB:       	0,
	})
	
	return &RedisCache{
		client: 	rdb,
		ttl:    	10 * time.Minute,
		logger: 	logger,
	}
}
func (c *RedisCache) GetPipeStatus(ctx context.Context, pipeId string) (*domain.PipeStatusDTO, error) {
	key := fmt.Sprintf("pipe:status:%s", pipeId)
	fmt.Println("ищем в редисе ключ:", key)
	
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return  nil, nil
		}
		c.logger.Errorf("redis get error %v",err)
		return nil, err
		
	}


	status := domain.PipeStatusDTO{}
	if err := json.Unmarshal([]byte(val), &status); err != nil {
		return nil, err
	}
	
	return &status, nil
}
func (c *RedisCache) SetPipeStatus(ctx context.Context, status domain.PipeStatusDTO) error {
	key := fmt.Sprintf("pipe:status:%s", status.PipeId)
	
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	
	
	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		c.logger.Errorf("redis set error %v",err)
		
		return err
	}
	fmt.Println("записали в редис:", key)
	return  nil
}

func (c *RedisCache) DeletePipeStatus(ctx context.Context, pipeId string) error {
	key := fmt.Sprintf("pipe:status:%s", pipeId)
	
	return c.client.Del(ctx, key).Err()
}
