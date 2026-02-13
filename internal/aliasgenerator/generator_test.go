package aliasgenerator

import (
	"strings"
	"testing"
)

func Test_md5AliasGenerator_GenerateByStr(t *testing.T) {
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

	generator := NewMd5AliasGenerator()

	strToAlias := make(map[string]string, 0)

	for _, tt := range tests {
		t.Run("generate for "+tt.name, func(t *testing.T) {
			strToAlias[tt.input] = generator.GenerateByStr(tt.input)

			if len(strToAlias[tt.input]) != aliasLen {
				t.Errorf("generator output len is %d, want %d", len(strToAlias[tt.input]), aliasLen)
			}

		})
	}

	for _, tt := range tests {
		t.Run("determinism for "+tt.name, func(t *testing.T) {
			first := strToAlias[tt.input]
			second := generator.GenerateByStr(tt.input)

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
	gen := NewMd5AliasGenerator()
	for i := 0; i < b.N; i++ {
		gen.GenerateByStr("https://example.com/very/long/path?with=query")
	}
}
