package trie

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vitelabs/go-vite/chain_db/database"
	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/common/types"
	"path/filepath"
	"testing"
)

func TestNewTrie(t *testing.T) {
	db := database.NewLevelDb(filepath.Join(common.GoViteTestDataDir(), "trie"))
	defer db.Close()

	pool := NewTrieNodePool()

	trie, ntErr := NewTrie(db, nil, pool)
	if ntErr != nil {
		t.Fatal(ntErr)
	}

	fmt.Println(1)
	trie.SetValue(nil, []byte("NilNilNilNilNil"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(2)
	trie.SetValue(nil, []byte("NilNilNilNilNil234"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(3)
	trie.SetValue([]byte("test"), []byte("value.hash"))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(4)
	trie.SetValue([]byte("tesa"), []byte("value.hash2"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(5)
	trie.SetValue([]byte("aofjas"), []byte("value.content1"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(6)
	trie.SetValue([]byte("aofjas"), []byte("value.content2value.content2value.content2value.content2value.content2value.content2value.content2value.content2"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(7)
	trie.SetValue([]byte("tesa"), []byte("value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(8)
	trie.SetValue([]byte("tesa"), []byte("value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value.hash3value09909"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(9)
	trie.SetValue([]byte("tesabcd"), []byte("value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4value.hash4"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(10)
	trie.SetValue([]byte("tesab"), []byte("value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println("10.1")
	trie.SetValue([]byte("tesab"), []byte("value.555"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(11)
	trie.SetValue([]byte("t"), []byte("somethinghiOkYesYourMyASDKJBNXA1239xnm.0j8n120k0k12nz$0231*&^$@!!())$S@@ST&&@@SDT&(OL<><:PP_+}}GC~@@@#$%^&&HCXZkasldjf1009100"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println("11.1")
	trie.SetValue([]byte("t"), []byte("abc"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(12)
	trie.SetValue([]byte("a"), []byte("a1230xm9"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println("12.1")
	trie.SetValue([]byte("a"), []byte("a10xm9"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(13)
	trie.SetValue([]byte("IamGood"), []byte("a1230xm90zm19ma"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGood")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(14)
	trie.SetValue([]byte("IamGood"), []byte("hahaheheh"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGood")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(15)
	trie.SetValue([]byte("IamGoo"), []byte("ijukh"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGoo")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(16)
	trie.SetValue([]byte("IamG"), []byte("ki10$%^%&@#!@#"))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGoo")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamG")))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()

	fmt.Println(17)
	trie.SetValue(nil, []byte("isNil"))
	fmt.Printf("%s\n", trie.GetValue([]byte("test")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", trie.GetValue([]byte("aofjas")))
	fmt.Printf("%s\n", trie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", trie.GetValue([]byte("t")))
	fmt.Printf("%s\n", trie.GetValue([]byte("a")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamGoo")))
	fmt.Printf("%s\n", trie.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", trie.GetValue(nil))
	fmt.Printf("%d\n", len(trie.unSavedRefValueMap))
	fmt.Println()
}

func TestTrieHash(t *testing.T) {
	db := database.NewLevelDb(filepath.Join(common.GoViteTestDataDir(), "trie"))
	defer db.Close()

	pool := NewTrieNodePool()

	trie, ntErr := NewTrie(db, nil, pool)
	if ntErr != nil {
		t.Fatal(ntErr)
	}

	trie.SetValue(nil, []byte("NilNilNilNilNil"))
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	trie.SetValue(nil, []byte("isNil"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("IamG"), []byte("ki10$%^%&@#!@#"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("IamGood"), []byte("a1230xm90zm19ma"))
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("tesab"), []byte("value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555"))
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("tesab"), []byte("value.555val"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("tesa"), []byte("vale....asdfasdfasdfvalue.555val"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("tes"), []byte("asdfvale....asdfasdfasdfvalue.555val"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("tesabcd"), []byte("asdfvale....asdfasdfasdfvalue.555val"))
	fmt.Println(trie.Hash())
	trie.SetValue([]byte("t"), []byte("asdfvale....asdfasdfasdfvalue.555valasd"))
	fmt.Println(trie.Hash())
}

func TestTrieSaveAndLoad(t *testing.T) {
	db := database.NewLevelDb(filepath.Join(common.GoViteTestDataDir(), "trie"))
	defer db.Close()

	pool := NewTrieNodePool()

	trie, ntErr := NewTrie(db, nil, pool)
	if ntErr != nil {
		t.Fatal(ntErr)
	}
	trie.SetValue(nil, []byte("NilNilNilNilNil"))
	trie.SetValue([]byte("IamG"), []byte("ki10$%^%&@#!@#"))
	trie.SetValue([]byte("IamGood"), []byte("a1230xm90zm19ma"))
	trie.SetValue([]byte("tesab"), []byte("value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555value.555"))

	trie.SetValue([]byte("tesab"), []byte("value.555val"))
	trie.SetValue([]byte("tesa"), []byte("vale....asdfasdfasdfvalue.555val"))
	trie.SetValue([]byte("tesa"), []byte("vale....asdfasdfasdfvalue.555val"))
	trie.SetValue([]byte("tes"), []byte("asdfvale....asdfasdfasdfvalue.555val"))
	trie.SetValue([]byte("tesabcd"), []byte("asdfvale....asdfasdfasdfvalue.555val"))
	trie.SetValue([]byte("t"), []byte("asdfvale....asdfasdfasdfvalue.555valasd"))
	fmt.Println(trie.Hash())
	fmt.Println()

	batch := new(leveldb.Batch)
	callback, _ := trie.Save(batch)
	db.Write(batch, nil)
	callback()

	trie = nil

	rootHash, _ := types.HexToHash("9df5e11da5cdaea43fa69991fb7cf575ec00c73e90d98cf67dbf2e9bdca9a998")
	newTrie, ntErr := NewTrie(db, &rootHash, pool)
	if ntErr != nil {
		t.Fatal(ntErr)
	}
	fmt.Printf("%s\n", newTrie.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTrie.GetValue([]byte("t")))
	fmt.Println(newTrie.Hash())
	fmt.Println()
	newTrie = nil

	newTri2, ntErr := NewTrie(db, &rootHash, pool)
	if ntErr != nil {
		t.Fatal(ntErr)
	}
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("t")))
	fmt.Println(newTri2.Hash())
	fmt.Println()

	newTri2.SetValue([]byte("tesab"), []byte("value.hahaha123"))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("t")))
	fmt.Println(newTri2.Hash())
	fmt.Println()

	newTri2.SetValue([]byte("IamGood"), []byte("Yes you are good."))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("t")))
	fmt.Println(newTri2.Hash())
	fmt.Println()

	batch2 := new(leveldb.Batch)
	callback2, _ := newTri2.Save(batch2)
	if err := db.Write(batch2, nil); err != nil {
		t.Fatal(err)
	}
	callback2()

	rootHash2 := newTri2.Hash()
	newTrie3, _ := NewTrie(db, rootHash2, pool)
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTrie3.GetValue([]byte("t")))
	fmt.Println(newTrie3.Hash())
	fmt.Println()

	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamG")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("IamGood")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesab")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesa")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tes")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("tesabcd")))
	fmt.Printf("%s\n", newTri2.GetValue([]byte("t")))
	fmt.Println(newTri2.Hash())
	fmt.Println()
}

func TestTrieCopy(t *testing.T) {

}
