#pragma once

#include <string>
#include <string_view>

// @bpl: pub type String
// @bpl: pub type StringView
using String = std::string;
using StringView = std::string_view;

// @bpl: pub to_string: forall ['a] 'a -> String
template <typename T>
inline String to_string(T value) { return std::to_string(value); }

// @bpl: pub String_::empty: String -> bool
// @bpl: pub String_::front: String -> i8
// @bpl: pub String_::size: String -> i64
// @bpl: pub String_::view: String -> StringView
// @bpl: pub String_::from_view: StringView -> String
// @bpl: pub String_::from_char: i8 -> String
// @bpl: pub String_::concat: (String, String) -> String
namespace String_ {

inline bool empty(String& s) { return s.empty(); }
inline char front(String& s) { return s.front(); }
inline int64_t size(String& s) { return s.size(); }
inline StringView view(const String& s) { return s; }
inline String from_view(StringView v) { return String(v); }
inline String from_char(char c) { return String(1, c); }
inline String concat(const String& a, const String& b) { return a + b; }

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
