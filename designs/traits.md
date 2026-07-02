# Traits for Bapel

This document proposes the design of traits (also known as typeclasses) for Bapel. Traits enable ad-hoc polymorphism, allowing different types to implement a common interface.

## Status

*   **Inherent Implementations**: Fully Implemented (including generic types).
*   **Trait Declarations & Implementations**: Fully Implemented.
*   **Coherence (Same-File Rule)**: Fully Implemented.
*   **Trait Bounds (Generic Constraints)**: Pending (Next Step).
*   **C++ Translation**: Fully Implemented for Trait/Inherent Impls. Pending for Trait Bounds.

## Motivation

Currently, Bapel has generic functions (parametric polymorphism) but lacks a way to constrain generic types to those that support certain operations. For example, we have `String_::size`, `Vector_::size`, and `Deque_::empty`, but we cannot write a generic function that works on any "container" that has a size or can be checked for emptiness.

Additionally, the proposed Runtime Linear Memory Model (LMM) requires a `Send` capability to ensure thread safety. Traits provide a natural way to model this capability.

## Proposed Design

### Inherent Implementation

An `impl` block without a trait name defines "inherent methods" directly for a type. In the scope of the `impl` block, `Self` refers to the type being implemented.

For the MVP, method call syntax (e.g., `x.size()`) is **not supported**. Methods must be called using their fully qualified names (e.g., `String::size &s`).

#### Inherent Methods for C++ Types (Opaque Types)

For types implemented in C++ (like `String` or `Vector`), we define the Bapel-side interface using `impl` blocks in `.bpl` files, and they call the underlying C++ implementations.

The C++ implementations are declared in the C++ header files using `@bpl` annotations, typically under a `TypeImpl` namespace to avoid conflicts with the Bapel method names.

For example, in `bapel/stl_vector.h`:

```cpp
template <typename T>
using Vector = std::vector<T>;

// @bpl: pub VectorImpl::mk: forall ['a] () -> Vector 'a
// @bpl: pub VectorImpl::push_back: forall ['a] (&Vector 'a, 'a) -> ()
```

And in `bapel/stl.bpl`:

```bapel
pub type Vector ['a]

impl ['a] (Vector 'a) {
  fn mk() -> Vector 'a {
    VectorImpl::mk ()
  }
  fn push_back(v: &Self, val: 'a) -> () {
    VectorImpl::push_back (v, val)
  }
}
```

In Bapel code, these are called as:

```bapel
let len: i64 = String::size &s;
let v_len: i64 = Vector::size &v;
```

### Trait Declaration

A trait defines a set of function signatures (methods) that a type must implement.

```bapel
trait Size {
  fn size(s: &Self) -> i64
}

trait Indexable ['elem] {
  fn get(v: &Self, index: i64) -> 'elem
}
```

### Trait Implementation

An `impl` block with a trait name defines the implementation of a trait for a specific type using the `for` keyword.

```bapel
impl Size for String {
  fn size(s: &Self) -> i64 {
    String::size s
  }
}

impl ['a] Size for (Vector 'a) {
  fn size(v: &Self) -> i64 {
    Vector::size v
  }
}

impl ['a] Indexable 'a for (Vector 'a) {
  fn get(v: &Self, index: i64) -> 'a {
    Vector::get (v, index)
  }
}
```


### Coherence (Same-File Rule)

To prevent conflicting implementations and ensure sound C++ code generation, Bapel enforces a strict coherence rule:

**A trait implementation (`impl Trait for Type`) must be defined in the same file as the `Type` itself.**

