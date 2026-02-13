package aliasgenerator

import (
	"math"
	"testing"
)

func Test_generateCharset(t *testing.T) {
	expected := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
	got := generateCharset()

	if len(got) != 63 {
		t.Errorf("charset length = %d, want 63", len(got))
	}

	charset := make(map[byte]int)

	for _, c := range expected {
		charset[c] = 0
	}

	for i, c := range got {
		counter, found := charset[c]
		if !found {
			t.Errorf("character %c is not found in desired charset", c)
		}
		if counter >= 1 {
			t.Errorf("duplicate character found: %c at %d", c, i)
		}
		charset[c] = counter + 1
	}

	for k, v := range charset {
		if v == 0 {
			t.Errorf("character %c is not in the charset", k)
		}
	}

}

func Test_encodeToCharset_length(t *testing.T) {
	tests := []uint64{0, 67, uint64(math.MaxUint64)}

	for _, n := range tests {
		got := encodeToCharset(n)

		if len(got) != aliasLen {
			t.Errorf("encodeToCharset(%d) returned length %d, want %d", n, len(got), aliasLen)
		}
	}
}

func Test_encodeToCharset_determinism(t *testing.T) {
	tests := []uint64{0, 67, uint64(math.MaxUint64)}

	for _, n := range tests {
		first := encodeToCharset(n)
		second := encodeToCharset(n)

		if string(first) != string(second) {
			t.Errorf("encodeToCharset(%d) not deterministic: '%s' vs '%s'", n, string(first), string(second))
		}
	}
}

func Test_generateCharset_uniqueness(t *testing.T) {
	seen := make(map[string]uint64)

	for n := range uint64(8192) {
		result := string(encodeToCharset(n))

		if prev, exists := seen[result]; exists {
			t.Errorf("Collision: %d and %d both encode to '%s'", prev, n, result)
		}
		seen[result] = n
	}

}
