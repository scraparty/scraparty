package routes

import (
	"fmt"
	"net/http"
	"os"

	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	clientId := os.Getenv("GITHUB_APP_ID")

	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=read:user&redirect_uri=http://%s%s/auth/callback",
		clientId,
		os.Getenv("CALLBACK_URL_HOST"),
		os.Getenv("CALLBACK_URL_PORT"),
	)

	http.Redirect(w, r, url, http.StatusFound)

	return nil
}
