
#include "conditionals_private.h"

std::monostate conditionals() {
  true;
  false;
  if (true) {
    true;
  };
  if (!true) {
    false;
  };
  if ((true) == (false)) {
    false;
  };
  bool v1;
  if ((true) == (false)) {
    v1 = false;
  } else {
    v1 = true;
  };
  if ((true) == (false)) {
    false;
  } else if ((false) == (true)) {
    true;
  };
  bool v2;
  if ((true) == (false)) {
    v2 = false;
  } else if ((false) == (true)) {
    v2 = true;
  };
  if ((true) == (false)) {
    false;
  } else if ((false) == (true)) {
    true;
  } else {
    false;
  };
  bool v3;
  if ((true) == (false)) {
    v3 = false;
  } else if ((false) == (true)) {
    v3 = true;
  } else {
    v3 = false;
  };
  return std::monostate();
}

bool ifLastTerm() {
  if (true) {
    return false;
  } else {
    return true;
  };
}

bool ftrue() { return true; }

bool conditionalsPolymorphic() {
  if (id<bool>(true)) {
    id<bool>(true);
  } else {
    id<bool>(false);
  };
  if (id<bool>(true)) {
    id<bool>(true);
  } else {
    id<bool>(false);
  };
  if (fconst<bool, std::monostate>(true, std::monostate())) {
    true;
  } else {
    false;
  };
  if (fconst<bool, std::monostate>(true, std::monostate())) {
    true;
  } else {
    false;
  };
  if ((id<bool>(true)) == (id<bool>(false))) {
    return (id<bool>(true)) == (id<bool>(false));
  } else {
    return (id<bool>(false)) == (id<bool>(true));
  };
}
