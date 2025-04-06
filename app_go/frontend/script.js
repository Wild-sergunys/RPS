// Ожидание полной загрузки html
document.addEventListener('DOMContentLoaded', function() {
    // Элементы DOM
    const elements = {
        arrayInput: document.getElementById('array-input'),
        sortBtn: document.getElementById('sort-btn'),
        saveBtn: document.getElementById('save-btn'),
        clearBtn: document.getElementById('clear-btn'),
        resultContainer: document.getElementById('result-container'),
        arraysList: document.getElementById('arrays-list'),
        inputError: document.getElementById('input-error')
    };

    // Делегирование событий для динамических кнопок
    elements.arraysList.addEventListener('click', function(e) {
        const target = e.target;
        if (target.classList.contains('load-btn')) {
            loadArray(target.dataset.id);
        } else if (target.classList.contains('sort-btn')) {
            sortAndSaveArray(target.dataset.id);
        } else if (target.classList.contains('delete-btn')) {
            deleteArray(target.dataset.id);
        }
    });

    // Обработчики кнопок
    elements.sortBtn.addEventListener('click', sortArray);
    elements.saveBtn.addEventListener('click', saveArray);
    elements.clearBtn.addEventListener('click', clearInput);

    // Загрузка данных при старте
    loadArrays();

    async function deleteArray(id) {
        if (!confirm('Вы уверены, что хотите удалить этот массив?')) {
            return;
        }
        
        try {
            const response = await fetch(`/arrays/delete?id=${id}`, { // http запрос к серверу
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            const data = await response.json(); // Обработка ответа сервера
            
            if (!response.ok) {
                throw new Error(data.message || 'Ошибка сервера');
            }
            
            alert('Массив успешно удален!');
            loadArrays();
        } catch (error) {
            console.error('Error:', error);
            showError(error.message);
        }
    }

    function sortArray() {
        try {
            // Получаем строку из input поля и преобразуем в массив
            const array = parseArray(elements.arrayInput.value);
            // Сортируем созданную копию массива ([...array])
            const sortedArray = selectionSort([...array]);
            
            elements.resultContainer.innerHTML = `
                <h3>Исходный массив:</h3>
                <p>[${array.join(', ')}]</p>
                <h3>Отсортированный массив:</h3>
                <p>[${sortedArray.join(', ')}]</p>
            `;

            clearError();
        } catch (e) {
            showError(e.message);
        }
    }

    function renderArrays(arrays) {
        elements.arraysList.innerHTML = '';
        
        if (!Array.isArray(arrays) || arrays.length === 0) {
            elements.arraysList.innerHTML = '<p>Нет сохраненных массивов</p>';
            return;
        }
        
        // forEach - выполнение для каждого элемента массива
        arrays.sort((a, b) => a.id - b.id).forEach(arr => {
            // div - division
            const arrayItem = document.createElement('div');
            arrayItem.className = 'array-item';
            // .inerHTML - св-во, позволяющее получить содержимое в виде строки или установиь новое
            arrayItem.innerHTML = `
                <h3>Массив #${arr.id}</h3>
                <p>${arr.array_data}</p>
                <p>Статус: ${arr.is_sorted ? 'Отсортирован' : 'Не отсортирован'}</p>
                <div class="array-actions">
                    <button data-id="${arr.id}" class="load-btn">Загрузить</button>
                    <button data-id="${arr.id}" class="sort-btn">Сортировать</button>
                    <button data-id="${arr.id}" class="delete-btn">Удалить</button>
                </div>
            `;
            // добавляем в DOM в конец дочерних эл-ов
            elements.arraysList.appendChild(arrayItem);
        });
    }

    async function saveArray() {
        try {
            // .trim() - удаляет пробелы в начале и коуе строки
            const input = elements.arrayInput.value.trim();

            // /.../ - ограничение выражения
            // ^ - начало строки
            // \d+ - одна или юолее цифр
            // (\s*,\s*\d+) - группа символов, (ноль или более пробелов , нибп одно или более цифр)
            // .test - проверка на соответствие
            if (!/^\d+(\s*,\s*\d+)*$/.test(input)) {
                throw new Error('Используйте формат: "123, 22, 111"');
            }
    
            // fetch - POST-запрос
            const response = await fetch('/arrays/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    array: input,
                    isSorted: false 
                })
            });
            
            const result = await response.json();
            
            if (!result.success) {
                throw new Error(result.message || 'Ошибка сохранения');
            }
            
            alert('Массив сохранен!');
            loadArrays();
        } catch (error) {
            console.error('Ошибка:', error);
            showError(error.message);
        }
    }

    async function loadArrays() {
        try {
            const response = await fetch('/arrays');
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.message || 'Ошибка сервера');
            }
            
            renderArrays(data.data);
        } catch (error) {
            console.error('Error:', error);
            showError('Ошибка загрузки массивов');
        }
    }

    async function loadArray(id) {
        try {
            const response = await fetch(`/arrays/load?id=${id}`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.message || 'Ошибка сервера');
            }
            
            const arrayData = data.array || data.array_data || data.data?.array;
            if (!arrayData) {
                throw new Error('Неверный формат данных массива');
            }
            
            elements.arrayInput.value = arrayData;
            elements.resultContainer.innerHTML = '';
            clearError();
        } catch (error) {
            console.error('Error:', error);
            showError(error.message);
        }
    }

    async function sortAndSaveArray(id) {
        try {
            const response = await fetch(`/arrays/sort?id=${id}`, { 
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.message || 'Ошибка сервера');
            }
            
            alert('Массив успешно отсортирован!');
            loadArrays();
        } catch (error) {
            console.error('Error:', error);
            showError(error.message);
        }
    }

    function clearInput() {
        elements.arrayInput.value = '';
        elements.resultContainer.innerHTML = '';
        clearError();
    }

    function showError(message) {
        elements.inputError.textContent = message;
        setTimeout(clearError, 3000); 
    }

    function clearError() {
        elements.inputError.textContent = '';
    }

    function parseArray(input) {
        // trim - удаляет пробелы в начале и конце
        if (!input.trim()) throw new Error('Введите числа для сортировки');
        
        return input.split(',')
            .map(item => item.trim()) // удаляем пробелы вокруг эл-ов
            .filter(item => item !== '') // удаляем пустые элементы
            .map(item => {
                const num = parseFloat(item); // строка в число
                if (isNaN(num)) throw new Error(`"${item}" не является числом`); // проверяем валидно ли число
                return num;
            });
    }

    function selectionSort(arr) {
        const n = arr.length;
        for (let i = 0; i < n - 1; i++) {
            let minIndex = i;
            for (let j = i + 1; j < n; j++) {
                if (arr[j] < arr[minIndex]) minIndex = j;
            }
            if (minIndex !== i) [arr[i], arr[minIndex]] = [arr[minIndex], arr[i]];
        }
        return arr;
    }
});