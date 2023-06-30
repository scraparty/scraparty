package db

import "gorm.io/gorm"

type Workflow struct {
	UserId int
	Id string
	SelectorPath string
}

func CreateWorkflow(workflow Workflow, db gorm.DB) {
	db.Create(workflow)

	return
}
