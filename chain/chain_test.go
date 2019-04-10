package chain

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"testing"

	"github.com/vitelabs/go-vite/chain/test_tools"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/config"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vm/quota"
	"math/rand"
	"sync"
)

var genesisConfigJson = "{  \"GenesisAccountAddress\": \"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\",  \"ForkPoints\": {  },  \"ConsensusGroupInfo\": {    \"ConsensusGroupInfoMap\":{      \"00000000000000000001\":{        \"NodeCount\":25,        \"Interval\":1,        \"PerCount\":3,        \"RandCount\":2,        \"RandRank\":100,        \"Repeat\":1,        \"CheckLevel\":0,        \"CountingTokenId\":\"tti_5649544520544f4b454e6e40\",        \"RegisterConditionId\":1,        \"RegisterConditionParam\":{          \"PledgeAmount\": 100000000000000000000000,          \"PledgeHeight\": 1,          \"PledgeToken\": \"tti_5649544520544f4b454e6e40\"        },        \"VoteConditionId\":1,        \"VoteConditionParam\":{},        \"Owner\":\"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\",        \"PledgeAmount\":0,        \"WithdrawHeight\":1      },      \"00000000000000000002\":{        \"NodeCount\":25,        \"Interval\":3,        \"PerCount\":1,        \"RandCount\":2,        \"RandRank\":100,        \"Repeat\":48,        \"CheckLevel\":1,        \"CountingTokenId\":\"tti_5649544520544f4b454e6e40\",        \"RegisterConditionId\":1,        \"RegisterConditionParam\":{          \"PledgeAmount\": 100000000000000000000000,          \"PledgeHeight\": 1,          \"PledgeToken\": \"tti_5649544520544f4b454e6e40\"        },        \"VoteConditionId\":1,        \"VoteConditionParam\":{},        \"Owner\":\"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\",        \"PledgeAmount\":0,        \"WithdrawHeight\":1      }    },    \"RegistrationInfoMap\":{      \"00000000000000000001\":{        \"s1\":{          \"NodeAddr\":\"vite_14edbc9214bd1e5f6082438f707d10bf43463a6d599a4f2d08\",          \"PledgeAddr\":\"vite_14edbc9214bd1e5f6082438f707d10bf43463a6d599a4f2d08\",          \"Amount\":100000000000000000000000,          \"WithdrawHeight\":7776000,          \"RewardTime\":1,          \"CancelTime\":0,          \"HisAddrList\":[\"vite_14edbc9214bd1e5f6082438f707d10bf43463a6d599a4f2d08\"]        },        \"s2\":{          \"NodeAddr\":\"vite_0acbb1335822c8df4488f3eea6e9000eabb0f19d8802f57c87\",          \"PledgeAddr\":\"vite_0acbb1335822c8df4488f3eea6e9000eabb0f19d8802f57c87\",          \"Amount\":100000000000000000000000,          \"WithdrawHeight\":7776000,          \"RewardTime\":1,          \"CancelTime\":0,          \"HisAddrList\":[\"vite_0acbb1335822c8df4488f3eea6e9000eabb0f19d8802f57c87\"]        }      }    }  },  \"MintageInfo\":{    \"TokenInfoMap\":{      \"tti_5649544520544f4b454e6e40\":{        \"TokenName\":\"Vite Token\",        \"TokenSymbol\":\"VITE\",        \"TotalSupply\":1000000000000000000000000000,        \"Decimals\":18,        \"Owner\":\"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\",        \"PledgeAmount\":0,        \"PledgeAddr\":\"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\",        \"WithdrawHeight\":0,        \"MaxSupply\":115792089237316195423570985008687907853269984665640564039457584007913129639935,        \"OwnerBurnOnly\":false,        \"IsReIssuable\":true      }    },    \"LogList\": [        {          \"Data\": \"\",          \"Topics\": [            \"3f9dcc00d5e929040142c3fb2b67a3be1b0e91e98dac18d5bc2b7817a4cfecb6\",            \"000000000000000000000000000000000000000000005649544520544f4b454e\"          ]        }      ]  },  \"AccountBalanceMap\": {    \"vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a\": {      \"tti_5649544520544f4b454e6e40\":899999000000000000000000000    },    \"vite_56fd05b23ff26cd7b0a40957fb77bde60c9fd6ebc35f809c23\": {      \"tti_5649544520544f4b454e6e40\":100000000000000000000000000    }  }}"

type MockConsensus struct{}

func (c *MockConsensus) VerifyAccountProducer(block *ledger.AccountBlock) (bool, error) {
	return true, nil
}

