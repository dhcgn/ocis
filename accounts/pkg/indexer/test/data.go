package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// User is a user.
type User struct {
	ID, UserName, Email string
	UID                 int
}

// Pet is a pet.
type Pet struct {
	ID, Kind, Color, Name string
	UID                   int
}

// Data mock data.
var Data = map[string][]interface{}{
	"users": {
		User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"},
		User{ID: "hijklmn-456", UserName: "frank", Email: "frank@example.com"},
		User{ID: "ewf4ofk-555", UserName: "jacky", Email: "jacky@example.com"},
		User{ID: "rulan54-777", UserName: "jones", Email: "jones@example.com"},
	},
	"pets": {
		Pet{ID: "rebef-123", Kind: "Dog", Color: "Brown", Name: "Waldo"},
		Pet{ID: "wefwe-456", Kind: "Cat", Color: "White", Name: "Snowy"},
		Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"},
		Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"},
	},
}

// WriteIndexTestData writes mock data to disk.
func WriteIndexTestData(m map[string][]interface{}, privateKey, dir string) (string, error) {
	rootDir, err := getRootDir(dir)
	if err != nil {
		return "", err
	}

	err = writeData(m, privateKey, rootDir)
	if err != nil {
		return "", err
	}

	return rootDir, nil
}

// getRootDir allows for some minimal behavior on destination on disk. Testing the cs3 api behavior locally means
// keeping track of where the cs3 data lives on disk, this function allows for multiplexing whether or not to use a
// temporary folder or an already defined one.
func getRootDir(dir string) (string, error) {
	var rootDir string
	var err error

	if dir != "" {
		rootDir = dir
	} else {
		rootDir, err = CreateTmpDir()
		if err != nil {
			return "", err
		}
	}
	return rootDir, nil
}

// writeData writes test data to disk on rootDir location Marshaled as json.
func writeData(m map[string][]interface{}, privateKey string, rootDir string) error {
	for dirName := range m {
		fileTypePath := path.Join(rootDir, dirName)

		if err := os.MkdirAll(fileTypePath, 0777); err != nil {
			return err
		}
		for _, u := range m[dirName] {
			data, err := json.Marshal(u)
			if err != nil {
				return err
			}

			pkVal := ValueOf(u, privateKey)
			if err := ioutil.WriteFile(path.Join(fileTypePath, pkVal), data, 0777); err != nil {
				return err
			}
		}
	}
	return nil
}
