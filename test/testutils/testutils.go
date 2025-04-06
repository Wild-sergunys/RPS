package testutils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Добавляем импорт драйвера MySQL
)

const (
	dbUser     = "sorting_user"
	dbPassword = "123"
	dbName     = "sorting_app"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GenerateRandomArray() string {
	size := rand.Intn(50) + 5 // Массивы от 5 до 55 элементов
	var arr []int
	for i := 0; i < size; i++ {
		arr = append(arr, rand.Intn(1000)-500) // Числа от -500 до 499
	}
	return IntArrayToString(arr)
}

func IntArrayToString(arr []int) string {
	var sb strings.Builder
	for i, num := range arr {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%d", num))
	}
	return sb.String()
}

func ClearDatabase(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE arrays")
	return err
}

func ParseArrayString(input string) ([]int, error) {
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

// Копия функции сортировки из handlers.go
func SelectionSort(arr []int) []int {
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