func NewChainInstance(dirName string, clear bool) (*chain, error) {
	dataDir := path.Join(test_tools.DefaultDataDir(), dirName)
	if clear {
		os.RemoveAll(dataDir)
	}
	genesisConfig := &config.Genesis{}
	json.Unmarshal([]byte(genesisConfigJson), genesisConfig)

	chainInstance := NewChain(dataDir, &config.Chain{}, genesisConfig)

	if err := chainInstance.Init(); err != nil {
		return nil, err
	}
	// mock consensus
	chainInstance.SetConsensus(&MockConsensus{})

	chainInstance.Start()
	return chainInstance, nil
}

func SetUp(t *testing.T, accountNum, txCount, snapshotPerBlockNum int) (*chain, map[types.Address]*Account, []*ledger.SnapshotBlock) {
	// test quota
	quota.InitQuotaConfig(true, true)

	chainInstance, err := NewChainInstance("unit_test", false)
	if err != nil {
		t.Fatal(err)
	}

	InsertSnapshotBlock(chainInstance)

	accounts := make(map[types.Address]*Account)
	unconfirmedBlocks := chainInstance.cache.GetUnconfirmedBlocks()
	for _, accountBlock := range unconfirmedBlocks {
		if _, ok := accounts[accountBlock.AccountAddress]; !ok {
			accounts[accountBlock.AccountAddress] = NewAccount(chainInstance, accountBlock.PublicKey, nil)
		}
	}

	if len(accounts) < accountNum {
		lackNum := accountNum - len(accounts)
		newAccounts := MakeAccounts(chainInstance, lackNum)
		for addr, account := range newAccounts {
			accounts[addr] = account
		}

	}
	var snapshotBlockList []*ledger.SnapshotBlock

	t.Run("InsertBlocks", func(t *testing.T) {
		//InsertAccountBlock(t, chainInstance, accounts, txCount, snapshotPerBlockNum)
		snapshotBlockList = InsertAccountBlock(t, chainInstance, accounts, txCount, snapshotPerBlockNum)
	})

	return chainInstance, accounts, snapshotBlockList
}

func TearDown(chainInstance *chain) {
	chainInstance.Stop()
	chainInstance.Destroy()
}

func TestSetup(t *testing.T) {
	SetUp(t, 100, 1231, 9)
}

func TestChain(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	chainInstance, accounts, snapshotBlockList := SetUp(t, 20, 100, 3)
	testChainAll(t, chainInstance, accounts, snapshotBlockList)

	snapshotBlockList = append(snapshotBlockList, InsertAccountBlock(t, chainInstance, accounts, rand.Intn(1232), rand.Intn(5))...)

	// test all
	testChainAll(t, chainInstance, accounts, snapshotBlockList)

	// test insert & delete
	snapshotBlockList = testInsertAndDelete(t, chainInstance, accounts, snapshotBlockList)

	// test panic
	//snapshotBlockList = testPanic(t, chainInstance, accounts, snapshotBlockList)

	TearDown(chainInstance)
}

func testChainAll(t *testing.T, chainInstance *chain, accounts map[types.Address]*Account, snapshotBlockList []*ledger.SnapshotBlock) {
	// account
	testAccount(t, chainInstance, accounts)

	// account block
	testAccountBlock(t, chainInstance, accounts)

	// on road
	testOnRoad(t, chainInstance, accounts)

	// snapshot block
	testSnapshotBlock(t, chainInstance, accounts, snapshotBlockList)

	// state
	testState(t, chainInstance, accounts, snapshotBlockList)

	// built-in contract
	testBuiltInContract(t, chainInstance, accounts, snapshotBlockList)
}

func testPanic(t *testing.T, chainInstance *chain, accounts map[types.Address]*Account, snapshotBlockList []*ledger.SnapshotBlock) []*ledger.SnapshotBlock {
	var wg sync.WaitGroup
	for j := 0; j < 3; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// test insert & delete
			go func() {
				panic("error")
			}()

			snapshotBlockList = testInsertAndDelete(t, chainInstance, accounts, snapshotBlockList)

		}()

		wg.Wait()

		// recover snapshotBlockList
		if len(snapshotBlockList) > 0 {
			maxIndex := len(snapshotBlockList) - 1
			i := maxIndex
			for ; i > 0; i-- {
				snapshotBlock := snapshotBlockList[i]
				if ok, err := chainInstance.IsSnapshotBlockExisted(snapshotBlock.Hash); err != nil {
					t.Fatal(err)
				} else if ok {
					break
				}
			}

			if i < maxIndex {
				for _, account := range accounts {
					account.DeleteSnapshotBlocks(accounts, snapshotBlockList[i+1:])
				}

				snapshotBlockList = snapshotBlockList[:i+1]
			}

		}

		// recover accounts
		for _, account := range accounts {
			for account.latestBlock != nil {
				if ok, err := chainInstance.IsAccountBlockExisted(account.latestBlock.Hash); err != nil {
					t.Fatal(err)
				} else if !ok {
					account.rollbackLatestBlock()
				} else {
					break
				}
			}
		}
	}
	return snapshotBlockList

}
