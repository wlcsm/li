package ansi

import (
	"bytes"
	"io"
	"testing"
)

func TestDecoder(t *testing.T) {
	for _, test := range []struct {
		in       []byte
		expected []Key
	}{
		{
			in:       []byte("hi"),
			expected: []Key{Key('h'), Key('i')},
		},
		{
			in:       []byte("\x1b[A"),
			expected: []Key{UpArrowKey},
		},
		{
			in:       []byte("\x1b[B"),
			expected: []Key{DownArrowKey},
		},
		{
			in:       []byte("\x1b[C"),
			expected: []Key{RightArrowKey},
		},
		{
			in:       []byte("\x1b[D"),
			expected: []Key{LeftArrowKey},
		},
		{
			in:       []byte("\x1b[F"),
			expected: []Key{EndKey},
		},
		{
			in:       []byte("\x1b[H"),
			expected: []Key{HomeKey},
		},
		{
			in:       []byte("\x1b[1~"),
			expected: []Key{HomeKey},
		},
		{
			in:       []byte("\x1b[7~"),
			expected: []Key{HomeKey},
		},
		{
			in:       []byte("\x1b[3~"),
			expected: []Key{DeleteKey},
		},
		{
			in:       []byte("\x1b[4~"),
			expected: []Key{EndKey},
		},
		{
			in:       []byte("\x1b[8~"),
			expected: []Key{EndKey},
		},
		{
			in:       []byte("\x1b[5~"),
			expected: []Key{PageUpKey},
		},
		{
			in:       []byte("\x1b[6~"),
			expected: []Key{PageDownKey},
		},
		{
			in:       []byte("\x1b[6~ab"),
			expected: []Key{PageDownKey, Key('a'), Key('b')},
		},
		{
			in:       []byte("\x1b[6~\x1b[B"),
			expected: []Key{PageDownKey, DownArrowKey},
		},
		{
			in:       []byte("\x1b[B\x1b[6~"),
			expected: []Key{DownArrowKey, PageDownKey},
		},
	} {
		d := NewDecoder(bytes.NewReader(test.in))
		for i := 0; i < len(test.expected); i++ {
			k, err := d.Decode()
			if err != nil {
				t.Fatal(err)
			}

			if exp := test.expected[i]; k != exp {
				t.Fatalf("index=%d expected=%d got=%d", i, exp, k)
			}
		}

		_, err := d.Decode()
		if err != io.EOF {
			t.Fatalf("at end of decoding, expected io.EOF, got: %v", err)
		}
	}
}
