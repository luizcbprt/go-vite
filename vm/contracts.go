package vm

import (
	"bytes"
	"github.com/vitelabs/go-vite/common/types"
	"math/big"
)

var (
	AddressRegister, _       = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	AddressVote, _           = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2})
	AddressMortgage, _       = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3})
	AddressConsensusGroup, _ = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4})
)

type precompiledContract interface {
	doSend(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error)
	doReceive(vm *VM, block VmAccountBlock) error
}

var simpleContracts = map[types.Address]precompiledContract{
	AddressRegister:       &register{},
	AddressVote:           &vote{},
	AddressMortgage:       &mortgage{},
	AddressConsensusGroup: &consensusGroup{},
}

func getPrecompiledContract(address types.Address) (precompiledContract, bool) {
	p, ok := simpleContracts[address]
	return p, ok
}

type register struct{}

var (
	DataRegister       = byte(1)
	DataCancelRegister = byte(2)
	DataReward         = byte(3)
)

func (p *register) doSend(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if len(block.Data()) == 11 && block.Data()[0] == DataRegister {
		return p.doSendRegister(vm, block, quotaLeft)
	} else if len(block.Data()) == 11 && block.Data()[0] == DataCancelRegister {
		return p.doSendCancelRegister(vm, block, quotaLeft)
	} else if len(block.Data()) >= 11 && block.Data()[0] == DataReward {
		return p.doSendReward(vm, block, quotaLeft)
	}
	return quotaLeft, ErrInvalidData
}

func (p *register) doSendRegister(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Cmp(big1e24) != 0 ||
		bytes.Equal(block.TokenId().Bytes(), viteTokenTypeId.Bytes()) ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	gid, _ := BytesToGid(block.Data()[1:11])
	if !vm.Db.IsExistGid(gid) {
		return quotaLeft, ErrInvalidData
	}
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) >= 72 && !allZero(old[0:32]) {
		return quotaLeft, ErrInvalidData
	}
	quotaLeft, err := useQuota(quotaLeft, registerGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}

func (p *register) doSendCancelRegister(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Sign() != 0 ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	gid, _ := BytesToGid(block.Data()[1:11])
	if !vm.Db.IsExistGid(gid) {
		return quotaLeft, ErrInvalidData
	}
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) < 72 || allZero(old[0:32]) ||
		vm.Db.SnapshotBlock(block.SnapshotHash()).Timestamp()-new(big.Int).SetBytes(old[32:40]).Int64() < registerLockTime {
		return quotaLeft, ErrInvalidData
	}

	quotaLeft, err := useQuota(quotaLeft, cancelRegisterGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}

func (p *register) doSendReward(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Sign() != 0 ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	if !bytes.Equal(block.Data()[1:11], snapshotGid.Bytes()) {
		return quotaLeft, ErrInvalidData
	}
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) < 72 {
		return quotaLeft, ErrInvalidData
	}
	// newRewardHeight := min(userDefined, currentSnapshotHeight-50, cancelSnapshotHeight)
	newRewardHeight := new(big.Int)
	if len(block.Data()) >= 43 {
		newRewardHeight.SetBytes(block.Data()[11:43])
	} else {
		newRewardHeight = BigMin(newRewardHeight.Sub(vm.Db.SnapshotBlock(block.SnapshotHash()).Height(), rewardHeightLimit), new(big.Int).SetBytes(old[40:72]))
	}
	if len(old) >= 104 && !allZero(old[72:104]) {
		newRewardHeight = BigMin(newRewardHeight, new(big.Int).SetBytes(old[72:104]))
	}
	oldRewardHeight := new(big.Int).SetBytes(old[40:72])
	if newRewardHeight.Cmp(oldRewardHeight) <= 0 {
		return quotaLeft, ErrInvalidData
	}
	heightGap := new(big.Int).Sub(newRewardHeight, oldRewardHeight)
	if heightGap.Cmp(rewardGapLimit) > 0 {
		return quotaLeft, ErrInvalidData
	}

	count := heightGap.Uint64()
	quotaLeft, err := useQuota(quotaLeft, rewardGas+count*calcRewardGasPerBlock)
	if err != nil {
		return quotaLeft, err
	}

	reward := calcReward(vm, block.AccountAddress().Bytes(), oldRewardHeight, count)
	block.SetData(joinBytes(block.Data()[0:11], leftPadBytes(newRewardHeight.Bytes(), 32), leftPadBytes(oldRewardHeight.Bytes(), 32), leftPadBytes(reward.Bytes(), 32)))
	return quotaLeft, nil
}

