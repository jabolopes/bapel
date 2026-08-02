module /* tests/testdata/in/conditionals.in:1 */conditionals

in "tests/testdata/in/conditionals.in" in lines 3-54 fn conditionals() -> () {
  true
  false
  ifthen (true, {
  true
})
  ifthen (! true, {
  false
})
  ifthen (== (true, false), {
  false
})
  let v1: bool = ifelse (== (true, false), {
  false
}, {
  true
})
  ifelse (== (true, false), {
  false
}, ifthen (== (false, true), {
  true
}))
  let v2: bool = ifelse (== (true, false), {
  false
}, ifthen (== (false, true), {
  true
}))
  ifelse (== (true, false), {
  false
}, ifelse (== (false, true), {
  true
}, {
  false
}))
  let v3: bool = ifelse (== (true, false), {
  false
}, ifelse (== (false, true), {
  true
}, {
  false
}))
  ()
}

in "tests/testdata/in/conditionals.in" in lines 56-62 fn ifLastTerm() -> bool {
  ifelse (true, {
  false
}, {
  true
})
}

in "tests/testdata/in/conditionals.in" in lines 64-66 fn ftrue() -> bool {
  true
}

in "tests/testdata/in/conditionals.in" in lines 68-70 fn id ['a](x: 'a) -> 'a {
  x
}

in "tests/testdata/in/conditionals.in" in lines 72-74 fn fconst ['a, 'b](x: 'a, y: 'b) -> 'a {
  x
}

in "tests/testdata/in/conditionals.in" in lines 76-106 fn conditionalsPolymorphic() -> bool {
  ifelse (id [bool] true, {
  id [bool] true
}, {
  id [bool] false
})
  ifelse (id true, {
  id true
}, {
  id false
})
  ifelse (fconst [bool] [()] (true, ()), {
  true
}, {
  false
})
  ifelse (fconst (true, ()), {
  true
}, {
  false
})
  ifelse (== (id true, id false), {
  == (id true, id false)
}, {
  == (id false, id true)
})
}
