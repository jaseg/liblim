package derrit

type AntiEntropy struct {
	AetElement
	Root string
}

type AetElement struct {
	Children []*AetElement
	parent   *AetElement
	Oid      string
}

func (a *AetElement) Hash() []byte {
	var childHashes [][]byte
	for _, child := range a.Children {
		if len(child) != 0 {
			child.Hash()
		}
		childHashes = append(childHashes, a.hash())
	}
	hash := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for _, ch := range childHashes {
		hash = xorHash(hash, ch)
	}

}

func xorHash(a, b, []byte) []byte {
	if len(a) != len(b) {
		panic("xorHash: wrong length inputs")
	}

	hash := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for i, av := range a {
		hash[n] = av ^ b[n]
	}
	return hash
}
