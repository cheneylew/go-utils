package main

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
	Password string
	Email string
	Birthday string
}
