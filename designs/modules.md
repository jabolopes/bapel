# Modules

A module is a collection of implementation files that share the same
(unexported) symbols.

A module can be a base module (file) or an implementation module (file).

An implementation module must begin with "implements MODULE_ID". This
establishes that this implementation module belongs to the base module with the
identifier MODULE_ID.

A module manages symbols. These symbols can be declared, defined, declared &
defined, imported and exported.

# All declared terms must be defined

All terms declared by a module must be defined within that module. In other
words, given a base module and one or more implementation modules that belong to
that base module, all declared terms must also be defined.

For example, if module A has:

    x: () -> ()

then there must be a definition of x either in A or in an implementation file of
A. For example, A_impl might have:

    fn x() -> () { ... }

# Imports

A module can import from other modules by importing their module identifier. For
example, if module `main` has:

    imports {
      core
    }

this imports all exported symbols from module `core` into `main`. This means all
exported types and terms from `core` are declared in `main`, and a build
dependency between `main` and `core` is established.
                                                                  
For example, if `core` defines `type T` and term `x`, then this `imports`
section translates to:
  
    import core type T
    import core x: () -> ()

## Fundamental imports

It's also possible to import fundamental symbols, such as, bool, i8, etc, using
this mechanism.

    import type bool

This states that bool is a type that exists and it is imported. This does not
establish an additional build dependency since there's no additional module
being imported.

# Impls

A module must have a base module and it can contain 0 or more implementation
files. The connection between a base module and the implementation modules is
established via the `impls` section in the base module, and the `implements`
clause in the implementation modules.

The `impls` section and the `implements` clauses must be consistent.

For example, if base module A has:

    impls {
      A_impl
    }

this is equivalent to importing all symbols (unexported and exported) from
module A_impl and also making module `A_impl` an implementation module file of
module A.

For example, if `A_impl` defines `type T` and term `x` with type () -> (), the
`impls` section translates to:

    impl A_impl.bpl type T
    impl A_impl.bpl x: () -> ()

# Exports

A module can export symbols to other modules.

## Type exports

A module can export a type (declaration), e.g.:

    export type T

This means T is a type and it is exported.


A module can also export a type definition, e.g.:

    export type T = (i8, i8)

This means type T is a tuple of i8s and it is exported.


The 2 options above are different, because if a module imports the type
declaration, all they know is that type T exists but they don't know the
internal implementation details of that type.

If the module imports the type definition, then they know the internal
implementation details.

The exporting module decides on whether the exported type is abstracted or
whether it is fully exported.

The exporting module can also export the abstracted type and keep the internal
implementation details for itself, e.g.:

    export type T

    type T = (i8, i8)

This declares type T as exported (abstracted) and defines type T as a tuple i8s,
without exporting the type definition.

A type declaration and type definition must always be consistent.

## Term exports

A module can export terms. This follows type exports.

### Export term declaration

A module can export a term (declaration), e.g.:

    export x: () -> ()

This means `x` is a term, it has type () -> (), and it is exported.

### Export term definition

A module can also export a term definition, e.g.:

    export fn x() -> () { ... }

This means `x` is a term, it has type () -> (), it has a function body defined
by the block { ... }, and it is exported.

### Hide term definition

A module can also export a term declaration and keep the term definition hidden
from the importing module, e.g.:

    export x: () -> ()

    fn x() -> () { ... }

This prevents the importing module from knowing the internal implementation
details of x, and therefore it cannot depend on x's function body.

A term declaration and term definition must always be consistent.

# Declarations

A module can declare symbols. For example,

    type T

declares that T is a type.


For example,

    x: () -> ()

declares that x is a term with type () -> ().


A module cannot declare a symbol that is already declared. Symbols that are
imported from other modules or from other module implementation files are also
already declared, and therefore cannot be further declared.

## Implicit declarations

A definition, type definition or term definition, is always implicitly declared.

For example, if a module has:

    type T = (i8, i8)
    
    fn x() -> () { ... }

The `type T` and term `x` will be implicitly declared, as if module A had been written as:

    type T
    
    x: () -> ()
    
    type T = (i8, i8)
    
    fn x() -> () { ... }

Implicit declarations are automatically added when a module defines a symbol
(type or a term) and does not additionally declare them explicitly.

This is done on a per-symbol basis, i.e., if module defines and declares `type T`
but only defines term `x`, then only term `x` will be implicitly declared.

Because implicit declarations are only added if the module does not have them,
this does not interfere with the requirement that symbols can be declared at
most once.

Implicit declarations are topologically sorted based on dependencies between
types and terms.

Implicit declarations are always moved above the first definition, i.e.,
implicit declarations. On the other hand, explicit declarations act in the order
in which they appear in the module.

An `import` or an `impl` clause are considered explicit declarations.

# Design

## New sources

New Source constructors:

```
ImportSource
  ModuleID string  // e.g., 'main'
  Decl IrDecl

ImplSource
  ModuleFilename string  // e.g., 'main_impl.bpl' or 'main_impl.cc'
  Decl IrDecl

ExportSource
  Decl IrDecl

DeclSource
  Decl IrDecl

```

Functions continue to be defined by the `FunctionSource`.

`DefSymbolSource` is to be replaced by `DeclSource`.

## Resolver

The `Resolver` is resolves imports and impls, by querying the modules referenced
in the `imports` section, and the implementation files referenced in the `impls`
section.

The imported and implemented symbols become new `ImportSource` and `ImplSource`
at the top of the module's body, with `ImportSource` before `ImplSource`.

The `Resolver` adds implicit declarations (if any). For any symbols defined in
the module's body (types or terms), an implicit declaration is produced if the
symbol is already not explicitly declared.

The implicit declarations become new `DeclSource` inserted in the module's body
right after the `ImplSource`.

## Symbol

The resolver builds a symbol table.

New `Symbol` type:

```
Origin =
| Import ModuleID
| Impl ModuleFilename
| Implicit IrDecl        // Symbol is defined and implicitly declared.
| ExplicitUndefined      // Symbol is explicitly declared and undefined.
| ExplicitDefined IrDecl // Symbol is explicitly declared an defined.

Symbol
  Decl IrDecl
  Origin Origin

SymbolTable = []Symbol
```

Any `DeclSource` for terms need to be matched with a function definition to
ensure there are no undefined terms.
