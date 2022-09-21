package json

import "strings"

func (jo *Object) IncludeFields(fields []string) *Object {
	if jo == nil || len(fields) == 0 {
		return nil
	}

	tree := fieldTree(fields)
	return jo.includeFields(tree)
}

func (jo *Object) ExcludeFields(fields []string) *Object {
	if jo == nil || len(fields) == 0 {
		return jo
	}

	tree := fieldTree(fields)
	return jo.excludeFields(tree)
}

type fieldTreeNode map[string]fieldTreeNode

func fieldTree(fields []string) fieldTreeNode {
	tree := make(map[string]fieldTreeNode)
	set := make(map[string]struct{})
	multi := make(map[string][]string)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		idx := strings.Index(field, ".")
		if idx == -1 {
			if _, ok := set[field]; ok {
				continue
			}

			tree[field] = nil
			set[field] = struct{}{}
			continue
		}

		l := field[:idx]
		r := field[idx+1:]

		if _, ok := set[l]; !ok {
			tree[l] = nil
		}

		set[l] = struct{}{}
		multi[l] = append(multi[l], r)
	}

	for field, fs := range multi {
		tree[field] = fieldTree(fs)
	}

	return tree
}

func (jo *Object) includeFields(tree fieldTreeNode) *Object {
	o := NewObject()

	for _, p := range jo.Properties {
		fs, ok := tree[p.Name()]
		if !ok {
			continue
		}

		if len(fs) == 0 {
			o.Set(p.Name(), p.Value().Value())
			continue
		}

		switch p.Value().Type() {
		case OBJECT:
			v, ok := p.Value().Object()
			if ok {
				o.Set(p.Name(), v.includeFields(fs))
			}
		case ARRAY:
			a, ok := p.Value().Array()
			if ok {
				o.Set(p.Name(), a.includeFields(fs))
			}
		}
	}

	return o
}

func (ja *Array) includeFields(tree fieldTreeNode) *Array {
	a := NewArray()

	for _, v := range ja.Values {
		switch v.Type() {
		case OBJECT:
			o, ok := v.Object()
			if ok {
				a.Values = append(a.Values, New(o.includeFields(tree)))
			}
		case ARRAY:
			t, ok := v.Array()
			if ok {
				a.Values = append(a.Values, New(t.includeFields(tree)))
			}
		}
	}

	return a
}

func (jo *Object) excludeFields(tree fieldTreeNode) *Object {
	o := NewObject()

	for _, p := range jo.Properties {
		fs, ok := tree[p.Name()]
		if !ok {
			o.Set(p.Name(), p.Value().Value())
			continue
		}

		if len(fs) == 0 {
			continue
		}

		switch p.Value().Type() {
		case OBJECT:
			v, ok := p.Value().Object()
			if ok {
				o.Set(p.Name(), v.excludeFields(fs))
			}
		case ARRAY:
			a, ok := p.Value().Array()
			if ok {
				o.Set(p.Name(), a.excludeFields(fs))
			}
		}
	}

	return o
}

func (ja *Array) excludeFields(tree fieldTreeNode) *Array {
	a := NewArray()

	for _, v := range ja.Values {
		switch v.Type() {
		case OBJECT:
			o, ok := v.Object()
			if ok {
				a.Values = append(a.Values, New(o.excludeFields(tree)))
			}
		case ARRAY:
			t, ok := v.Array()
			if ok {
				a.Values = append(a.Values, New(t.excludeFields(tree)))
			}
		}
	}

	return a
}

/*
func (jo *Object) setFields(tree fieldTreeNode, include bool) *Object {
	o := NewObject()

	for _, p := range jo.Properties {
		fs, ok := tree[p.Name()]
		if (ok && !include) || (!ok && include) {
			continue
		}

		if len(fs) == 0 {
			o.Set(p.Name(), p.Value().Value())
			continue
		}

		switch p.Value().Type() {
		case OBJECT:
			v, ok := p.Value().Object()
			if ok {
				o.Set(p.Name(), v.setFields(fs, include))
			}
		case ARRAY:
			a, ok := p.Value().Array()
			if ok {
				o.Set(p.Name(), a.setFields(fs, include))
			}
		}
	}

	return o
}

func (ja *Array) setFields(tree fieldTreeNode, include bool) *Array {
	a := NewArray()

	for _, v := range ja.Values {
		switch v.Type() {
		case OBJECT:
			o, ok := v.Object()
			if ok {
				a.Values = append(a.Values, New(o.setFields(tree, include)))
			}
		case ARRAY:
			t, ok := v.Array()
			if ok {
				a.Values = append(a.Values, New(t.setFields(tree, include)))
			}
		}
	}

	return a
}
*/
