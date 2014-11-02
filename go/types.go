package liblim

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"reflect"
	"time"
)

var (
	ErrInvalidType = errors.New("type is not supported")
	ErrMergeCompat = errors.New("cannot merge incompatible Types")
)

// Vector Clock Style Version
type Version struct {
	// Maps host identifications to version number
	// Should only contain some hosts not all, but always the local version.
	Perspective map[string]uint64
}

type Crdt struct {
	Oid        string  `json:"oid"`
	Version    Version `json:"-"`
	VersionInt uint    `json:"version,int"`
	mtime      time.Time
}

func newCrdt() *Crdt {
	return &Crdt{
		Oid:     uuid.New(),
		Version: Version{Perspective: make(map[string]uint64)},
		mtime:   time.Now(),
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

type Immutable struct {
	Crdt
	Val interface{}
}

func sameKind(l, r interface{}) bool {
	lk := reflect.TypeOf(l).Kind()
	rk := reflect.TypeOf(r).Kind()
	return lk == rk
}
