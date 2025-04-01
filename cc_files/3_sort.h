#pragma once

#include "0_input_check.h"

class Sort {
public:
  void selection_sort(std::vector<int> &vec) {
    size_t min_index;         //О(1)
    size_t size = vec.size(); //О(1)
    for (size_t i = 0; i < size - 1; i++) {
      //О(N) - внешний цикл проходит по каждому элементу вектора, кроме последнего
      min_index = i;          //О(1)
      for (size_t j = i + 1; j < size; j++)
      //О(N) - внутренний цикл проходит по всем элементам вектора, начиная с i + 1
        if (vec[j] < vec[min_index])     //О(1)
          min_index = j;                 //О(1)
      std::swap(vec[i], vec[min_index]); //О(1)
    }
    std::cout << "\nОтсортированный массив:\n";
    for (size_t i = 0; i < size; i++)
      std::cout << vec[i] << " ";
    std::cout << "\n";
  }
};
