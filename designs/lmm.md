# Runtime Linear Memory Model for Bapel

This document proposes a design for a runtime linear memory model for Bapel. This model aims to track memory allocations and enforce linearity constraints at runtime, detecting errors such as use-after-free, double-free, and use-after-move.

Since Bapel compiles to C++, this model is designed to be implemented as a C++ runtime library that replaces the current raw pointer implementation.

## Core Concepts

1.  **Owned Pointer (`OwnedPtr<T>`)**: A wrapper around a raw pointer that tracks ownership and state. It is linear (move-only) and manages the resource lifetime.
2.  **Borrow Pointer (`BorrowPtr<T>`)**: A temporary pointer that points to a resource owned by an `OwnedPtr<T>`. It does not own the resource and cannot free it. Multiple borrows can exist.
3.  **Resource State**: A resource managed by `OwnedPtr` can be in one of the following states:
    *   `VALID`: The resource is owned and can be accessed.
    *   `MOVED`: Ownership of the resource has been transferred to another `OwnedPtr`. Accessing this pointer is an error.
    *   `FREED`: The resource has been deallocated. Accessing this pointer is an error.
    *   `UNINITIALIZED`: The pointer has been default-constructed and does not yet point to a valid resource. Accessing this pointer is an error.
4.  **Borrow Tracking**: We track the number of active borrows for each `OwnedPtr`. Moving or freeing a resource while it has active borrows is a runtime error.
5.  **Nullability**: Pointers in Bapel are non-nullable by default. Optional pointers are represented using Bapel's `Maybe` variant type (e.g., `Maybe (Ptr T)`), ensuring null-safety at compile-time via pattern matching. The C++ runtime pointer classes do not support a null state.
6.  **Stack Borrows**: Borrows of stack-allocated (local) variables are allowed. The Bapel compiler statically guarantees that these borrows do not escape their lexical scope. Consequently, the C++ runtime does not track lifetimes of stack borrows (they have a null owner) and relies on compiler safety guarantees.

## Syntactic Sugar

To make the memory model easier to use, Bapel introduces syntactic sugar for these pointer types:

*   **`&&T`** denotes `OwnedPtr T`.
*   **`&T`** denotes `BorrowPtr T`.

For example:
*   `let p: &&i32 = Ptr_::alloc 10;` (declares an owned pointer to `i32`).
*   `let b: &i32 = Ptr_::borrow p;` (declares a borrow of `p`).

The Bapel compiler will map `&&T` to C++ `OwnedPtr<T>` and `&T` to C++ `BorrowPtr<T>`.


## Runtime Library Design (C++)

We replace the raw pointer implementation of `Ptr<T>` with two smart pointer classes: `OwnedPtr<T>` (representing ownership) and `BorrowPtr<T>` (representing a borrow).

### State Diagram (for `OwnedPtr`)

```mermaid
stateDiagram-v2
    [*] --> VALID : Allocation (Ptr_::alloc)
    VALID --> MOVED : Move / Transfer Ownership
    VALID --> FREED : Free / Deallocation (Ptr_::free)
    MOVED --> [*] : Destructor (no-op)
    FREED --> [*] : Destructor (no-op)
```

### C++ Implementation (`owned_ptr.h` and `borrow_ptr.h`)

