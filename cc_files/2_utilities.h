#pragma once

#include "0_input_check.h"

#define FIRST_STAGE  "[1] - Заполнить массив вручную\n"            \
                     "[2] - Заполнить массив рандомными числами\n" \
                     "[3] - Выход из программы\n";


class Utilities {
  InputCheck _input;
public:
  int choice_check(const int start, const int finish, int &res) {
    do {
      std::cout << "\n\nВыберите пункт из меню\n";
      std::cout << "\n" << FIRST_STAGE;
      res = _input.get_int();
    } while((start > res)||(res > finish));
    return res;
  }

  size_t array_size(size_t &res) {
    std::cout << "Введите размер массива: ";
    int num_of_objects = _input.get_int();
    while (num_of_objects < 1) {
      std::cout << "Введите число больше 0:";
      num_of_objects = _input.get_int();
    }
    res = static_cast<size_t>(num_of_objects);
    return res;
  }

  void save_array(size_t &number_of_objects, std::vector<int> &vec) {
    int choice = 0;
    while ((choice < 1) || (choice > 2)) {
      std::cout << "Сохранить массив в файл? [1] - Да, [2] - Нет\n";
      choice = _input.get_int();
    }

    if (choice == 1) {
      std::ofstream output_file("sorted_array.txt");
      output_file << "Отсортированный массив:\n";
      for (size_t i = 0; i < number_of_objects; i++)
        output_file << vec[i] << " ";
      output_file.close();
      std::cout << "\nМассив сохранен в файл sorted_array.txt\n";
    }
    else std::cout << "\nМассив не был сохранен.\n";
  }
};
