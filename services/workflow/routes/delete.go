package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	scrapartyDb "github.com/scraparty/scraparty/workflow/db"
	"github.com/scraparty/scraparty/workflow/session"
	"gorm.io/gorm"
)

func Delete(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	sid, err := r.Cookie("scraparty")

	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return err
	}

	sess, err := session.GetSession(sid.Value)	

	if errors.Is(err, session.ErrNotAuthenticated) {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return err
	}

	if err != nil {
		http.Error(w, "there was an error authenticating you", http.StatusInternalServerError)

		return err
	}

	workflowId := r.URL.Query().Get("workflow_id")

	w.WriteHeader(200)
	w.Write([]byte("your workflow was created successfully"))

	return nil
}
