module bapel.stl

imports {
  bapel.core
}

impls {
  "stl_deque.h"
  "stl_filesystem.h"
  "stl_fstream.h"
  "stl_map.h"
  "stl_memory.h"
  "stl_optional.h"
  "stl_set.h"
  "stl_sstream.h"
  "stl_string.h"
  "stl_unordered_map.h"
  "stl_unordered_set.h"
  "stl_vector.h"
}

pub type String
pub type StringView
pub type Vector ['a]
pub type Optional ['a]
pub type Deque ['a]
pub type Map ['k, 'v]
pub type Set ['a]
pub type UnorderedMap ['k, 'v]
pub type UnorderedSet ['a]
pub type Ofstream
pub type Ifstream
pub type OStringStream
pub type SharedPtr ['a]
pub type UniquePtr ['a]

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
    StringImpl::from_view s
  }
  fn starts_with(s: Self, pref: StringView) -> bool {
    StringViewImpl::starts_with (s, pref)
  }
  fn ends_with(s: Self, suff: StringView) -> bool {
    StringViewImpl::ends_with (s, suff)
  }
  fn remove_prefix(s: &Self, n: i64) -> () {
    StringViewImpl::remove_prefix (s, n)
  }
  fn remove_suffix(s: &Self, n: i64) -> () {
    StringViewImpl::remove_suffix (s, n)
  }
  fn trim_prefix(s: &Self, pref: StringView) -> bool {
    if StringView::starts_with (*s, pref) {
      StringView::remove_prefix (s, StringView::size pref);
      return true
    }
    false
  }
  fn trim_suffix(s: &Self, suff: StringView) -> bool {
    if StringView::ends_with (*s, suff) {
      StringView::remove_suffix (s, StringView::size suff);
      return true
    }
    false
  }
}

impl String {
  fn empty(s: &Self) -> bool {
    StringImpl::empty s
  }
  fn front(s: &Self) -> i8 {
    StringImpl::front s
  }
  fn at(s: &Self, i: i64) -> i8 {
    StringImpl::at (s, i)
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
  fn append(s: &Self, other: StringView) -> () {
    StringImpl::append (s, other)
  }
  fn push_back(s: &Self, c: i8) -> () {
    StringImpl::push_back (s, c)
  }
  fn reserve(s: &Self, cap: i64) -> () {
    StringImpl::reserve (s, cap)
  }
  fn substr(s: &Self, pos: i64, len: i64) -> String {
    StringImpl::substr (s, pos, len)
  }
  fn find(s: &Self, target: &String, pos: i64) -> i64 {
    StringImpl::find (s, target, pos)
  }
  // TODO: rfind should return an Optional i64 instead.
  fn rfind(s: &Self, target: &String) -> i64 {
    StringImpl::rfind (s, target)
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
  fn trim_prefix(s: &Self, pref: &String) -> bool {
    if String::starts_with (s, pref) {
      let sv: StringView = String::view s;
      StringView::remove_prefix (&sv, String::size pref);
      Ptr::set (s, StringView::to_string sv);
      return true
    }
    false
  }
  fn trim_suffix(s: &Self, suff: &String) -> bool {
    if String::ends_with (s, suff) {
      let sv: StringView = String::view s;
      StringView::remove_suffix (&sv, String::size suff);
      Ptr::set (s, StringView::to_string sv);
      return true
    }
    false
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
  fn sort(v: &Self) -> () {
    VectorImpl::sort v
  }
  fn dedup(v: &Self) -> () {
    VectorImpl::dedup v
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

impl OStringStream {
  fn mk() -> OStringStream {
    OStringStreamImpl::mk ()
  }
  fn write(oss: &Self, s: StringView) -> () {
    OStringStreamImpl::write (oss, s)
  }
  fn write_char(oss: &Self, c: i8) -> () {
    OStringStreamImpl::write_char (oss, c)
  }
  fn str(oss: &Self) -> String {
    OStringStreamImpl::str oss
  }
}

impl ['k, 'v] (Map 'k 'v) {
  fn mk() -> Map 'k 'v {
    MapImpl::mk ()
  }
  fn insert(m: &Self, key: 'k, val: 'v) -> () {
    MapImpl::insert (m, key, val)
  }
  fn remove(m: &Self, key: &'k) -> bool {
    MapImpl::remove (m, key)
  }
  fn clear(m: &Self) -> () {
    MapImpl::clear m
  }
  fn size(m: &Self) -> i64 {
    MapImpl::size m
  }
  fn empty(m: &Self) -> bool {
    MapImpl::empty m
  }
  fn contains(m: &Self, key: &'k) -> bool {
    MapImpl::contains (m, key)
  }
  fn get(m: &Self, key: &'k) -> Optional 'v {
    MapImpl::get (m, key)
  }
  fn keys(m: &Self) -> Vector 'k {
    MapImpl::keys m
  }
  fn entries(m: &Self) -> Vector ('k, 'v) {
    MapImpl::entries m
  }
}

impl ['a] (Set 'a) {
  fn mk() -> Set 'a {
    SetImpl::mk ()
  }
  fn insert(s: &Self, val: 'a) -> bool {
    SetImpl::insert (s, val)
  }
  fn remove(s: &Self, val: &'a) -> bool {
    SetImpl::remove (s, val)
  }
  fn clear(s: &Self) -> () {
    SetImpl::clear s
  }
  fn size(s: &Self) -> i64 {
    SetImpl::size s
  }
  fn empty(s: &Self) -> bool {
    SetImpl::empty s
  }
  fn contains(s: &Self, val: &'a) -> bool {
    SetImpl::contains (s, val)
  }
  fn values(s: &Self) -> Vector 'a {
    SetImpl::values s
  }
}

impl ['k, 'v] (UnorderedMap 'k 'v) {
  fn mk() -> UnorderedMap 'k 'v {
    UnorderedMapImpl::mk ()
  }
  fn insert(m: &Self, key: 'k, val: 'v) -> () {
    UnorderedMapImpl::insert (m, key, val)
  }
  fn size(m: &Self) -> i64 {
    UnorderedMapImpl::size m
  }
  fn empty(m: &Self) -> bool {
    UnorderedMapImpl::empty m
  }
  fn contains(m: &Self, key: &'k) -> bool {
    UnorderedMapImpl::contains (m, key)
  }
  fn get(m: &Self, key: &'k) -> Optional 'v {
    UnorderedMapImpl::get (m, key)
  }
}

impl ['a] (UnorderedSet 'a) {
  fn mk() -> UnorderedSet 'a {
    UnorderedSetImpl::mk ()
  }
  fn insert(s: &Self, val: 'a) -> bool {
    UnorderedSetImpl::insert (s, val)
  }
  fn remove(s: &Self, val: &'a) -> bool {
    UnorderedSetImpl::remove (s, val)
  }
  fn clear(s: &Self) -> () {
    UnorderedSetImpl::clear s
  }
  fn size(s: &Self) -> i64 {
    UnorderedSetImpl::size s
  }
  fn empty(s: &Self) -> bool {
    UnorderedSetImpl::empty s
  }
  fn contains(s: &Self, val: &'a) -> bool {
    UnorderedSetImpl::contains (s, val)
  }
  fn values(s: &Self) -> Vector 'a {
    UnorderedSetImpl::values s
  }
}

impl ['a] (SharedPtr 'a) {
  fn make_shared(val: 'a) -> SharedPtr 'a {
    SharedPtrImpl::make_shared val
  }
  fn null() -> SharedPtr 'a {
    SharedPtrImpl::null ()
  }
  fn is_null(ptr: &Self) -> bool {
    SharedPtrImpl::is_null ptr
  }
  fn get(ptr: &Self) -> 'a {
    SharedPtrImpl::get ptr
  }
  fn use_count(ptr: &Self) -> i64 {
    SharedPtrImpl::use_count ptr
  }
}

impl ['a] (UniquePtr 'a) {
  fn make_unique(val: 'a) -> UniquePtr 'a {
    UniquePtrImpl::make_unique val
  }
  fn null() -> UniquePtr 'a {
    UniquePtrImpl::null ()
  }
  fn is_null(ptr: &Self) -> bool {
    UniquePtrImpl::is_null ptr
  }
  fn get(ptr: &Self) -> 'a {
    UniquePtrImpl::get ptr
  }
  fn reset(ptr: &Self) -> () {
    UniquePtrImpl::reset ptr
  }
}



