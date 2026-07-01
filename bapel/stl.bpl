module bapel.stl

imports {
  bapel.core
}

impls {
  "stl_deque.h"
  "stl_filesystem.h"
  "stl_fstream.h"
  "stl_sstream.h"
  "stl_string.h"
  "stl_vector.h"
}

pub type String
pub type StringView

pub trait String_ {
  fn empty(s: Self) -> bool
  fn front(s: Self) -> i8
  fn size(s: Self) -> i64
  fn view(s: Self) -> StringView
  fn from_view(v: StringView) -> Self
  fn from_char(c: i8) -> Self
  fn concat(a: Self, b: Self) -> Self
  fn find(s: Self, target: Self, pos: i64) -> i64
  fn replace(s: Self, pos: i64, count: i64, to: Self) -> Self
}

impl String_ for String {
  fn empty(s: String) -> bool {
    StringImpl::empty s
  }
  fn front(s: String) -> i8 {
    StringImpl::front s
  }
  fn size(s: String) -> i64 {
    StringImpl::size s
  }
  fn view(s: String) -> StringView {
    StringImpl::view s
  }
  fn from_view(v: StringView) -> String {
    StringImpl::from_view v
  }
  fn from_char(c: i8) -> String {
    StringImpl::from_char c
  }
  fn concat(a: String, b: String) -> String {
    StringImpl::concat (a, b)
  }
  fn find(s: String, target: String, pos: i64) -> i64 {
    StringImpl::find (s, target, pos)
  }
  fn replace(s: String, pos: i64, count: i64, to: String) -> String {
    StringImpl::replace (s, pos, count, to)
  }
}
