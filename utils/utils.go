package utils

import (
	"fmt"
	"net/http"
	"personalblog/middleware"
	"time"
)

// допоміжна структура для валідації
type PostData struct {
	User    *middleware.Claims
	Title   string
	Content string
	Date    string
}

// функція валідації данних
func Validate(r *http.Request) (*PostData, error) {
	// перевірка авторизації користувача
	user := middleware.GetUserFromContext(r)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// получение значень
	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	// валідація
	if title == "" || len(title) > 100 {
		return nil, fmt.Errorf("invalid title")
	}

	if content == "" {
		return nil, fmt.Errorf("invalid content")
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid data")
	}
	// відправка значень та помилки
	return &PostData{
		User:    user,
		Title:   title,
		Content: content,
		Date:    date,
	}, nil
}
