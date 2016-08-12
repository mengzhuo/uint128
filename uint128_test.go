package uint128

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

var (
	target = map[string]struct{ H, L uint64 }{
		"456":               {0x0, 0x456},
		"10000000000000456": {0x1, 0x456},
		"e0000000000000009": {0xe, 0x9},
		"10000000000000000": {0x1, 0x0},
	}
)

func TestEncode(t *testing.T) {

	for k, v := range target {
		u := new(Uint128)
		u.H = v.H
		u.L = v.L
		if u.HexString() != k {
			z := []rune(u.HexString())
			for i, r := range k {
				if r != z[i] {
					t.Error("missmatch expect\n", k, "\n", u.HexString(), "\n",
						strings.Repeat(" ", i), "^")
					break
				}
			}
		}
	}
}

func TestLoadFromByte(t *testing.T) {
	for k, v := range target {
		u := &Uint128{}
		b, _ := hex.DecodeString(fmt.Sprintf("%032s", k))
		err := binary.Read(bytes.NewReader(b),
			binary.BigEndian, u)

		if err != nil || u.H != v.H || u.L != v.L {
			i := new(big.Int)
			i.SetString(fmt.Sprintf("0x%032s", k), 0)
			t.Error("missmatch ", fmt.Sprintf("%032s", k), i.Text(16), u, err)
		}
	}
}

func TestXor(t *testing.T) {
	xor := []struct{ s, x, cmp string }{
		{"1", "1", "0"},
		{"2", "1", "3"},
		{"e0000000000000009", "f0000000000000000", "10000000000000009"},
	}
	for _, entry := range xor {
		s, err := NewFromString(entry.s)
		if err != nil {
			t.Error(err)
		}
		x, _ := NewFromString(entry.x)
		if err != nil {
			t.Error(err)
		}
		cmp, _ := NewFromString(entry.cmp)
		if err != nil {
			t.Error(err)
		}

		t.Log(s, x, cmp)

		s.Xor(x)
		if s.Compare(cmp) != 0 {
			t.Error("failed xor at", entry, s, x, cmp)
		}
	}
}
