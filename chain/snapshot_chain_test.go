package chain

import (
	"fmt"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vm_context"
	"math/big"
	"testing"
	"time"
)

//func TestGetNeedSnapshotContent(t *testing.T) {
//	chainInstance := getChainInstance()
//	content := chainInstance.GetNeedSnapshotContent()
//	for addr, item := range content {
//		fmt.Printf("%s: %+v\n", addr.String(), item)
//	}
//}
//
//func TestInsertSnapshotBlock(t *testing.T) {
//
//}
//
//func TestGetSnapshotBlocksByHash(t *testing.T) {
//	chainInstance := getChainInstance()
//	blocks, err := chainInstance.GetSnapshotBlocksByHash(nil, 100, true, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for index, block := range blocks {
//		fmt.Printf("%d: %+v\n", index, block)
//	}
//
//	blocks2, err2 := chainInstance.GetSnapshotBlocksByHash(nil, 100, true, true)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	for index, block := range blocks2 {
//		fmt.Printf("%d: %+v\n", index, block)
//	}
//
//	blocks3, err3 := chainInstance.GetSnapshotBlocksByHash(nil, 100, false, true)
//	if err3 != nil {
//		t.Fatal(err3)
//	}
//	for index, block := range blocks3 {
//		fmt.Printf("%d: %+v\n", index, block)
//	}
//}
//
//func TestGetSnapshotBlocksByHeight(t *testing.T) {
//	chainInstance := getChainInstance()
//	blocks, err := chainInstance.GetSnapshotBlocksByHeight(2, 10, false, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for index, block := range blocks {
//		fmt.Printf("%d: %+v\n", index, block)
//	}
//}
//
//func TestGetSnapshotBlockByHeight(t *testing.T) {
//	chainInstance := getChainInstance()
//	block, err := chainInstance.GetSnapshotBlockByHeight(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", block)
//
//	block2, err2 := chainInstance.GetSnapshotBlockByHeight(2)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	fmt.Printf("%+v\n", block2)
//}
//
//func TestGetSnapshotBlockByHash(t *testing.T) {
//	chainInstance := getChainInstance()
//	block, err := chainInstance.GetSnapshotBlockByHash(&GenesisMintageSendBlock.Hash)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", block)
//
//	hash2, _ := types.HexToHash("f34e00c283f11728e28ccf2cf2138a7976b9ed7daaf7dbcc2ca598f66139f80d")
//	block2, err2 := chainInstance.GetSnapshotBlockByHash(&hash2)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	fmt.Printf("%+v\n", block2)
//}
//
//func TestGetLatestSnapshotBlock(t *testing.T) {
//	chainInstance := getChainInstance()
//	block := chainInstance.GetLatestSnapshotBlock()
//	fmt.Printf("%+v\n", block)
//}
//
//func TestGetGenesisSnapshotBlock(t *testing.T) {
//	chainInstance := getChainInstance()
//	block := chainInstance.GetGenesisSnapshotBlock()
//	fmt.Printf("%+v\n", block)
//}
//
//func TestGetConfirmBlock(t *testing.T) {
//	chainInstance := getChainInstance()
//	block, err := chainInstance.GetConfirmBlock(&GenesisMintageSendBlock.Hash)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", block)
//
//	hash, _ := types.HexToHash("8d9cef33f1c053f976844c489fc642855576ccd535cf2648412451d783147394")
//	block2, err2 := chainInstance.GetConfirmBlock(&hash)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	fmt.Printf("%+v\n", block2)
//
//	block3, err3 := chainInstance.GetConfirmBlock(&GenesisMintageBlock.Hash)
//	if err3 != nil {
//		t.Fatal(err3)
//	}
//	fmt.Printf("%+v\n", block3)
//}
//
//func TestGetConfirmTimes(t *testing.T) {
//	chainInstance := getChainInstance()
//	times1, err := chainInstance.GetConfirmTimes(&GenesisMintageSendBlock.Hash)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", times1)
//
//	hash, _ := types.HexToHash("8d9cef33f1c053f976844c489fc642855576ccd535cf2648412451d783147394")
//	times2, err2 := chainInstance.GetConfirmTimes(&hash)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	fmt.Printf("%+v\n", times2)
//
//	times3, err3 := chainInstance.GetConfirmTimes(&GenesisMintageBlock.Hash)
//	if err3 != nil {
//		t.Fatal(err3)
//	}
//	fmt.Printf("%+v\n", times3)
//}
//
//// TODO
//func TestGetSnapshotBlockBeforeTime(t *testing.T) {
//	chainInstance := getChainInstance()
//	time1 := time.Now()
//	block, err := chainInstance.GetSnapshotBlockBeforeTime(&time1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", block)
//
//	time2 := time.Unix(1535209021, 0)
//	block2, err2 := chainInstance.GetSnapshotBlockBeforeTime(&time2)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//	fmt.Printf("%+v\n", block2)
//
//	time3 := GenesisSnapshotBlock.Timestamp.Add(time.Second * 100)
//	block3, err3 := chainInstance.GetSnapshotBlockBeforeTime(&time3)
//	if err3 != nil {
//		t.Fatal(err3)
//	}
//	fmt.Printf("%+v\n", block3)
//}
//
//func TestGetConfirmAccountBlock(t *testing.T) {
//	chainInstance := getChainInstance()
//	block, err := chainInstance.GetConfirmAccountBlock(GenesisSnapshotBlock.Height, &GenesisMintageSendBlock.AccountAddress)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", block)
//
//	block2, err2 := chainInstance.GetConfirmAccountBlock(GenesisSnapshotBlock.Height, &GenesisRegisterBlock.AccountAddress)
//	if err2 != nil {
//		t.Fatal(err2)
//	}
//
//	fmt.Printf("%+v\n", block2)
//
//	block3, err3 := chainInstance.GetConfirmAccountBlock(GenesisSnapshotBlock.Height+10, &GenesisMintageSendBlock.AccountAddress)
//	if err3 != nil {
//		t.Fatal(err3)
//	}
//	fmt.Printf("%+v\n", block3)
//
//	block4, err4 := chainInstance.GetConfirmAccountBlock(0, &GenesisMintageSendBlock.AccountAddress)
//	if err4 != nil {
//		t.Fatal(err4)
//	}
//	fmt.Printf("%+v\n", block4)
//}

