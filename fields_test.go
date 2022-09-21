package json

import (
	"fmt"
	"strings"
	"testing"
)

func TestFields(t *testing.T) {
	s := `{"id":1,"name":{"first":"Yuri","last":"Metelkin","nick":"YM"},"history":[{"date":"2022-09-21","action":"test"},{"date":"2022-09-20","action":"test"}]}`
	full, err := ParseObject(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Run("include many fields", func(t *testing.T) {
		jo := full.IncludeFields([]string{"id", "name.nick", "history.date"})
		fmt.Println(jo.String())

		if id, _ := jo.GetInt("id"); id != 1 {
			t.Fail()
		}

		if name, ok := jo.GetObject("name"); !ok {
			t.Fail()
		} else if nick, _ := name.GetString("nick"); nick != "YM" {
			t.Fail()
		} else if _, ok := name.Get("first"); ok {
			t.Fail()
		}

		if history, ok := jo.GetObjects("history"); !ok {
			t.Fail()
		} else if len(history) != 2 {
			t.Fail()
		} else if date, _ := history[1].GetString("date"); date != "2022-09-20" {
			t.Fail()
		} else if _, ok := history[1].Get("action"); ok {
			t.Fail()
		}

	})

	t.Run("include one full array field", func(t *testing.T) {
		jo := full.IncludeFields([]string{"history"})
		fmt.Println(jo.String())

		if _, ok := jo.GetString("id"); ok {
			t.Fail()
		}

		if _, ok := jo.GetObject("name"); ok {
			t.Fail()
		}

		if history, ok := jo.GetObjects("history"); !ok {
			t.Fail()
		} else if len(history) != 2 {
			t.Fail()
		} else if s, _ := history[0].GetString("action"); s != "test" {
			t.Fail()
		} else if s, _ := history[1].GetString("date"); s != "2022-09-20" {
			t.Fail()
		}
	})

	t.Run("include array of arrays fields", func(t *testing.T) {
		jo := NewObject(
			P("people", A(
				A(O(P("id", 1), P("name", "YM"))),
				A(O(P("id", 2), P("name", "SV"))),
			),
			),
		)

		jo = jo.IncludeFields([]string{"people.name"})
		fmt.Println(jo.String())

		people, ok := jo.GetArray("people")
		if !ok || len(people.Values) != 2 {
			t.Fail()
		}

		ja, ok := people.Values[1].Array()
		if !ok || len(ja.Values) != 1 {
			t.Fail()
		}

		o, ok := ja.Values[0].Object()
		if !ok {
			t.Fail()
		}

		if s, _ := o.GetString("name"); s != "SV" {
			t.Fail()
		}
	})

	t.Run("bad include fields input", func(t *testing.T) {
		jo := full.IncludeFields(nil)
		if jo != nil {
			t.Fail()
		}

		jo = full.IncludeFields([]string{"id", " ", "name", "id"})
		fmt.Println(jo.String())

		if id, _ := jo.GetInt("id"); id != 1 {
			t.Fail()
		}
	})

	t.Run("exclude fields", func(t *testing.T) {
		jo := full.ExcludeFields([]string{"id", "name.nick", "history.date"})
		fmt.Println(jo.String())

		if _, ok := jo.GetInt("id"); ok {
			t.Fail()
		}

		if name, ok := jo.GetObject("name"); !ok {
			t.Fail()
		} else if _, ok := name.GetString("nick"); ok {
			t.Fail()
		} else if s, _ := name.GetString("first"); s != "Yuri" {
			t.Fail()
		} else if s, _ := name.GetString("last"); s != "Metelkin" {
			t.Fail()
		}

		if history, ok := jo.GetObjects("history"); !ok {
			t.Fail()
		} else if len(history) != 2 {
			t.Fail()
		} else if s, _ := history[1].GetString("action"); s != "test" {
			t.Fail()
		} else if _, ok := history[1].Get("date"); ok {
			t.Fail()
		}

		jo = full.ExcludeFields(nil)
		if jo == nil {
			t.Fail()
		}
		fmt.Println(jo.String())
	})

	t.Run("exclude array of arrays fields", func(t *testing.T) {
		jo := NewObject(
			P("people", A(
				A(O(P("id", 1), P("name", "YM"))),
				A(O(P("id", 2), P("name", "SV"))),
			),
			),
		)

		jo = jo.ExcludeFields([]string{"people.id"})
		fmt.Println(jo.String())

		people, ok := jo.GetArray("people")
		if !ok || len(people.Values) != 2 {
			t.Fail()
		}

		ja, ok := people.Values[1].Array()
		if !ok || len(ja.Values) != 1 {
			t.Fail()
		}

		o, ok := ja.Values[0].Object()
		if !ok {
			t.Fail()
		}

		if s, _ := o.GetString("name"); s != "SV" {
			t.Fail()
		}
	})

	t.Run("include exclude fields", func(t *testing.T) {
		jo := full.IncludeFields([]string{"name"}).ExcludeFields([]string{"name.nick"})
		fmt.Println(jo.String())

		if name, ok := jo.GetObject("name"); !ok {
			t.Fail()
		} else if _, ok := name.GetString("nick"); ok {
			t.Fail()
		} else if s, _ := name.GetString("first"); s != "Yuri" {
			t.Fail()
		} else if s, _ := name.GetString("last"); s != "Metelkin" {
			t.Fail()
		}
	})
}
