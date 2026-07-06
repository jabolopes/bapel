module main

imports {
  bapel.args
  bapel.core
  bapel.os
  bapel.stl
}



type BazelTarget = struct {
  kind: String,
  name: String,
  srcs: Vector String,
  hdrs: Vector String,
  deps: Vector String
}


type MatchResult = struct {
  found: bool,
  path: String,
  prefixLength: i64
}

type PackageMapping = struct {
  is_prefix: bool,
  name: String,
  path: String
}

type SourceFileInfo = struct {
  importModules: Vector String,
  implFiles: Vector String
}


fn copyFile(src: &String, dst: &String) -> bool {
  let dstDir: String = fs::parent_path (*dst);
  if !fs::create_directories dstDir {
     return false
  }
  fs::remove (*dst);
  fs::copy (*src, *dst)
}

fn execInDir(cmd: &String, args: &Vector String, dir: &String) -> (i64, String) {
  let origPath: String = fs::current_path ();
  if !fs::set_current_path (*dir) {
     let res: (i64, String) = (-1, "chdir failed");
     return res
  }
  let res: (i64, String) = os::exec (*cmd, *args);
  if !fs::set_current_path origPath {
     core::print [String] "Warning: failed to restore CWD";
     ()
  }
  res
}

fn replaceSeparator(s: String, from: &String, to: &String) -> String {
  let pos: i64 = String::find (&s, from, 0);
  let from_len: i64 = String::size from;
  let to_len: i64 = String::size to;
  
  for pos != -1 {
    s <- String::replace (&s, pos, from_len, to);
    pos <- pos + to_len;
    pos <- String::find (&s, from, pos);
  };
  s
}

fn resolveMappedPath(path: &String, moduleID: &String) -> String {
  let dot: String = ".";
  let slash: String = "/";
  let relPath: String = replaceSeparator (*moduleID, &dot, &slash);
  let relPathWithExt: String = String::concat (&relPath, &".bpl");
  fs::join (*path, relPathWithExt)
}

fn isPrefixOf(pref: &String, s: &String) -> bool {
  if *s == *pref {
     return true
  }
  let p: String = String::concat (pref, &".");
  let p_len: i64 = String::size &p;
  let s_len: i64 = String::size s;
  if s_len < p_len {
     return false
  }
  let s_view: StringView = String::view s;
  let sub_view: StringView = StringView::substr (s_view, 0, p_len);
  let sub: String = String::from_view sub_view;
  sub == p
}

fn findBestMatch(
    mappings: &Vector PackageMapping,
    moduleID: &String,
    index: i64,
    currentBest: &MatchResult) -> MatchResult {
  
  if index >= Vector::size mappings {
     return *currentBest
  }
  
  let mapping: PackageMapping = Vector::get (mappings, index);
  
  if mapping.is_prefix {
     let mapping_name: String = mapping.name;
     if isPrefixOf (&mapping_name, moduleID) {
        let prefixLen: i64 = String::size &mapping_name;
        if prefixLen > (*currentBest).prefixLength {
           let mapping_path: String = mapping.path;
           let resolvedPath: String = resolveMappedPath (&mapping_path, moduleID);
           let newBest: MatchResult = struct { found = true, path = resolvedPath, prefixLength = prefixLen };
           return findBestMatch (mappings, moduleID, index + 1, &newBest)
        }
     }
  } else {
     if mapping.name == *moduleID {
        let mapping_path: String = mapping.path;
        let exactPath: String = resolveMappedPath (&mapping_path, moduleID);
        let newBest: MatchResult = struct { found = true, path = exactPath, prefixLength = 999999 };
        return findBestMatch (mappings, moduleID, index + 1, &newBest)
     }
  }
  
  findBestMatch (mappings, moduleID, index + 1, currentBest)
}

fn processWorkspaceLine(line: &String, mappings: &Vector PackageMapping) -> () {
  if String::size line == 0 {
     return ()
  }
  let line_iss: IStringStream = IStringStream::mk (*line);
  let type_str: String = "";
  let name: String = "";
  let path: String = "";
  
  if !IStringStream::read (&line_iss, &type_str) {
     return ()
  }
  if !IStringStream::read (&line_iss, &name) {
     return ()
  }
  if !IStringStream::read (&line_iss, &path) {
     return ()
  }
  
  let is_prefix: bool = type_str == "PREFIX";
  let mapping: PackageMapping = struct {
     is_prefix = is_prefix,
     name = name,
     path = path
  };
  Vector::push_back [PackageMapping] (mappings, mapping);
  ()
}

