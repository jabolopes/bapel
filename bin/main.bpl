module main

imports {
  bapel.core
  bapel.stl
}

impls {
  "main_helper.h"
  "main_helper.cc"
}

type BazelTarget = struct {
  kind String,
  name String,
  srcs core::Vector String,
  hdrs core::Vector String,
  deps core::Vector String
}


type cli::MatchResult = struct {
  found bool,
  path String,
  prefixLength i64
}

fn copyFile(src: String, dst: String) -> bool {
  let dstDir: String = fs::parent_path dst;
  if !fs::create_directories dstDir {
     return false
  }
  fs::remove dst;
  fs::copy (src, dst)
}

fn execInDir(cmd: String, args: core::Vector String, dir: String) -> (i64, String) {
  let origPath: String = fs::current_path ();
  if !fs::set_current_path dir {
     let res: (i64, String) = (-1, "chdir failed");
     return res
  }
  let res: (i64, String) = cli::exec (cmd, args);
  if !fs::set_current_path origPath {
     core::print [String] "Warning: failed to restore CWD";
     ()
  }
  res
}

fn resolveMappedPath(path: String, moduleID: String) -> String {
  let relPath: String = cli::replaceSeparator (moduleID, ".", "/");
  let relPathWithExt: String = String_::concat (relPath, ".bpl");
  fs::join (path, relPathWithExt)
}

fn findBestMatch(
    mappings: & (core::Vector cli::PackageMapping),
    moduleID: String,
    index: i64,
    currentBest: cli::MatchResult) -> cli::MatchResult {
  
  if index >= core::vec_size mappings {
     return currentBest
  }
  
  let mapping: cli::PackageMapping = core::get (mappings, index);
  
  if mapping.is_prefix {
     if cli::isPrefixOf (mapping.name, moduleID) {
        let prefixLen: i64 = String_::size mapping.name;
        if prefixLen > currentBest.prefixLength {
           let resolvedPath: String = resolveMappedPath (mapping.path, moduleID);
           let newBest: cli::MatchResult = struct { found = true, path = resolvedPath, prefixLength = prefixLen };
           return findBestMatch (mappings, moduleID, index + 1, newBest)
        }
     }
  } else {
     if mapping.name == moduleID {
        let exactPath: String = resolveMappedPath (mapping.path, moduleID);
        let newBest: cli::MatchResult = struct { found = true, path = exactPath, prefixLength = 999999 };
        return findBestMatch (mappings, moduleID, index + 1, newBest)
     }
  }
  
  findBestMatch (mappings, moduleID, index + 1, currentBest)
}

fn resolveModule(moduleID: String) -> String {
  let wsFile: String = "workspace.bpl";
  if fs::exists wsFile {
     let args: core::Vector String = core::mk [String] ();
     core::add [String] (&args, "-workspace");
     core::add [String] (&args, "-format=flat");
     core::add [String] (&args, wsFile);
     let res: (i64, String) = cli::exec ("bootstrap/parser", args);
     if res.0 == 0 {
        let mappings: core::Vector cli::PackageMapping = cli::parseWorkspaceFlat res.1;
        let emptyBest: cli::MatchResult = struct { found = false, path = "", prefixLength = 0 };
        let best: cli::MatchResult = findBestMatch (&mappings, moduleID, 0, emptyBest);
        if best.found {
           return best.path
        }
     }
  }
  let relPath: String = cli::replaceSeparator (moduleID, ".", "/");
  String_::concat (relPath, ".bpl")
}

fn vecContains(v: & (core::Vector String), s: String, index: i64) -> bool {
  if index >= core::vec_size v {
     return false
  }
  if core::get (v, index) == s {
     return true
  }
  vecContains (v, s, index + 1)
}

fn appendVectors(dst: & (core::Vector String), src: & (core::Vector String), index: i64) -> () {
  if index >= core::vec_size src {
     return ()
  }
  core::add [String] (dst, core::get (src, index));
  appendVectors (dst, src, index + 1)
}

