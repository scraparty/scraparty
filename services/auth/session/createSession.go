package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	Id string `json:"cookie"`
	Session map[string]interface{} `json:"session"`
}

func CreateSession(id string, accessToken string, githubId int) error {
	ctx := context.Background()
	 
	ip := os.Getenv("SESSION_STORAGE_SERVICE_HOST")
	port := os.Getenv("SESSION_STORAGE_SERVICE_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", ip, port),
		Password: "",
		DB: 0,
	})

	session := Session{
		Session: map[string]interface{}{
			"access_token": accessToken,
			"github_id": githubId,
		},
	}

	stringified, err := json.Marshal(session)

	if err != nil {
		return err
	}

	err = rdb.Set(ctx, id, string(stringified), 0).Err()

	if err != nil {
		return err
	}

	return nil
}
