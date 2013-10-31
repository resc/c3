package c3

import "sort"

type Sorter struct {
	List
	Lesser
}

func (s *Sorter) Less(i, j int) bool {
	iv, _ := s.Get(i)
	jv, _ := s.Get(j)
	return s.Lesser(iv, jv)
}

func (s *Sorter) Sort() {
	sort.Sort(s)
}
