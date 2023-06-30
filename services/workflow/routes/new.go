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

type WorkflowObject struct {
	SelectorPath string `json:"selector_path"`
}

func New(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
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

	workflowObject := WorkflowObject{}

	bytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "there was an error creating your workflow", http.StatusInternalServerError)

		return err
	}

	err = json.Unmarshal(bytes, &workflowObject)

	if err != nil {
		http.Error(w, "there was an error creating your workflow", http.StatusInternalServerError)

		return err
	}

	workflow := scrapartyDb.Workflow{
		UserId: sess.Session["github_id"].(int),
		Id: uuid.NewString(),
		SelectorPath: workflowObject.SelectorPath,
	}

	scrapartyDb.CreateWorkflow(workflow, db)

	w.WriteHeader(200)
	w.Write([]byte("your workflow was created successfully"))

	return nil
}
