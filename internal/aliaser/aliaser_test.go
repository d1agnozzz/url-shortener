package aliaser

import (
	"strings"
	"testing"
)

func Test_md5Aliaser(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"simple string", "test case string"},
		{"different string", "different test case string"},
		{"with special chars", "hello!@#$%^&*()"},
		{"unicode", "Hello 世界"},
		{"long string", strings.Repeat("a", 1000)},
		{"similar 1", "hello"},
		{"similar 2", "helloo"},
		{"similar 3", "hellO"},
	}

	aliasgen := NewMd5Aliaser()

	strToAlias := make(map[string]string, 0)

	for _, tt := range tests {
		t.Run("generate for "+tt.name, func(t *testing.T) {
			strToAlias[tt.input] = aliasgen.GenerateByStr(tt.input).String()

			if len(strToAlias[tt.input]) != AliasLen {
				t.Errorf("generator output len is %d, want %d", len(strToAlias[tt.input]), AliasLen)
			}

		})
	}

	for _, tt := range tests {
		t.Run("determinism for "+tt.name, func(t *testing.T) {
			first := strToAlias[tt.input]
			second := aliasgen.GenerateByStr(tt.input).String()

			if first != second {
				t.Errorf("generator is not deterministic: given '%s' outputs '%s' and '%s'", tt.input, first, second)
			}
		})
	}

	t.Run("uniqueness", func(t *testing.T) {
		aliasToStr := make(map[string]string)

		for input, alias := range strToAlias {
			existing, found := aliasToStr[alias]

			if found {
				t.Errorf("generator collision: given '%s' and '%s' returns the same '%s'", input, existing, alias)
			}

			aliasToStr[alias] = input
		}
	})

}

func BenchmarkGenerateByStr(b *testing.B) {
	gen := NewMd5Aliaser()
	for b.Loop() {
		gen.GenerateByStr("https://example.com/very/long/path?with=query")
	}
}
