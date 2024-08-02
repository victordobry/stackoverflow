package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const CACHE_KEY = "value"

func main() {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, "postgres://postgres:123@db:5432/postgres")
	if err != nil {
		log.Fatalf("DB connection failed: %v\n", err)
	}
	fmt.Println("DB connection succeeded.")
	defer db.Close(ctx)
	cache := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		value, err := cache.Get(ctx, CACHE_KEY).Result()
		cacheHit := false
		switch {
		case err == redis.Nil:
			err = db.QueryRow(r.Context(), `SELECT 'Hello world!!!'`).Scan(&value)
			if err != nil {
				log.Fatalf("Failed to read value from DB: %v\n", err)
			}
		case err != nil:
			log.Fatalf("Failed to read cached value: %v\n", err)
		default:
			cacheHit = true
		}
		err = cache.Set(ctx, CACHE_KEY, value, 0).Err()
		if err != nil {
			log.Fatalf("Failed to cache value: %v\n", err)
		}
		fmt.Fprintf(w, "Value is '%s' (cacheHit=%t)", value, cacheHit)
	})
	http.ListenAndServe(":8089", nil)
}
