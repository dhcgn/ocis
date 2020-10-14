package disk

import (
	"errors"
	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NonUniqueIndex is able to index an document by a key which might contain non-unique values
//
// /var/tmp/testfiles-395764020/index.disk/PetByColor/
// ├── Brown
// │   └── rebef-123 -> /var/tmp/testfiles-395764020/pets/rebef-123
// ├── Green
// │    ├── goefe-789 -> /var/tmp/testfiles-395764020/pets/goefe-789
// │    └── xadaf-189 -> /var/tmp/testfiles-395764020/pets/xadaf-189
// └── White
//     └── wefwe-456 -> /var/tmp/testfiles-395764020/pets/wefwe-456
type NonUniqueIndex struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string
}

func init() {
	registry.IndexConstructorRegistry["disk"]["non_unique"] = NewNonUniqueIndexWithOptions
}

// NewNonUniqueIndexWithOptions instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewNonUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &NonUniqueIndex{
		indexBy:      opts.IndexBy,
		typeName:     opts.TypeName,
		filesDir:     opts.FilesDir,
		indexBaseDir: path.Join(opts.DataDir, "index.disk"),
		indexRootDir: path.Join(path.Join(opts.DataDir, "index.disk"), strings.Join([]string{"non_unique", opts.TypeName, opts.IndexBy}, ".")),
	}
}

// Init initializes a unique index.
func (idx NonUniqueIndex) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

// Lookup exact lookup by value.
func (idx NonUniqueIndex) Lookup(v string) ([]string, error) {
	searchPath := path.Join(idx.indexRootDir, v)
	fi, err := ioutil.ReadDir(searchPath)
	if os.IsNotExist(err) {
		return []string{}, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	if err != nil {
		return []string{}, err
	}

	var ids []string = nil
	for _, f := range fi {
		ids = append(ids, f.Name())
	}

	if len(ids) == 0 {
		return []string{}, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return ids, nil
}

// Add adds a value to the index, returns the path to the root-document
func (idx NonUniqueIndex) Add(id, v string) (string, error) {
	oldName := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v, id)

	if err := os.MkdirAll(path.Join(idx.indexRootDir, v), 0777); err != nil {
		return "", err
	}

	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err

}

// Remove a value v from an index.
func (idx NonUniqueIndex) Remove(id string, v string) error {
	res, err := filepath.Glob(path.Join(idx.indexRootDir, "/*/", id))
	if err != nil {
		return err
	}

	for _, p := range res {
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	// Remove value directory if it is empty
	valueDir := path.Join(idx.indexRootDir, v)
	fi, err := ioutil.ReadDir(valueDir)
	if err != nil {
		return err
	}

	if len(fi) == 0 {
		if err := os.RemoveAll(valueDir); err != nil {
			return err
		}
	}

	return nil
}

// Update index from <oldV> to <newV>.
func (idx NonUniqueIndex) Update(id, oldV, newV string) (err error) {
	oldDir := path.Join(idx.indexRootDir, oldV)
	oldPath := path.Join(oldDir, id)
	newDir := path.Join(idx.indexRootDir, newV)
	newPath := path.Join(newDir, id)

	if _, err = os.Stat(oldPath); os.IsNotExist(err) {
		return &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: oldV}
	}

	if err != nil {
		return
	}

	if err = os.MkdirAll(newDir, 0777); err != nil {
		return
	}

	if err = os.Rename(oldPath, newPath); err != nil {
		return
	}

	di, err := ioutil.ReadDir(oldDir)
	if err != nil {
		return err
	}

	if len(di) == 0 {
		err = os.RemoveAll(oldDir)
		if err != nil {
			return
		}
	}

	return

}

// Search allows for glob search on the index.
func (idx NonUniqueIndex) Search(pattern string) ([]string, error) {
	paths, err := filepath.Glob(path.Join(idx.indexRootDir, pattern, "*"))
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: pattern}
	}

	return paths, nil
}

// IndexBy undocumented.
func (idx NonUniqueIndex) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx NonUniqueIndex) TypeName() string {
	return idx.typeName
}

// FilesDir undocumented.
func (idx NonUniqueIndex) FilesDir() string {
	return idx.filesDir
}