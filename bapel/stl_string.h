#pragma once

#include <string>
#include <string_view>

// @bpl: pub type String
// @bpl: pub type StringView

using String = std::string;
using StringView = std::string_view;

// @bpl: pub String_::empty: String -> bool
// @bpl: pub String_::front: String -> i8
// @bpl: pub String_::size: String -> i64
namespace String_ {

inline bool empty(String s) { return s.empty(); }
inline char front(String s) { return s.front(); }
inline int64_t size(String s) { return s.size(); }

}  // namespace String_

// @bpl: pub StringView_::at: (StringView, i64) -> i8
// @bpl: pub StringView_::empty: StringView -> bool
// @bpl: pub StringView_::front: StringView -> i8
// @bpl: pub StringView_::size: StringView -> i64
// @bpl: pub StringView_::substr: (StringView, i64, i64) -> StringView
namespace StringView_ {

inline char at(StringView s, int64_t i) { return s.at(i); }
inline bool empty(StringView s) { return s.empty(); }
inline char front(StringView s) { return s.front(); }
inline int64_t size(StringView s) { return s.size(); }

inline StringView substr(StringView s, int64_t pos, int64_t size) {
  if (pos > s.size()) {
	  return StringView();
  }
	return s.substr(pos, size);
}

}  // namespace StringView_
