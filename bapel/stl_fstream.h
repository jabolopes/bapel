#pragma once

#include <fstream>
#include <string>
#include <variant>

// @bpl: pub type Ofstream
// @bpl: pub Ofstream::open: String -> Ofstream
// @bpl: pub Ofstream::is_open: &Ofstream -> bool
// @bpl: pub Ofstream::close: &Ofstream -> ()
// @bpl: pub Ofstream::write: (&Ofstream, String) -> ()

using Ofstream = std::ofstream;
struct Ofstream_ {
  Ofstream_() = delete;

  static inline std::ofstream open(const std::string& filename) {
    return std::ofstream(filename);
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
