package db

import "gorm.io/gorm"

type Workflow struct {
	UserId int
	Id string
	SelectorPath string
}

func GetWorkflow(workflowId string, db gorm.DB) Workflow {
	workflow := Workflow{}

  db.First(&workflow, workflowId)	

	return workflow
}

func CreateWorkflow(workflow Workflow, db gorm.DB) {
	db.Create(workflow)

	return
}