```cpp
#include <iostream>
#include <cassert>

#ifndef NDEBUG
#define BAPEL_MEM_CHECK(expr) expr
#else
#define BAPEL_MEM_CHECK(expr)
#endif

#ifdef BAPEL_TEST_THROW
#include <stdexcept>
class BapelMemoryError : public std::runtime_error {
public:
    explicit BapelMemoryError(const char* msg) : std::runtime_error(msg) {}
};
inline void runtime_error(const char* msg) {
    throw BapelMemoryError(msg);
}
#else
inline void runtime_error(const char* msg) {
    std::cerr << "Bapel Runtime Memory Error: " << msg << std::endl;
    std::abort();
}
#endif

template <typename T>
class BorrowPtr;

template <typename T>
class OwnedPtr {
public:
#ifndef NDEBUG
    enum State { VALID, MOVED, FREED, UNINITIALIZED };
#endif

    OwnedPtr() : ptr_(nullptr) {
        BAPEL_MEM_CHECK(state_ = UNINITIALIZED; borrows_ = 0;);
    }
    
    explicit OwnedPtr(T* ptr) : ptr_(ptr) {
        BAPEL_MEM_CHECK(state_ = VALID; borrows_ = 0;);
    }

    // Move constructor
    OwnedPtr(OwnedPtr&& other) noexcept {
#ifndef NDEBUG
        if (other.borrows_ > 0) {
            runtime_error("Cannot move a borrowed resource");
        }
        state_ = other.state_;
        borrows_ = 0;
        if (other.state_ == VALID) {
            other.state_ = MOVED;
        }
#endif
        ptr_ = other.ptr_;
        other.ptr_ = nullptr;
    }

    // Move assignment
    OwnedPtr& operator=(OwnedPtr&& other) noexcept {
        if (this != &other) {
#ifndef NDEBUG
            if (state_ == VALID) {
                runtime_error("Resource leaked: overwritten without consumption");
            }
            if (other.borrows_ > 0) {
                runtime_error("Cannot move a borrowed resource");
            }
            state_ = other.state_;
            borrows_ = 0;
            if (other.state_ == VALID) {
                other.state_ = MOVED;
            }
#endif
            ptr_ = other.ptr_;
            other.ptr_ = nullptr;
        }
        return *this;
    }

    // Copy is disabled
    OwnedPtr(const OwnedPtr&) = delete;
    OwnedPtr& operator=(const OwnedPtr&) = delete;

    ~OwnedPtr() {
#ifndef NDEBUG
        if (state_ == VALID) {
            if (borrows_ > 0) {
                runtime_error("Owner destroyed while borrowed");
            }
            // Strict linearity: Leaks are errors.
            runtime_error("Resource leaked: went out of scope without being consumed");
        }
#endif
    }

    T& get_ref() const {
        check_valid();
        return *ptr_;
    }

    T* operator->() const {
        check_valid();
        return ptr_;
    }

    T& operator*() const {
        check_valid();
        return *ptr_;
    }

    void deallocate() {
#ifndef NDEBUG
        if (state_ == FREED) {
            runtime_error("Double free detected");
        }
        if (state_ == MOVED) {
            runtime_error("Attempting to free a moved resource");
        }
        if (borrows_ > 0) {
            runtime_error("Cannot free a borrowed resource");
        }
        state_ = FREED;
#endif
        delete ptr_;
        ptr_ = nullptr;
    }

    void add_borrow() const {
        check_valid();
        BAPEL_MEM_CHECK(borrows_++;);
    }

    void remove_borrow() const {
        BAPEL_MEM_CHECK(borrows_--;);
    }

    bool is_valid() const {
#ifndef NDEBUG
        return state_ == VALID;
#else
        return true;
#endif
    }

private:
    T* ptr_;
#ifndef NDEBUG
    State state_;
    mutable int borrows_;
#endif

    void check_valid() const {
#ifndef NDEBUG
        if (state_ == UNINITIALIZED) {
            runtime_error("Attempted to access uninitialized pointer");
        }
        if (state_ == FREED) {
            runtime_error("Use after free detected");
        }
        if (state_ == MOVED) {
            runtime_error("Use after move (ownership transferred) detected");
        }
#endif
    }    
    friend class BorrowPtr<T>;
};

template <typename T>
class BorrowPtr {
public:
    BorrowPtr() : ptr_(nullptr) {
        BAPEL_MEM_CHECK(owner_ = nullptr;);
    }

    // Construct borrow from stack reference (trusted static analysis)
    explicit BorrowPtr(T& stack_ref) : ptr_(&stack_ref) {
        BAPEL_MEM_CHECK(owner_ = nullptr;);
    }

    // Construct borrow from owner
    explicit BorrowPtr(const OwnedPtr<T>& owner) : ptr_(owner.ptr_) {
        BAPEL_MEM_CHECK(
            owner_ = &owner;
            owner_->add_borrow();
        );
    }

    // Copy constructor (create another borrow)
    BorrowPtr(const BorrowPtr& other) : ptr_(other.ptr_) {
        BAPEL_MEM_CHECK(
            owner_ = other.owner_;
            if (owner_) {
                owner_->add_borrow();
            }
        );
    }

    // Copy assignment
    BorrowPtr& operator=(const BorrowPtr& other) {
        if (this != &other) {
            BAPEL_MEM_CHECK(
                if (owner_) {
                    owner_->remove_borrow();
                }
            );
            ptr_ = other.ptr_;
            BAPEL_MEM_CHECK(
                owner_ = other.owner_;
                if (owner_) {
                    owner_->add_borrow();
                }
            );
        }
        return *this;
    }

    // Move constructor (transfer borrow)
    BorrowPtr(BorrowPtr&& other) noexcept : ptr_(other.ptr_) {
        BAPEL_MEM_CHECK(
            owner_ = other.owner_;
            other.owner_ = nullptr;
        );
        other.ptr_ = nullptr;
    }

    // Move assignment
    BorrowPtr& operator=(BorrowPtr&& other) noexcept {
        if (this != &other) {
            BAPEL_MEM_CHECK(
                if (owner_) {
                    owner_->remove_borrow();
                }
                owner_ = other.owner_;
                other.owner_ = nullptr;
            );
            ptr_ = other.ptr_;
            other.ptr_ = nullptr;
        }
        return *this;
    }

    ~BorrowPtr() {
        BAPEL_MEM_CHECK(
            if (owner_) {
                owner_->remove_borrow();
            }
        );
    }

    T& get_ref() const {
        check_valid();
        return *ptr_;
    }

    T* get_ptr() const {
        check_valid();
        return ptr_;
    }

    T* operator->() const {
        check_valid();
        return ptr_;
    }

    T& operator*() const {
        check_valid();
        return *ptr_;
    }

private:
    T* ptr_;
#ifndef NDEBUG
    const OwnedPtr<T>* owner_;
#endif

    void check_valid() const {
#ifndef NDEBUG
        if (owner_ && !owner_->is_valid()) {
            runtime_error("Use of invalid borrow (owner moved or freed)");
        }
#endif
    }
};
```

