package main

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title           string
	Content         string
	Summary         string
	FeaturedImage   string
	Slug            string
	MetaTitle       string
	MetaDescription string
	UserId          uint
	User            User
}
