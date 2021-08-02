package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func fillUsers(w http.ResponseWriter, r *http.Request) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	user := User{FirstName: "Jinzhu", LastName: "Bag"}

	result := db.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		fmt.Println(result.Error) // returns error
	}

	fmt.Printf("ID: %v, Rows affected: %v", user.ID, result.RowsAffected)
}

func routes() {
	http.HandleFunc("/users/fill", fillUsers)
	http.HandleFunc("/posts/fill", fillPosts)
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/post/", getPost) // to use route params (/post/:id) use MUX or CHI package
	log.Fatal(http.ListenAndServe(":10000", nil))
}
