#pragma once

template <typename T>
struct Vector;

// @bpl: pub vector_get: forall ['a] (&Vector 'a, i64) -> 'a
template <typename T>
T vector_get(Vector<T>*, int64_t);
