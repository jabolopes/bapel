implements bin.main

type PackageMapping = struct {
  is_prefix: bool,
  name: String,
  path: String
}

type ModuleFinder = struct {
  modules_by_name: UnorderedMap String String,
  modules_by_prefix: UnorderedMap String String
}

type SourceFileQuery = struct {
  import_modules: Vector String,
  impl_files: Vector String,
  flag_files: Vector String,
  declarations: Vector String,
  trait_implementations: Vector String
}

type ModuleQuery = struct {
  import_modules: Vector String,
  impl_files: Vector String,
  flag_files: Vector String,
  declarations: Vector String,
  trait_implementations: Vector String
}

fn process_workspace_line(
    line: &String,
    modules_by_prefix: &UnorderedMap String String,
    modules_by_name: &UnorderedMap String String) -> () {
  if line.size == 0 {
     return ()
  }
  let line_iss: IStringStream = IStringStream::mk (*line);
  let type_str: String = "".to_string;
  let name: String = "".to_string;
  let path: String = "".to_string;
  if !line_iss.read &type_str {
     return ()
  }
  if !line_iss.read &name {
     return ()
  }
  if !line_iss.read &path {
     return ()
  }
  if type_str == "PREFIX".to_string {
     UnorderedMap::insert [String, String] (modules_by_prefix, name, path);
     ()
  } else if type_str == "MODULE".to_string {
     UnorderedMap::insert [String, String] (modules_by_name, name, path);
     ()
  } else {
     ()
  }
}

fn mk_module_finder() -> ModuleFinder {
  let modules_by_name: UnorderedMap String String = UnorderedMap::mk [String, String] ();
  let modules_by_prefix: UnorderedMap String String = UnorderedMap::mk [String, String] ();
  let ws_file: String = "workspace.bpl".to_string;
  if fs::exists ws_file {
     let args: Vector String = Vector::mk [String] ();
     Vector::push_back [String] (&args, "-workspace".to_string);
     Vector::push_back [String] (&args, "-format=flat".to_string);
     Vector::push_back [String] (&args, ws_file);
     let res: (i64, String) = os::exec ("bootstrap/parser".to_string, args);
     if res.0 == 0 {
        let flat_text: String = res.1;
        let iss: IStringStream = IStringStream::mk flat_text;
        let line: String = "".to_string;
        for getline (&iss, &line) {
           process_workspace_line (&line, &modules_by_prefix, &modules_by_name);
        };
     }
  }
  struct {
     modules_by_name = modules_by_name,
     modules_by_prefix = modules_by_prefix
  }
}

fn lookup_by_prefix_step(finder: &ModuleFinder, name: &String) -> (bool, String) {
  let target: String = ".".to_string;
  let index: i64 = name.rfind &target;
  if index == -1 {
     return (false, "".to_string)
  }
  let sv: StringView = String::view name;
  let prefix_sv: StringView = sv.substr (0, index);
  let prefix_str: String = prefix_sv.to_string;
  if UnorderedMap::contains [String, String] (&(*finder).modules_by_prefix, &prefix_str) {
     let opt_val: Optional String = UnorderedMap::get [String, String] (&(*finder).modules_by_prefix, &prefix_str);
     return (true, Optional::get_value [String] &opt_val)
  }
  lookup_by_prefix_step (finder, &prefix_str)
}

fn lookup_module(finder: &ModuleFinder, mod_id: &String) -> (bool, String) {
  if UnorderedMap::contains [String, String] (&(*finder).modules_by_name, mod_id) {
    let opt_val: Optional String = UnorderedMap::get [String, String] (&(*finder).modules_by_name, mod_id);
    return (true, Optional::get_value [String] &opt_val)
  }
  lookup_by_prefix_step (finder, mod_id)
}

fn base_filename(finder: &ModuleFinder, mod_id: &String) -> String {
  let package_name: String = "".to_string;
  let res: (bool, String) = lookup_module (finder, mod_id);
  if res.0 {
    package_name <- res.1
  };
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let rel_path: String = replaceSeparator (*mod_id, &dot, &slash);
  let full_path: String = fs::join (package_name, rel_path);
  full_path.concat &".bpl".to_string
}

fn impl_filename(base_file: &String, rel_impl: &String) -> String {
  let parent: String = fs::parent_path (*base_file);
  fs::join (parent, *rel_impl)
}

