package liblim

import (
	"time"
)

type Register struct {
	*Crdt
	Val   interface{} `json:"val"`
	Mtime int64       `json:"mtime,int"`
}

func NewRegister(val interface{}) *Register {
	c := newCrdt()
	return &Register{
		Crdt:  c,
		Val:   val,
		Mtime: time.Now().UnixNano(),
	}
}

func (r *Register) Get() interface{} {
	return r.Val
}

func (local *Register) Merge(remote *Register) error {
	if !sameKind(local.Val, remote.Val) {
		return ErrMergeCompat
	}
	if local.Mtime < remote.Mtime {
		local.Val = remote.Val
	} else if local.Mtime == remote.Mtime {
		panic("Mtimes are the same! Send a Postcard to Jaseg!")
	}
	return nil
}
