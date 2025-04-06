package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")                   // Разрешаем все домены
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Разрешение методов, options для предварительных CORS-запросов
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Разрешенные заголовки запросов
}

type ArrayRequest struct {
	Array    string `json:"array"`
	IsSorted bool   `json:"isSorted"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"` //omitempty - пропуск поля при нулевом значении
	Data    interface{} `json:"data,omitempty"`
}

// Обработчик HTTP-запроса, для получения данных массивов
// w http.ResponseWriter - формирование HTTP-ответа
// r *http.Request - инофрмация об HTTP-запросе
func arraysHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w) // настройка CORS

	if r.Method != "GET" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	arrays, err := getAllArrays()
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при получении массивов: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, Response{
		Success: true,
		Data:    arrays,
	}, http.StatusOK)
}

func saveArrayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req ArrayRequest
	// Декодирование json в ArrayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка декодирования запроса: %v", err),
		}, http.StatusBadRequest)
		return
	}

	// Преобразуем строку в массив чисел с валидацией
	numbers, err := parseArrayString(req.Array)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Неверный формат массива: %v", err),
		}, http.StatusBadRequest)
		return
	}

	id, err := saveArrayToDB(numbers, req.IsSorted)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при сохранении: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	// Выполняем переиндексацию
	if err := reindexArrays(); err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка переиндексации: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	// Получаем обновленный список
	arrays, err := getAllArrays()
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при получении массивов: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, Response{
		Success: true,
		Data:    arrays,
		Message: fmt.Sprintf("Массив сохранен. База переиндексирована. Новый ID: %d", id),
	}, http.StatusCreated)
}

func loadArrayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем id из url
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: "Неверный ID массива",
		}, http.StatusBadRequest)
		return
	}

	arrayData, err := getArrayByID(id)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при загрузке массива: %v", err),
		}, http.StatusNotFound)
		return
	}

	jsonResponse(w, Response{
		Success: true,
		Data:    map[string]string{"array": arrayData},
	}, http.StatusOK)
}

func sortArrayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: "Неверный ID массива",
		}, http.StatusBadRequest)
		return
	}

	// Загружаем массив из БД
	arrayData, err := getArrayByID(id)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при загрузке массива: %v", err),
		}, http.StatusNotFound)
		return
	}

	// Парсим и сортируем массив
	numbers, err := parseArrayString(arrayData)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Неверный формат массива: %v", err),
		}, http.StatusBadRequest)
		return
	}

	sortedNumbers := selectionSort(numbers)

	// Сохраняем отсортированный массив
	_, err = saveArrayToDB(sortedNumbers, true)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при сохранении отсортированного массива: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	// Выполняем переиндексацию
	if err := reindexArrays(); err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка переиндексации: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	// Получаем обновленный список
	arrays, err := getAllArrays()
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при получении массивов: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, Response{
		Success: true,
		Data:    arrays,
		Message: "Массив успешно отсортирован и база переиндексирована",
	}, http.StatusOK)
}


// Вспомогательные функции

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json") // Устанавливаем Content-Type
	w.WriteHeader(statusCode) // Устанавливаем http статус
	json.NewEncoder(w).Encode(data) // Петеровдим данные в json формат
}

func parseArrayString(input string) ([]int, error) {
	var numbers []int
	items := strings.Split(input, ",")

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		num, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("элемент '%s' не является числом", item)
		}
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("массив не может быть пустым")
	}

	return numbers, nil
}

func selectionSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
	return arr
}

func deleteArrayHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method != "DELETE" && r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: "Неверный ID массива",
		}, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM arrays WHERE id = ?", id)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка при удалении массива: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	if err := reindexArrays(); err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка переиндексации: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, Response{
		Success: true,
		Message: "Массив успешно удален",
	}, http.StatusOK)
}

func reindexArraysHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Переиндексация в одном запросе (по старшинству создания)
	_, err := db.Exec(`
			SET @new_id = 0;
			UPDATE arrays SET id = (@new_id := @new_id + 1) ORDER BY id ASC;
	`)
	if err != nil {
		jsonResponse(w, Response{
			Success: false,
			Message: fmt.Sprintf("Ошибка переиндексации: %v", err),
		}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный статус
	jsonResponse(w, Response{
		Success: true,
		Message: "Массивы успешно переиндексированы",
	}, http.StatusOK)
}