### Bapel API Changes (`core_pointer.h`)

We redefine the pointer operations to use `OwnedPtr` and `BorrowPtr`.

```cpp
#pragma once
#include "bapel/owned_ptr.h"
#include "bapel/borrow_ptr.h"

// @bpl: pub type Ptr ['a]
template <typename A>
using Ptr = OwnedPtr<A>;

// @bpl: pub type Ref ['a]
template <typename A>
using Ref = BorrowPtr<A>;

// @bpl: pub Ptr_::alloc: forall ['a] 'a -> Ptr 'a
template <typename A>
Ptr<A> alloc(A a) {
    return Ptr<A>(new A(std::move(a)));
}

// @bpl: pub Ptr_::free: forall ['a] Ptr 'a -> ()
template <typename A>
std::monostate free(Ptr<A>& p) {
    p.deallocate();
    return std::monostate();
}

// @bpl: pub Ptr_::get: forall ['a] Ptr 'a -> 'a
template <typename A>
A& get(const Ptr<A>& ptr) {
    return ptr.get_ref();
}

// @bpl: pub Ref_::get: forall ['a] Ref 'a -> 'a
template <typename A>
A& get(const Ref<A>& ptr) {
    return ptr.get_ref();
}

// @bpl: pub Ptr_::borrow: forall ['a] Ptr 'a -> Ref 'a
template <typename A>
Ref<A> borrow(const Ptr<A>& ptr) {
    return Ref<A>(ptr);
}
```

## Single vs. Dual Pointer Type Design

An important design decision is whether to use a single C++ class (`LinearPtr`) for both owned and borrowed pointers, or to split them into two distinct C++ classes (e.g., `OwnedPtr` and `BorrowPtr`).

### Option A: Unified `LinearPtr` (Current Proposal)

In this design, a single class handles both roles, distinguished by a runtime flag (`is_borrow_`).

*   **Pros**:
    *   **Simpler Type Mapping**: Bapel's `Ptr T` maps 1:1 to C++ `LinearPtr<T>`. The Bapel compiler does not need to distinguish between owned and borrowed pointer types in the generated C++ code.
    *   **Easier Generic Code**: Functions that accept pointers can simply take `LinearPtr<T>` regardless of whether they take ownership or just borrow.
