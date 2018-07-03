package sack

// ParentKey the key in sack indicate it's parent, don't use it directly
const ParentKey = "+_-PARENT-_+"

// Sack is basically a map from string to anything
// it's nil-safe, means all function call won't panic if s == nil
type Sack map[string]interface{}

// Parent get the parent
func (s Sack) Parent() (p Sack) {
	if s == nil {
		return nil
	}
	if p, ok := s[ParentKey].(Sack); ok {
		return p
	}
	return nil
}

// SetParent set the parent of a Sack
func (s Sack) SetParent(p Sack) {
	if s == nil {
		return
	}
	s[ParentKey] = p
}

// Value get a value from Sack, if not found, find in it's parent
func (s Sack) Value(k string) (v interface{}) {
	if s == nil {
		return nil
	}
	v = s[k]
	if v == nil {
		p := s.Parent()
		if p != nil {
			v = p.Value(k)
		}
	}
	return
}

// Set set a value for key
func (s Sack) Set(k string, v interface{}) {
	if s == nil {
		return
	}
	s[k] = v
}
