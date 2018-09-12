package vm

import (
	"bytes"
	"github.com/vitelabs/go-vite/common/types"
	"math/big"
	"testing"
)

func TestCodeAnalysis(t *testing.T) {
	tests := []struct {
		code   []byte
		result []byte
	}{
		{[]byte{byte(PUSH1), 0x01, 0x01, 0x01}, []byte{0x40, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH1), byte(PUSH1), byte(PUSH1), byte(PUSH1)}, []byte{0x50, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), 0x01, 0x01, 0x01}, []byte{0x7f, 0x80, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH8), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x7f, 0x80, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x01, 0x01, 0x01, 0x01, 0x01, byte(PUSH2), byte(PUSH2), byte(PUSH2), 0x01, 0x01, 0x01}, []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x01, 0x01, 0x01, 0x01, 0x01, byte(PUSH2), 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH3), 0x01, 0x01, 0x01, byte(PUSH1), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x74, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH3), 0x01, 0x01, 0x01, byte(PUSH1), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x74, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x01, byte(PUSH8), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x3f, 0xc0, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x01, byte(PUSH8), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x3f, 0xc0, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH16), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x7f, 0xff, 0x80, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH16), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x7f, 0xff, 0x80, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH16), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x7f, 0xff, 0x80, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH8), 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, byte(PUSH1), 0x01}, []byte{0x7f, 0xa0, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH8), 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, byte(PUSH1), 0x01}, []byte{0x7f, 0xa0, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{byte(PUSH32)}, []byte{0x7f, 0xff, 0xff, 0xff, 0x80}},
		{[]byte{byte(PUSH32)}, []byte{0x7f, 0xff, 0xff, 0xff, 0x80}},
		{[]byte{byte(PUSH32)}, []byte{0x7f, 0xff, 0xff, 0xff, 0x80}},
	}
	for _, test := range tests {
		ret := codeBitmap(test.code)
		if !bytes.Equal(test.result, ret) {
			t.Fatalf("analysis fail, got %v, expected %v", ret, test.result)
		}
	}
}

func TestHas(t *testing.T) {
	tests := []struct {
		code   []byte
		dest   *big.Int
		result bool
	}{
		{[]byte{byte(PUSH1), byte(JUMPDEST)}, big.NewInt(1), false},
		{[]byte{byte(PUSH1), 0, byte(JUMPDEST)}, big.NewInt(2), true},
		{[]byte{byte(PUSH1), 0, byte(JUMPDEST)}, big.NewInt(1), false},
		{[]byte{byte(PUSH32), 0, byte(JUMPDEST)}, big.NewInt(2), false},
	}
	for _, test := range tests {
		d := make(destinations)
		result := d.has(types.Address{}, test.code, test.dest)
		if result != test.result {
			t.Fatalf("analysis result error, code: [%v], dest: %v, expected: %v, got: %v", test.code, test.dest, test.result, result)
		}
	}
}