fn parseWorkspaceFlat(text: &String) -> Vector PackageMapping {
  let mappings: Vector PackageMapping = Vector::mk [PackageMapping] ();
  let iss: IStringStream = IStringStream::mk (*text);
  let line: String = "";
  
  for getline (&iss, &line) {
     processWorkspaceLine (&line, &mappings);
  };
  mappings
}

fn parseSourceFileFlat(text: &String) -> SourceFileInfo {
  let importModules: Vector String = Vector::mk [String] ();
  let implFiles: Vector String = Vector::mk [String] ();
  let iss: IStringStream = IStringStream::mk (*text);
  let line: String = "";
  
  for getline (&iss, &line) {
    if String::size &line > 0 {
      let line_iss: IStringStream = IStringStream::mk line;
      let type_str: String = "";
      let value: String = "";
      if IStringStream::read (&line_iss, &type_str) {
        if IStringStream::read (&line_iss, &value) {
          if type_str == "IMPORT" {
            Vector::push_back [String] (&importModules, value);
            ()
          } else if type_str == "IMPL" {
            Vector::push_back [String] (&implFiles, value);
            ()
          }
        }
      }
    }
  };
  
  let info: SourceFileInfo = struct {
    importModules = importModules,
    implFiles = implFiles
  };
  info
}

fn resolveModule(moduleID: &String) -> String {
  let wsFile: String = "workspace.bpl";
  if fs::exists wsFile {
     let args: Vector String = Vector::mk [String] ();
     Vector::push_back [String] (&args, "-workspace");
     Vector::push_back [String] (&args, "-format=flat");
     Vector::push_back [String] (&args, wsFile);
     let res: (i64, String) = os::exec ("bootstrap/parser", args);
     if res.0 == 0 {
        let flatText: String = res.1;
        let mappings: Vector PackageMapping = parseWorkspaceFlat (&flatText);
        let emptyBest: MatchResult = struct { found = false, path = "", prefixLength = 0 };
        let best: MatchResult = findBestMatch (&mappings, moduleID, 0, &emptyBest);
        if best.found {
           return best.path
        }
     }
  }
  let dot: String = ".";
  let slash: String = "/";
  let relPath: String = replaceSeparator (*moduleID, &dot, &slash);
  String::concat (&relPath, &".bpl")
}

fn vecContains(v: &Vector String, s: &String, index: i64) -> bool {
  if index >= Vector::size v {
     return false
  }
  if Vector::get (v, index) == *s {
     return true
  }
  vecContains (v, s, index + 1)
}

fn appendVectors(dst: &Vector String, src: &Vector String, index: i64) -> () {
  if index >= Vector::size src {
     return ()
  }
  Vector::push_back [String] (dst, Vector::get (src, index));
  appendVectors (dst, src, index + 1)
}

fn buildImpls(
    implFiles: &Vector String,
    moduleID: &String,
    baseFileDir: &String,
    index: i64,
    srcs: &Vector String,
    hdrs: &Vector String) -> i64 {
  
  if index >= Vector::size implFiles {
     return 0
  }
  let implFile: String = Vector::get (implFiles, index);
  let fullImplPath: String = fs::join (*baseFileDir, implFile);
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let baseName: String = fs::stem implFile;
     let dot: String = ".";
     let slash: String = "/";
     let baseOutputBasename: String = replaceSeparator (*moduleID, &dot, &slash);
     let implOutBasename: String = String::concat (&(String::concat (&baseOutputBasename, &"-")), &baseName);
     let outPath: String = fs::join ("out", implOutBasename);
     let outCcPath: String = String::concat (&outPath, &".cc");
     
     if !fs::create_directories (fs::parent_path outCcPath) {
        core::print [String] (String::concat (&"Failed to create directory: ", &(fs::parent_path outCcPath)));
        return 1
     }
     
     let ccArgs: Vector String = Vector::mk [String] ();
     Vector::push_back [String] (&ccArgs, "-o");
     Vector::push_back [String] (&ccArgs, outCcPath);
     Vector::push_back [String] (&ccArgs, fullImplPath);
     
     let ccRes: (i64, String) = os::exec ("bootstrap/compiler", ccArgs);
     if ccRes.0 != 0 {
        core::print [String] (String::concat (&"Failed to compile impl: ", &fullImplPath));
        core::print [String] ccRes.1;
        return ccRes.0
     }
     
     Vector::push_back [String] (srcs, String::concat (&implOutBasename, &".cc"));
     
  } else {
     let dst: String = fs::join ("out", fullImplPath);
     if !copyFile (&fullImplPath, &dst) {
        core::print [String] (String::concat (&"Failed to copy impl: ", &fullImplPath));
        return 1
     }
     
     if ext == ".cc" {
        Vector::push_back [String] (srcs, fullImplPath);
        ()
     } else if ext == ".cpp" {
        Vector::push_back [String] (srcs, fullImplPath);
        ()
     } else if ext == ".h" {
        Vector::push_back [String] (hdrs, fullImplPath);
        ()
     }
  }
  
  buildImpls (implFiles, moduleID, baseFileDir, index + 1, srcs, hdrs)
}