fn buildImpls(
    implFiles: & (core::Vector String),
    moduleID: String,
    baseFileDir: String,
    index: i64,
    srcs: & (core::Vector String),
    hdrs: & (core::Vector String)) -> i64 {
  
  if index >= core::vec_size implFiles {
     return 0
  }
  let implFile: String = core::get (implFiles, index);
  let fullImplPath: String = fs::join (baseFileDir, implFile);
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let baseName: String = fs::stem implFile;
     let baseOutputBasename: String = cli::replaceSeparator (moduleID, ".", "/");
     let implOutBasename: String = String_::concat (String_::concat (baseOutputBasename, "-"), baseName);
     let outPath: String = fs::join ("out", implOutBasename);
     let outCcPath: String = String_::concat (outPath, ".cc");
     
     if !fs::create_directories (fs::parent_path outCcPath) {
        core::print [String] (String_::concat ("Failed to create directory: ", fs::parent_path outCcPath));
        return 1
     }
     
     let ccArgs: core::Vector String = core::mk [String] ();
     core::add [String] (&ccArgs, "-o");
     core::add [String] (&ccArgs, outCcPath);
     core::add [String] (&ccArgs, fullImplPath);
     
     let ccRes: (i64, String) = cli::exec ("bootstrap/compiler", ccArgs);
     if ccRes.0 != 0 {
        core::print [String] (String_::concat ("Failed to compile impl: ", fullImplPath));
        core::print [String] ccRes.1;
        return ccRes.0
     }
     
     core::add [String] (srcs, String_::concat(implOutBasename, ".cc"));
     
  } else {
     let dst: String = fs::join ("out", fullImplPath);
     if !copyFile (fullImplPath, dst) {
        core::print [String] (String_::concat ("Failed to copy impl: ", fullImplPath));
        return 1
     }
     
     if ext == ".cc" {
        core::add [String] (srcs, fullImplPath);
        ()
     } else if ext == ".cpp" {
        core::add [String] (srcs, fullImplPath);
        ()
     } else if ext == ".h" {
        core::add [String] (hdrs, fullImplPath);
        ()
     }
  }
  
  buildImpls (implFiles, moduleID, baseFileDir, index + 1, srcs, hdrs)
}

fn mergeUnique(src: & (core::Vector String), dst: & (core::Vector String), index: i64) -> () {
  if index >= core::vec_size src {
     return ()
  }
  let item: String = core::get (src, index);
  if !vecContains (dst, item, 0) {
     core::add [String] (dst, item);
     ()
  }
  mergeUnique (src, dst, index + 1)
}

fn collectImplImports(
    implFiles: & (core::Vector String),
    baseFileDir: String,
    index: i64,
    importsList: & (core::Vector String)) -> i64 {
  
  if index >= core::vec_size implFiles {
     return 0
  }
  let implFile: String = core::get (implFiles, index);
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let fullImplPath: String = fs::join (baseFileDir, implFile);
     let args: core::Vector String = core::mk [String] ();
     core::add [String] (&args, "-format=flat");
     core::add [String] (&args, fullImplPath);
     let res: (i64, String) = cli::exec ("bootstrap/parser", args);
     if res.0 != 0 {
        core::print [String] (String_::concat ("Failed to parse impl for imports: ", fullImplPath));
        return res.0
     }
     let info: cli::SourceFileInfo = cli::parseSourceFileFlat res.1;
     let implImports: core::Vector String = info.importModules;
     mergeUnique (&implImports, importsList, 0);
     ()
  }
  
  collectImplImports (implFiles, baseFileDir, index + 1, importsList)
}

fn buildImports(
    importModules: & (core::Vector String),
    builtModules: & (core::Vector String),
    deps: & (core::Vector String),
    index: i64,
    targets: & (core::Vector BazelTarget)) -> i64 {
  
  if index >= core::vec_size importModules {
     return 0
  }
  let imp: String = core::get (importModules, index);
  let err: i64 = buildModule (imp, builtModules, false, targets);
  if err != 0 {
     return err
  }
  
  let sanitized: String = cli::replaceSeparator (cli::replaceSeparator (imp, ".", "_"), "/", "_");
  let depTarget: String = String_::concat (":", sanitized);
  core::add [String] (deps, depTarget);
  
  buildImports (importModules, builtModules, deps, index + 1, targets)
}

