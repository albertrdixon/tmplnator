package stack

import (
	"testing"
)

var stacktests = []struct {
	list []interface{}
}{
	{[]interface{}{"one", "two", "three"}},
	{[]interface{}{1, 2, 3}},
}

func TestStack(t *testing.T) {
	for _, st := range stacktests {
		i, s := 0, NewStack()
		for _, in := range st.list {
			if s.Len() != i {
				t.Errorf("stack.Len(): %d, want %d", s.Len(), i)
			}
			s.Push(in)
			i++
		}
		for j := 2; j >= 0; j-- {
			item := s.Pop()
			if st.list[j] != item {
				t.Errorf("stack.Pop(): %v, want %v", item, st.list[j])
			}
		}
	}
}