func calcReward(vm *VM, producer []byte, startHeight *big.Int, count uint64) *big.Int {
	var rewardCount uint64
	for count >= 0 {
		var list []VmSnapshotBlock
		if count < dbPageSize {
			list = vm.Db.SnapshotBlockList(startHeight, count, true)
			count = 0
		} else {
			list = vm.Db.SnapshotBlockList(startHeight, dbPageSize, true)
			count = count - dbPageSize
		}
		for _, block := range list {
			if bytes.Equal(block.Producer().Bytes(), producer) {
			}
			rewardCount++
		}
	}
	return new(big.Int).Mul(rewardPerBlock, new(big.Int).SetUint64(rewardCount))
}

func (p *register) doReceive(vm *VM, block VmAccountBlock) error {
	if len(block.Data()) == 11 && block.Data()[0] == DataRegister {
		return p.doReceiveRegister(vm, block)
	} else if len(block.Data()) == 11 && block.Data()[0] == DataCancelRegister {
		return p.doReceiveCancelRegister(vm, block)
	} else if len(block.Data()) == 107 && block.Data()[0] == DataReward {
		return p.doReceiveReward(vm, block)
	}
	return ErrInvalidData
}
func (p *register) doReceiveRegister(vm *VM, block VmAccountBlock) error {
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) >= 72 && !allZero(old[0:32]) {
		return ErrInvalidData
	}
	snapshotBlock := vm.Db.SnapshotBlock(block.SnapshotHash())
	rewardHeight := leftPadBytes(snapshotBlock.Height().Bytes(), 32)
	if len(old) >= 72 && !allZero(old[40:72]) {
		rewardHeight = old[40:72]
	}
	registerInfo := joinBytes(leftPadBytes(block.Amount().Bytes(), 32),
		leftPadBytes(new(big.Int).SetInt64(snapshotBlock.Timestamp()).Bytes(), 32),
		rewardHeight,
		emptyWord)
	vm.Db.SetStorage(block.ToAddress(), locHash, registerInfo)
	return nil
}
func (p *register) doReceiveCancelRegister(vm *VM, block VmAccountBlock) error {
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) < 72 || allZero(old[0:32]) {
		return ErrInvalidData
	}
	amount := new(big.Int).SetBytes(old[0:32])
	snapshotBlock := vm.Db.SnapshotBlock(block.SnapshotHash())
	registerInfo := joinBytes(emptyWord,
		leftPadBytes(new(big.Int).SetInt64(snapshotBlock.Timestamp()).Bytes(), 8),
		old[40:72],
		leftPadBytes(snapshotBlock.Height().Bytes(), 32))
	vm.Db.SetStorage(block.ToAddress(), locHash, registerInfo)
	refundBlock := vm.createBlock(block.ToAddress(), block.AccountAddress(), BlockTypeSendCall, block.Depth()+1)
	refundBlock.SetAmount(amount)
	refundBlock.SetTokenId(viteTokenTypeId)
	refundBlock.SetHeight(new(big.Int).Add(block.Height(), big1))
	vm.blockList = append(vm.blockList, refundBlock)
	return nil
}
func (p *register) doReceiveReward(vm *VM, block VmAccountBlock) error {
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) < 72 || !bytes.Equal(old[40:72], block.Data()[43:75]) {
		return ErrInvalidData
	}
	if len(old) >= 104 && bytes.Equal(old[72:104], block.Data()[11:43]) {
		vm.Db.SetStorage(block.ToAddress(), locHash, []byte{})
	} else {
		var registerInfo []byte
		if len(old) >= 104 {
			registerInfo = joinBytes(old[0:40], block.Data()[11:43], old[72:104])
		} else {
			registerInfo = joinBytes(old[0:40], block.Data()[11:43])
		}
		vm.Db.SetStorage(block.ToAddress(), locHash, registerInfo)
	}
	refundBlock := vm.createBlock(block.ToAddress(), block.AccountAddress(), BlockTypeSendReward, block.Depth()+1)
	refundBlock.SetAmount(new(big.Int).SetBytes(block.Data()[75:107]))
	refundBlock.SetTokenId(viteTokenTypeId)
	refundBlock.SetHeight(new(big.Int).Add(block.Height(), big1))
	vm.blockList = append(vm.blockList, refundBlock)
	return nil
}

type vote struct{}

var (
	DataVote       = byte(1)
	DataCancelVote = byte(2)
)

