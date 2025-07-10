# Modules

A module is the unit of abstraction. It is a collection of source
files that share the management of a set of symbols.
                                     
A module decides which symbols exist in that module, and which are visible or
hidden to the other modules.

## Related designs

* [Module discovery](module-discovery.md)

## Source files

A module is a collection of source files.

A source file can be a base file or an implementation file.

A module must have exactly 1 base file, an it can have 0 or more implementation
files.

## Module identifiers

A module identifier (`$MODULE_ID`) is a collection of identifier components
separated by `.`, e.g., `mymodule`, `main`, `bapel.core`,
`bapel.internal.utils`, etc.

A source file filename (`$MODULE_FILENAME`) is, e.g., `"myfile.bpl"`,
`"myfile.cc"`, `"main.bpl"`, `"main_impl.bpl"`, `"bapel/core.bpl"`,
`"bapel/core_impl.bpl"`, etc.

## `module` clause

A base file must begin with the `module` clause on the first line.

The `module` clause is `module $MODULE_ID`, e.g., `module bapel.core`.

The `module` clause is what identifies a file as being a base file, and the
resulting module as being a actual module.

## `implements` clause

An implementation file must being with the `implements` clause on the first
line.

The `implements` clause is `implements $MODULE_ID`, e.g., `implements bapel.core`.

The `implements` clause identifies that a file is an implementation file, and it
must be consistent with the `impls` section (see below).

## `imports` section

A source file (base and implementation) can have at most one `imports` section
with 1 or more module identifiers (`$MODULE_ID`), e.g.:

    imports {
      core
      vec
    }

The module identifiers in the `imports` section must be distinct and must be
lexicographically sorted. Otherwise, it is an error.

The module names given in the `imports` section are "imported" into the current
module. What "import" means is that the symbols exported by the imported modules
are made visible to the current module, to be used. That's all it means.

Importing some other module does not mean to copy it's module text into the
current module.

## `impls` section

A base file can have at most one `impls` section with 1 or more module filenames
(`$MODULE_FILENAME`), e.g.:

    impls {
      "myfile1.bpl"
      "myfile2.ccm"
    }

An implementation file cannot have an `impls` section.

The `impls` section and the `implements` clauses must be consistent, i.e., the
base file must declare all its implementation files and all the implementation
files must declare the same `$MODULE_ID`. Otherwise, it is an error.

TODO: Does the `impls` section need to be sorted? See
https://github.com/jabolopes/bapel/issues/7.

## Symbol management

A module manages symbols. These symbols:
* can be imported from other modules.
* can be imported from implementation files.
* can be locally declared by the current source file
* can be locally defined by the current source file.
* can be locally declared & defined by the current source file.

Symbols can also be exported.

## Symbol frequency

A symbol can be imported from other modules at most once.

A symbol can be imported from implementation files at most once.

A symbol can be locally declared at most once.

A symbol can be locally defined at most once.

If a symbol is both locally declared & defined, it can be locally declared &
defined at most once.

## Symbol shadowing

Symbols imported from other modules can be shadowed by symbols imported from
implementation files, and local symbols. The reason for this is to avoid the
following problem: let's say our module M imports some module C, and module C
happened to add a new function that module M already defined. Without shadowing,
module M would no longer compile. In other words, any module should be allowed
to add new symbols without breaking other modules that depend on it.

## All declared terms must be defined

All terms declared by a module must be defined within that module. In other
words, given a base file and 0 or more implementation files that belong to the
same module, all declared terms must also be defined in that module.

For example, if module `A` has:

    x: () -> ()

then there must be a definition of `x` in `A`'s base file or one of its
implementation files. A definition could be for example:

    fn x() -> () { ... }

## Type abstraction

A module can export a type (declaration), e.g.:

    export type T

A module can also export a type (declaration) and have a hidden type definition,
e.g.:

    export type T

    type T = (i8, i8)

In both cases, the type `T` is exported but its internal representation is
hidden from any modules that import it, but still visible to implementation
files in the same module where it is declared / defined.

A module can also export a type definition, e.g.:

    export type T = (i8, i8)

In this case, type `T` is a tuple of `i8`s and it is exported. Type `T`'s
internal representation is made visible to any modules that import it.

The options above differ because the first 2 options hide the type's internal
representation, whereas the last option makes the internal representation
visible. This is how to employ abstraction.