*   **Cons**:
    *   **Higher Memory Overhead**: Every pointer must carry the `is_borrow_` flag and the `owner_` pointer, even if it is an owned pointer (where `owner_` is always null).
    *   **Runtime Branching**: Operations like destruction and moves must check `is_borrow_` at runtime to determine their behavior, adding CPU overhead.
    *   **Deferred Safety Checks**: The C++ compiler cannot prevent invalid operations on borrows (like calling `deallocate()`); these can only be caught at runtime.

### Option B: Split `OwnedPtr` and `BorrowPtr`

In this design, we define two separate templates: `OwnedPtr<T>` (which owns the resource and manages borrow counts) and `BorrowPtr<T>` (which points to the resource and references the parent `OwnedPtr`).

*   **Pros**:
    *   **Reduced Memory Overhead**:
        *   `OwnedPtr<T>` only needs `ptr_`, `state_`, and `borrows_` (typically **16 bytes**, saving 16 bytes per owned pointer).
        *   `BorrowPtr<T>` only needs `ptr_` and `owner_` (typically **16 bytes**).
    *   **Compile-Time Safety in C++**: The C++ compiler enforces that `BorrowPtr` cannot be freed (it lacks a `deallocate()` method) and cannot be moved in a way that transfers ownership.
    *   **Fewer Runtime Branches**: No need to check `is_borrow_` at runtime; the type system guarantees the behavior.
*   **Cons**:
    *   **Compiler Complexity**: The Bapel compiler must distinguish between ownership and borrowing at the type level and generate different C++ code accordingly.
        *   Bapel `Ptr T` (owned) -> C++ `OwnedPtr<T>`.
        *   Bapel `&T` (borrowed) -> C++ `BorrowPtr<T>`.
    *   **Signature Proliferation**: Functions must explicitly declare if they accept `OwnedPtr` or `BorrowPtr`.

### Summary Comparison

| Feature | Option A: Unified `LinearPtr` | Option B: Split `OwnedPtr`/`BorrowPtr` |
| :--- | :--- | :--- |
| **C++ Type Safety** | Low (Runtime checks only) | **High (Compile-time enforcement)** |
| **Memory Size (Owned)**| 32 bytes | **16 bytes** |
| **Memory Size (Borrow)**| 32 bytes | **16 bytes** |
| **Runtime Branching** | Yes (checks `is_borrow_`) | **No (determined by type)** |
| **Bapel Compiler Work** | Low | High (needs type distinction) |

On balance, **Option B (Split Types)** is superior for a production-ready compiler as it significantly reduces memory overhead and leverages the C++ compiler to enforce safety rules, but **Option A (Unified Type)** is easier to prototype.

## Integrating with Bapel Compiler

To make this model work seamlessly, the Bapel compiler needs some adjustments:

1.  **Move Semantics by Default**: When a pointer is assigned or passed to a function, the compiler should generate a C++ move (`std::move`).
    ```bapel
    let p2 = p1; // Compiler generates: Ptr p2 = std::move(p1);
    ```
2.  **Borrowing Syntax**: The compiler should automatically insert `Ptr_::borrow` when a pointer is passed to a function that doesn't take ownership (borrows).
    *   Currently, Bapel uses `&vec` in `push_back(&vec, 10)`.
    *   If `vec` is `Vector i8`, then `&vec` is `Ptr (Vector i8)`.
    *   If `&` operator always creates a borrow, then `&vec` should call `Ptr_::borrow` if `vec` is already a pointer, or allocate/wrap if it is a local variable.
    *   If `vec` is a local variable, `&vec` creates a temporary pointer. We need to ensure this temporary pointer doesn't outlive `vec`.

### Compiler Code Generation Logic

To generate `std::move` and manage `OwnedPtr`/`BorrowPtr` transitions, the Bapel compiler uses type information available during the code generation phase (after typechecking).

#### 1. Variable Assignment (`let x = y` or `x <- y`)

