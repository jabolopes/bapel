# Module discovery

Document how modules are organized and are discovered by the bapel toolchain.

## Workspace

The bapel toolchain must always be executed in the context of a workspace and in
particular, it must be executed in the root of that workspace.

The workspace is the filesystem tree that starts at the directory that contains
the file `workspace.bpl`.

This file is optional and can be absent. If this file is absent, the workspace
root is assumed to be the working directory of the toolchain.

The bapel toolchain DOES NOT perform workspace discovery, i.e., it does not
traverse the filesystem upwards to find the `workspace.bpl` file.

That means it is the developer's responsibility to ensure that the toolchain is
executed at the root of the workspace.

## Module names and filenames

Modules are identified by their name, as proposed in [P1634R0 - Naming
guidelines for modules](https://isocpp.org/files/papers/P1634R0.html).

Module names have a relationship with the filesystem, namely, a module `main`
refers by default to the module with the base file `$WORKSPACE/main.bpl`, and a
module `bapel.core` refers by default to module with the base file
`$WORKSPACE/bapel/core.bpl`.

Because of this relationship, it's not possible to follow the recommendation
from that naming proposal, and only allow module names to have 2 levels, as that
would result in too large directories for large organizations.

For this reason, module names with 1 or more levels are allowed, e.g., `main` (1
level), `bapel.core` (2 levels), `bapel.internal.util` (3 levels), etc.

## Module names and namespaces

There are no relationship requirements between module names and namespaces,
i.e., the modules can use any namespace.

However, if the module name has more than 1 level, it's recommended that the
namespace be derived from the topmost component of that module name.

For example, if the module name is `mycompany.myproject`, then it's recommended
that the namespace contain at least `mycompany`, e.g., `mycompany.MyType` or
`mycompany.myfunc()`, `mycompany.myproject.MyType`, etc. This recommendation is
not enforced.

## Resolve module base files from imports

Given that module names map to the filesystem, imports are resolved relative to
the workspace. For example, `import mymodule` resolved by default to
`$WORKSPACE/mymodule.bpl`, and `import myproject.mymodule` resolved by default
to `$WORKSPACE/myproject/mymodule.bpl`.

## Resolve module implementation files from imports

Modules are a collection of files. These files are the module base file and any
module implementation files specified in the `impls` section of that module base
file.

The module implementation files are specified by their filename, relative to the
resolved location of the module base file.

For example, if `import myproject.mymodule` resolved to the module base file
`$WORKSPACE/myproject/mymodule.bpl`, and `mydir/mymodule_impl.bpl` is in `impls`
section of that base module file, then it resolves to
`$WORKSPACE/myproject/mydir/mymodule_impl.bpl`.

## Module redirection

With all of the above, it's already possible to run the toolchain and build any
project. However, it's not convenient to require all the source code to be in
the same workspace because of the following reasons.

Some modules belong to standard libraries, e.g., `bapel.core`, that are not part
of the developer's source. It makes more sense to "point" the developer's
workspace to the intended standard library rather than require the developer to
copy the standard library to their project.

Some of those modules, e.g., standard libraries, are shipped with the toolchain
installation, and that means they are files stored together with toolchain, and
have a particular version. Again, it makes sense to "point" to those installed
versions rather than keep a copy in the developer's source.

Some modules are managed by other organizations, so the developer won't have
source code for those in their workspace.

Some modules are available on the internet, so they need to be downloaded during
the build process, and as such they are not available in the source tree, rather
only on the build tree.

Some modules are available in the source tree, but the developer may wish to
point those modules to another local copy for experimental purposes. For
example, a developer might have a published version of a project in some code
hosting website, but locally they would point to their local copy of that
project since it has some local, experimental customizations, that are not
published to that hosting website.

For these reasons, it is necessary to have a module name redirection
mechanism. This redirection mechanism is essential a mapping between modules and
their location.

For example, the mapping can point:
* a specific module to a specific file location or internet location, e.g.,
  module `myproject.mymodule` in `/home/.../myproject/mymodule`.
* a module prefix to a specific file location, e.g., prefix `myproject` in
  `/home/.../myproject`.

This mapping is managed by the `workspace.bpl` file.

## The `workspace.bpl` file.

The `workspace.bpl` contains workspace configuration and module redirection
configuration.

For example, a `workspace.bpl` might look like:

```
workspace {
  packages {
    prefix bapel in "/home/jose/Projects/bapel"
    module bapel.core in "/home/jose/Projects/experimental"
  }
}
```

The `module bapel.core` is a module package, and it maps the module `bapel.core`
to some source. It maps no other module.

There can be 0 or more module packages but the module name they map must be
unique, otherwise it's an error.

The rule `prefix bapel` is a prefix package, and it maps to some source all
modules that have `bapel` as their first component. For example, this prefix
package would map the modules `bapel.vector`, `bapel.internal.utils`, etc.

The `prefix` package matches at module name component boundaries and does not
match partial components. This means that `prefix bapel` does not map the
modules `bape`, or `bapel`, or `bapeli`, etc.

The order of the packages does not matter. The rules are always applied from
more specific to less specific, i.e.:
1. If a `module` package matches exactly, that rule is used.
2. Otherwise, if one or more `prefix` package applies, then the longest matching
   prefix is used.
3. Otherwise, the module name is relative to the workspace.

For example, using the `workspace.bpl` example from above:
* module `bapel.core` resolves to `/home/jose/Projects/bapel/bapel/core.bpl`
* module `bapel.vector` resolves to `/home/jose/Projects/experimental/bapel/vector.bpl`

It's important to observe that `bapel.core` resolves to `bapel/core.bpl` inside
`/home/jose/Projects/bapel`.

It's not the case that `bapel.core` resolves to `core.bpl` inside
`/home/jose/Projects/bapel`.

In the future, the `workspace.bpl` file may contain more workspaces, to allow
developers to switch between workspaces easily, and it may also contain other
configuration data.

In the future, the `workspace.bpl` file may also be allowed to contain custom
data that developers can use to store arbitrary workspace related data, e.g.,
build configuration data.