fn mergeUnique(src: &Vector String, dst: &Vector String, index: i64) -> () {
  if index >= Vector::size src {
     return ()
  }
  let item: String = Vector::get (src, index);
  if !vecContains (dst, &item, 0) {
     Vector::push_back [String] (dst, item);
     ()
  }
  mergeUnique (src, dst, index + 1)
}

fn collectImplImports(
    implFiles: &Vector String,
    baseFileDir: &String,
    index: i64,
    importsList: &Vector String) -> i64 {
  
  if index >= Vector::size implFiles {
     return 0
  }
  let implFile: String = Vector::get (implFiles, index);
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let fullImplPath: String = fs::join (*baseFileDir, implFile);
     let args: Vector String = Vector::mk [String] ();
     Vector::push_back [String] (&args, "-format=flat");
     Vector::push_back [String] (&args, fullImplPath);
     let res: (i64, String) = os::exec ("bootstrap/parser", args);
     if res.0 != 0 {
        core::print [String] (String::concat (&"Failed to parse impl for imports: ", &fullImplPath));
        return res.0
     }
     let flatText: String = res.1;
     let info: SourceFileInfo = parseSourceFileFlat (&flatText);
     let implImports: Vector String = info.importModules;
     mergeUnique (&implImports, importsList, 0);
     ()
  }
  
  collectImplImports (implFiles, baseFileDir, index + 1, importsList)
}

fn buildImports(
    importModules: &Vector String,
    builtModules: &Vector String,
    deps: &Vector String,
    index: i64,
    targets: &Vector BazelTarget) -> i64 {
  
  if index >= Vector::size importModules {
     return 0
  }
  let imp: String = Vector::get (importModules, index);
  let err: i64 = buildModule (&imp, builtModules, false, targets);
  if err != 0 {
     return err
  }
  
  let dot: String = ".";
  let slash: String = "/";
  let under: String = "_";
  let tempName: String = replaceSeparator (imp, &dot, &under);
  let sanitized: String = replaceSeparator (tempName, &slash, &under);
  let depTarget: String = String::concat (&":", &sanitized);
  Vector::push_back [String] (deps, depTarget);
  
  buildImports (importModules, builtModules, deps, index + 1, targets)
}

