package liblim

import (
	"time"
)

type Register struct {
	*Crdt
	Value interface{} `json:"val"`
	Mtime int64       `json:"mtime,int"`
}

type LwwRegister interface {
	Val() interface{}
	SetVal(interface{})
	Compare(interface{}) bool
	Merge(remote LwwRegister) error
}

func NewRegister(val interface{}) *Register {
	c := newCrdt()
	return &Register{
		Crdt:  c,
		Value: val,
		Mtime: time.Now().UnixNano(),
	}
}

func (r *Register) Val() interface{} {
	return r.Value
}

func (r *Register) SetVal(val interface{}) {
	r.Mtime = time.Now().UnixNano()
	r.Value = val
}

func (r *Register) Compare(remote *Register) bool {
	return r.Mtime <= remote.Mtime
}

func (local *Register) Merge(remote *Register) error {
	if !sameKind(local.Val, remote.Val) {
		return ErrMergeCompat
	}
	if local.Compare(remote) {
		local.SetVal(remote.Val())
	}
	return nil
}