A type declaration and type definition must always be consistent.

A type declaration must come before its definition. It is an error to declare a
type that is already defined.

## Term abstraction

A module can export a term and its type, e.g.:

    export f: () -> ()
    
    fn f() -> () { () }

A module can export a term and its definition, e.g.:

    export fn() -> () { () }

Given that all declared terms must also be defined, there's no 3rd option like
there was for types.

A term declaration and its term definition must always be consistent.

A term declaration must come before its definition. It is an error to declare a
term that is already defined.

# Imports

A source file can import from other modules by importing their module identifier.

For example, module `core` has:

    export type T1
    export f1: () -> ()
    export type T2 = (i8, i8)
    export fn f2() -> () { () }

The type `T1` is exported but its internal representation is hidden.

The term `f1` is exported with its type, but its internal representation is hidden.

The type `T2` is exported as well as its internal representation.

The term `f2` is exported with its type, as well as its internal representation.

For example, base file `main.bpl` has:

    imports {
      core
    }

This is translated to:
  
    import core type T1
    import core f1: () -> ()
    import core type T2 = (i8, i8)
    import core f2: () -> ()

As a result, base file `main.bpl` can:
* use type `T1` in an abstracted way, without referring to its internal
  representation.
* use term `f1` and its type in an abstracted way, without referring to its
  internal representation (function body).
* use type `T2` including its internal representation, e.g., construct and
  destruct the tuple type.
* use term `f2` and its type, including its internal representation (function
  type), e.g., inline its function body.

# Impls

Implementation files are primarily a way of splitting a module across several
files to avoid having very large files. For this reason, implementation files
must feel the same as if the symbols they declare and define were defined in the
base file.

For example, implementation file `main_impl.bpl` has:

    type T1
    f1: () -> ()
    type T2 = (i8, i8)
    fn f2() -> () { () }
    
    type T1 = (i8, i8)
    fn f1() -> () { () }

And, base file `main.bpl` has:

    impls {
      main_impl.bpl
    }

This is translated to:

    impl main_impl.bpl type T1 = (i8, i8)
    impl main_impl.bpl f1: () -> () { () }
    impl main_impl.bpl type T2 = (i8, i8)
    impl main_impl.bpl fn f2() -> () { () }
    
As a result, base file `main.bpl` can:
* use type `T1` and its internal representation, e.g., construct and destruct
  the tuple type.
* use term `f1` and its type, including its internal representation (function
  body).
* use type `T2` and its internal representation, e.g., construct and destruct
  the tuple type.
* use term `f2` and its type, including its internal representation (function
  type).

The result would be the same if `main_impl.bpl` exported any of those
symbols. In other words, a base file sees all the exported and unexported
symbols from the implementation files in the same module.

## Implicit declarations

All locally defined types and terms are also declared, either implicitly or
explicitly.

An explicit declaration is when the source program contains a type declaration
(e.g., `type T`) or a term declaration (e.g., `f: () -> ()`).

An implicit declaration is a declaration that is automatically inserted when a
type or a term are defined but not explicitly declared.

For example, `main.bpl` has:

    type T1
    type T1 = (i8, i8)
    
    type T2 = (i8, i8)

    f1: () -> ()
    fn f1() -> ()

    fn f2() -> () { () }

This is translated to:

    type T1
    type T2
    f1: () -> ()
    f2: () -> ()

    type T1 = (i8, i8)
   
    fn f1() -> ()

    type T2 = (i8, i8)

    fn f2() -> () { () }

The result is:
* type `T1` was already explicitly declared, so it has no implicit declaration.
* type `T2` was not explicitly declared, so it has an implicit declaration.
* term `f1` was already explicitly declared, so it has no implicit declaration.
* term `f2` was not explicitly declared, so it has an implicit declaration.

Implicit declarations must obey the requirements of symbols being declared /
defined at most once. In fact, implicit declarations do not interfere with that
requirement.

The main purpose of implicit declarations is to solve the following problem: the
order of declarations and definitions is important since symbols are only
entered into the context once they are declared or defined. Implicit
declarations allow the types and terms to refer to other types and terms that
are defined in non-sequential order without requiring the programmer to
explicitly declare them.

All declarations (implicit and explicit) are automatically topologically sorted
based on dependencies between types and terms.

Imported symbols either from other modules or implementation files are
considered explicit declarations.