When the compiler generates C++ code for an assignment:
*   **Static Type Check**: The compiler checks the Bapel type of the Right-Hand Side (RHS) expression.
*   **Linearity Check**: If the type of RHS is an owned pointer (`&&T`), it is marked as a linear resource.
*   **Code Generation**:
    *   If RHS is `&&T`, generate: `x = std::move(y);`
    *   This calls the C++ move assignment operator of `OwnedPtr<T>`, which invalidates `y` (sets its state to `MOVED` and pointer to `nullptr`).
    *   If RHS is not linear (e.g., a primitive type or a borrow `&T`), generate standard copy: `x = y;`

#### 2. Function Calls (`f(y)`)

When passing arguments to a function:
*   **Signature Matching**: The compiler looks at the signature of the function being called.
*   **Ownership Transfer**: If the function parameter is declared as `&&T` (owned pointer), the function takes ownership.
    *   Generate C++: `f(std::move(y))`
*   **Borrowing**: If the function parameter is declared as `&T` (borrow), the function only borrows the resource.
    *   If `y` is an owned pointer (`&&T`), the compiler must generate a borrow from it.
    *   Generate C++: `f(BorrowPtr<T>(y))` (which calls the borrow constructor, incrementing the borrow count on `y`).
    *   If `y` is already a borrow (`&T`), it can be copied (borrows are copyable).
    *   Generate C++: `f(y)` (calls C++ copy constructor of `BorrowPtr<T>`, which increments the borrow count on the original owner).

#### 3. Address-of Operator (`&x`)

In Bapel, `&x` creates a borrow.
*   If `x` is a local variable (stack allocated):
    *   Generate C++: `BorrowPtr<T>(x)` (assuming we have a C++ constructor that can wrap stack references safely, or we restrict this to heap).
*   If `x` is an owned pointer (`&&T`):
    *   In Bapel, `&p` where `p` is `&&T` creates a borrow.
    *   Generate C++: `BorrowPtr<T>(p)` (calls borrow constructor).

#### 4. Member Access (`ptr.field`)

In Bapel, member access on pointers uses the same dot operator (`.`) as value types. The compiler automatically dereferences the pointer.
*   If `p` is of type `&&T` (owned pointer) or `&T` (borrow pointer):
    *   Generate C++: `p->field` (leveraging the overloaded `operator->` in `OwnedPtr` and `BorrowPtr`).

## Thread Safety and Concurrency

To ensure thread safety and prevent data races, Bapel will restrict how pointers and borrows are transferred across threads.

### The `Send` Capability

We introduce a static classification of types based on whether they can be safely transferred to another thread (the `Send` capability):

1.  **`Send` Types**: Types that can be moved to another thread.
    *   Primitive types (e.g., `i32`, `bool`) are `Send`.
    *   Owned pointers (`&&T` / `Ptr T`) are `Send` (assuming `T` is `Send`). Moving an `OwnedPtr` transfers ownership to the new thread. This is safe because it is only allowed when there are no active borrows.
2.  **Non-`Send` Types**: Types that cannot cross thread boundaries.
    *   Borrow pointers (`&T` / `Ref T`) are **never** `Send` because they refer to memory owned by another thread.

### Lambda Captures and Thread Spawning

When spawning a thread (e.g., via `Thread_::spawn(lambda)`):
*   The function/lambda passed to `spawn` must have the `Send` capability.
*   A lambda is `Send` **if and only if** all of its captured variables are `Send`.
*   Consequently, a lambda cannot capture a `BorrowPtr` (`&T`) if it is to be run on another thread. It can capture `OwnedPtr` (`&&T`), which will be moved into the thread.

### Status

This static checking model is the preferred long-term solution for Bapel concurrency safety. However, it requires a trait-like system (or similar capability tracking) in the Bapel compiler. 

**This feature is currently blocked on Bapel adding support for traits.**

## Detection Examples

### Use-After-Free

```bapel
let p = Ptr_::alloc (10 [i32]);
Ptr_::free p;
let x = Ptr_::get p; // Runtime Error: Use after free detected
```

### Double-Free

```bapel
let p = Ptr_::alloc (10 [i32]);
Ptr_::free p;
Ptr_::free p; // Runtime Error: Double free detected
```

### Use-After-Move

```bapel
let p1 = Ptr_::alloc (10 [i32]);
let p2 = p1; // p1 is moved to p2
let x = Ptr_::get p1; // Runtime Error: Use after move detected
```

