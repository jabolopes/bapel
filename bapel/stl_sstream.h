#pragma once
#include <sstream>
#include <string>

// @bpl: pub type IStringStream
using IStringStream = std::istringstream;

// @bpl: pub IStringStream::read: forall ['a] (& IStringStream, & 'a) -> bool
// @bpl: pub IStringStream::mk: String -> IStringStream
struct IStringStream_ {
  IStringStream_() = delete;

  static inline IStringStream mk(const std::string& s) {
    return IStringStream(s);
  }

  template <typename T>
  static inline bool read(IStringStream* iss, T* val) {
    return static_cast<bool>(*iss >> *val);
  }
};
