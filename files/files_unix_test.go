// +build darwin freebsd linux netbsd openbsd

package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestExistsPerms(t *testing.T) {
	startDir, err := ioutil.TempDir("/tmp", "filestest_existsperms")
	check(t, err)
	base := startDir + "/TestExistsPerms"
	err = os.MkdirAll(base, 0777)
	check(t, err)

	fileName := base + "/bar"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0660)
	check(t, err)
	defer f.Close()

	e := Exists(fileName)
	if !e {
		t.Errorf(`%s : expected exists but got false`, fileName)
	}

	execCommand(t, "/bin/chmod", "000", fileName)
	execCommand(t, "/bin/chmod", "000", base)

	e2 := Exists(fileName)
	if e2 {
		t.Errorf(`%s : expected exists returning FALSE for an unattainable file`, fileName)
	}
}

var isSymlinkData = []maybeln{
	{"", false},
	{"   ", false},
	{"?|!", false},
	{".notfound", false},
	{".", false},
	{"testdata", false},
	{"testdata/", false},
	{"testdata/files", false},
	{"testdata/files/", false},
	{"testdata/files/01.txt", false},
	{"testdata/files/01.txt/foo", false},
	{"testdata/files/linkto01", true},
	{"testdata/files/../files/linkto01", true},
	{"testdata/files/./linkto01", true},
	{"testdata/files/sub", false},
	{"testdata/files/sub/", false},
}

func TestIsSymlink(t *testing.T) {
	for _, data := range isSymlinkData {
		is := IsSymlink(data.path)
		if is != data.isln {
			t.Errorf(`Expected IsSymlink=%t for path "%s"`, data.isln, data.path)
		}
	}
}
