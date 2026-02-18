implements program

fn conditionals() -> () {
  true
  false

  if true {
    true
  }

  if !true {
    false
  }

  if true == false {
    false
  }

  let v1: bool = if true == false {
    false
  } else {
    true
  }

  if true == false {
    false
  } else if false == true {
    true
  }

  let v2: bool = if true == false {
    false
  } else if false == true {
    true
  }

  if true == false {
    false
  } else if false == true {
    true
  } else {
    false
  }

  let v3: bool = if true == false {
    false
  } else if false == true {
    true
  } else {
    false
  }

  ()
}

fn ftrue() -> bool {
  true
}

fn id['a](x: 'a) -> 'a {
  x
}

fn fconst['a, 'b](x: 'a, y: 'b) -> 'a {
  x
}

fn conditionalsPolymorphic() -> bool {
  if id [bool] true {
    id [bool] true
  } else {
    id [bool] false
  }

  if id true {
    id true
  } else {
    id false
  }

  if fconst [bool, ()] (true, ()) {
    true
  } else {
    false
  }

  if fconst (true, ()) {
    true
  } else {
    false
  }

  if id true == id false {
    id true == id false
  } else {
    id false == id true
  }
}