func randomSendViteBlock(snapshotBlockHash types.Hash, addr1 *types.Address, addr2 *types.Address) ([]*vm_context.VmAccountBlock, []types.Address, error) {
	chainInstance := getChainInstance()
	now := time.Now()

	if addr1 == nil {
		accountAddress, _, _ := types.CreateAddress()
		addr1 = &accountAddress
	}
	if addr2 == nil {
		accountAddress, _, _ := types.CreateAddress()
		addr2 = &accountAddress
	}

	vmContext, err := vm_context.NewVmContext(chainInstance, nil, nil, addr1)
	if err != nil {
		return nil, nil, err
	}
	latestBlock, _ := chainInstance.GetLatestAccountBlock(addr1)
	nextHeight := uint64(1)
	var prevHash types.Hash
	if latestBlock != nil {
		nextHeight = latestBlock.Height + 1
		prevHash = latestBlock.Hash
	}

	sendAmount := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e9))
	var sendBlock = &ledger.AccountBlock{
		PrevHash:       prevHash,
		BlockType:      ledger.BlockTypeSendCall,
		AccountAddress: *addr1,
		ToAddress:      *addr2,
		Amount:         sendAmount,
		TokenId:        ledger.ViteTokenId,
		Height:         nextHeight,
		Fee:            big.NewInt(0),
		//PublicKey:      publicKey,
		SnapshotHash: snapshotBlockHash,
		Timestamp:    &now,
		Nonce:        []byte("test nonce test nonce"),
		Signature:    []byte("test signature test signature test signature"),
	}

	vmContext.AddBalance(&ledger.ViteTokenId, sendAmount)

	sendBlock.StateHash = *vmContext.GetStorageHash()
	sendBlock.Hash = sendBlock.ComputeHash()
	return []*vm_context.VmAccountBlock{{
		AccountBlock: sendBlock,
		VmContext:    vmContext,
	}}, []types.Address{*addr1, *addr2}, nil
}

func newReceiveBlock(snapshotBlockHash types.Hash, accountAddress types.Address, fromHash types.Hash) ([]*vm_context.VmAccountBlock, error) {
	chainInstance := getChainInstance()
	latestBlock, _ := chainInstance.GetLatestAccountBlock(&accountAddress)
	nextHeight := uint64(1)
	var prevHash types.Hash
	if latestBlock != nil {
		nextHeight = latestBlock.Height + 1
		prevHash = latestBlock.Hash
	}

	now := time.Now()

	vmContext, err := vm_context.NewVmContext(chainInstance, nil, nil, &accountAddress)
	if err != nil {
		return nil, err
	}

	var receiveBlock = &ledger.AccountBlock{
		PrevHash:       prevHash,
		BlockType:      ledger.BlockTypeReceive,
		AccountAddress: accountAddress,
		FromBlockHash:  fromHash,
		Height:         nextHeight,
		Fee:            big.NewInt(0),
		SnapshotHash:   snapshotBlockHash,
		Timestamp:      &now,
		Nonce:          []byte("test nonce test nonce"),
		Signature:      []byte("test signature test signature test signature"),
	}

	vmContext.AddBalance(&ledger.ViteTokenId, big.NewInt(100))

	receiveBlock.StateHash = *vmContext.GetStorageHash()
	receiveBlock.Hash = receiveBlock.ComputeHash()

	return []*vm_context.VmAccountBlock{{
		AccountBlock: receiveBlock,
		VmContext:    vmContext,
	}}, nil
}

