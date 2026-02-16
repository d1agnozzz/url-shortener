package aliaser

func generateCharset() []byte {
	b := make([]byte, 0, 63)

	for i := byte('0'); i <= '9'; i++ {
		b = append(b, i)
	}

	for i := byte('a'); i <= 'z'; i++ {
		b = append(b, i)
	}

	for i := byte('A'); i <= 'Z'; i++ {
		b = append(b, i)
	}

	b = append(b, '_')

	return b
}

func encodeToCharset(n uint64) []byte {
	encoded := make([]byte, AliasLen)
	base := len(charset)

	for i := range encoded {
		idx := n % uint64(base)
		encoded[i] = charset[idx]
		n /= uint64(base)
	}

	return encoded
}