fn buildModule(
    moduleID: &String,
    builtModules: &Vector String,
    isRoot: bool,
    targets: &Vector BazelTarget) -> i64 {
  if vecContains (builtModules, moduleID, 0) {
     return 0
  }
  
  let baseFile: String = resolveModule moduleID;
  if !fs::exists baseFile {
     core::print [String] (String::concat (&"File not found: ", &baseFile));
     return 1
  }
  
  let args: Vector String = Vector::mk [String] ();
  Vector::push_back [String] (&args, "-format=flat");
  Vector::push_back [String] (&args, baseFile);
  let res: (i64, String) = os::exec ("bootstrap/parser", args);
  if res.0 != 0 {
     core::print [String] (String::concat (&"Failed to parse: ", &baseFile));
     core::print [String] res.1;
     return res.0
  }
  
  let flatText: String = res.1;
  let info: SourceFileInfo = parseSourceFileFlat (&flatText);
  
  let importsList: Vector String = info.importModules;
  let baseFileDir: String = fs::parent_path baseFile;
  let implsList: Vector String = info.implFiles;
  let err_impls: i64 = collectImplImports (&implsList, &baseFileDir, 0, &importsList);
  if err_impls != 0 {
     return err_impls
  }
  
  let deps: Vector String = Vector::mk [String] ();
  let err: i64 = buildImports (&importsList, builtModules, &deps, 0, targets);
  if err != 0 {
     return err
  }
  
  let dot: String = ".";
  let slash: String = "/";
  let under: String = "_";
  let baseOutputBasename: String = replaceSeparator (*moduleID, &dot, &slash);
  let outPath: String = fs::join ("out", baseOutputBasename);
  let outHeader: String = String::concat (&outPath, &".h");
  
  if !fs::create_directories (fs::parent_path outHeader) {
     core::print [String] (String::concat (&"Failed to create directory: ", &(fs::parent_path outHeader)));
     return 1
  }
  
  let ccArgs: Vector String = Vector::mk [String] ();
  Vector::push_back [String] (&ccArgs, "-o");
  Vector::push_back [String] (&ccArgs, outHeader);
  Vector::push_back [String] (&ccArgs, baseFile);
  
  let ccRes: (i64, String) = os::exec ("bootstrap/compiler", ccArgs);
  if ccRes.0 != 0 {
     core::print [String] (String::concat (&"Failed to compile: ", &baseFile));
     core::print [String] ccRes.1;
     return ccRes.0
  }
  
  let srcs: Vector String = Vector::mk [String] ();
  let hdrs: Vector String = Vector::mk [String] ();
  
  Vector::push_back [String] (&srcs, String::concat (&baseOutputBasename, &".cc"));
  Vector::push_back [String] (&hdrs, String::concat (&baseOutputBasename, &".h"));
  Vector::push_back [String] (&hdrs, String::concat (&baseOutputBasename, &"_private.h"));
  
  let err2: i64 = buildImpls (&implsList, moduleID, &baseFileDir, 0, &srcs, &hdrs);
  if err2 != 0 {
     return err2
  }
  
  let tempName: String = replaceSeparator (*moduleID, &dot, &under);
  let targetName: String = replaceSeparator (tempName, &slash, &under);
  if isRoot {
     appendVectors (&srcs, &hdrs, 0);
     let emptyHdrs: Vector String = Vector::mk [String] ();
     let target: BazelTarget = struct {
        kind = "cc_binary",
        name = targetName,
        srcs = srcs,
        hdrs = emptyHdrs,
        deps = deps
     };
     Vector::push_back [BazelTarget] (targets, target);
     ()
  } else {
     let target: BazelTarget = struct {
        kind = "cc_library",
        name = targetName,
        srcs = srcs,
        hdrs = hdrs,
        deps = deps
     };
     Vector::push_back [BazelTarget] (targets, target);
     ()
  }
  
  Vector::push_back [String] (builtModules, *moduleID);
  0
}

fn writeVector(f: & Ofstream, label: &String, v: &Vector String) -> () {
  if Vector::size v == 0 {
     return ()
  }
  Ofstream::write (f, *label);
  Ofstream::write (f, " = [\n");
  writeVectorElems (f, v, 0);
  Ofstream::write (f, "    ],\n");
  ()
}

fn writeVectorElems(f: & Ofstream, v: &Vector String, index: i64) -> () {
  if index >= Vector::size v {
     return ()
  }
  Ofstream::write (f, "        \"");
  Ofstream::write (f, Vector::get (v, index));
  Ofstream::write (f, "\",\n");
  writeVectorElems (f, v, index + 1)
}

fn writeTargets(f: & Ofstream, targets: &Vector BazelTarget, index: i64) -> () {
  if index >= Vector::size targets {
     return ()
  }
  let target: BazelTarget = Vector::get (targets, index);
  Ofstream::write (f, target.kind);
  Ofstream::write (f, "(\n");
  
  Ofstream::write (f, "    name = \"");
  Ofstream::write (f, target.name);
  Ofstream::write (f, "\",\n");
  
  let srcs: Vector String = target.srcs;
  let srcsLabel: String = "    srcs";
  writeVector (f, &srcsLabel, &srcs);
  let hdrs: Vector String = target.hdrs;
  let hdrsLabel: String = "    hdrs";
  writeVector (f, &hdrsLabel, &hdrs);
  
  Ofstream::write (f, "    copts = [\n");
  Ofstream::write (f, "        \"-std=c++17\",\n");
  Ofstream::write (f, "        \"-Xassembler\",\n");
  Ofstream::write (f, "        \"--gsframe=no\",\n");
  Ofstream::write (f, "    ],\n");
  
  let deps: Vector String = target.deps;
  let depsLabel: String = "    deps";
  writeVector (f, &depsLabel, &deps);
  
  Ofstream::write (f, ")\n\n");
  
  writeTargets (f, targets, index + 1)
}