func newSnapshotBlock() (*ledger.SnapshotBlock, error) {
	chainInstance := getChainInstance()

	latestBlock := chainInstance.GetLatestSnapshotBlock()
	now := time.Now()
	snapshotBlock := &ledger.SnapshotBlock{
		Height:    latestBlock.Height + 1,
		PrevHash:  latestBlock.Hash,
		Timestamp: &now,
	}

	content := chainInstance.GetNeedSnapshotContent()
	snapshotBlock.SnapshotContent = content

	trie, err := chainInstance.GenStateTrie(latestBlock.StateHash, content)
	if err != nil {
		return nil, err
	}

	snapshotBlock.StateTrie = trie
	snapshotBlock.StateHash = *trie.Hash()
	snapshotBlock.Hash = snapshotBlock.ComputeHash()

	return snapshotBlock, err
}

func TestDeleteSnapshotBlocksToHeight(t *testing.T) {
	chainInstance := getChainInstance()

	snapshotBlock, err0 := newSnapshotBlock()
	if err0 != nil {
		t.Fatal(err0)
	}

	chainInstance.InsertSnapshotBlock(snapshotBlock)

	blocks, addressList, err := randomSendViteBlock(snapshotBlock.Hash, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	chainInstance.InsertAccountBlocks(blocks)

	blocks2, addressList2, err2 := randomSendViteBlock(snapshotBlock.Hash, nil, nil)
	if err2 != nil {
		t.Fatal(err)
	}
	chainInstance.InsertAccountBlocks(blocks2)

	snapshotBlock2, err3 := newSnapshotBlock()
	if err3 != nil {
		t.Fatal(err3)
	}

	chainInstance.InsertSnapshotBlock(snapshotBlock2)

	receiveBlock, _ := newReceiveBlock(snapshotBlock2.Hash, addressList[1], blocks[0].AccountBlock.Hash)
	chainInstance.InsertAccountBlocks(receiveBlock)

	receiveBlock2, _ := newReceiveBlock(snapshotBlock2.Hash, addressList2[1], blocks2[0].AccountBlock.Hash)
	chainInstance.InsertAccountBlocks(receiveBlock2)

	snapshotBlock3, _ := newSnapshotBlock()
	chainInstance.InsertSnapshotBlock(snapshotBlock3)

	var display = func() {
		dBlocks1, _ := chainInstance.GetAccountBlocksByHeight(blocks[0].AccountBlock.AccountAddress, 0, 10, true)
		for _, block := range dBlocks1 {
			fmt.Printf("%+v\n", block)
		}
		dBlocks2, _ := chainInstance.GetAccountBlocksByHeight(blocks2[0].AccountBlock.AccountAddress, 0, 10, true)
		for _, block := range dBlocks2 {
			fmt.Printf("%+v\n", block)
		}
		dBlocks3, _ := chainInstance.GetAccountBlocksByHeight(receiveBlock[0].AccountBlock.AccountAddress, 0, 10, true)
		for _, block := range dBlocks3 {
			fmt.Printf("%+v\n", block)
		}
		dBlocks4, _ := chainInstance.GetAccountBlocksByHeight(receiveBlock2[0].AccountBlock.AccountAddress, 0, 10, true)
		for _, block := range dBlocks4 {
			fmt.Printf("%+v\n", block)
		}

		latestBlock := chainInstance.GetLatestSnapshotBlock()
		fmt.Printf("%+v\n", latestBlock)

	}
	display()

	fmt.Println()

	blockMeta, _ := chainInstance.ChainDb().Ac.GetBlockMeta(&blocks[0].AccountBlock.Hash)
	fmt.Printf("%+v\n", blockMeta)

	blockMeta1, _ := chainInstance.ChainDb().Ac.GetBlockMeta(&blocks2[0].AccountBlock.Hash)
	fmt.Printf("%+v\n", blockMeta1)

	fmt.Println()
	sbList, abMap, deleteErr := chainInstance.DeleteSnapshotBlocksToHeight(3)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
	for _, sb := range sbList {
		fmt.Printf("%+v\n", sb)
	}

	for addr, abs := range abMap {
		fmt.Printf("%s\n", addr.String())
		for _, ab := range abs {
			fmt.Printf("%+v\n", ab)
		}
	}
	fmt.Println()

	display()
	fmt.Println()

	needContent := chainInstance.GetNeedSnapshotContent()
	for addr, content := range needContent {
		fmt.Printf("%s: %+v\n", addr.String(), content)
	}
	fmt.Println()
	blockMeta2, _ := chainInstance.ChainDb().Ac.GetBlockMeta(&blocks[0].AccountBlock.Hash)
	fmt.Printf("%+v\n", blockMeta2)

	blockMeta3, _ := chainInstance.ChainDb().Ac.GetBlockMeta(&blocks2[0].AccountBlock.Hash)
	fmt.Printf("%+v\n", blockMeta3)

}
