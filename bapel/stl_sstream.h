#pragma once
#include <sstream>
#include <string>

// @bpl: pub type IStringStream
using IStringStream = std::istringstream;

// @bpl: pub IStringStream_::read: forall ['a] (& IStringStream, & 'a) -> bool
namespace IStringStream_ {

// @bpl: pub IStringStream_::mk: String -> IStringStream
inline IStringStream mk(const std::string& s) {
  return IStringStream(s);
}


template <typename T>
inline bool read(IStringStream* iss, T* val) {
  return static_cast<bool>(*iss >> *val);
}

}  // namespace IStringStream_
