package main


import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",              
		DB:       0,                
	})

	// PING
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis PING:", pong)

	// SET a key
	err = rdb.Set(ctx, "username", "sibi", 10*time.Minute).Err()
	if err != nil {
		log.Fatalf("Failed to SET key: %v", err)
	}
	fmt.Println("Key 'username' set successfully")

	// GET the key
	val, err := rdb.Get(ctx, "username").Result()
	if err != nil {
		log.Fatalf("Failed to GET key: %v", err)
	}
	fmt.Println("GET 'username':", val)

	// Delete the key
	err = rdb.Del(ctx, "username").Err()
	if err != nil {
		log.Fatalf("Failed to DELETE key: %v", err)
	}
	fmt.Println("Key 'username' deleted")
}
