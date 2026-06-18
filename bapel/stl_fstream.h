#pragma once

#include <fstream>
#include <string>
#include <variant>

// @bpl: pub type Ofstream
// @bpl: pub Ofstream_::open: String -> Ofstream
// @bpl: pub Ofstream_::is_open: &Ofstream -> bool
// @bpl: pub Ofstream_::close: &Ofstream -> ()
// @bpl: pub Ofstream_::write: (&Ofstream, String) -> ()

using Ofstream = std::ofstream;
namespace Ofstream_ {

inline std::ofstream open(const std::string& filename) {
  return std::ofstream(filename);
}

inline bool is_open(std::ofstream* f) {
  return f->is_open();
}

inline std::monostate close(std::ofstream* f) {
  f->close();
  return std::monostate();
}

inline std::monostate write(std::ofstream* f, const std::string& s) {
  *f << s;
  return std::monostate();
}

}  // namespace Ofstream_
