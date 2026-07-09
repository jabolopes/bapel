implements program

imports {
  bapel.core
  bapel.stl
}

fn testUnorderedMap() -> () {
  let m: UnorderedMap String i64 = UnorderedMap::mk [String, i64] ();
  core::print [bool] (UnorderedMap::empty [String, i64] &m);
  core::print [i64] (UnorderedMap::size [String, i64] &m);

  let key1: String = "hello".to_string;
  UnorderedMap::insert [String, i64] (&m, key1, 42);

  core::print [bool] (UnorderedMap::empty [String, i64] &m);
  core::print [i64] (UnorderedMap::size [String, i64] &m);

  let key2: String = "hello".to_string;
  core::print [bool] (UnorderedMap::contains [String, i64] (&m, &key2));

  let key3: String = "world".to_string;
  core::print [bool] (UnorderedMap::contains [String, i64] (&m, &key3));

  let opt1: Optional i64 = UnorderedMap::get [String, i64] (&m, &key2);
  core::print [bool] (Optional::has_value [i64] &opt1);
  core::print [i64] (Optional::get_value [i64] &opt1);

  let opt2: Optional i64 = UnorderedMap::get [String, i64] (&m, &key3);
  core::print [bool] (Optional::has_value [i64] &opt2);
  ()
}
