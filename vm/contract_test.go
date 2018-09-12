package vm

import (
	"bytes"
	"errors"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func SetTestLogContext() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(strings.Replace(dir, "\\", "/", -1), "runlog")
	if err := os.MkdirAll(dir, 0777); err != nil {
		return
	}
	dir = filepath.Join(dir, "test.log")
	log15.Info(dir)
	log15.Root().SetHandler(
		log15.LvlFilterHandler(log15.LvlInfo, log15.Must.FileHandler(dir, log15.TerminalFormat())),
	)
}

func TestRun(t *testing.T) {
	SetTestLogContext()
	tests := []struct {
		input, result          []byte
		err                    error
		quotaLeft, quotaRefund uint64
		summary                string
	}{
		{[]byte{byte(PUSH1), 1, byte(PUSH1), 2, byte(ADD), byte(PUSH1), 32, byte(DUP1), byte(SWAP2), byte(SWAP1), byte(MSTORE), byte(PUSH1), 32, byte(SWAP1), byte(RETURN)}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}, nil, 999964, 0, "return 1+2"},
		{[]byte{0xFA}, []byte{}, errors.New(""), 1000000, 0, "invalid opcode"},
		{[]byte{byte(PUSH1), 32, byte(PUSH17), 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, byte(MLOAD)}, []byte{}, errGasUintOverflow, 999994, 0, "memory size overflow"},
		{[]byte{byte(POP)}, []byte{}, errors.New(""), 1000000, 0, "make stack error"},
		{[]byte{byte(PUSH1), 32, byte(JUMP)}, []byte{}, errors.New(""), 999989, 0, "execution error"},
		{[]byte{byte(CALLVALUE), byte(DUP1), byte(ISZERO), byte(PUSH2), 0, 11, byte(JUMPI), byte(PUSH1), 0, byte(DUP1), byte(REVERT), byte(JUMPDEST), byte(PUSH1), 32, byte(PUSH1), 0, byte(DUP2), byte(DUP2), byte(MSTORE), byte(RETURN)}, []byte{}, ErrExecutionReverted, 999973, 0, "execution revert"},
		{[]byte{byte(CALLVALUE), byte(DUP1), byte(ISZERO), byte(NOT), byte(PUSH2), 0, 12, byte(JUMPI), byte(PUSH1), 0, byte(DUP1), byte(REVERT), byte(JUMPDEST), byte(PUSH1), 32, byte(PUSH1), 0, byte(DUP2), byte(DUP2), byte(MSTORE), byte(RETURN)}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32}, nil, 999957, 0, "jumpi"},
	}
	for _, test := range tests {
		vm := &VM{Db: NewNoDatabase(), instructionSet: simpleInstructionSet}
		vm.Debug = true
		receiveCallBlock := &ledger.AccountBlock{AccountAddress: types.Address{}, ToAddress: types.Address{}, BlockType: ledger.BlockTypeReceive}
		receiveCallBlock.Amount = big.NewInt(10)
		c := newContract(receiveCallBlock.AccountAddress, receiveCallBlock.ToAddress, receiveCallBlock, 1000000, 0)
		c.setCallCode(types.Address{}, test.input)
		ret, err := c.run(vm)
		if bytes.Compare(ret, test.result) != 0 || c.quotaLeft != test.quotaLeft || c.quotaRefund != test.quotaRefund || (err == nil && test.err != nil) || (err != nil && test.err == nil) {
			t.Fatalf("contract run failed, summary: %v", test.summary)
		}
	}
}
