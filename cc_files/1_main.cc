#include "2_utilities.h"
#include "3_sort.h"

#define GREETINGS    "******************************************\n"         \
                     "* Здравствуйте!                          *\n"         \
                     "* Лабораторная работа №2                 *\n"         \
                     "* Предмет: Разработка программных систем *\n"         \
                     "* Работу выполнили студенты группы 434   *\n"         \
                     "* Щербаков С.И. & Журавлева А.В.         *\n"         \
                     "******************************************\n";

enum FIRST_CHOICE {MANUAL_INPUT = 1, RANDOM_INPUT, EXIT};

int main(void)
{
  InputCheck _input;
  Utilities _utils;
  Sort _sort;

  std::vector<int> array_for_sort;
  std::cout << GREETINGS;

  const int start_stage = 1, finish_stage = 3;

  for (;;) {
    int choice = 0;
    size_t num_of_objects = 0;
    _utils.choice_check(start_stage, finish_stage, choice);

    switch (choice)
    {
      case MANUAL_INPUT:
      {
        num_of_objects = _utils.array_size(num_of_objects);
        array_for_sort.resize(num_of_objects);
        for (size_t i = 0; i < array_for_sort.size(); i++) {
          std::cout << "Введите [" << (i + 1) << "] элемент массива: ";
          array_for_sort[i] = _input.get_int();
        }

        _sort.selection_sort(array_for_sort);
        array_for_sort.clear();
        _utils.save_array(num_of_objects, array_for_sort);

        break;
      }
      case RANDOM_INPUT:
      {
        num_of_objects = _utils.array_size(num_of_objects);
        array_for_sort.resize(num_of_objects);
        std::cout << "\nВведите диапазон значений для заполения массива (введите целые числа)\n"
                     "Максимальное значение: ";
        int max_num = _input.get_int();
        std::cout << "Минимальное значение: ";
        int min_num = _input.get_int();
        while ((max_num < min_num) || (max_num == min_num)
                                   || (max_num - min_num == 1)) {
          std::cout << "Введите допустимый диапазон!\n"
                       "Максимальное значение: ";
          max_num = _input.get_int();
          std::cout << "Минимальное значение: ";
          min_num = _input.get_int();
        }

        time_t system_time = time(0);
        std::default_random_engine generated_number(system_time);
        std::uniform_int_distribution<int> range(min_num, max_num);
        for (size_t arr_element_index = 0; arr_element_index < num_of_objects; arr_element_index++)
          array_for_sort[arr_element_index] = range(generated_number);

        std::cout << "\n\nИсходный массив:\n";
        for (size_t arr_element_index = 0; arr_element_index < num_of_objects; arr_element_index++)
          std::cout << array_for_sort[arr_element_index] << " ";
        std::cout << "\n";

        _sort.selection_sort(array_for_sort);
        _utils.save_array(num_of_objects, array_for_sort);
        array_for_sort.clear();

        break;
      }
    }
    if (static_cast<FIRST_CHOICE>(choice) == EXIT) exit(0);
  }
  return 0;
}
