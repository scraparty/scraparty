package db

import "gorm.io/gorm"

type User struct {
	GithubId int
	Workflows []string
}

func GetUser(githubId int, db gorm.DB) (User) {
	user := User{}

	db.First(&user, githubId)

	return user
}
