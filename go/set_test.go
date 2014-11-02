package liblim

import (
	"reflect"
	"testing"
)

func TestSetInt(t *testing.T) {
	c1oid, c1 := NewSetElement(147)
	c2oid, c2 := NewSetElement(258)
	c3oid, c3 := NewSetElement(369)
	c4oid, c4 := NewSetElement(741)
	c5oid, c5 := NewSetElement(852)
	c6oid, c6 := NewSetElement(963)

	s1 := NewSet()
	s1.Add(c1oid, c1)
	s1.Add(c2oid, c2)
	s1.Add(c3oid, c3)
	s1.Add(c4oid, c4)
	s1.Add(c5oid, c5)
	s1.Add(c6oid, c6)
	s1oid1, s1c1 := NewSetElement(123)
	s1oid2, s1c2 := NewSetElement(456)
	s1oid3, s1c3 := NewSetElement(789)
	s1.Add(s1oid1, s1c1)
	s1.Add(s1oid2, s1c2)
	s1.Add(s1oid3, s1c3)

	s2 := NewSet()
	s2.Add(c1oid, c1)
	s2.Add(c2oid, c2)
	s2.Add(c3oid, c3)
	s2.Add(c4oid, c4)
	s2.Add(c5oid, c5)
	s2.Add(c6oid, c6)
	s2oid1, s2c1 := NewSetElement(321)
	s2oid2, s2c2 := NewSetElement(654)
	s2oid3, s2c3 := NewSetElement(987)
	s2.Add(s2oid1, s2c1)
	s2.Add(s2oid2, s2c2)
	s2.Add(s2oid3, s2c3)

	s3 := NewSet()

	s3.Add(c1oid, c1)
	s3.Add(c2oid, c2)
	s3.Add(c3oid, c3)
	s3.Add(c4oid, c4)
	s3.Add(c5oid, c5)
	s3.Add(c6oid, c6)
	s3.Add(s1oid1, s1c1)
	s3.Add(s1oid2, s1c2)
	s3.Add(s1oid3, s1c3)
	s3.Add(s2oid1, s2c1)
	s3.Add(s2oid2, s2c2)
	s3.Add(s2oid3, s2c3)

	s2.Merge(s1)

	num := len(s3.Elements)
	cNum := 0
	for oid, val := range s3.Elements {
		cVal, exists := s2.Elements[oid]
		if !exists {
			t.Log("Values does not exist")
			t.Fail()
		}
		if cVal != val {
			t.Fail()
		}
		cNum = cNum + 1
	}
	if cNum != num {
		t.Log("number of matches does not match")
		t.Fail()
	}

	if num != len(s2.Elements) {
		t.Log("Wrong post merge length")
	}
}

var setTestStrings1 []string = []string{
	"Hello World",
	"Cruel World",
	"Happs World",
	"asdfg",
}

var setTestStrings2 []string = []string{
	"Hello Space",
	"Cruel Space",
	"Happs Space",
	"qwertz",
}

var setTestStrings3 []string = []string{
	"pre set 1",
	"pre set 3",
	"pre set 4",
	"yxcvb",
}

func TestSetString(t *testing.T) {
	refMap := make(map[string]interface{})
	for _, testString := range setTestStrings3 {
		oid, val := NewSetElement(testString)
		refMap[oid] = val
	}

	set1 := NewSet()
	set1Map := make(map[string]interface{})
	for _, testString := range setTestStrings1 {
		oid, val := NewSetElement(testString)
		set1.Add(oid, interface{}(val))
		set1Map[oid] = val
	}
	for oid, val := range refMap {
		set1.Add(oid, val)
	}

	set2 := NewSet()
	set2Map := make(map[string]interface{})
	for _, testString := range setTestStrings2 {
		oid, val := NewSetElement(testString)
		set2.Add(oid, val)
		set1Map[oid] = val
	}
	for oid, val := range refMap {
		set2.Add(oid, val)
	}

	for oid, val := range set1Map {
		refMap[oid] = val
	}
	for oid, val := range set2Map {
		refMap[oid] = val
	}

	err := set2.Merge(set1)
	if err != nil {
		t.Log("Merge Failed")
		t.Fail()
	}
	if !reflect.DeepEqual(refMap, set2.Elements) {
		t.Log("Maps are not equal")
		t.Fail()
	}

}
