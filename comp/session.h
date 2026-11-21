#pragma once

#include "comp/module_finder.h"
#include "comp/querier.h"

namespace comp {

// @bpl: pub type comp::Session
struct Session {
  Querier querier;

  Session() : querier(ModuleFinder({}, {{"", "."}})) {}
  explicit Session(Querier q) : querier(std::move(q)) {}
};

// @bpl: pub comp::mk_session: () -> comp::Session
inline Session mk_session() {
  return Session();
}

} // namespace comp
