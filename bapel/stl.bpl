module bapel.stl

imports {
  bapel.core
}

impls {
  "stl_deque.h"
  "stl_filesystem.h"
  "stl_fstream.h"
  "stl_optional.h"
  "stl_sstream.h"
  "stl_string.h"
  "stl_vector.h"
}

pub type String
pub type StringView
pub type Vector ['a]
pub type Optional ['a]
pub type Deque ['a]
pub type Ofstream
pub type Ifstream

impl String {
  fn empty(s: &Self) -> bool {
    StringImpl::empty s
  }
  fn front(s: &Self) -> i8 {
    StringImpl::front s
  }
  fn size(s: &Self) -> i64 {
    StringImpl::size s
  }
  fn view(s: &Self) -> StringView {
    StringImpl::view s
  }
  fn from_view(v: StringView) -> String {
    StringImpl::from_view v
  }
  fn from_char(c: i8) -> String {
    StringImpl::from_char c
  }
  fn concat(a: &Self, b: &String) -> String {
    StringImpl::concat (a, b)
  }
  fn find(s: &Self, target: &String, pos: i64) -> i64 {
    StringImpl::find (s, target, pos)
  }
  fn replace(s: &Self, pos: i64, count: i64, to: &String) -> String {
    StringImpl::replace (s, pos, count, to)
  }
  fn starts_with(s: &Self, pref: &String) -> bool {
    StringImpl::starts_with (s, pref)
  }
  fn ends_with(s: &Self, suff: &String) -> bool {
    StringImpl::ends_with (s, suff)
  }
}

impl StringView {
  fn at(s: Self, i: i64) -> i8 {
    StringViewImpl::at (s, i)
  }
  fn empty(s: Self) -> bool {
    StringViewImpl::empty s
  }
  fn front(s: Self) -> i8 {
    StringViewImpl::front s
  }
  fn size(s: Self) -> i64 {
    StringViewImpl::size s
  }
  fn substr(s: Self, pos: i64, size: i64) -> StringView {
    StringViewImpl::substr (s, pos, size)
  }
  fn to_string(s: Self) -> String {
    String::from_view s
  }
  fn starts_with(s: Self, pref: StringView) -> bool {
    StringViewImpl::starts_with (s, pref)
  }
  fn ends_with(s: Self, suff: StringView) -> bool {
    StringViewImpl::ends_with (s, suff)
  }
}

impl ['a] (Vector 'a) {
  fn mk() -> Vector 'a {
    VectorImpl::mk ()
  }
  fn push_back(v: &Self, val: 'a) -> () {
    VectorImpl::push_back (v, val)
  }
  fn size(v: &Self) -> i64 {
    VectorImpl::size v
  }
  fn get(v: &Self, index: i64) -> 'a {
    VectorImpl::get (v, index)
  }
  fn set_at(v: &Self, index: i64, val: 'a) -> () {
    VectorImpl::set (v, index, val)
  }
}

impl ['a] (Optional 'a) {
  fn none() -> Optional 'a {
    OptionalImpl::none ()
  }
  fn make_optional(val: 'a) -> Optional 'a {
    OptionalImpl::make_optional val
  }
  fn has_value(opt: &Self) -> bool {
    OptionalImpl::has_value opt
  }
  fn get_value(opt: &Self) -> 'a {
    OptionalImpl::get_value opt
  }
}

impl ['a] (Deque 'a) {
  fn mk() -> Deque 'a {
    DequeImpl::mk ()
  }
  fn push_back(d: &Self, val: 'a) -> () {
    DequeImpl::push_back (d, val)
  }
  fn pop_front(d: &Self) -> () {
    DequeImpl::pop_front d
  }
  fn front(d: &Self) -> 'a {
    DequeImpl::front d
  }
  fn empty(d: &Self) -> bool {
    DequeImpl::empty d
  }
}

impl Ofstream {
  fn open(filename: &String) -> Ofstream {
    OfstreamImpl::open filename
  }
  fn is_open(f: &Self) -> bool {
    OfstreamImpl::is_open f
  }
  fn close(f: &Self) -> () {
    OfstreamImpl::close f
  }
  fn write(f: &Self, s: String) -> () {
    OfstreamImpl::write (f, s)
  }
}

impl Ifstream {
  fn open(filename: &String) -> Ifstream {
    IfstreamImpl::open filename
  }
  fn is_open(f: &Self) -> bool {
    IfstreamImpl::is_open f
  }
  fn close(f: &Self) -> () {
    IfstreamImpl::close f
  }
  fn read['a](f: &Self, val: &'a) -> bool {
    IfstreamImpl::read ['a] (f, val)
  }
}

