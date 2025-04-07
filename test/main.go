package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"test/testutils"
)

func main() {
	db, err := testutils.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		fmt.Println("\nТесты производительности базы данных")
		fmt.Println("1. Запуск всех тестов")
		fmt.Println("2. Тесты вставки")
		fmt.Println("3. Тесты сортировки")
		fmt.Println("4. Тесты очистки")
		fmt.Println("0. Выход")

		var choice int
		fmt.Print("Выберите вариант: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода, попробуйте еще раз")
			continue
		}

		switch choice {
		case 1:
			runAllTests(db)
		case 2:
			runInsertTests(db)
		case 3:
			runSortTests(db)
		case 4:
			runClearTests(db)
		case 0:
			fmt.Println("Выход из программы")
			return
		default:
			fmt.Println("Неверный выбор, попробуйте еще раз")
		}
	}
}

func runAllTests(db *sql.DB) {
	fmt.Println("\n=== Запуск всех тестов ===")
	runInsertTests(db)
	runSortTests(db)
	runClearTests(db)
}

func runInsertTests(db *sql.DB) {
	fmt.Println("\n=== Тесты вставки ===")
	sizes := []int{100, 1000, 10000}
	for _, size := range sizes {
		runInsertTest(db, size)
	}
	if err := testutils.ClearDatabase(db); err != nil {
		log.Printf("Ошибка очистки базы после тестов сортировки: %v", err)
	} else {
		fmt.Println("\nБаза данных очищена после тестов сортировки")
	}
}

func runInsertTest(db *sql.DB, count int) {
	fmt.Printf("\nТест вставки %d массивов:\n", count)

	// Очищаем базу перед тестом
	if err := testutils.ClearDatabase(db); err != nil {
		log.Printf("Ошибка очистки базы: %v", err)
		return
	}

	start := time.Now()
	success := true

	for i := 0; i < count; i++ {
		arr := testutils.GenerateRandomArray()
		_, err := db.Exec("INSERT INTO arrays (array_data, is_sorted) VALUES (?, ?)", arr, false)
		if err != nil {
			log.Printf("Ошибка вставки: %v", err)
			success = false
			break
		}
	}

	duration := time.Since(start)
	fmt.Printf("Успешно: %v\n", success)
	fmt.Printf("Общее время: %v\n", duration)
	fmt.Printf("Среднее время на массив: %v\n", duration/time.Duration(count))
}

func runSortTests(db *sql.DB) {
	fmt.Println("\n=== Тесты сортировки ===")
	dbSizes := []int{100, 1000, 10000}
	for _, size := range dbSizes {
		runSortTest(db, size)
	}
	if err := testutils.ClearDatabase(db); err != nil {
		log.Printf("Ошибка очистки базы после тестов сортировки: %v", err)
	} else {
		fmt.Println("\nБаза данных очищена после тестов сортировки")
	}
}

func runSortTest(db *sql.DB, dbSize int) {
	fmt.Printf("\nТест сортировки (база из %d записей):\n", dbSize)

	// Подготовка тестовых данных
	if err := testutils.ClearDatabase(db); err != nil {
		log.Printf("Ошибка очистки базы: %v", err)
		return
	}

	// Заполняем базу
	for i := 0; i < dbSize; i++ {
		arr := testutils.GenerateRandomArray()
		if _, err := db.Exec("INSERT INTO arrays (array_data, is_sorted) VALUES (?, ?)", arr, false); err != nil {
			log.Printf("Ошибка заполнения базы: %v", err)
			return
		}
	}

	start := time.Now()
	success := true
	processed := 0
	var totalSortTime time.Duration

	rows, err := db.Query("SELECT id, array_data FROM arrays ORDER BY RAND() LIMIT 100")
	if err != nil {
		log.Printf("Ошибка выборки: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var arrayData string
		if err := rows.Scan(&id, &arrayData); err != nil {
			log.Printf("Ошибка чтения: %v", err)
			success = false
			break
		}

		// Парсим массив
		numbers, err := testutils.ParseArrayString(arrayData)
		if err != nil {
			log.Printf("Ошибка парсинга: %v", err)
			success = false
			break
		}

		// Сортируем (используем копию функции из handlers.go)
		sortStart := time.Now()
		sorted := testutils.SelectionSort(numbers)
		sortedStr := testutils.IntArrayToString(sorted)
		totalSortTime += time.Since(sortStart)

		// Сохраняем результат
		_, err = db.Exec("UPDATE arrays SET array_data = ?, is_sorted = ? WHERE id = ?",
			sortedStr, true, id)
		if err != nil {
			log.Printf("Ошибка обновления: %v", err)
			success = false
			break
		}
		processed++
	}

	duration := time.Since(start)
	fmt.Printf("Успешно: %v\n", success)
	fmt.Printf("Обработано массивов: %d\n", processed)
	fmt.Printf("Общее время: %v\n", duration)
	if processed > 0 {
		fmt.Printf("Среднее время сортировки: %v\n", totalSortTime/time.Duration(processed))
	}
}

func runClearTests(db *sql.DB) {
	fmt.Println("\n=== Тесты очистки ===")
	sizes := []int{100, 1000, 10000}
	for _, size := range sizes {
		runClearTest(db, size)
	}
}

func runClearTest(db *sql.DB, count int) {
	fmt.Printf("\nТест очистки (%d записей):\n", count)

	// Подготовка тестовых данных
	if err := testutils.ClearDatabase(db); err != nil {
		log.Printf("Ошибка очистки базы: %v", err)
		return
	}

	// Заполняем базу
	startFill := time.Now()
	for i := 0; i < count; i++ {
		arr := testutils.GenerateRandomArray()
		if _, err := db.Exec("INSERT INTO arrays (array_data) VALUES (?)", arr); err != nil {
			log.Printf("Ошибка заполнения базы: %v", err)
			return
		}
	}
	fmt.Printf("База заполнена за: %v\n", time.Since(startFill))

	// Тестируем очистку
	start := time.Now()
	err := testutils.ClearDatabase(db)
	duration := time.Since(start)

	fmt.Printf("Успешно: %v\n", err == nil)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	fmt.Printf("Время очистки: %v\n", duration)
}