### Freeing Borrowed Resource

```bapel
let p1 = Ptr_::alloc (10 [i32]);
let b = Ptr_::borrow p1;
Ptr_::free p1; // Runtime Error: Cannot free a borrowed resource
```

## Overhead Analysis

The overhead of the memory model depends on the build configuration.

### Debug Builds (Checks Enabled)

In debug builds (without `NDEBUG` defined), tracking fields are present and validity checks are executed.

#### Memory Overhead

By splitting the pointer into `OwnedPtr` and `BorrowPtr`, we reduce the memory overhead compared to a unified wrapper:

*   **Owned Pointer (`OwnedPtr<T>`)**:
    *   `ptr_`: 8 bytes.
    *   `state_`: 4 bytes (enum).
    *   `borrows_`: 4 bytes (int).
    *   **Total Size**: 16 bytes (100% overhead vs 8-byte raw pointer).
*   **Borrow Pointer (`BorrowPtr<T>`)**:
    *   `ptr_`: 8 bytes.
    *   `owner_`: 8 bytes (pointer to `OwnedPtr`).
    *   **Total Size**: 16 bytes (100% overhead vs 8-byte raw pointer).

This represents a **100% increase** in memory usage for storing pointers (8 bytes of overhead per pointer).

#### CPU Overhead (Time)

1.  **Dereference (`get_ref()`, `operator->`, `operator*`)**:
    *   **Cost**: Validity checks.
        *   For `OwnedPtr`: Checks if `state_ == VALID`.
        *   For `BorrowPtr`: Checks if `owner_->is_valid()`.
    *   **Impact**: High in tight loops. Every dereference incurs branch checks.
2.  **Move (Transfer Ownership)**:
    *   **Cost**: Checks on `borrows_` (for `OwnedPtr`) and field updates.
    *   **Impact**: Medium.
3.  **Borrowing (`Ptr_::borrow`)**:
    *   **Cost**: Construction of `BorrowPtr`, registering with `OwnedPtr` (incrementing `borrows_`).
    *   **Impact**: Medium.
4.  **Borrow Destruction**:
    *   **Cost**: Decrementing the owner's `borrows_` count.
    *   **Impact**: Low-Medium.
5.  **Allocation/Deallocation (`alloc`/`free`)**:
    *   **Cost**: Dominated by C++ `new`/`delete`. The wrapper overhead is negligible.

### Production Builds (Checks Disabled, `NDEBUG` defined)

In production builds, all tracking fields and checks are conditionally compiled out.

#### Memory Overhead

*   **Owned Pointer (`OwnedPtr<T>`)**:
    *   `ptr_`: 8 bytes.
    *   **Total Size**: 8 bytes (**0% overhead** vs 8-byte raw pointer).
*   **Borrow Pointer (`BorrowPtr<T>`)**:
    *   `ptr_`: 8 bytes.
    *   **Total Size**: 8 bytes (**0% overhead** vs 8-byte raw pointer).

#### CPU Overhead (Time)

All checks and tracking operations are optimized away by the C++ compiler (inlined empty functions or omitted code blocks). The CPU overhead is **0%** compared to using raw pointers.

## Limitations of Runtime Verification

*   **Performance Cost**: Best suited for **Debug/Testing builds**.
*   **Late Detection**: Errors are only detected if the problematic code path is executed.

## Out of Scope

*   **Inheritance and Polymorphism**: Bapel does not support inheritance, so polymorphism and related pointer casting (e.g., derived-to-base) are out of scope for this memory model.
*   **Const Correctness / Mutability Distinction**: Bapel does not distinguish between mutable and immutable borrows at the language level. Therefore, enforcing read/write exclusivity (aliasing XOR mutability) at runtime is out of scope. All borrows are treated as potentially mutable, and multiple active borrows are allowed as long as the resource is not moved or freed.

## Alternative: ECS-based Memory Model (using EnTT)

If Bapel adopts the "ECS is the heap" philosophy, using a library like **EnTT** can further optimize the memory model.

### How it Improves the Proposal

1.  **Zero Memory Overhead for Handles**:
    *   Pointers compile to integer IDs (4 or 8 bytes), matching or beating raw pointer size.
2.  **Built-in Versioning**:
    *   Uses entity versioning for use-after-free detection.