fn writeBuildFile(targets: &Vector BazelTarget) -> bool {
  let f: Ofstream = Ofstream::open "out/BUILD";
  if !Ofstream::is_open &f {
     return false
  }
  Ofstream::write (&f, "load(\"@rules_cc//cc:defs.bzl\", \"cc_binary\", \"cc_library\")\n\n");
  writeTargets (&f, targets, 0);
  Ofstream::close &f;
  true
}

fn ensureWorkspaceSetup() -> bool {
  if !fs::create_directories "out" {
     return false
  }
  
  let ws: Ofstream = Ofstream::open "out/WORKSPACE";
  if !Ofstream::is_open &ws {
     return false
  }
  Ofstream::write (&ws, "workspace(name = \"bapel_out\")\n");
  Ofstream::close &ws;

  let mod: Ofstream = Ofstream::open "out/MODULE.bazel";
  if !Ofstream::is_open &mod {
     return false
  }
  Ofstream::write (&mod, "module(name = \"bapel_out\")\n");
  Ofstream::write (&mod, "bazel_dep(name = \"rules_cc\", version = \"0.2.17\")\n");
  Ofstream::close &mod;

  true
}

fn build(moduleID: &String) -> i64 {
  if !ensureWorkspaceSetup () {
     core::print [String] "Failed to setup workspace";
     return 1
  }
  
  let builtModules: Vector String = Vector::mk [String] ();
  let targets: Vector BazelTarget = Vector::mk [BazelTarget] ();
  let err: i64 = buildModule (moduleID, &builtModules, true, &targets);
  if err != 0 {
     return err
  }
  
  if !writeBuildFile (&targets) {
     core::print [String] "Failed to write BUILD file";
     return 1
  }
  
  let dot: String = ".";
  let slash: String = "/";
  let under: String = "_";
  let tempName: String = replaceSeparator (*moduleID, &dot, &under);
  let targetName: String = replaceSeparator (tempName, &slash, &under);
  
  // Safe construction of "//:" to avoid parser comment bugs
  let slash2: String = "/";
  let doubleSlash: String = String::concat (&slash2, &slash2);
  let bazelTarget: String = String::concat (&(String::concat (&doubleSlash, &":")), &targetName);
  
  let bazelArgs: Vector String = Vector::mk [String] ();
  Vector::push_back [String] (&bazelArgs, "build");
  Vector::push_back [String] (&bazelArgs, bazelTarget);
  
  let bazelCmd: String = "bazel";
  let outDir: String = "out";
  let bazelRes: (i64, String) = execInDir (&bazelCmd, &bazelArgs, &outDir);
  if bazelRes.0 != 0 {
     core::print [String] "Bazel build failed";
     core::print [String] bazelRes.1;
     return bazelRes.0
  }
  
  let bazelBinPath: String = fs::join (fs::join ("out", "bazel-bin"), targetName);
  let outputPath: String = fs::join ("out", *moduleID);
  
  if !copyFile (&bazelBinPath, &outputPath) {
     core::print [String] "Failed to copy built binary";
     return 1
  }
  
  0
}

fn getSubArgs(args: &Vector String, start: i64) -> Vector String {
  let sub: Vector String = Vector::mk [String] ();
  sliceArgs (args, start, &sub);
  sub
}

fn sliceArgs(args: &Vector String, index: i64, dst: &Vector String) -> () {
  if index >= Vector::size args {
     return ()
  }
  Vector::push_back [String] (dst, Vector::get (args, index));
  sliceArgs (args, index + 1, dst)
}

pub fn main(argc: args::Argc, argv: args::Argv) -> i32 {
  args::init (argc, argv);
  let args: Vector String = args::get_args ();
  let count: i64 = Vector::size &args;
  
  if count < 2 {
     core::print [String] "expected subcommand, e.g., 'parse', 'cc', 'build', 'query'";
     return 1
  }
  
  let command: String = Vector::get (&args, 1);
  
  if command == "cc" {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/compiler", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "build" {
     if count < 3 {
        core::print [String] "usage: bpl build <module>";
        return 1
     }
     let buildTarget: String = Vector::get (&args, 2);
     let err: i64 = build (&buildTarget);
     return core::i64_to_i32 err
  }
  
  if command == "query" {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/querier", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "parse" {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/parser", subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  core::print [String] (String::concat (&"unknown command: ", &command));
  1
}
