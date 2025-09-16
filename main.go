package main

import (
	"fmt"
	"net/http"
	"os"
	api "personalblog/handlers"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	if _, err := os.Stat("posts"); os.IsNotExist(err) {
		os.Mkdir("posts", 0755)
	}
}
func main() {

	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/posts/", api.ArticleHandler)

	http.HandleFunc("/dashboard", api.DashBoardWithAuth())
	http.HandleFunc("/new", api.CreatePostWithAuth())
	http.HandleFunc("/edit/", api.UpdatePostWithAuth())
	http.HandleFunc("/delete/", api.DeletePostWithAuth())
	http.HandleFunc("/search", api.SearchPostHanler)

	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/logout", api.LogoutHandler)

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	fmt.Println("http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

}