### How it Worsens the Proposal

1.  **Borrow Safety is Harder**:
    *   Does not track raw C++ references to components; still requires wrappers for safe access.
2.  **Paradigm Shift**:
    *   Requires compiler to target ECS registry for all allocations.

### Summary Comparison

| Metric | Custom Split Ptrs (Option B) | EnTT (ECS-based) |
| :--- | :--- | :--- |
| **Pointer Size** | 16 bytes | **4 or 8 bytes** |
| **Memory Overhead** | Medium (+100%) | **Minimal to None** |
| **Use-After-Free** | Runtime checks | **Entity Version Check** |
| **Double-Free** | Runtime checks | **Entity Validation Check** |
| **Borrow Safety** | Easier to track (via `BorrowPtr`) | Harder (requires component wrappers) |
| **Complexity** | Medium (Library + syntax) | High (Compiler integration) |

## Changes to Bapel System Libraries

To support the new `OwnedPtr`/`BorrowPtr` model, the system libraries in the `bapel/` directory must be audited and updated. 

Most library functions that accept pointers do not take ownership of the resource; they only need to read or write to it temporarily. Therefore, their signatures must be updated to use Bapel's borrow syntax (`&T`) instead of the owned pointer syntax (`&&T` or `Ptr T`), and their C++ implementations must be updated to accept `BorrowPtr<T>` instead of raw pointers (`T*`).

### Required Updates by File

1.  **`bapel/core_pointer.h`**:
    *   Define `OwnedPtr` and `BorrowPtr` (as shown in the Runtime Library Design).
    *   Map Bapel `&&T` (or `Ptr T`) to `OwnedPtr<T>`.
    *   Map Bapel `&T` (or `Ref T`) to `BorrowPtr<T>`.
    *   Add `Ref_::get` and `Ref_::set` to allow accessing and modifying borrowed resources.

2.  **`bapel/stl_vector.h`**:
    *   Update signatures in comments to use borrow syntax:
        *   `Vector_::push_back: forall ['a] (&Vector 'a, 'a) -> ()` (already uses `&` in comment, but must map to `BorrowPtr` in C++).
        *   `Vector_::size: forall ['a] &Vector 'a -> i64`
        *   `Vector_::get: forall ['a] (&Vector 'a, i64) -> 'a`
        *   `Vector_::set: forall ['a] (&Vector 'a, i64, 'a) -> ()`
    *   Update C++ implementation to accept `BorrowPtr<std::vector<T>>` instead of `std::vector<T>*`:
        ```cpp
        template <typename T>
        inline std::monostate push_back(BorrowPtr<std::vector<T>> v, T value) {
            v.get_ref().push_back(std::move(value));
            return std::monostate();
        }
        // Apply similar changes to size, get_ref, and set.
        ```

3.  **`bapel/stl_deque.h`**:
    *   Apply similar changes as `stl_vector.h` to use `BorrowPtr<std::deque<T>>`.

4.  **`bapel/stl_fstream.h`**:
    *   Update C++ functions (`is_open`, `close`, `write`) to accept `BorrowPtr<std::ofstream>` instead of `std::ofstream*`.

5.  **`bapel/stl_sstream.h`**:
    *   Update `IStringStream_::read` to accept `BorrowPtr<std::istringstream>` and `BorrowPtr<T>`.

6.  **`bapel/stl_string.h`**:
    *   Update `getline` signature in comments to use borrows:
        `// @bpl: pub getline: forall ['s] (& 's, & String) -> bool`
    *   Update C++ `getline` implementation:
        ```cpp
        template <typename Stream>
        inline bool getline(BorrowPtr<Stream> is, BorrowPtr<String> s) {
            return static_cast<bool>(std::getline(is.get_ref(), s.get_ref()));
        }
        ```

### Migration Strategy

1.  **Phase 1: Compiler Rebuild**: Rebuild the compiler with the new grammar supporting `&&` and `&` mapping to `OwnedPtr` and `BorrowPtr`.
2.  **Phase 2: Library Update**: Update the header files in `bapel/` to match the new design.
3.  **Phase 3: Program Update**: Update user programs (like `program_vector.bpl`) to use the new syntax.


