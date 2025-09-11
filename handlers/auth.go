package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"personalblog/middleware"
	"text/template"
)

type LoginRequest struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

type LoginResponse struct {
	Token   string `json: "token"`
	Message string `json: "message"`
	Success bool   `json: "success"`
}

// фунція для установки Cookie
func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

// функція для хешування пароля
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// хендлер логина(авторизації)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// перевірка запроса і парсинг html
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}
	// получення значень
	username := r.FormValue("username")
	password := r.FormValue("password")

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// валідація
	if adminUsername == "" {
		adminUsername = "admin"
	}
	if adminPassword == "" {
		adminPassword = "123"
	}

	// валідація данних авторизації
	if username != adminUsername || hashPassword(password) != hashPassword(adminPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		tmpl := template.Must(template.ParseFiles("templates/login_error.html"))
		tmpl.Execute(w, nil)
		return
	}

	// генерація JWT токена
	token, err := middleware.GenerateJWT(username)
	if err != nil {
		http.Error(w, `{"error":"falid generate token"}`, http.StatusInternalServerError)
		return
	}

	// встановлення cookie
	SetAuthCookie(w, token)
	// перенапревлення
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// хендлер лагаут(виход з акаунту)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// очистка cookie
	middleware.ClearAuthCookie(w)
	// перенаправлення
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
