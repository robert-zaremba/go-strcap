package strcap

import (
	"fmt"
	"strings"
)

// Builder builds a string with length not bigger then Limit
type Builder struct {
	Limit   int
	Elems   []string
	length  int
	stopped bool
}

// AddStr adds new string to the Builder
func (lb *Builder) Add(ls ...string) bool {
	if lb.stopped {
		return false
	}
	for _, s := range ls {
		n := len(s)
		if lb.length+n > lb.Limit {
			lb.Elems = append(lb.Elems, s[:lb.Limit-lb.length], "...")
			lb.stopped = true
			return false
		}
		lb.Elems = append(lb.Elems, s)
		lb.length += n
	}
	return true
}

// AddAny adds new object to the Builder using fmt.Sprint serializer
func (lb *Builder) AddAny(ls ...interface{}) bool {
	for i := range ls {
		if !lb.Add(fmt.Sprint(ls[i])) {
			return false
		}
	}
	return true
}

// AddSlice adds all elements from the string slice
func (lb *Builder) AddSlice(ls []string) bool {
	if !lb.Add("[") {
		return false
	}
	var first = true
	for i := range ls {
		if first {
			first = false
		} else {
			lb.Add(", ")
		}
		if !lb.Add(ls[i]) {
			break
		}
	}
	lb.Elems = append(lb.Elems, "]")
	lb.length++
	return lb.stopped
}

// AddMapKeys adds all keys from the string map
func (lb *Builder) AddMapKeys(ls map[string]string) bool {
	if !lb.Add("[") {
		return false
	}
	var first = true
	for k := range ls {
		if first {
			first = false
		} else {
			lb.Add(", ")
		}
		if !lb.Add(k) {
			break
		}
	}
	lb.Elems = append(lb.Elems, "]")
	lb.length++
	return lb.stopped
}

// Join joins all elements into a single string
func (lb Builder) Join() string {
	return strings.Join(lb.Elems, "")
}

// String implements Stringer interface by calling Join method
func (lb Builder) String() string {
	return lb.Join()
}
