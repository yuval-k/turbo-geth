package trie

func CompressWitnessKey(nibbles []byte, previousNibbles []byte) []byte {
	commonPrefixLen := 0
	for i := 0; i < len(previousNibbles) && i < len(nibbles); i++ {
		if nibbles[i] == previousNibbles[i] {
			commonPrefixLen++
		} else {
			break
		}
	}

	nibbles = nibbles[commonPrefixLen:]
	key := keyNibblesToBytes(nibbles)
	if commonPrefixLen < 15 {
		key[0] += byte(commonPrefixLen) << 4
	} else {
		prefix := []byte{key[0], 0}
		prefix[0] += byte(15) << 4
		commonPrefixLen -= 15
		prefix[1] = byte(commonPrefixLen)

		remainder := key[1:]

		key = append(prefix, remainder...)
	}

	return key
}

func UncompressWitnessKey(compressedKey []byte, previousNibbles []byte) []byte {
	keyByte := compressedKey[0]

	commonPrefixLen := int(keyByte >> 4)

	compressedKey[0] = compressedKey[0] & 0b00001111 // erase prefix len

	if commonPrefixLen == 15 {
		commonPrefixLen += int(compressedKey[1])
		prefix := []byte{compressedKey[0]}
		remainder := compressedKey[2:]
		compressedKey = append(prefix, remainder...)
	}

	nibbles := keyBytesToNibbles(compressedKey)
	if commonPrefixLen == 0 {
		return nibbles
	}

	return append(previousNibbles[:commonPrefixLen], nibbles...)
}
