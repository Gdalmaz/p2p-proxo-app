package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func ConnectRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(pong, err)
		return err
	}

	return nil
}
func SaveClickCountToRedis(foodID uint, clickCount int) error {
	ctx := context.Background()

	if rdb == nil {
		return errors.New("Redis client is not initialized")
	}

	key := fmt.Sprintf("food:%d:clicks", foodID)
	value := strconv.Itoa(clickCount)

	err := rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		log.Printf("Redis'e kaydedilemedi: %v", err)
		return err
	}

	return nil
}
func GetAllClickCountsFromRedis() (map[string]int, error) {
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "food:*:clicks", 0).Iterator()
	result := make(map[string]int)

	for iter.Next(ctx) {
		key := iter.Val()
		value, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		clickCount, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		result[key] = clickCount
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteClickCountFromRedis(foodID uint) error {
	ctx := context.Background()

	if rdb == nil {
		return errors.New("Redis client is not initialized")
	}

	key := fmt.Sprintf("food:%d:clicks", foodID)

	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
