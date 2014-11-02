package liblim

import (
	"code.google.com/p/go-uuid/uuid"
)

func NewSetElement(val interface{}) (string, interface{}) {
	return uuid.New(), val
}

type GrowSet interface {
	Add(oid string, val interface{})
	Lookup(string) bool
	Compare(remote *GrowSet) bool
	Merge(remote *GrowSet) error
}

// Default Grow only set.
type Set struct {
	*Crdt
	// Contains all elements with indexed by uuid.
	Elements map[string]interface{} `json:"elements"`
}

func NewSet() *Set {
	c := newCrdt()
	return &Set{
		Crdt:     c,
		Elements: make(map[string]interface{}),
	}
}

func (s *Set) Add(oid string, val interface{}) {
	s.Elements[oid] = val
}

func (s *Set) Lookup(oid string) bool {
	_, ok := s.Elements[oid]
	return ok
}

func (s *Set) Element(oid string) interface{} {
	val, ok := s.Elements[oid]
	if ok {
		return val
	}
	return nil
}

func (s *Set) Merge(remote *Set) error {
	for oid, element := range remote.Elements {
		if !s.Lookup(oid) {
			s.Add(oid, element)
		}
	}
	return nil
}
