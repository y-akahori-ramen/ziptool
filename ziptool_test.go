package ziptool_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/y-akahori-ramen/ziptool"
)

func TestArchive(t *testing.T) {

	err := ziptool.Archive(filepath.Join("tmp", "singleArchive.zip"), filepath.Join("testdata", "singlefilesample.txt"))
	if err != nil {
		t.Fatal(err)
	}

	err = ziptool.Unarchive("tmp", filepath.Join("tmp", "singleArchive.zip"))
	if err != nil {
		t.Fatal(err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = ziptool.Archive(filepath.Join("tmp", "fileArchiveAbs.zip"), filepath.Join(curDir, "testdata"))
	if err != nil {
		t.Fatal(err)
	}

	err = ziptool.Unarchive(filepath.Join("tmp", "unarchive"), filepath.Join("tmp", "fileArchiveAbs.zip"))
	if err != nil {
		t.Fatal(err)
	}

}

func TestNotExistSrcArchive(t *testing.T) {
	err := ziptool.Archive("notExistSrc", "dst.zip")
	if err == nil {
		t.Fatal()
	}
	os.Remove("dst.zip")
}
