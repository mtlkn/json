# JSON parser and object model

## The Parser

The parser does not fully validate parsed values. The validation happens in "getter" methods: GetString(), GetInt(), etc.

This approach makes parsing 3 times faster than standard GO parser:

> | Benchmarks `1`           | ops     | ns/op   | B/op    | allocs/op |
> | :---                     |    ---: |    ---: |     --: |      ---: |
> | Standard small file `2`  | 143,935 | 8,303   | 2,744   | 89        |
> | This parser small file   | 480,003 | 2,532   | 2,528   | 60        |
> | Standard large file `3`  | 1,932   | 632,392 | 58,125  | 5,154     |
> | This parser large file   | 5,436   | 224,327 | 197,816 | 4,590     |
>
> `1` cpu: 11th Gen Intel(R) Core(TM) i7-11850H @ 2.50GHz 
>
> `2`  1KB
>
> `3` 90KB



```
jo, err := ParseObjectString(`{ "name": "YM", "age": 27 }`)
ja, err := ParseArrayString(`[ "xyz", 123 ]`)
```

## The Model

```
// create new object and add two properites
jo := New().Add("name", "YM").Add("age", 27)

// add more properties
jo.Add("pi", 3.14)
jo.Add("ok", true)

// remove properties
jo.Remove("pi")

// get value
s, ok := jo.GetString("name")
i, ok := jo.GetInt("age")

// create new array
ja := NewArray([]string{ "abc", "xyz" })

// add more values
ja.Add(3.14)
ja.Add(true)

// remove value
ja.Remove(3)

// get value
s, ok := ja.GetString(0)
i, ok := ja.GetFloat(2)

```