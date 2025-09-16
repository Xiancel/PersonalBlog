package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"personalblog/middleware"
	"personalblog/model"
	"personalblog/utils"
	"strconv"
	"strings"
	"text/template"
)

// функція для парсингу html
func parseTemplate(templateName string) *template.Template {
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		panic(fmt.Sprintf("Error parsing templates %s, %v", templateName, err))
	}
	return tmpl
}

// функція для получення всіх постів
func getPosts() []model.Posts {
	files, _ := os.ReadDir("posts")
	var posts []model.Posts

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		data, _ := os.ReadFile(filepath.Join("posts", f.Name()))
		var pos model.Posts
		json.Unmarshal(data, &pos)
		posts = append(posts, pos)
	}

	return posts
}

// функція для получення посту по ID
func GetPostByID(id int) *model.Posts {
	filepath := fmt.Sprintf("posts/post%d.json", id)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	var a model.Posts
	json.Unmarshal(data, &a)
	return &a
}

// хендлер головної сторінки
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts := getPosts()
	tmpl := parseTemplate("home.html")
	tmpl.Execute(w, posts)
}

// хендлер посту
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// получення айді посту з путі
	idStr := r.URL.Path[len("/posts/"):]
	// конвертація в int, валідація айді
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	// получення айді посту і валідація айді
	post := GetPostByID(id)
	if post == nil {
		http.Error(w, `{"error":"no found"}`, http.StatusNotFound)
		return
	}

	// парсинг
	tmpl := parseTemplate("postpage.html")
	tmpl.Execute(w, post)
}

// хендлер Dashboardа
func DashBoardHandler(w http.ResponseWriter, r *http.Request) {
	post := getPosts()
	tmpl := parseTemplate("dashboard.html")
	tmpl.Execute(w, post)
}

// хендлер створення посту
func NewHandler(w http.ResponseWriter, r *http.Request) {
	// перевірка методі і парсинг html
	if r.Method == "GET" {
		tmpl := parseTemplate("newPost.html")
		tmpl.Execute(w, nil)
		return
	}

	//validation
	postData, err := utils.Validate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	//Create Posts
	// визначення айді для постів
	files, _ := os.ReadDir("posts")
	maxID := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			var id int
			fmt.Sscanf(file.Name(), "post%d.json", &id)
			if id > maxID {
				maxID = id
			}
		}
	}

	var a model.Posts

	a.Title = postData.Title
	a.Content = postData.Content
	a.Date = postData.Date
	a.Author = postData.User.Username
	a.ID = maxID + 1

	filePath := fmt.Sprintf("posts/post%d.json", a.ID)
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// хендлер редагування посту
func EditHandler(w http.ResponseWriter, r *http.Request) {
	// получення айді посту з шляху
	idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
	// конвертація в int, валідація айді
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		// получення айді посту
		post := GetPostByID(id)
		if post == nil {
			http.Error(w, `{"error":"no found"}`, http.StatusNotFound)
			return
		}

		// прарсинг html
		tmpl := parseTemplate("updatePost.html")
		tmpl.Execute(w, post)
		return
	}

	//validation
	postData, err := utils.Validate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// редагування посту
	var p model.Posts

	p.ID = id
	p.Title = postData.Title
	p.Content = postData.Content
	p.Date = postData.Date
	p.Author = postData.User.Username

	filepath := fmt.Sprintf("posts/post%d.json", id)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	file, _ := os.Create(filepath)
	defer file.Close()
	json.NewEncoder(file).Encode(p)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// хендлер видалення посту
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// получення айді з шляху
	idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
	// конвертація в int, валідація айді
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	// видалення посту
	filepath := fmt.Sprintf("posts/post%d.json", id)
	if err := os.Remove(filepath); err != nil {
		http.Error(w, `{"error":"failed to delete post"}`, http.StatusInternalServerError)
		return
	}

	// перенаправлення
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// хендлер пошуку
func SearchPostHanler(w http.ResponseWriter, r *http.Request) {
	// перевірка методу
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// пошук поста по назві
	query := r.FormValue("search")
	post := SearchPosts(query)

	// валідація
	if post == nil {
		http.Error(w, `{"error":"no found"}`, http.StatusNotFound)
		return
	}

	// парсинг
	tmpl := parseTemplate("postpage.html")
	tmpl.Execute(w, post)
}

// функція пошуку посту
func SearchPosts(query string) *model.Posts {
	// получення всіх постів
	posts := getPosts()
	// пошук посту
	for _, p := range posts {
		if strings.EqualFold(p.Title, query) {
			return &p
		}
	}

	return nil
}

// перевірка на авторизацію(на автора/адміна )
func CreatePostWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(NewHandler)
}

func DeletePostWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(DeleteHandler)
}
func UpdatePostWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(EditHandler)
}
func DashBoardWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(DashBoardHandler)
}
