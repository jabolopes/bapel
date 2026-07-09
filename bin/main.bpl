module bin.main

imports {
  bapel.args
  bapel.core
  bapel.os
  bapel.stl
}

impls {
  "query.bpl"
}



type BazelTarget = struct {
  kind: String,
  name: String,
  srcs: Vector String,
  hdrs: Vector String,
  deps: Vector String
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
     let res: (i64, String) = (-1, "chdir failed".to_string);
     return res
  }
  let res: (i64, String) = os::exec (*cmd, *args);
  if !fs::set_current_path origPath {
     core::print [String] "Warning: failed to restore CWD".to_string;
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

fn resolveModule(moduleID: &String) -> String {
  let finder: ModuleFinder = mk_module_finder ();
  base_filename (&finder, moduleID)
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
  
  if ext == ".bpl".to_string {
     let baseName: String = fs::stem implFile;
     let dot: String = ".".to_string;
     let slash: String = "/".to_string;
     let baseOutputBasename: String = replaceSeparator (*moduleID, &dot, &slash);
     let implOutBasename: String = (baseOutputBasename.concat &"-".to_string).concat &baseName;
     let outPath: String = fs::join ("out".to_string, implOutBasename);
     let outCcPath: String = outPath.concat &".cc".to_string;
     
     if !fs::create_directories (fs::parent_path outCcPath) {
        core::print [String] (("Failed to create directory: ".to_string).concat &(fs::parent_path outCcPath));
        return 1
     }
     
     let ccArgs: Vector String = Vector::mk [String] ();
     ccArgs.push_back "-o".to_string;
     ccArgs.push_back outCcPath;
     ccArgs.push_back fullImplPath;
     
     let ccRes: (i64, String) = os::exec ("bootstrap/compiler".to_string, ccArgs);
     if ccRes.0 != 0 {
        core::print [String] (("Failed to compile impl: ".to_string).concat &fullImplPath);
        core::print [String] ccRes.1;
        return ccRes.0
     }
     
     srcs.push_back (implOutBasename.concat &".cc".to_string);
     
  } else {
     let dst: String = fs::join ("out".to_string, fullImplPath);
     if !copyFile (&fullImplPath, &dst) {
        core::print [String] (("Failed to copy impl: ".to_string).concat &fullImplPath);
        return 1
     }
     
     if ext == ".cc".to_string {
        srcs.push_back fullImplPath;
        ()
     } else if ext == ".cpp".to_string {
        srcs.push_back fullImplPath;
        ()
     } else if ext == ".h".to_string {
        hdrs.push_back fullImplPath;
        ()
     }
  }
  
  buildImpls (implFiles, moduleID, baseFileDir, index + 1, srcs, hdrs)
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
  
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let under: String = "_".to_string;
  let tempName: String = replaceSeparator (imp, &dot, &under);
  let sanitized: String = replaceSeparator (tempName, &slash, &under);
  let depTarget: String = (":".to_string).concat &sanitized;
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
  
  let finder: ModuleFinder = mk_module_finder ();
  let mod_query: ModuleQuery = query_module (&finder, moduleID);
  let importsList: Vector String = mod_query.import_modules;
  let implsList: Vector String = mod_query.impl_files;
  let baseFile: String = base_filename (&finder, moduleID);
  if !fs::exists baseFile {
     core::print [String] (("File not found: ".to_string).concat &baseFile);
     return 1
  }
  let baseFileDir: String = fs::parent_path baseFile;
  
  let deps: Vector String = Vector::mk [String] ();
  let err: i64 = buildImports (&importsList, builtModules, &deps, 0, targets);
  if err != 0 {
     return err
  }
  
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let under: String = "_".to_string;
  let baseOutputBasename: String = replaceSeparator (*moduleID, &dot, &slash);
  let outPath: String = fs::join ("out".to_string, baseOutputBasename);
  let outHeader: String = outPath.concat &".h".to_string;
  
  if !fs::create_directories (fs::parent_path outHeader) {
     core::print [String] (("Failed to create directory: ".to_string).concat &(fs::parent_path outHeader));
     return 1
  }
  
  let ccArgs: Vector String = Vector::mk [String] ();
  ccArgs.push_back "-o".to_string;
  ccArgs.push_back outHeader;
  ccArgs.push_back baseFile;
  
  let ccRes: (i64, String) = os::exec ("bootstrap/compiler".to_string, ccArgs);
  if ccRes.0 != 0 {
     core::print [String] (("Failed to compile: ".to_string).concat &baseFile);
     core::print [String] ccRes.1;
     return ccRes.0
  }
  
  let srcs: Vector String = Vector::mk [String] ();
  let hdrs: Vector String = Vector::mk [String] ();
  
  srcs.push_back (baseOutputBasename.concat &".cc".to_string);
  hdrs.push_back (baseOutputBasename.concat &".h".to_string);
  hdrs.push_back (baseOutputBasename.concat &"_private.h".to_string);
  
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
        kind = "cc_binary".to_string,
        name = targetName,
        srcs = srcs,
        hdrs = emptyHdrs,
        deps = deps
     };
     targets.push_back target;
     ()
  } else {
     let target: BazelTarget = struct {
        kind = "cc_library".to_string,
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
  f.write " = [\n".to_string;
  writeVectorElems (f, v, 0);
  f.write "    ],\n".to_string;
  ()
}

fn writeVectorElems(f: & Ofstream, v: &Vector String, index: i64) -> () {
  if index >= v.size {
     return ()
  }
  f.write "        \"".to_string;
  f.write (v.get index);
  f.write "\",\n".to_string;
  writeVectorElems (f, v, index + 1)
}

fn writeTargets(f: & Ofstream, targets: &Vector BazelTarget, index: i64) -> () {
  if index >= targets.size {
     return ()
  }
  let target: BazelTarget = targets.get index;
  f.write target.kind;
  f.write "(\n".to_string;
  
  f.write "    name = \"".to_string;
  f.write target.name;
  f.write "\",\n".to_string;
  
  let srcs: Vector String = target.srcs;
  let srcsLabel: String = "    srcs".to_string;
  writeVector (f, &srcsLabel, &srcs);
  let hdrs: Vector String = target.hdrs;
  let hdrsLabel: String = "    hdrs".to_string;
  writeVector (f, &hdrsLabel, &hdrs);
  
  f.write "    copts = [\n".to_string;
  f.write "        \"-std=c++17\",\n".to_string;
  f.write "        \"-Xassembler\",\n".to_string;
  f.write "        \"--gsframe=no\",\n".to_string;
  f.write "    ],\n".to_string;
  
  let deps: Vector String = target.deps;
  let depsLabel: String = "    deps".to_string;
  writeVector (f, &depsLabel, &deps);
  
  f.write ")\n\n".to_string;
  
  writeTargets (f, targets, index + 1)
}

fn writeBuildFile(targets: &Vector BazelTarget) -> bool {
  let f: Ofstream = Ofstream::open &"out/BUILD".to_string;
  if !f.is_open {
     return false
  }
  f.write "load(\"@rules_cc//cc:defs.bzl\", \"cc_binary\", \"cc_library\")\n\n".to_string;
  writeTargets (&f, targets, 0);
  f.close;
  true
}

fn ensureWorkspaceSetup() -> bool {
  if !fs::create_directories "out".to_string {
     return false
  }
  
  let ws: Ofstream = Ofstream::open &"out/WORKSPACE".to_string;
  if !ws.is_open {
     return false
  }
  ws.write "workspace(name = \"bapel_out\")\n".to_string;
  ws.close;

  let mod: Ofstream = Ofstream::open &"out/MODULE.bazel".to_string;
  if !mod.is_open {
     return false
  }
  mod.write "module(name = \"bapel_out\")\n".to_string;
  mod.write "bazel_dep(name = \"rules_cc\", version = \"0.2.17\")\n".to_string;
  mod.close;

  true
}

fn build(moduleID: &String) -> i64 {
  if !ensureWorkspaceSetup () {
     core::print [String] "Failed to setup workspace".to_string;
     return 1
  }
  
  let builtModules: Vector String = Vector::mk [String] ();
  let targets: Vector BazelTarget = Vector::mk [BazelTarget] ();
  let err: i64 = buildModule (moduleID, &builtModules, true, &targets);
  if err != 0 {
     return err
  }
  
  if !writeBuildFile (&targets) {
     core::print [String] "Failed to write BUILD file".to_string;
     return 1
  }
  
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let under: String = "_".to_string;
  let tempName: String = replaceSeparator (*moduleID, &dot, &under);
  let targetName: String = replaceSeparator (tempName, &slash, &under);
  
  // Safe construction of "//:" to avoid parser comment bugs
  let slash2: String = "/".to_string;
  let doubleSlash: String = slash2.concat &slash2;
  let bazelTarget: String = (doubleSlash.concat &":".to_string).concat &targetName;
  
  let bazelArgs: Vector String = Vector::mk [String] ();
  bazelArgs.push_back "build".to_string;
  bazelArgs.push_back bazelTarget;
  
  let bazelCmd: String = "bazel".to_string;
  let outDir: String = "out".to_string;
  let bazelRes: (i64, String) = execInDir (&bazelCmd, &bazelArgs, &outDir);
  if bazelRes.0 != 0 {
     core::print [String] "Bazel build failed".to_string;
     core::print [String] bazelRes.1;
     return bazelRes.0
  }
  
  let bazelBinPath: String = fs::join (fs::join ("out".to_string, "bazel-bin".to_string), targetName);
  let outputPath: String = fs::join ("out".to_string, *moduleID);
  
  if !copyFile (&bazelBinPath, &outputPath) {
     core::print [String] "Failed to copy built binary".to_string;
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
     core::print [String] "expected subcommand, e.g., 'parse', 'cc', 'build', 'query'".to_string;
     return 1
  }
  
  let command: String = args.get 1;
  
  if command == "cc".to_string {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/compiler".to_string, subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "build".to_string {
     if count < 3 {
        core::print [String] "usage: bpl build <module>".to_string;
        return 1
     }
     let buildTarget: String = args.get 2;
     let err: i64 = build (&buildTarget);
     return core::i64_to_i32 err
  }
  
  if command == "query".to_string {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/querier".to_string, subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  if command == "parse".to_string {
     let subArgs: Vector String = getSubArgs (&args, 2);
     let res: (i64, String) = os::exec ("bootstrap/parser".to_string, subArgs);
     core::print [String] res.1;
     return core::i64_to_i32 res.0
  }
  
  core::print [String] (("unknown command: ".to_string).concat &command);
  1
}