fn buildModule(
    moduleID: String,
    builtModules: & (core::Vector String),
    isRoot: bool,
    targets: & (core::Vector BazelTarget)) -> i64 {
  if vecContains (builtModules, moduleID, 0) {
     return 0
  }
  
  let baseFile: String = resolveModule moduleID;
  if !fs::exists baseFile {
     core::print [String] (String_::concat ("File not found: ", baseFile));
     return 1
  }
  
  let args: core::Vector String = core::mk [String] ();
  core::add [String] (&args, "-format=flat");
  core::add [String] (&args, baseFile);
  let res: (i64, String) = cli::exec ("bootstrap/parser", args);
  if res.0 != 0 {
     core::print [String] (String_::concat ("Failed to parse: ", baseFile));
     core::print [String] res.1;
     return res.0
  }
  
  let info: cli::SourceFileInfo = cli::parseSourceFileFlat res.1;
  
  let importsList: core::Vector String = info.importModules;
  let baseFileDir: String = fs::parent_path baseFile;
  let implsList: core::Vector String = info.implFiles;
  let err_impls: i64 = collectImplImports (&implsList, baseFileDir, 0, &importsList);
  if err_impls != 0 {
     return err_impls
  }
  
  let deps: core::Vector String = core::mk [String] ();
  let err: i64 = buildImports (&importsList, builtModules, &deps, 0, targets);
  if err != 0 {
     return err
  }
  
  let baseOutputBasename: String = cli::replaceSeparator (moduleID, ".", "/");
  let outPath: String = fs::join ("out", baseOutputBasename);
  let outHeader: String = String_::concat (outPath, ".h");
  
  if !fs::create_directories (fs::parent_path outHeader) {
     core::print [String] (String_::concat ("Failed to create directory: ", fs::parent_path outHeader));
     return 1
  }
  
  let ccArgs: core::Vector String = core::mk [String] ();
  core::add [String] (&ccArgs, "-o");
  core::add [String] (&ccArgs, outHeader);
  core::add [String] (&ccArgs, baseFile);
  
  let ccRes: (i64, String) = cli::exec ("bootstrap/compiler", ccArgs);
  if ccRes.0 != 0 {
     core::print [String] (String_::concat ("Failed to compile: ", baseFile));
     core::print [String] ccRes.1;
     return ccRes.0
  }
  
  let srcs: core::Vector String = core::mk [String] ();
  let hdrs: core::Vector String = core::mk [String] ();
  
  core::add [String] (&srcs, String_::concat(baseOutputBasename, ".cc"));
  core::add [String] (&hdrs, String_::concat(baseOutputBasename, ".h"));
  core::add [String] (&hdrs, String_::concat(baseOutputBasename, "_private.h"));
  
  let err2: i64 = buildImpls (&implsList, moduleID, baseFileDir, 0, &srcs, &hdrs);
  if err2 != 0 {
     return err2
  }
  
  let targetName: String = cli::replaceSeparator (cli::replaceSeparator (moduleID, ".", "_"), "/", "_");
  if isRoot {
     appendVectors (&srcs, &hdrs, 0);
     let emptyHdrs: core::Vector String = core::mk [String] ();
     let target: BazelTarget = struct {
        kind = "cc_binary",
        name = targetName,
        srcs = srcs,
        hdrs = emptyHdrs,
        deps = deps
     };
     core::add [BazelTarget] (targets, target);
     ()
  } else {
     let target: BazelTarget = struct {
        kind = "cc_library",
        name = targetName,
        srcs = srcs,
        hdrs = hdrs,
        deps = deps
     };
     core::add [BazelTarget] (targets, target);
     ()
  }
  
  core::add [String] (builtModules, moduleID);
  0
}

fn writeVector(f: & Ofstream, label: String, v: & (core::Vector String)) -> () {
  if core::vec_size v == 0 {
     return ()
  }
  Ofstream_::write (f, label);
  Ofstream_::write (f, " = [\n");
  writeVectorElems (f, v, 0);
  Ofstream_::write (f, "    ],\n");
  ()
}

fn writeVectorElems(f: & Ofstream, v: & (core::Vector String), index: i64) -> () {
  if index >= core::vec_size v {
     return ()
  }
  Ofstream_::write (f, "        \"");
  Ofstream_::write (f, core::get (v, index));
  Ofstream_::write (f, "\",\n");
  writeVectorElems (f, v, index + 1)
}

fn writeTargets(f: & Ofstream, targets: & (core::Vector BazelTarget), index: i64) -> () {
  if index >= core::vec_size targets {
     return ()
  }
  let target: BazelTarget = core::get (targets, index);
  Ofstream_::write (f, target.kind);
  Ofstream_::write (f, "(\n");
  
  Ofstream_::write (f, "    name = \"");
  Ofstream_::write (f, target.name);
  Ofstream_::write (f, "\",\n");
  
  let srcs: core::Vector String = target.srcs;
  writeVector (f, "    srcs", &srcs);
  let hdrs: core::Vector String = target.hdrs;
  writeVector (f, "    hdrs", &hdrs);
  
  Ofstream_::write (f, "    copts = [\n");
  Ofstream_::write (f, "        \"-std=c++17\",\n");
  Ofstream_::write (f, "        \"-Xassembler\",\n");
  Ofstream_::write (f, "        \"--gsframe=no\",\n");
  Ofstream_::write (f, "    ],\n");
  
  let deps: core::Vector String = target.deps;
  writeVector (f, "    deps", &deps);
  
  Ofstream_::write (f, ")\n\n");
  
  writeTargets (f, targets, index + 1)
}