This implies:
1.  You can implement a foreign trait for a local type (since the `impl` is in the local type's file).
2.  You **cannot** implement a local trait for a foreign type (since you cannot modify the foreign type's file).
3.  Implementing a foreign trait for a foreign type is disallowed.

*   **Primitives and Tuples Exception (MVP):** Since primitive types (e.g., `i32`, `i64`) and tuples do not have a defining source file, they cannot have manual `impl` blocks. For the MVP, any trait implementations for primitive types and tuples are handled internally by the compiler (either as built-ins or via automatic derivation). Users cannot implement custom traits for primitives or tuples.

This simplifies compiler resolution and ensures that a type's capabilities are always bundled with its definition.


### Name Uniqueness

To support clean C++ translation without prefixes, Bapel enforces that traits, structs, variants, and type aliases share the same type namespace within a module. You cannot define a trait and a struct with the same name in the same module.

### Trait Bounds (Generic Constraints)

Generic functions can be constrained using trait bounds.

```bapel
fn printSize['t: Size](x: & 't) -> () {
  core::print [i64] (Size::size x)
}
```

If multiple bounds are needed:

```bapel
fn printElementIfLarge['t: Size + Indexable 'elem, 'elem](x: & 't, index: i64) -> () {
  if Size::size x > 10 {
    core::print ['elem] (Indexable::get (x, index))
  }
}
```


## C++ Translation (C++17)

Bapel targets C++17. Traits are translated to C++ using template specialization and SFINAE (Substitution Failure Is Not An Error) to enforce constraints.

### Inherent Methods Translation

For types implemented in C++ (opaque types), we cannot define static methods directly on the type if it is a C++ type alias (e.g., `using String = std::string`).

To avoid naming conflicts in C++, the inherent methods are defined in a helper `struct` named with a `_` suffix (e.g., `String_`, `Vector_`).

#### Non-templated Types
For non-templated types like `String`, the helper is a non-templated struct:

```cpp
struct String_ {
  String_() = delete;
  static inline bool empty(const String& s) { ... }
  static inline int64_t size(const String& s) { ... }
};
```

#### Templated Types
For templated types like `Vector`, the helper is a templated struct:

```cpp
template <typename T>
struct Vector_ {
  Vector_() = delete;
  static inline int64_t size(const Vector<T>* v) { ... }
};
```

#### Compiler Mapping
The Bapel compiler automatically translates Bapel `Type::method` calls to C++ `::Type_::method` calls.
*   For non-templated types: `String::size &s` -> `::String_::size(s)`
*   For templated types: `Vector::size &v` -> `::Vector_<T>::size(v)` (the compiler automatically splits the type arguments, applying the type-level arguments to the `Vector_` struct template).

### Trait Representation

Each trait translates to a C++ template struct. The struct contains static member functions for each method in the trait.

```cpp
// Trait Size (implicit Self)
template <typename Self>
struct Size;
```

### Trait Implementation

An `impl` block translates to a template specialization of the trait struct for the concrete types.

```cpp
// impl Size for String
template <>
struct Size<String> {
  static int64_t size(const String& s) {
    return String_::size(s);
  }
};

// impl ['a] Size for (Vector 'a)
template <typename T>
struct Size<std::vector<T>> {
  static int64_t size(const std::vector<T>& v) {
    return Vector_::size(v);
  }
};
```

### Trait Bounds (SFINAE)

Generic functions with trait bounds use `std::enable_if_t` to ensure the template is only instantiated if the type implements the trait.

To enforce the constraint at the function signature level (for better error messages and overloading), we can use SFINAE:

```cpp
// Helper to check if trait is implemented (via completeness of Size)
template <typename T, typename = void>
struct Size_is_implemented : std::false_type {};

template <typename T>
struct Size_is_implemented<T, std::void_t<
  decltype(sizeof(Size<T>))
>> : std::true_type {};

template <typename T>
inline constexpr bool Size_is_implemented_v = Size_is_implemented<T>::value;

// Generic function with bound: fn printSize['t: Size](x: & 't)
template <typename T, typename = std::enable_if_t<Size_is_implemented_v<T>>>
void printSize(const T& x) {
  print(Size<T>::size(x));
}
```

If we don't need complex SFINAE for overloading, we can just let it fail during instantiation:

```cpp
template <typename T>
void printSize(const T& x) {
  // If Size<T> is not specialized, this will fail to compile.
  print(Size<T>::size(x));
}
```

## Alternatives Considered: C++20 Concepts

If Bapel were to target C++20 in the future, traits could map directly to concepts, which are much cleaner and provide better compiler errors.

```cpp
template <typename T>
concept Size = requires(const T& t) {
  { size(t) } -> std::same_as<int64_t>;
};

template <Size T>
void printSize(const T& x) {
  print(size(x));
}
```

We choose the C++17 approach for now to maximize compatibility with existing C++ toolchains.

## Future Work (Post-MVP)

The following features are deferred to post-MVP:

1.  **Associated Types**: Introducing associated types (like Rust's `Iterator::Item`) to simplify signatures.
    *   For example, instead of `trait Indexable ['elem] { get: (&Self, i64) -> 'elem }`, we could have `trait Indexable { type Elem; get: (&Self, i64) -> Elem }`.
2.  **Deriving**: Supporting `deriving` for common traits (e.g., `Eq`, `Show`) to reduce boilerplate.
3.  **Marker Traits and Auto-Derivation**: Supporting marker traits (capabilities like `Send`) and their automatic derivation based on structural conformance (e.g., a struct is `Send` if all its fields are `Send`), as well as explicit opt-out (e.g., `impl !Send`).
4.  **Blanket and Overlapping Implementations**: Supporting blanket implementations (implementing a trait for all types that implement another trait, e.g., `impl ['t: Foo] Bar for 't`) and the complex coherence checking required to prevent overlaps.
5.  **Generic Methods in Traits**: Supporting trait methods that have their own type parameters, distinct from the trait's parameters (e.g., `fn map['u](s: &Self, f: (Self) -> 'u) -> 'u`).

## Open Questions and Refinements

During the review of this design, several open questions and areas for refinement were identified:

### 1. Trait Method Invocation Syntax
*   **Accepted Design (MVP):**
    *   All trait and implementation functions have a **fully qualified name only** (e.g., `Size::size`, `String::size`).
    *   To call these functions, the fully qualified name must be used: `Size::size x` or `String::size x`.
    *   **Generic Fully Qualified Calls:** If a trait is parameterized, the fully qualified call treats all type parameters (the implicit `'self` and trait-level parameters) as a flat list of type arguments passed to the function: `Trait::method ['self, 'trait_params...] (receiver, args...)`.
        *   Example: `Indexable::get [Vector i8, i8] (v, 0)`
        *   Most of the time, these type arguments can be omitted as the compiler can infer them from the value arguments (e.g., `Indexable::get (v, 0)`).
    *   **Method call syntax (`x.method(args)`) is deferred to post-MVP.**
*   **Discussion:** Deferring method syntax significantly simplifies the MVP compiler by eliminating the need for complex method resolution rules and automatic receiver adjustment (borrowing/dereferencing). It ensures that all function calls are explicit and unambiguous.

### 2. Method Resolution and Conflict Handling
*   **Accepted Design (MVP):**
    *   Because method call syntax is deferred, there is **no conflict resolution needed at call sites** in the MVP.
    *   Inherent methods (e.g., `String::size`) and trait methods (e.g., `Size::size`) are distinct fully qualified names.
    *   If a type has an inherent method and implements a trait method with the same name, they are called explicitly: `String::size s` vs `Size::size s`.
    *   Resolution rules for desugaring will be defined post-MVP when method syntax is introduced.

### 3. Trait Scope and Imports
*   **Accepted Design (MVP):**
    *   Because of the **Same-File Rule**, trait implementations are automatically loaded whenever the type they are implemented for is in scope. There is no need to import implementations separately.
    *   To call a trait method using its fully qualified name (e.g., `Size::size x`), the **Trait itself (e.g., `Size`) must be in scope** (explicitly imported or defined in the current file).

### 4. Multi-parameter and Parameterized Traits Translation
*   **Accepted Design:**
    *   Parameterized traits (like `Indexable ['elem]`) translate to C++ template structs with multiple template parameters, representing the implicit `Self` and the explicit trait parameters.
    *   **Trait Representation:**
        ```cpp
        template <typename Self, typename Elem>
        struct Indexable;
        ```
    *   **Trait Implementation (Specialization):**
        An `impl` block translates to a template specialization binding both the type and the trait parameters:
        ```cpp
        template <typename T>
        struct Indexable<std::vector<T>, T> {
          static T get(const std::vector<T>& v, int64_t index) {
            return Vector_::get(v, index);
          }
        };
        ```
    *   **SFINAE Helpers:**
        The helper must also be parameterized by the trait parameters to verify the implementation exists for that specific combination:
        ```cpp
        template <typename T, typename Elem, typename = void>
        struct Indexable_is_implemented : std::false_type {};
 
        template <typename T, typename Elem>
        struct Indexable_is_implemented<T, Elem, std::void_t<
          decltype(sizeof(Indexable<T, Elem>))
        >> : std::true_type {};

        template <typename T, typename Elem>
        inline constexpr bool Indexable_is_implemented_v = Indexable_is_implemented<T, Elem>::value;
        ```
    *   **Generic Function Translation:**
        Multiple bounds are combined using `std::enable_if_t` and the helper variables:
        ```cpp
        template <
            typename T, 
            typename Elem, 
            typename = std::enable_if_t<
                Size_is_implemented_v<T> && 
                Indexable_is_implemented_v<T, Elem>
            >
        >
        void printElementIfLarge(const T& x, int64_t index) {
          if (Size<T>::size(x) > 10) {
            print(Indexable<T, Elem>::get(x, index));
          }
        }
        ```

### 5. Static vs. Dynamic Dispatch
*   **Accepted Design:**
    *   For the MVP, traits in Bapel are **strictly for static dispatch (monomorphization)**.
    *   All generic functions using trait bounds are monomorphized at compile time.
    *   There is no support for trait objects or runtime dynamic dispatch.
    *   This ensures zero runtime overhead, matching Bapel's performance goals and simplifying the C++ translation (which maps directly to C++ templates).


### 6. Compiler-Side Verification vs. C++ SFINAE
*   **Accepted Design:**
    *   The Bapel compiler will perform **full eager verification** of trait bounds and method resolution.
    *   It will not rely on the C++ compiler to catch trait-related type errors via SFINAE or template instantiation failures.
    *   If a type does not implement a required trait, the Bapel compiler will emit a clean, readable, Bapel-specific error message pointing to the source code.
    *   This is made highly feasible for the MVP due to other simplifying constraints (Same-File Rule, no blanket impls, static dispatch only), which reduce trait resolution to a simple lookup in the type's defining file.


