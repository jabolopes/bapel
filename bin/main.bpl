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
  let pos: i64 = s.find (from, 0);
  let from_len: i64 = from.size;
  let to_len: i64 = to.size;
  
  for pos != -1 {
    s <- s.replace (pos, from_len, to);
    pos <- pos + to_len;
    pos <- s.find (from, pos);
  };
  s
}

fn resolveMappedPath(path: &String, moduleID: &String) -> String {
  let dot: String = ".";
  let slash: String = "/";
  let relPath: String = replaceSeparator (*moduleID, &dot, &slash);
  let relPathWithExt: String = relPath.concat &".bpl";
  fs::join (*path, relPathWithExt)
}

fn isPrefixOf(pref: &String, s: &String) -> bool {
  if *s == *pref {
     return true
  }
  let p: String = pref.concat &".";
  let p_len: i64 = p.size;
  let s_len: i64 = s.size;
  if s_len < p_len {
     return false
  }
  let s_view: StringView = s.view;
  let sub_view: StringView = s_view.substr (0, p_len);
  let sub: String = String::from_view sub_view;
  sub == p
}

fn findBestMatch(
    mappings: &Vector PackageMapping,
    moduleID: &String,
    index: i64,
    currentBest: &MatchResult) -> MatchResult {
  
  if index >= mappings.size {
     return *currentBest
  }
  
  let mapping: PackageMapping = mappings.get index;
  
  if mapping.is_prefix {
     let mapping_name: String = mapping.name;
     if isPrefixOf (&mapping_name, moduleID) {
        let prefixLen: i64 = mapping_name.size;
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
  if line.size == 0 {
     return ()
  }
  let line_iss: IStringStream = IStringStream::mk (*line);
  let type_str: String = "";
  let name: String = "";
  let path: String = "";
  
  if !line_iss.read &type_str {
     return ()
  }
  if !line_iss.read &name {
     return ()
  }
  if !line_iss.read &path {
     return ()
  }
  
  let is_prefix: bool = type_str == "PREFIX";
  let mapping: PackageMapping = struct {
     is_prefix = is_prefix,
     name = name,
     path = path
  };
  mappings.push_back mapping;
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
    if line.size > 0 {
      let line_iss: IStringStream = IStringStream::mk line;
      let type_str: String = "";
      let value: String = "";
      if line_iss.read &type_str {
        if line_iss.read &value {
          if type_str == "IMPORT" {
            importModules.push_back value;
            ()
          } else if type_str == "IMPL" {
            implFiles.push_back value;
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
  relPath.concat &".bpl"
}

fn vecContains(v: &Vector String, s: &String, index: i64) -> bool {
  if index >= v.size {
     return false
  }
  if v.get index == *s {
     return true
  }
  vecContains (v, s, index + 1)
}

fn appendVectors(dst: &Vector String, src: &Vector String, index: i64) -> () {
  if index >= src.size {
     return ()
  }
  dst.push_back (src.get index);
  appendVectors (dst, src, index + 1)
}

fn buildImpls(
    implFiles: &Vector String,
    moduleID: &String,
    baseFileDir: &String,
    index: i64,
    srcs: &Vector String,
    hdrs: &Vector String) -> i64 {
  
  if index >= implFiles.size {
     return 0
  }
  let implFile: String = implFiles.get index;
  let fullImplPath: String = fs::join (*baseFileDir, implFile);
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let baseName: String = fs::stem implFile;
     let dot: String = ".";
     let slash: String = "/";
     let baseOutputBasename: String = replaceSeparator (*moduleID, &dot, &slash);
     let implOutBasename: String = (baseOutputBasename.concat &"-").concat &baseName;
     let outPath: String = fs::join ("out", implOutBasename);
     let outCcPath: String = outPath.concat &".cc";
     
     if !fs::create_directories (fs::parent_path outCcPath) {
        core::print [String] (("Failed to create directory: " [String]).concat &(fs::parent_path outCcPath));
        return 1
     }
     
     let ccArgs: Vector String = Vector::mk [String] ();
     ccArgs.push_back "-o";
     ccArgs.push_back outCcPath;
     ccArgs.push_back fullImplPath;
     
     let ccRes: (i64, String) = os::exec ("bootstrap/compiler", ccArgs);
     if ccRes.0 != 0 {
        core::print [String] (("Failed to compile impl: " [String]).concat &fullImplPath);
        core::print [String] ccRes.1;
        return ccRes.0
     }
     
     srcs.push_back (implOutBasename.concat &".cc");
     
  } else {
     let dst: String = fs::join ("out", fullImplPath);
     if !copyFile (&fullImplPath, &dst) {
        core::print [String] (("Failed to copy impl: " [String]).concat &fullImplPath);
        return 1
     }
     
     if ext == ".cc" {
        srcs.push_back fullImplPath;
        ()
     } else if ext == ".cpp" {
        srcs.push_back fullImplPath;
        ()
     } else if ext == ".h" {
        hdrs.push_back fullImplPath;
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
  
  if index >= implFiles.size {
     return 0
  }
  let implFile: String = implFiles.get index;
  let ext: String = fs::extension implFile;
  
  if ext == ".bpl" {
     let fullImplPath: String = fs::join (*baseFileDir, implFile);
     let args: Vector String = Vector::mk [String] ();
     args.push_back "-format=flat";
     args.push_back fullImplPath;
     let res: (i64, String) = os::exec ("bootstrap/parser", args);
     if res.0 != 0 {
        core::print [String] (("Failed to parse impl for imports: " [String]).concat &fullImplPath);
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
  
  if index >= importModules.size {
     return 0
  }
  let imp: String = importModules.get index;
  let err: i64 = buildModule (&imp, builtModules, false, targets);
  if err != 0 {
     return err
  }
  
  let dot: String = ".";
  let slash: String = "/";
  let under: String = "_";
  let tempName: String = replaceSeparator (imp, &dot, &under);
  let sanitized: String = replaceSeparator (tempName, &slash, &under);
  let depTarget: String = (":" [String]).concat &sanitized;
  deps.push_back depTarget;
  
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
     core::print [String] (("File not found: " [String]).concat &baseFile);
     return 1
  }
  
  let args: Vector String = Vector::mk [String] ();
  args.push_back "-format=flat";
  args.push_back baseFile;
  let res: (i64, String) = os::exec ("bootstrap/parser", args);
  if res.0 != 0 {
     core::print [String] (("Failed to parse: " [String]).concat &baseFile);
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
  let outHeader: String = outPath.concat &".h";
  
  if !fs::create_directories (fs::parent_path outHeader) {
     core::print [String] (("Failed to create directory: " [String]).concat &(fs::parent_path outHeader));
     return 1
  }
  
  let ccArgs: Vector String = Vector::mk [String] ();
  ccArgs.push_back "-o";
  ccArgs.push_back outHeader;
  ccArgs.push_back baseFile;
  
  let ccRes: (i64, String) = os::exec ("bootstrap/compiler", ccArgs);
  if ccRes.0 != 0 {
     core::print [String] (("Failed to compile: " [String]).concat &baseFile);
     core::print [String] ccRes.1;
     return ccRes.0
  }
  
  let srcs: Vector String = Vector::mk [String] ();
  let hdrs: Vector String = Vector::mk [String] ();
  
  srcs.push_back (baseOutputBasename.concat &".cc");
  hdrs.push_back (baseOutputBasename.concat &".h");
  hdrs.push_back (baseOutputBasename.concat &"_private.h");
  
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
     targets.push_back target;
     ()
  } else {
     let target: BazelTarget = struct {
        kind = "cc_library",
        name = targetName,
        srcs = srcs,
        hdrs = hdrs,
        deps = deps
     };
     targets.push_back target;
     ()
  }
  
  builtModules.push_back (*moduleID);
  0
}

fn writeVector(f: & Ofstream, label: &String, v: &Vector String) -> () {
  if v.size == 0 {
     return ()
  }
  f.write (*label);
  f.write " = [\n";
  writeVectorElems (f, v, 0);
  f.write "    ],\n";
  ()
}

fn writeVectorElems(f: & Ofstream, v: &Vector String, index: i64) -> () {
  if index >= v.size {
     return ()
  }
  f.write "        \"";
  f.write (v.get index);
  f.write "\",\n";
  writeVectorElems (f, v, index + 1)
}

fn writeTargets(f: & Ofstream, targets: &Vector BazelTarget, index: i64) -> () {
  if index >= targets.size {
     return ()
  }
  let target: BazelTarget = targets.get index;
  f.write target.kind;
  f.write "(\n";
  
  f.write "    name = \"";
  f.write target.name;
  f.write "\",\n";
  
  let srcs: Vector String = target.srcs;
  let srcsLabel: String = "    srcs";
  writeVector (f, &srcsLabel, &srcs);
  let hdrs: Vector String = target.hdrs;
  let hdrsLabel: String = "    hdrs";
  writeVector (f, &hdrsLabel, &hdrs);
  
  f.write "    copts = [\n";
  f.write "        \"-std=c++17\",\n";
  f.write "        \"-Xassembler\",\n";
  f.write "        \"--gsframe=no\",\n";
  f.write "    ],\n";
  
  let deps: Vector String = target.deps;
  let depsLabel: String = "    deps";
  writeVector (f, &depsLabel, &deps);
  
  f.write ")\n\n";
  
  writeTargets (f, targets, index + 1)
}

fn writeBuildFile(targets: &Vector BazelTarget) -> bool {
  let f: Ofstream = Ofstream::open "out/BUILD";
  if !f.is_open {
     return false
  }
  f.write "load(\"@rules_cc//cc:defs.bzl\", \"cc_binary\", \"cc_library\")\n\n";
  writeTargets (&f, targets, 0);
  f.close;
  true
}

fn ensureWorkspaceSetup() -> bool {
  if !fs::create_directories "out" {
     return false
  }
  
  let ws: Ofstream = Ofstream::open "out/WORKSPACE";
  if !ws.is_open {
     return false
  }
  ws.write "workspace(name = \"bapel_out\")\n";
  ws.close;

  let mod: Ofstream = Ofstream::open "out/MODULE.bazel";
  if !mod.is_open {
     return false
  }
  mod.write "module(name = \"bapel_out\")\n";
  mod.write "bazel_dep(name = \"rules_cc\", version = \"0.2.17\")\n";
  mod.close;

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
  let doubleSlash: String = slash2.concat &slash2;
  let bazelTarget: String = (doubleSlash.concat &":").concat &targetName;
  
  let bazelArgs: Vector String = Vector::mk [String] ();
  bazelArgs.push_back "build";
  bazelArgs.push_back bazelTarget;
  
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
  if index >= args.size {
     return ()
  }
  dst.push_back (args.get index);
  sliceArgs (args, index + 1, dst)
}

pub fn main(argc: args::Argc, argv: args::Argv) -> i32 {
  args::init (argc, argv);
  let args: Vector String = args::get_args ();
  let count: i64 = args.size;
  
  if count < 2 {
     core::print [String] "expected subcommand, e.g., 'parse', 'cc', 'build', 'query'";
     return 1
  }
  
  let command: String = args.get 1;
  
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
     let buildTarget: String = args.get 2;
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
  
  core::print [String] (("unknown command: " [String]).concat &command);
  1
}
