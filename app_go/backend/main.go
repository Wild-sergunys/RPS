package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql" // Драйвер MySQL для работы с базой данных
)

var db *sql.DB // Глобальная переменная для хранения соединения с базой данных

func main() {
	// Инициализация БД
	var err error
	// ИНнициализация соединения с DB (тип данных, connection string)
	db, err = sql.Open("mysql", "sorting_user:123@tcp(127.0.0.1:3306)/sorting_app")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Закрытие соединения с БД при завершении функции main

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err) // Логирование ошибки подключения
	}

	// Получаем абсолютный путь к директории frontend
	frontendPath := filepath.Join(getProjectRoot(), "frontend")

	// Определение маршрутов
	http.HandleFunc("/arrays", arraysHandler)
	http.HandleFunc("/arrays/save", saveArrayHandler)
	http.HandleFunc("/arrays/load", loadArrayHandler)
	http.HandleFunc("/arrays/sort", sortArrayHandler)
	http.HandleFunc("/arrays/delete", deleteArrayHandler)
	http.HandleFunc("/arrays/reindex", reindexArraysHandler)

	// fs - файловый сервер
	// Настройка статического сервера для обслуживания файлов из директории frontend
	fs := http.FileServer(http.Dir(frontendPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Удаление префикса /static/ из URL

	// Главная страница приложения (index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendPath, "index.html")) // Отправка файла index.html клиенту
	})

	fmt.Println("Сервер запущен на http://localhost:8080") // Сообщение о запуске сервера
	log.Fatal(http.ListenAndServe(":8080", nil))           // Запуск HTTP-сервера на порту 8080 и логирование ошибок
}

// Функция для получения корневой директории проекта
func getProjectRoot() string {
	dir, err := os.Getwd() // Получение текущей рабочей директории
	if err != nil {
		log.Fatal(err) // Логирование ошибки при получении директории
	}
	return filepath.Dir(dir) // Возврат родительской директории (корень проекта)
}
