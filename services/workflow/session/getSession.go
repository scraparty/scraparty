package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	Id string `json:"cookie"`
	Session map[string]interface{} `json:"session"`
}

var ErrNotAuthenticated = errors.New("this user is not authenticated")

func GetSession(id string) (Session, error) {
	ctx := context.Background()

	ip := os.Getenv("SESSION_STORAGE_SERVICE_HOST")
	port := os.Getenv("SESSION_STORAGE_SERVICE_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", ip, port),
		Password: "",
		DB: 0,
	})

	result, err := rdb.Get(ctx, id).Result()

	if err == redis.Nil {
		return Session{}, ErrNotAuthenticated
	}

	if err != nil { 
		return Session{}, err
	}

	session := Session{} 

	err = json.Unmarshal([]byte(result), &session) 	

	if err != nil {
		return Session{}, nil
	}
	
	return session, nil
}
