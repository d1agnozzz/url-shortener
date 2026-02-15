package urlsanitizer

import "testing"

func Test_urlSanitizer_Sanitize(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		raw     string
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			raw:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "not a url",
			raw:     "some text which is not valid url",
			want:    "",
			wantErr: true,
		},
		{
			name:    "simple url",
			raw:     "https://youtube.com/",
			want:    "https://youtube.com/",
			wantErr: false,
		},
		{
			name:    "plain path",
			raw:     "/home/foo/bar",
			want:    "",
			wantErr: true,
		},
		{
			name:    "corrupted scheme",
			raw:     "https:/youtube.com",
			want:    "",
			wantErr: true,
		},
		{
			name:    "leading and trailing whitespaces",
			raw:     "  \t  https://youtube.com/  \t   \n",
			want:    "https://youtube.com/",
			wantErr: false,
		},
		{
			name:    "path to directory",
			raw:     "https://example.com/foobar/",
			want:    "https://example.com/foobar",
			wantErr: false,
		},
		{
			name:    "path to file",
			raw:     "https://example.com/foobar",
			want:    "https://example.com/foobar",
			wantErr: false,
		},
		{
			name:    "many trailing slashes",
			raw:     "https://example.com/foobar//////////",
			want:    "https://example.com/foobar",
			wantErr: false,
		},
		{
			name:    "with query params",
			raw:     "https://example.com/foo/bar?key=value&thing=another",
			want:    "https://example.com/foo/bar?key=value&thing=another",
			wantErr: false,
		},
		{
			name:    "query params sorting",
			raw:     "https://example.com/foo/bar?b=true&a=false",
			want:    "https://example.com/foo/bar?a=false&b=true",
			wantErr: false,
		},
		{
			name:    "case sensitiveness",
			raw:     "HTTPS://EXAMPLE.COM/Foo/BAR?paramA=True&paramb=false",
			want:    "https://example.com/Foo/BAR?paramA=True&paramb=false",
			wantErr: false,
		},
		{
			name:    "stripping user info",
			raw:     "https://user:pass@example.com/foo",
			want:    "https://example.com/foo",
			wantErr: false,
		},
		{
			name:    "default schema",
			raw:     "example.com",
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "opaque url",
			raw:     "mailto:a@b.com",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUrlSanitizer()
			got, gotErr := s.Sanitize(tt.raw)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Sanitize() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Sanitize() succeeded unexpectedly")
			}
			if got.String() != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
