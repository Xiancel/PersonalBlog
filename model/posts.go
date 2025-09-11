package model

// структура поста
type Posts struct {
	ID      int    `json:"id"`                                        // id поста
	Title   string `json:"title" validate:"required, min=5, max=200"` // назва посту
	Content string `json:"content" validate:"required, min=50"`       // вміст посту
	Date    string `json:"date" validate:"required"`                  // дата створення
	Author  string `json:"author,omitempty"`                          // автор
}

// я не пом'ятаю для чого ми його створювали але видаляти його не хочеться
// оно конечно нечего не ламает я проверял но ОН НУЖЕН ДЛЯ БАЛАНСА ВСЕЛЕННОЙ
// Інтерфейс в якому методи для:
type PostRepository interface {
	Save(post Posts) error          // зберегання постів
	FindByID(id int) (Posts, error) // пошуку поста за його айді
}
