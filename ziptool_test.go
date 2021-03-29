package ziptool_test

import (
	"os"
	"testing"

	"github.com/y-akahori-ramen/ziptool"
)

func TestArchive(t *testing.T) {

	err := ziptool.Archive("tmp/singleArchive.zip", "testdata/singlefilesample.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = ziptool.Unarchive("tmp", "tmp/singleArchive.zip")
	if err != nil {
		t.Fatal(err)
	}

	err = ziptool.Archive("tmp/fileArchiveAbs.zip", "/Volumes/Data/Program/ziptool/testdata")
	if err != nil {
		t.Fatal(err)
	}

	err = ziptool.Unarchive("tmp/unarchive", "tmp/fileArchiveAbs.zip")
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
