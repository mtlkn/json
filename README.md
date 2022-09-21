# JSON
JSON parser and structure 

## Usage
`jo, err := ParseObject(strings.NewReader("{\"id\":1, \"name\":\"YM\"}"))`

`fmt.Println(jo.String())`

```
{"id":1,"name":"YM"}
```

`jo.Set("name", "Yuri Metelkin")`

`jo.Set("history", A(O(P("date","2022-09-21"), P("action", "test")), O(P("date","2022-09-22"), P("action", "test"))))`

`jo = jo.ExcludeFields([]string{"id"})`

`fmt.Println(jo.String())`

```
{"name":"Yuri Metelkin","history":[{"date":"2022-09-21","action"test"},{"date":"2022-09-22","action"test"}]}
```

