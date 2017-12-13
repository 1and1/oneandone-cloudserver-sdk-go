package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_bs       sync.Once
	test_bs_name string
	test_bs_desc string
	test_bs      *BlockStorage
)

func create_block_storage() *BlockStorage {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_bs_name = fmt.Sprintf("BlockStorage_%d", rint)
	test_bs_desc = fmt.Sprintf("BlockStorage_%d description", rint)
	req := BlockStorageRequest{
		Name:         test_bs_name,
		Description:  test_bs_desc,
		Size:         Int2Pointer(20),
		DatacenterID: "908DC2072407C94C8054610AD5A53B8C",
	}
	fmt.Printf("Creating new block storage '%s'...\n", test_bs_name)
	bs_id, bs, err := api.CreateBlockStorage(&req)

	if err != nil {
		fmt.Printf("Unable to create a block storage. Error: %s", err.Error())
		return nil
	}

	if bs_id == "" || bs.Id == "" {
		fmt.Printf("Unable to create block storage '%s'.", test_bs_name)
		return nil
	}

	api.WaitForState(bs, "ACTIVE", 10, 30)
	return bs
}

func set_block_storage() {
	test_bs = create_block_storage()
}

func TestListBlockStorages(t *testing.T) {
	set_bs.Do(set_block_storage)

	res, err := api.ListBlockStorages()
	if err != nil {
		t.Errorf("ListBlockStorages failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No block storage found.")
	}
}

func TestDeleteBlockStorage(t *testing.T) {
	set_bs.Do(set_block_storage)

	fmt.Printf("Deleting block storage '%s'...\n", test_bs.Name)
	bs, err := api.DeleteBlockStorage(test_bs.Id)
	if err != nil {
		t.Errorf("DeleteSharedStorage failed. Error: " + err.Error())
		return
	} else {
		api.WaitUntilDeleted(bs)
	}

	bs, err = api.GetBlockStorage(bs.Id)

	if bs != nil {
		t.Errorf("Unable to delete the block storage.")
	} else {
		test_bs = nil
	}
}
