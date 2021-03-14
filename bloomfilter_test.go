package bloomfilter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		n          int
		e          float64
		bitsetLen  int
		hashFnsLen int
		err        error
	}{
		{
			n:          100,
			e:          .001,
			bitsetLen:  1438,
			hashFnsLen: 10,
			err:        nil,
		},
		{
			n:          110,
			e:          .001,
			bitsetLen:  1582,
			hashFnsLen: 10,
			err:        nil,
		},
		{
			n:          0,
			e:          .001,
			bitsetLen:  0,
			hashFnsLen: 0,
			err:        ErrCapacity,
		},
		{
			n:          1,
			e:          0,
			bitsetLen:  0,
			hashFnsLen: 0,
			err:        ErrErrorRate,
		},
		{
			n:          1,
			e:          100,
			bitsetLen:  0,
			hashFnsLen: 0,
			err:        ErrErrorRate,
		},
	}

	for _, c := range cases {
		got, err := New(c.n, c.e)

		if err != c.err {
			w := fmt.Sprintf(`New(%v, %v)`, c.n, c.e)
			t.Fatal(w, "expect", c.err, "got", err)
		}

		if got == nil {
			continue
		}

		bn := len(got.Bitset)
		if bn != c.bitsetLen {
			w := fmt.Sprintf(`len(New(%v, %v).Bitset)`, c.n, c.e)
			t.Fatal(w, "expect", c.bitsetLen, "got", bn)
		}

		hn := len(got.HashFns)
		if hn != c.hashFnsLen {
			w := fmt.Sprintf(`len(New(%v, %v).HashFns)`, c.n, c.e)
			t.Fatal(w, "expect", c.hashFnsLen, "got", hn)
		}
	}
}

func TestFindK(t *testing.T) {
	cases := []struct {
		e   float64
		out int
	}{
		{
			e:   0.9,
			out: 1,
		},
		{
			e:   0.75,
			out: 1,
		},
		{
			e:   0.5,
			out: 1,
		},
		{
			e:   0.25,
			out: 2,
		},
		{
			e:   0.1,
			out: 4,
		},
		{
			e:   0.01,
			out: 7,
		},
		{
			e:   0.001,
			out: 10,
		},
	}

	for _, c := range cases {
		got := findK(c.e)

		if got != c.out {
			w := fmt.Sprintf(`findK(%v)`, c.e)
			t.Fatal(w, "expect", c.out, "got", got)
		}
	}
}

func TestFindM(t *testing.T) {
	cases := []struct {
		n   int
		e   float64
		out int
	}{
		{
			n:   100,
			e:   0.9,
			out: 22,
		},
		{
			n:   100,
			e:   0.75,
			out: 60,
		},
		{
			n:   100,
			e:   0.5,
			out: 145,
		},
		{
			n:   100,
			e:   0.25,
			out: 289,
		},
		{
			n:   100,
			e:   0.1,
			out: 480,
		},
		{
			n:   100,
			e:   0.01,
			out: 959,
		},
		{
			n:   100,
			e:   0.001,
			out: 1438,
		},
	}

	for _, c := range cases {
		got := findM(c.n, c.e)

		if got != c.out {
			w := fmt.Sprintf(`findM(%v, %v)`, c.n, c.e)
			t.Fatal(w, "expect", c.out, "got", got)
		}
	}
}

func TestCreateHashFns(t *testing.T) {
	cases := []struct {
		k   int
		out []hash.Hash
	}{
		{k: 1, out: []hash.Hash{hmac.New(sha256.New, []byte("key-0"))}},
		{k: 2, out: []hash.Hash{
			hmac.New(sha256.New, []byte("key-0")),
			hmac.New(sha256.New, []byte("key-1")),
		}},
	}

	for _, c := range cases {
		got := createHashFns(c.k)

		if len(got) != c.k {
			w := fmt.Sprintf(`len(createHashFns(%v))`, c.k)
			t.Fatal(w, "expect", c.k, "got", len(got))
		}

		for i, g := range got {
			b := []byte("test")
			c.out[i].Write(b)
			g.Write(b)

			ch := c.out[i].Sum(nil)
			gh := g.Sum(nil)

			if !bytes.Equal(ch, gh) {
				w := fmt.Sprintf(`createHashFns(%v)[%d]`, c.k, i)
				t.Fatal(w, "\nexpect", ch, "\ngot   ", gh)
			}
		}
	}
}
