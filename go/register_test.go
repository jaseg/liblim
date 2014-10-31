package liblim

import (
	"bytes"
	"testing"
	"time"
)

// Give back two registers that are created after another. I promise!
func twoNewRegisters(l, r interface{}) (*Register, *Register) {
	rl := NewRegister(l)
	time.Sleep(500)
	return rl, NewRegister(r)
}

func TestRegisterMergeInt(t *testing.T) {
	l, r := twoNewRegisters(123456, 654321)
	l.Merge(r)
	if l.Val.(int) != 654321 {
		t.Fail()
	}
}

func TestRegisterMergeFloat(t *testing.T) {
	l, r := twoNewRegisters(float32(123456.654321), float32(654321.123456))
	l.Merge(r)
	if l.Val.(float32) != 654321.123456 {
		t.Fail()
	}
	l, r = twoNewRegisters(123456.654321, 654321.123456)
	l.Merge(r)
	if l.Val.(float64) != 654321.123456 {
		t.Fail()
	}
}

func TestRegisterMergeString(t *testing.T) {
	l, r := twoNewRegisters("Welt Hallo", "Hallo Welt")
	l.Merge(r)
	if l.Val.(string) != "Hallo Welt" {
		t.Fail()
	}
}

func TestRegisterMergeByteSlice(t *testing.T) {
	l, r := twoNewRegisters([]byte("Welt Hallo"), []byte("Hallo Welt"))
	l.Merge(r)
	if !bytes.Equal(l.Val.([]byte), []byte("Hallo Welt")) {
		t.Fail()
	}
}
