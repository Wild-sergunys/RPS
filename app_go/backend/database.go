package main

import (
	"fmt"
	"strconv"
	"strings"
)

func saveArrayToDB(numbers []int, isSorted bool) (int64, error) {
	var sb strings.Builder
	for i, num := range numbers {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Itoa(num))
	}
	arrayStr := sb.String()

	res, err := db.Exec("INSERT INTO arrays (array_data, is_sorted) VALUES (?, ?)", arrayStr, isSorted)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func getAllArrays() ([]map[string]interface{}, error) {
	// Сортируем по текущим ID
	// Query отправляет запрос к бд, rows - итератор для доступа к рез sql запроса
	rows, err := db.Query("SELECT id, array_data, is_sorted FROM arrays ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arrays []map[string]interface{}

	// for rows.Next() возращает true, если строка доступна для чтения
	for rows.Next() {
		var id int
		var arrayData string
		var isSorted bool

		err = rows.Scan(&id, &arrayData, &isSorted)
		if err != nil {
			return nil, err
		}

		arrays = append(arrays, map[string]interface{}{
			"id":         id,
			"array_data": arrayData,
			"is_sorted":  isSorted,
		})
	}

	return arrays, nil
}

func getArrayByID(id int) (string, error) {
	var arrayData string
	err := db.QueryRow("SELECT array_data FROM arrays WHERE id = ?", id).Scan(&arrayData)
	return arrayData, err
}

func reindexArrays() error {
	// Проверяем соединение с БД
	if err := db.Ping(); err != nil {
		return fmt.Errorf("проверка соединения с БД не удалась: %v", err)
	}

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Временная таблица для переиндексации
	// _ - игнорируем результат (кол-во строк)
	_, err = tx.Exec(`
		CREATE TEMPORARY TABLE IF NOT EXISTS temp_reindex AS 
		SELECT id, ROW_NUMBER() OVER (ORDER BY created_at) as new_id 
		FROM arrays
	`)
	if err != nil {
		return fmt.Errorf("ошибка создания временной таблицы: %v", err)
	}

	// Обновляем ID
	_, err = tx.Exec(`
			UPDATE arrays a
			JOIN temp_reindex t ON a.id = t.id
			SET a.id = t.new_id
	`)
	if err != nil {
		return fmt.Errorf("ошибка обновления ID: %v", err)
	}

	// Удаляем временную таблицу
	_, err = tx.Exec("DROP TEMPORARY TABLE temp_reindex")
	if err != nil {
		return fmt.Errorf("ошибка удаления временной таблицы: %v", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	return nil
}
