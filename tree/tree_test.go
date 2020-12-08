package tree

import (
	"bytes"
	"strings"
	"testing"
)

const testDir = "../testdata"

func TestDeep(t *testing.T) {
	tests := []struct {
		name      string
		targetDir string
		isAll     bool
		deepLevel int
		buf       *bytes.Buffer

		want string
	}{
		{
			name:      "success 1",
			targetDir: testDir,
			isAll:     true,
			deepLevel: 1,
			buf:       &bytes.Buffer{},

			want: genTestTree(`
├── a
└── .a
`),
		},
		{
			name:      "success 2",
			targetDir: testDir,
			isAll:     true,
			deepLevel: 2,
			buf:       &bytes.Buffer{},

			want: genTestTree(`
├── a
│   ├── b
│   └── .z
└── .a
`),
		},
		{
			name:      "success 3",
			targetDir: testDir,
			isAll:     true,
			deepLevel: 3,
			buf:       &bytes.Buffer{},

			want: genTestTree(`
├── a
│   ├── b
│   │   ├── c
│   │   └── x
│   └── .z
└── .a
`),
		},
		{
			name:      "success 4",
			targetDir: testDir,
			isAll:     false,
			deepLevel: 3,
			buf:       &bytes.Buffer{},

			want: genTestTree(`
└── a
    └── b
        ├── c
        └── x
`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			deep(tt.targetDir, tt.isAll, tt.deepLevel, 0, tt.buf)
			got := tt.buf.String()
			if got != tt.want {
				t.Errorf("\nwant: \n%s\ngot: \n%s\n", tt.want, got)
			}
		})
	}
}

func genTestTree(t string) string {
	return strings.TrimPrefix(t, "\n")
}
