#pragma once

template <typename T>
struct Ptr;

template <typename T>
struct Vector;

// @bpl: pub vector_get: forall ['a] (&Vector 'a, i64) -> 'a
template <typename T>
T vector_get(Ptr<Vector<T>>, int64_t);
