package liblim

import (
	"code.google.com/p/go-uuid/uuid"
	"time"
)

// Vector Clock Style Version
type Version struct {
	// Maps host identifications to version number
	// Should only contain some hosts not all, but always the local version.
	Perspective map[string]uint64
}

type Crdt struct {
	Oid        string  `json:"oid,string"`
	Version    Version `json:"-"`
	VersionInt uint    `json:"version,int"`
	timestamp  time.Time
}

func newCrdt() *Crdt {
	return &Crdt{
		Oid:       uuid.NewRandom().String(),
		Version:   Version{Perspective: make(map[string]uint64)},
		timestamp: time.Now(),
	}
}

type Set struct {
	Crdt
	Elements map[string]interface{}
}

func (s *Set) Add(oid string, val interface{}) {
	s.Elements[oid] = val
}

func (s *Set) Get(oid string) interface{} {
	val, ok := s.Elements[oid]
	if ok {
		return val
	}
	return nil
}

type Register struct {
	*Crdt
	Val interface{}
}

func NewRegister(val interface{}) *Register {
	c := newCrdt()
	return &Register{
		Crdt: c,
		Val:  val,
	}
}

func (r *Register) Get() interface{} {
	return r.Val
}

type Immutable struct {
	Crdt
	Val interface{}
}
