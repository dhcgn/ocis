package indexer

import (
	"os"
	"path"
	"testing"

	"github.com/owncloud/ocis/accounts/pkg/indexer/option"

	"github.com/owncloud/ocis/accounts/pkg/config"
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/cs3"
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/disk"
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/stretchr/testify/assert"
)

const cs3RootFolder = "/var/tmp/ocis/storage/users/data"

func TestIndexer_CS3_AddWithUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", cs3RootFolder)
	assert.NoError(t, err)
	indexer := createCs3Indexer()

	err = indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique", nil, false)
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_CS3_AddWithNonUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", cs3RootFolder)
	assert.NoError(t, err)
	indexer := createCs3Indexer()

	err = indexer.AddIndex(&User{}, "UserName", "ID", "users", "non_unique", nil, false)
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_FindByWithUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique", nil, false)
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	res, err := indexer.FindBy(User{}, "UserName", "mikey")
	assert.NoError(t, err)
	t.Log(res)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_AddWithUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique", nil, false)
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_AddWithNonUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&Pet{}, "Kind", "ID", "pets", "non_unique", nil, false)
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindBy(Pet{}, "Kind", "Hog")
	assert.NoError(t, err)

	t.Log(res)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_AddWithAutoincrementIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&User{}, "UID", "ID", "users", "autoincrement", &option.Bound{Lower: 5}, false)
	assert.NoError(t, err)

	res1, err := indexer.Add(Data["users"][0])
	assert.NoError(t, err)
	assert.Equal(t, "UID", res1[0].Field)
	assert.Equal(t, "5", path.Base(res1[0].Value))

	res2, err := indexer.Add(Data["users"][1])
	assert.NoError(t, err)
	assert.Equal(t, "UID", res2[0].Field)
	assert.Equal(t, "6", path.Base(res2[0].Value))

	resFindBy, err := indexer.FindBy(User{}, "UID", "6")
	assert.NoError(t, err)
	assert.Equal(t, "hijklmn-456", resFindBy[0])
	t.Log(resFindBy)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_DeleteWithNonUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&Pet{}, "Kind", "ID", "pets", "non_unique", nil, false)
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	err = indexer.Delete(pet2)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_SearchWithNonUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&Pet{}, "Name", "ID", "pets", "non_unique", nil, false)
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindByPartial(pet2, "Name", "*ky")
	assert.NoError(t, err)

	t.Log(res)
	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_UpdateWithUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique", nil, false)
	assert.NoError(t, err)

	err = indexer.AddIndex(&User{}, "Email", "ID", "users", "unique", nil, false)
	assert.NoError(t, err)

	user1 := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	user2 := &User{ID: "hijklmn-456", UserName: "frank", Email: "frank@example.com"}

	_, err = indexer.Add(user1)
	assert.NoError(t, err)

	_, err = indexer.Add(user2)
	assert.NoError(t, err)

	err = indexer.Update(user1, &User{
		ID:       "abcdefg-123",
		UserName: "mikey-new",
		Email:    "mikey@example.com",
	})
	assert.NoError(t, err)
	v, err1 := indexer.FindBy(&User{}, "UserName", "mikey-new")
	assert.NoError(t, err1)
	assert.Len(t, v, 1)
	v, err2 := indexer.FindBy(&User{}, "UserName", "mikey")
	assert.NoError(t, err2)
	assert.Len(t, v, 0)

	err1 = indexer.Update(&User{
		ID:       "abcdefg-123",
		UserName: "mikey-new",
		Email:    "mikey@example.com",
	}, &User{
		ID:       "abcdefg-123",
		UserName: "mikey-newest",
		Email:    "mikey-new@example.com",
	})
	assert.NoError(t, err1)
	fbUserName, err2 := indexer.FindBy(&User{}, "UserName", "mikey-newest")
	assert.NoError(t, err2)
	assert.Len(t, fbUserName, 1)
	fbEmail, err3 := indexer.FindBy(&User{}, "Email", "mikey-new@example.com")
	assert.NoError(t, err3)
	assert.Len(t, fbEmail, 1)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_Disk_UpdateWithNonUniqueIndex(t *testing.T) {
	dataDir, err := WriteIndexTestData(Data, "ID", "")
	assert.NoError(t, err)
	indexer := createDiskIndexer(dataDir)

	err = indexer.AddIndex(&Pet{}, "Name", "ID", "pets", "non_unique", nil, false)
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func createCs3Indexer() *Indexer {
	return CreateIndexer(&config.Config{
		Repo: config.Repo{
			CS3: config.CS3{
				ProviderAddr: "0.0.0.0:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
	})
}

func createDiskIndexer(dataDir string) *Indexer {
	return CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})
}
