package main

import (
	"fmt"
	"net/http"
	"regexp"

	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func strip_tags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	return re.ReplaceAllString(content, "")
}

func fillPosts(w http.ResponseWriter, r *http.Request) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var posts = []Post{{Title: "Title 1", Content: "Content 1 is just 1", Summary: "A summary for title 1", FeaturedImage: "image_500x500.jpg", Slug: "title-1", MetaTitle: "Title 1 - Meta", MetaDescription: "Description for Meta", UserId: 5},
		{Title: "Title 2", Content: "Content 2 is just 2", Summary: "A summary for title 2", FeaturedImage: "image_500x500.jpg", Slug: "title-2", MetaTitle: "Title 2 - Meta", MetaDescription: "Description for Meta 2", UserId: 1}}

	result := db.CreateInBatches(posts, 2) // pass pointer of data to Create

	if result.Error != nil {
		fmt.Println(result.Error) // returns error
	}

	fmt.Printf("Rows affected: %v", result.RowsAffected)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var posts []Post

	// Get all records
	result := db.Find(&posts)
	// SELECT * FROM users;

	for _, post := range posts {
		fmt.Printf("%d %v\n", post.ID, post.Title)
	}

	if result.Error != nil {
		fmt.Println(result.Error) // returns error
	}

	fmt.Printf("Rows affected: %v", result.RowsAffected)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	id := strings.TrimPrefix(r.URL.Path, "/post/")

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	switch r.Method {
	case "GET":
		result := db.First(&post, id)

		if result.Error != nil {
			fmt.Println(result.Error) // returns error
		}
		fmt.Printf("Post: %v", post)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)

		post = Post{
			Title:           r.FormValue("title"),
			Content:         r.FormValue("content"),
			Summary:         r.FormValue("summary"),
			FeaturedImage:   r.FormValue("featured_image"),
			Slug:            r.FormValue("slug"),
			MetaTitle:       r.FormValue("title"),
			MetaDescription: string(strip_tags(r.FormValue("content"))[0:155]),
			UserId:          1,
		}

		result := db.Create(&post)

		if result.Error != nil {
			fmt.Println(result.Error) // returns error
		}

		fmt.Printf("Rows affected: %v", result.RowsAffected)
	case "PUT":

		db.Model(&post).Where("id = ?", id).Updates(Post{
			Title:           r.FormValue("title"),
			Content:         r.FormValue("content"),
			Summary:         r.FormValue("summary"),
			FeaturedImage:   r.FormValue("featured_image"),
			Slug:            r.FormValue("slug"),
			MetaTitle:       r.FormValue("title"),
			MetaDescription: string(strip_tags(r.FormValue("content"))[0:155]),
		})

		fmt.Printf("Post: %v", post)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
