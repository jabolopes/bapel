#pragma once
#include <sstream>
#include <string>

// @bpl: pub type IStringStream
using IStringStream = std::istringstream;

// @bpl: pub IStringStream::read: forall ['a] (& IStringStream, & 'a) -> bool
namespace IStringStream_ {

template <typename T>
inline bool read(IStringStream* iss, T* val) {
  return static_cast<bool>(*iss >> *val);
}

}  // namespace IStringStream_
