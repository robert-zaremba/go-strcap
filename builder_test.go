package strcap

import (
	. "github.com/scale-it/checkers"
	. "gopkg.in/check.v1"
)

type builderSuite struct{}

func (s *builderSuite) TestAdd(c *C) {
	lb := Builder{Limit: 7}

	c.Check(lb.Join(), Equals, "")
	lb.Add("1", "2 3")
	c.Check(lb.Join(), Equals, "12 3")
	lb.Add("will break")
	c.Assert(lb.Join(), Equals, "12 3wil...")
	lb.Add("another try")
	c.Assert(lb.Join(), Equals, "12 3wil...")
	lb.Add("another", "another try")
	c.Assert(lb.Join(), Equals, "12 3wil...")

	lb = Builder{Limit: 3}
	lb.Add("123")
	c.Check(lb.Join(), Equals, "123")
}

func (s *builderSuite) TestAddAny(c *C) {
	lb := Builder{Limit: 4}
	lb.AddAny(12, 3)
	c.Check(lb.Join(), Equals, "123")
	lb.AddAny(456)
	c.Check(lb.Join(), Equals, "1234...")

	lb = Builder{Limit: 4}
	lb.AddAny(123456, 78)
	c.Check(lb.Join(), Equals, "1234...")
}

func (s *builderSuite) checkStrSlice(limit int, base string, slice []string, expected string, c *C) {
	lb := Builder{Limit: limit}
	lb.Add(base)
	lb.AddSlice(slice)
	c.Check(lb.Join(), Equals, expected)
}

func (s *builderSuite) TestAddSlice(c *C) {
	empty := []string{}
	s.checkStrSlice(1, "", empty, "[]", c)
	s.checkStrSlice(1, "a", empty, "a...", c)
	s.checkStrSlice(2, "a", empty, "a[]", c)
	s.checkStrSlice(3, "a", empty, "a[]", c)
	s.checkStrSlice(5, "a", empty, "a[]", c)

	ls := []string{"1b"}
	s.checkStrSlice(2, "a", ls, "a[...]", c)
	s.checkStrSlice(3, "a", ls, "a[1...]", c)
	s.checkStrSlice(4, "a", ls, "a[1b]", c)
	s.checkStrSlice(4, "a", []string{"1bc"}, "a[1b...]", c)

	ls = []string{"1b", "2b"}
	s.checkStrSlice(4, "a", ls, "a[1b...]", c)
	s.checkStrSlice(5, "a", ls, "a[1b,...]", c)
	s.checkStrSlice(6, "a", ls, "a[1b, ...]", c)
	s.checkStrSlice(7, "a", ls, "a[1b, 2...]", c)
	s.checkStrSlice(10, "a", ls, "a[1b, 2b]", c)

	lb := Builder{Limit: 13}
	lb.AddSlice(ls)
	lb.Add("12", " 34", "5")
	c.Check(lb.Join(), Equals, "[1b, 2b]12 34...")
}

func (s *builderSuite) checkStrMap(limit int, base string, m map[string]string, c *C, expected ...string) {
	lb := Builder{Limit: limit}
	lb.Add(base)
	lb.AddMapKeys(m)
	c.Check(lb.Join(), IsIn, expected)
}

func (s *builderSuite) TestAddMap(c *C) {
	empty := map[string]string{}
	s.checkStrMap(1, "", empty, c, "[]")
	s.checkStrMap(1, "a", empty, c, "a...")
	s.checkStrMap(2, "a", empty, c, "a[]")
	s.checkStrMap(3, "a", empty, c, "a[]")
	s.checkStrMap(5, "a", empty, c, "a[]")

	ls := map[string]string{"1b": "not important"}
	s.checkStrMap(2, "a", ls, c, "a[...]")
	s.checkStrMap(3, "a", ls, c, "a[1...]")
	s.checkStrMap(4, "a", ls, c, "a[1b]")
	ls = map[string]string{"1bc": "not important"}
	s.checkStrMap(4, "a", ls, c, "a[1b...]")

	ls = map[string]string{"1b": "not important", "2b": "not important"}
	s.checkStrMap(4, "a", ls, c, "a[1b...]", "a[2b...]")
	s.checkStrMap(5, "a", ls, c, "a[1b,...]", "a[2b,...]")
	s.checkStrMap(6, "a", ls, c, "a[1b, ...]", "a[2b, ...]")
	s.checkStrMap(7, "a", ls, c, "a[1b, 2...]", "a[2b, 1...]")
	s.checkStrMap(10, "a", ls, c, "a[1b, 2b]", "a[2b, 1b]")

	lb := Builder{Limit: 13}
	lb.AddMapKeys(ls)
	lb.Add("12", " 34", "5")
	c.Check(lb.Join(), IsIn, []string{"[1b, 2b]12 34...", "[2b, 1b]12 34..."})
}
