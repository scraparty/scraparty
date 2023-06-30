package db

import "gorm.io/gorm"

type User struct {
	GithubId int
	Workflows []string
}

func CreateUser(githubId int, db gorm.DB) {
	user := User{}

	user.GithubId = githubId
	user.Workflows = []string{}

	db.Create(user)

	return
}
