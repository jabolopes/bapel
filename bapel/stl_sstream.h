#pragma once
#include <sstream>
#include <string>
#include <string_view>
#include <variant>

// @bpl: pub type IStringStream
using IStringStream = std::istringstream;

// @bpl: pub IStringStream::read: forall ['a] (& IStringStream, & 'a) -> bool
// @bpl: pub IStringStream::mk: String -> IStringStream
namespace inherents {
struct IStringStream {
  IStringStream() = delete;

  static inline ::IStringStream mk(const std::string& s) {
    return ::IStringStream(s);
  }

  template <typename T>
  static inline bool read(::IStringStream* iss, T* val) {
    return static_cast<bool>(*iss >> *val);
  }
};
}

using OStringStream = std::ostringstream;

// @bpl: pub OStringStreamImpl::mk: () -> OStringStream
// @bpl: pub OStringStreamImpl::write: (&OStringStream, StringView) -> ()
// @bpl: pub OStringStreamImpl::write_char: (&OStringStream, i8) -> ()
// @bpl: pub OStringStreamImpl::str: &OStringStream -> String
struct OStringStreamImpl {
  OStringStreamImpl() = delete;

  static inline OStringStream mk() {
    return OStringStream();
  }

  static inline std::monostate write(OStringStream* oss, std::string_view s) {
    *oss << s;
    return std::monostate();
  }

  static inline std::monostate write_char(OStringStream* oss, char c) {
    oss->put(c);
    return std::monostate();
  }

  static inline std::string str(const OStringStream* oss) {
    return oss->str();
  }
};
