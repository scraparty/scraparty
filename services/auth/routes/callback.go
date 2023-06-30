package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/scraparty/scraparty/auth/session"
	scrapartyDb "github.com/scraparty/scraparty/auth/db"
	"gorm.io/gorm"
)

type GithubResponse struct {
	AccessToken string `json:"access_token"`
}

type GithubUser struct {
	Id int `json:"id"`
}

func Callback(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	code := r.URL.Query().Get("code")

	token, err := getToken(code)

	if err != nil {
		return err
	}

	sid := uuid.NewString() 

	cookie := http.Cookie{
		Name: "scraparty",
		Value: sid,
		Domain: os.Getenv("COOKIE_DOMAIN"),
		Path: "/",
	}

	http.SetCookie(w, &cookie)

	user, err := getUser(token)

	if err != nil {
		return err
	}

	err = session.CreateSession(sid, token, user.Id)

	if err != nil {
		return err
	}

	scrapartyUser := scrapartyDb.User{}

	db.First(&scrapartyUser, user.Id)

	if scrapartyUser.GithubId != 0 {
		w.WriteHeader(200)
		w.Write([]byte("you have been logged in successfully"))

		return nil
	}

	scrapartyDb.CreateUser(user.Id, db)

	w.WriteHeader(200)
	w.Write([]byte("your account has been created successfully"))

	return nil
}

func getUser(token string) (GithubUser, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)

	if err != nil {
		return GithubUser{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)

	if err != nil {
		return GithubUser{}, err
	}

	user := GithubUser{}

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return GithubUser{}, err
	}

	err = json.Unmarshal(bytes, &user)

	if err != nil {
		return GithubUser{}, err
	}

	return user, err
}

func getToken(code string) (string, error) {
	client := http.Client{}

	clientId := os.Getenv("GITHUB_APP_ID")
	clientSecret := os.Getenv("GITHUB_APP_SECRET")

	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", 
		clientId,
		clientSecret,
		code,
	)

	req, err := http.NewRequest(
		"POST",
		url,
		nil,
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	githubResponse := GithubResponse{}

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, &githubResponse)

	if err != nil {
		return "", err
	}

	return githubResponse.AccessToken, nil
}