func (p *vote) doSend(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if len(block.Data()) == 31 && block.Data()[0] == DataVote {
		return p.doSendVote(vm, block, quotaLeft)
	} else if len(block.Data()) == 11 && block.Data()[0] == DataCancelVote {
		return p.doSendCancelVote(vm, block, quotaLeft)
	}
	return quotaLeft, ErrInvalidData
}
func (p *vote) doSendVote(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Sign() != 0 ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	gid, _ := BytesToGid(block.Data()[1:11])
	if !vm.Db.IsExistGid(gid) {
		return quotaLeft, ErrInvalidData
	}
	address, _ := types.BytesToAddress(block.Data()[11:31])
	if !vm.Db.IsExistAddress(address) {
		return quotaLeft, ErrInvalidData
	}
	quotaLeft, err := useQuota(quotaLeft, voteGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}
func (p *vote) doSendCancelVote(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Sign() != 0 ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	gid, _ := BytesToGid(block.Data()[1:11])
	if !vm.Db.IsExistGid(gid) {
		return quotaLeft, ErrInvalidData
	}
	quotaLeft, err := useQuota(quotaLeft, cancelVoteGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}
func (p *vote) doReceive(vm *VM, block VmAccountBlock) error {
	if len(block.Data()) == 11 && block.Data()[0] == DataVote {
		return p.doReceiveVote(vm, block)
	} else if len(block.Data()) == 11 && block.Data()[0] == DataCancelVote {
		return p.doReceiveCancelVote(vm, block)
	}
	return nil
}
func (p *vote) doReceiveVote(vm *VM, block VmAccountBlock) error {
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	vm.Db.SetStorage(block.ToAddress(), locHash, block.Data()[11:31])
	return nil
}
func (p *vote) doReceiveCancelVote(vm *VM, block VmAccountBlock) error {
	gid, _ := BytesToGid(block.Data()[1:11])
	locHash := getKey(block.AccountAddress(), gid)
	vm.Db.SetStorage(block.ToAddress(), locHash, nil)
	return nil
}

type mortgage struct{}

var (
	DataMortgage       = byte(1)
	DataCancelMortgage = byte(2)
)

func (p *mortgage) doSend(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if len(block.Data()) == 29 && block.Data()[0] == DataMortgage {
		return p.doSendMortgage(vm, block, quotaLeft)
	} else if len(block.Data()) == 53 && block.Data()[0] == DataCancelMortgage {
		return p.doSendCancelMortgage(vm, block, quotaLeft)
	}
	return quotaLeft, ErrInvalidData
}
func (p *mortgage) doSendMortgage(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if block.Amount().Sign() == 0 ||
		!bytes.Equal(block.TokenId().Bytes(), viteTokenTypeId.Bytes()) ||
		!isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	address, _ := types.BytesToAddress(block.Data()[1:21])
	if !vm.Db.IsExistAddress(address) {
		return quotaLeft, ErrInvalidData
	}
	withdrawTime := new(big.Int).SetBytes(block.Data()[21:29]).Int64()
	if withdrawTime < vm.Db.SnapshotBlock(block.SnapshotHash()).Timestamp()+mortgageTime {
		return quotaLeft, ErrInvalidData
	}
	quotaLeft, err := useQuota(quotaLeft, mortgageGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}
func (p *mortgage) doSendCancelMortgage(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	if !isUserAccount(vm.Db, block.AccountAddress()) {
		return quotaLeft, ErrInvalidData
	}
	address, _ := types.BytesToAddress(block.Data()[1:21])
	if !vm.Db.IsExistAddress(address) {
		return quotaLeft, ErrInvalidData
	}
	if new(big.Int).SetBytes(block.Data()[21:53]).Sign() != 1 {
		return quotaLeft, ErrInvalidData
	}
	quotaLeft, err := useQuota(quotaLeft, cancelMortgageGas)
	if err != nil {
		return quotaLeft, err
	}
	return quotaLeft, nil
}

func (p *mortgage) doReceive(vm *VM, block VmAccountBlock) error {
	if len(block.Data()) == 11 && block.Data()[0] == DataVote {
		return p.doReceiveMortgage(vm, block)
	} else if len(block.Data()) == 11 && block.Data()[0] == DataCancelVote {
		return p.doReceiveCancelMortgage(vm, block)
	}
	return nil
}

func (p *mortgage) doReceiveMortgage(vm *VM, block VmAccountBlock) error {
	locHash := types.DataHash(append(block.Data()[1:21], types.DataHash(block.AccountAddress().Bytes()).Bytes()...))
	old := vm.Db.Storage(block.ToAddress(), locHash)
	newWithdrawTime := new(big.Int).SetBytes(block.Data()[21:29]).Int64() + vm.Db.SnapshotBlock(block.SnapshotHash()).Timestamp()
	newAmount := new(big.Int)
	if len(old) >= 40 {
		if newWithdrawTime < new(big.Int).SetBytes(old[32:40]).Int64() {
			return ErrInvalidData
		}
		newAmount.SetBytes(old[0:32])
	}
	newAmount.Add(newAmount, block.Amount())
	vm.Db.SetStorage(block.ToAddress(), locHash, joinBytes(leftPadBytes(newAmount.Bytes(), 32), leftPadBytes(new(big.Int).SetInt64(newWithdrawTime).Bytes(), 8)))

	locHashAmount := types.DataHash(block.Data()[1:21])
	oldAmount := vm.Db.Storage(block.ToAddress(), locHashAmount)
	newMortgageAmount := new(big.Int)
	if len(oldAmount) >= 32 {
		newMortgageAmount.SetBytes(oldAmount[0:32])
	}
	newMortgageAmount.Add(newMortgageAmount, block.Amount())
	vm.Db.SetStorage(block.ToAddress(), locHashAmount, leftPadBytes(newMortgageAmount.Bytes(), 32))
	return nil
}
func (p *mortgage) doReceiveCancelMortgage(vm *VM, block VmAccountBlock) error {
	locHash := types.DataHash(append(block.Data()[1:21], types.DataHash(block.AccountAddress().Bytes()).Bytes()...))
	old := vm.Db.Storage(block.ToAddress(), locHash)
	if len(old) < 40 {
		return ErrInvalidData
	}
	withdrawTime := new(big.Int).SetBytes(old[32:40]).Int64()
	if withdrawTime > vm.Db.SnapshotBlock(block.SnapshotHash()).Timestamp() {
		return ErrInvalidData
	}
	leftAmount := new(big.Int).SetBytes(old[0:32])
	withdrawAmount := new(big.Int).SetBytes(block.Data()[21:53])
	if leftAmount.Cmp(withdrawAmount) < 0 {
		return ErrInvalidData
	}
	leftAmount.Sub(leftAmount, withdrawAmount)

	locHashAmount := types.DataHash(block.Data()[1:21])
	oldLeftMortgage := vm.Db.Storage(block.ToAddress(), locHashAmount)
	if len(oldLeftMortgage) < 32 {
		return ErrInvalidData
	}
	newMortgageAmount := new(big.Int).SetBytes(oldLeftMortgage[0:32])
	if newMortgageAmount.Cmp(withdrawAmount) < 0 {
		return ErrInvalidData
	}
	newMortgageAmount.Sub(newMortgageAmount, withdrawAmount)

	if leftAmount.Sign() == 0 {
		vm.Db.SetStorage(block.ToAddress(), locHash, nil)
	} else {
		vm.Db.SetStorage(block.ToAddress(), locHash, joinBytes(leftPadBytes(leftAmount.Bytes(), 32), old[32:40]))
	}

	if newMortgageAmount.Sign() == 0 {
		vm.Db.SetStorage(block.ToAddress(), locHashAmount, nil)
	} else {
		vm.Db.SetStorage(block.ToAddress(), locHashAmount, leftPadBytes(newMortgageAmount.Bytes(), 32))
	}

	// append refund block
	refundBlock := vm.createBlock(block.ToAddress(), block.AccountAddress(), BlockTypeSendCall, block.Depth()+1)
	refundBlock.SetAmount(withdrawAmount)
	refundBlock.SetTokenId(viteTokenTypeId)
	refundBlock.SetHeight(new(big.Int).Add(block.Height(), big1))
	vm.blockList = append(vm.blockList, refundBlock)
	return nil
}

type consensusGroup struct{}

func (p *consensusGroup) doSend(vm *VM, block VmAccountBlock, quotaLeft uint64) (uint64, error) {
	// TODO
	return 0, nil
}
func (p *consensusGroup) doReceive(vm *VM, block VmAccountBlock) error {
	// TODO
	return nil
}

func isUserAccount(db VmDatabase, addr types.Address) bool {
	return len(db.ContractCode(addr)) == 0
}

func getKey(addr types.Address, gid Gid) types.Hash {
	var data = types.Hash{}
	copy(data[2:12], gid[:])
	copy(data[12:], addr[:])
	return data
}
