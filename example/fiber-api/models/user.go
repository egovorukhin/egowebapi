package models

import (
	"errors"
	"fmt"
)

var users Users

type User struct {
	Id         int        `json:"id,omitempty" jsonschema:"description=ID user" ewa:"query:desc:Id пользователя"`
	Firstname  string     `json:"firstname" jsonschema:"description=Firstname" ewa:"query:desc:Имя пользователя"`
	Lastname   string     `json:"lastname" jsonschema:"description=Lastname" ewa:"query:desc:Фамилия пользователя,array=User1&User2&User3"`
	Department Department `json:"department"`
}

type Department struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Users map[int]*User
type UserArray []*User

const (
	ModelUser       = "Aw.User"
	ModelDepartment = "Aw.Department"
)

func GetUser(id int) *User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func GetUsers() Users {
	if users == nil {
		return nil
	}
	return users
}

func (u User) Set() {
	if users == nil {
		users = map[int]*User{}
	}
	users[u.Id] = &u
}

func (u User) Update(id int) error {
	if users == nil {
		return nil
	}
	if _, ok := users[id]; ok {
		users[id] = &u
		return nil
	}
	return errors.New(fmt.Sprintf("Запись не найдена - %d", id))
}

func (u User) Delete() {
	if users == nil {
		return
	}
	delete(users, u.Id)
}
