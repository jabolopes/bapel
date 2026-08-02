
#include "number_literals_private.h"

std::monostate callWithIDs() {
  int8_t i = static_cast<int8_t>(0);
  int8_t j = static_cast<int8_t>(0);
  (i) + (j);
  return std::monostate();
}

std::monostate callWithIDAndLiterals() {
  int8_t i = static_cast<int8_t>(0);
  (i) + (static_cast<int8_t>(1));
  return std::monostate();
}

std::monostate callWithLiterals() {
  int8_t i = static_cast<int8_t>(0);
  (i) + (static_cast<int8_t>(1));
  return std::monostate();
}

std::monostate letWithIDs() {
  int8_t i = static_cast<int8_t>(0);
  int8_t j = static_cast<int8_t>(0);
  int8_t x = (i) + (j);
  return std::monostate();
}

std::monostate letWithIDAndLiterals() {
  int8_t i = static_cast<int8_t>(0);
  int8_t x = (i) + (static_cast<int8_t>(1));
  return std::monostate();
}

std::monostate letWithLiterals() {
  int8_t x = (static_cast<int8_t>(1)) + (static_cast<int8_t>(2));
  return std::monostate();
}

std::monostate assignWithIDs() {
  int8_t x = static_cast<int8_t>(0);
  int8_t i = static_cast<int8_t>(0);
  int8_t j = static_cast<int8_t>(0);
  x = (i) + (j);
  return std::monostate();
}

std::monostate assignWithIDAndLiterals() {
  int8_t x = static_cast<int8_t>(0);
  int8_t i = static_cast<int8_t>(0);
  x = (i) + (static_cast<int8_t>(1));
  return std::monostate();
}

std::monostate assignWithLiterals() {
  int8_t x = static_cast<int8_t>(0);
  x = (static_cast<int8_t>(1)) + (static_cast<int8_t>(2));
  return std::monostate();
}

int8_t returnWithIDs() {
  int8_t i = static_cast<int8_t>(0);
  int8_t j = static_cast<int8_t>(0);
  return (i) + (j);
}

int8_t returnWithIDAndLiterals() {
  int8_t i = static_cast<int8_t>(0);
  return (i) + (static_cast<int8_t>(1));
}

int8_t returnWithLiterals() {
  return (static_cast<int8_t>(1)) + (static_cast<int8_t>(2));
}