fn writeBuildFile(targets: & (core::Vector BazelTarget)) -> bool {
  let f: Ofstream = Ofstream_::open "out/BUILD";
  if !Ofstream_::is_open &f {
     return false
  }
  Ofstream_::write (&f, "load(\"@rules_cc//cc:defs.bzl\", \"cc_binary\", \"cc_library\")\n\n");
  writeTargets (&f, targets, 0);
  Ofstream_::close &f;
  true
}

fn ensureWorkspaceSetup() -> bool {
  if !fs::create_directories "out" {
     return false
  }
  
  let ws: Ofstream = Ofstream_::open "out/WORKSPACE";
  if !Ofstream_::is_open &ws {
     return false
  }
  Ofstream_::write (&ws, "workspace(name = \"bapel_out\")\n");
  Ofstream_::close &ws;

  let mod: Ofstream = Ofstream_::open "out/MODULE.bazel";
  if !Ofstream_::is_open &mod {
     return false
  }
  Ofstream_::write (&mod, "module(name = \"bapel_out\")\n");
  Ofstream_::write (&mod, "bazel_dep(name = \"rules_cc\", version = \"0.2.17\")\n");
  Ofstream_::close &mod;

  true
}

fn build(moduleID: String) -> i64 {
  if !ensureWorkspaceSetup () {
     core::print [String] "Failed to setup workspace";
     return 1
  }
  
  let builtModules: core::Vector String = core::mk [String] ();
  let targets: core::Vector BazelTarget = core::mk [BazelTarget] ();
  let err: i64 = buildModule (moduleID, &builtModules, true, &targets);
  if err != 0 {
     return err
  }
  
  if !writeBuildFile (&targets) {
     core::print [String] "Failed to write BUILD file";
     return 1
  }
  
  let targetName: String = cli::replaceSeparator (cli::replaceSeparator (moduleID, ".", "_"), "/", "_");
  
  // Safe construction of "//:" to avoid parser comment bugs
  let slash: String = "/";
  let doubleSlash: String = String_::concat (slash, slash);
  let bazelTarget: String = String_::concat (String_::concat (doubleSlash, ":"), targetName);
  
  let bazelArgs: core::Vector String = core::mk [String] ();
  core::add [String] (&bazelArgs, "build");
  core::add [String] (&bazelArgs, bazelTarget);
  
  let bazelRes: (i64, String) = execInDir ("bazel", bazelArgs, "out");
  if bazelRes.0 != 0 {
     core::print [String] "Bazel build failed";
     core::print [String] bazelRes.1;
     return bazelRes.0
  }
  
  let bazelBinPath: String = fs::join (fs::join ("out", "bazel-bin"), targetName);
  let outputPath: String = fs::join ("out", moduleID);
  
  if !copyFile (bazelBinPath, outputPath) {
     core::print [String] "Failed to copy built binary";
     return 1
  }
  
  0
}

pub fn bapel_main() -> i32 {
  let count: i64 = cli::getArgCount ();
  if count < 2 {
     core::print [String] "expected subcommand, e.g., 'parse', 'cc', 'build', 'query'";
     return 1
  }
  
  let command: String = cli::getArg 1;
  
  if command == "cc" {
     let subArgs: core::Vector String = cli::getSubArgs 2;
     let res: (i64, String) = cli::exec ("bootstrap/compiler", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "build" {
     if count < 3 {
        core::print [String] "usage: bpl build <module>";
        return 1
     }
     let err: i64 = build (cli::getArg 2);
     return core::i64_to_i32 err
  }
  
  if command == "query" {
     let subArgs: core::Vector String = cli::getSubArgs 2;
     let res: (i64, String) = cli::exec ("bootstrap/querier", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "parse" {
     let subArgs: core::Vector String = cli::getSubArgs 2;
     let res: (i64, String) = cli::exec ("bootstrap/parser", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  core::print [String] (String_::concat ("unknown command: ", command));
  1
}