fn query_annotation_file(path: &String) -> SourceFileQuery {
  let import_modules: Vector String = Vector::mk [String] ();
  let impl_files: Vector String = Vector::mk [String] ();
  let flag_files: Vector String = Vector::mk [String] ();
  let declarations: Vector String = Vector::mk [String] ();
  let trait_implementations: Vector String = Vector::mk [String] ();

  let f: Ifstream = Ifstream::open path;
  if f.is_open {
    let line: String = "".to_string;
    for getline (&f, &line) {
      if line.size > 0 {
        let import_pref: String = "import ".to_string;
        let import_part: String = "import :".to_string;
        let semi: String = ";".to_string;
        if String::starts_with (&line, &import_pref) {
          if !String::starts_with (&line, &import_part) {
            if String::ends_with (&line, &semi) {
              let sv: StringView = String::view &line;
              StringView::remove_prefix (&sv, import_pref.size);
              StringView::remove_suffix (&sv, semi.size);
              Vector::push_back [String] (&import_modules, StringView::to_string sv);
              ()
            }
          }
        };
        let bpl_pref: String = "// @bpl: ".to_string;
        if String::starts_with (&line, &bpl_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, bpl_pref.size);
          Vector::push_back [String] (&declarations, StringView::to_string sv);
          ()
        }
      }
    };
    f.close;
    ()
  };

  struct {
    import_modules = import_modules,
    impl_files = impl_files,
    flag_files = flag_files,
    declarations = declarations,
    trait_implementations = trait_implementations
  }
}

fn query_bpl_file(path: &String) -> SourceFileQuery {
  let import_modules: Vector String = Vector::mk [String] ();
  let impl_files: Vector String = Vector::mk [String] ();
  let flag_files: Vector String = Vector::mk [String] ();
  let declarations: Vector String = Vector::mk [String] ();
  let trait_implementations: Vector String = Vector::mk [String] ();

  let args: Vector String = Vector::mk [String] ();
  Vector::push_back [String] (&args, "-format=flat".to_string);
  Vector::push_back [String] (&args, *path);
  let res: (i64, String) = os::exec ("bootstrap/parser".to_string, args);
  if res.0 == 0 {
    let flat_text: String = res.1;
    let iss: IStringStream = IStringStream::mk flat_text;
    let line: String = "".to_string;
    for getline (&iss, &line) {
      if line.size > 0 {
        let import_pref: String = "IMPORT ".to_string;
        let impl_pref: String = "IMPL ".to_string;
        let flag_pref: String = "FLAG ".to_string;
        let decl_pref: String = "DECL ".to_string;
        let trait_impl_pref: String = "TRAIT_IMPL ".to_string;

        if String::starts_with (&line, &import_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, import_pref.size);
          Vector::push_back [String] (&import_modules, StringView::to_string sv);
          ()
        } else if String::starts_with (&line, &impl_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, impl_pref.size);
          Vector::push_back [String] (&impl_files, StringView::to_string sv);
          ()
        } else if String::starts_with (&line, &flag_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, flag_pref.size);
          Vector::push_back [String] (&flag_files, StringView::to_string sv);
          ()
        } else if String::starts_with (&line, &decl_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, decl_pref.size);
          let raw_str: String = StringView::to_string sv;
          let unescaped: String = replaceSeparator (raw_str, &"\\n".to_string, &"\n".to_string);
          Vector::push_back [String] (&declarations, unescaped);
          ()
        } else if String::starts_with (&line, &trait_impl_pref) {
          let sv: StringView = String::view &line;
          StringView::remove_prefix (&sv, trait_impl_pref.size);
          let raw_str: String = StringView::to_string sv;
          let unescaped: String = replaceSeparator (raw_str, &"\\n".to_string, &"\n".to_string);
          Vector::push_back [String] (&trait_implementations, unescaped);
          ()
        } else {
          ()
        }
      }
    };
  } else {
    core::print [String] (("Failed to parse file: ".to_string).concat path);
    core::print [String] res.1;
    ()
  };

  struct {
    import_modules = import_modules,
    impl_files = impl_files,
    flag_files = flag_files,
    declarations = declarations,
    trait_implementations = trait_implementations
  }
}

fn query_source_file(path: &String) -> SourceFileQuery {
  if fs::extension (*path) == ".bpl".to_string {
    return query_bpl_file path
  }
  query_annotation_file path
}

