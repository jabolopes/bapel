#pragma once

#include <fstream>
#include <string>
#include <variant>

// @bpl: pub OfstreamImpl::open: &String -> Ofstream
// @bpl: pub OfstreamImpl::is_open: &Ofstream -> bool
// @bpl: pub OfstreamImpl::close: &Ofstream -> ()
// @bpl: pub OfstreamImpl::write: (&Ofstream, String) -> ()

// @bpl: pub IfstreamImpl::open: &String -> Ifstream
// @bpl: pub IfstreamImpl::is_open: &Ifstream -> bool
// @bpl: pub IfstreamImpl::close: &Ifstream -> ()
// @bpl: pub IfstreamImpl::read: forall ['a] (&Ifstream, &'a) -> bool

using Ofstream = std::ofstream;
using Ifstream = std::ifstream;

struct OfstreamImpl {
  OfstreamImpl() = delete;

  static inline std::ofstream open(const std::string* filename) {
    return std::ofstream(*filename);
  }

  static inline bool is_open(std::ofstream* f) {
    return f->is_open();
  }

  static inline std::monostate close(std::ofstream* f) {
    f->close();
    return std::monostate();
  }

  static inline std::monostate write(std::ofstream* f, const std::string& s) {
    *f << s;
    return std::monostate();
  }
};

struct IfstreamImpl {
  IfstreamImpl() = delete;

  static inline std::ifstream open(const std::string* filename) {
    return std::ifstream(*filename);
  }

  static inline bool is_open(std::ifstream* f) {
    return f->is_open();
  }

  static inline std::monostate close(std::ifstream* f) {
    f->close();
    return std::monostate();
  }

  template <typename T>
  static inline bool read(std::ifstream* ifs, T* val) {
    return static_cast<bool>(*ifs >> *val);
  }
};

