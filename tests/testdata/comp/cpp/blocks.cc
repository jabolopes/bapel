
#include "blocks_private.h"

std::monostate blocks1() {
  {
    std::monostate();
  };
  {
    {
      std::monostate();
    };
    {
      std::monostate();
    };
  };
  return std::monostate();
}

std::monostate blocks2() {
  {
    std::monostate();
  };
  {
    {
      std::monostate();
    };
    {
      return std::monostate();
    };
  };
}
