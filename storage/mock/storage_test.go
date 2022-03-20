package mock

import (
	"context"
	"testing"
)

func TestStorageWrite(t *testing.T) {
	cases := []struct {
		desc       string
		data       string
		writeError bool
		readError  bool
		contents   []byte
		ok         bool
	}{
		{
			desc:       "simple",
			data:       "",
			writeError: false,
			readError:  false,
			contents:   []byte("foo"),
			ok:         true,
		},
		{
			desc:       "write error",
			data:       "",
			writeError: true,
			readError:  false,
			contents:   []byte("foo"),
			ok:         false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewStorage(tc.data, tc.writeError, tc.readError)
			err := s.Write(context.Background(), tc.contents)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err: %s", err)
			}
			if !tc.ok && err == nil {
				t.Fatal("expected to return an error, but no error")
			}

			if tc.ok {
				got := []byte(s.data)
				if err != nil {
					t.Fatalf("failed to read contents: %s", err)
				}
				if string(got) != string(tc.contents) {
					t.Errorf("got: %s, want: %s", string(got), string(tc.contents))
				}
			}
		})
	}
}

func TestStorageRead(t *testing.T) {
	cases := []struct {
		desc       string
		data       string
		writeError bool
		readError  bool
		contents   []byte
		ok         bool
	}{
		{
			desc:       "simple",
			data:       "foo",
			writeError: false,
			readError:  false,
			contents:   []byte("foo"),
			ok:         true,
		},
		{
			desc:       "read error",
			data:       "foo",
			writeError: false,
			readError:  true,
			contents:   nil,
			ok:         false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewStorage(tc.data, tc.writeError, tc.readError)
			got, err := s.Read(context.Background())
			if tc.ok && err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}
			if !tc.ok && err == nil {
				t.Fatal("expected to return an error, but no error")
			}

			if tc.ok {
				if string(got) != string(tc.contents) {
					t.Errorf("got: %s, want: %s", string(got), string(tc.contents))
				}
			}
		})
	}
}
