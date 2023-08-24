package entity

import "fmt"

type User struct {
	ID       int
	Name     string
	Password string
	Email    string
}

func (u User) print() {
	fmt.Println("Name: ", u.Name, "Email: ", u.Email, "Password: ", u.Password)
}
