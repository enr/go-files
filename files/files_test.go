package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type maybedir struct {
	path  string
	isdir bool
}

var isDirData = []maybedir{
	{"", false},
	{"   ", false},
	{"?|!", false},
	{".notfound", false},
	{".", true},
	{"testdata", true},
	{"testdata/", true},
	{"testdata/files", true},
	{"testdata/files/", true},
	{"testdata/files/01.txt", false},
	{"testdata/files/linkto01", false},
	{"testdata/files/sub", true},
	{"testdata/files/sub/", true},
}

func TestIsDir(t *testing.T) {
	for _, data := range isDirData {
		is := IsDir(data.path)
		if is != data.isdir {
			t.Errorf(`Expected IsDir=%t for path "%s"`, data.isdir, data.path)
		}
	}
}

type testfile struct {
	path    string
	sha1sum string
}

var testfiles = []testfile{
	{"testdata/files/01.txt", "89c47433ed8741caf3b6747c18e0d242b0d39993"},
	{"testdata/files/02.txt", "45981845bb1ab6c784bfd781bddde5fb70b57151"},
	{"testdata/files/sub/03.txt", "c51fce748bb1654be53575aa244de59fcf63f18c"},
}

func TestSha1Sum(t *testing.T) {
	for _, data := range testfiles {
		sha1sum, err := Sha1Sum(data.path)
		if err != nil {
			t.Errorf("error in Sha1Sum(%s): %s %s", data.path, reflect.TypeOf(err), err.Error())
		}
		if sha1sum != data.sha1sum {
			t.Errorf(`%s : expected sha1sum "%s" but got "%s"`, data.path, data.sha1sum, sha1sum)
		}
	}
}

func TestCopy(t *testing.T) {
	outputDir, err := ioutil.TempDir("/tmp", "filestest_copy")
	check(t, err)
	for _, data := range testfiles {
		of := fmt.Sprintf("%s/%s", outputDir, filepath.Base(data.path))
		deleteFile(of, t)
		err := Copy(data.path, of)
		if err != nil {
			t.Errorf("error in Copy(%s, %s): %s %s", data.path, of, reflect.TypeOf(err), err.Error())
		}
		sha1sum, err := Sha1Sum(of)
		if err != nil {
			t.Errorf("error in Sha1Sum(%s): %s %s", of, reflect.TypeOf(err), err.Error())
		}
		if sha1sum != data.sha1sum {
			t.Errorf(`%s : expected sha1sum "%s" but got "%s"`, of, data.sha1sum, sha1sum)
		}
	}
}

type copyerrorargs struct {
	source      string
	destination string
}

var copyErrorData = []copyerrorargs{
	{"", ""},
	{"   ", ""},
	{"not_here", "test.txt"},
	{"testdata/not_here.txt", "test.txt"},
	{"testdata/files/01.txt", "not/a/dir/01.txt"},
	{"testdata", "testdata.txt"},
}

func TestCopyError(t *testing.T) {
	for _, data := range copyErrorData {
		err := Copy(data.source, data.destination)
		if err == nil {
			t.Errorf("expected error in Copy(%s, %s) but got nil", data.source, data.destination)
		}
		if Exists(data.destination) {
			t.Errorf("Copy(%s, %s) created destination file", data.source, data.destination)
		}
	}
}

func TestCopyInDir(t *testing.T) {
	sourceFile := "testdata/files/01.txt"
	destinationDir, err := ioutil.TempDir("/tmp", "filestest_copyindir")
	check(t, err)
	expectedFile := fmt.Sprintf("%s/01.txt", destinationDir)

	err = Copy(sourceFile, destinationDir)
	if err != nil {
		t.Errorf("error in Copy(%s, %s)", sourceFile, destinationDir)
	}
	if !Exists(expectedFile) {
		t.Errorf("Copy(%s, %s) : no expected file %s", sourceFile, destinationDir, expectedFile)
	}
}

type maybeexists struct {
	path   string
	exists bool
}

var existsData = []maybeexists{
	{"", false},
	{"   ", false},
	{"?|!", false},
	{".notfound", false},
	{".", true},
	{"    . ", true},
	{"testdata", true},
	{"testdata/", true},
	{"testdata/files", true},
	{"testdata/files/", true},
	{"testdata/files/01.txt", true},
	{"testdata/files/02.txt", true},
	{"testdata/files/linkto01", true},
	{"testdata/files/sub", true},
	{"testdata/files/sub/", true},
	{"testdata/files/sub/03", false},
	{"testdata/files/sub/03.txt", true},
	{"./testdata/files/sub/03.txt", true},
	{"../files/./testdata/files/sub/03.txt", true},
}

func TestExists(t *testing.T) {
	for _, data := range existsData {
		e := Exists(data.path)
		if e != data.exists {
			t.Errorf(`%s : expected exists "%t"`, data.path, data.exists)
		}
	}
}

func TestExistsPerms(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	startDir, err := ioutil.TempDir("/tmp", "filestest_existsperms")
	check(t, err)
	base := startDir + "/TestExistsPerms"
	err = os.MkdirAll(base, 0777)
	check(t, err)

	fileName := base + "/bar"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0660)
	check(t, err)
	defer f.Close()

	execCommand(t, "/bin/chmod", "777", fileName)
	execCommand(t, "/bin/chmod", "777", base)

	e := Exists(fileName)
	if !e {
		t.Errorf(`%s : expected exists but got false`, fileName)
	}

	execCommand(t, "/bin/chmod", "000", fileName)
	execCommand(t, "/bin/chmod", "000", base)

	e2 := Exists(fileName)
	if !e2 {
		t.Errorf(`%s : expected exists but got false`, fileName)
	}
}

func TestIsAccessible(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	startDir, err := ioutil.TempDir("/tmp", "filestest_isaccessible")
	check(t, err)
	base := startDir + "/TestIsAccessible"
	err = os.MkdirAll(base, 0777)
	check(t, err)

	fileName := base + "/bar"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0660)
	check(t, err)
	defer f.Close()

	execCommand(t, "/bin/chmod", "777", fileName)
	execCommand(t, "/bin/chmod", "777", base)
	a1 := IsAccessible(fileName)
	if !a1 {
		t.Errorf(`%s : expected accessible but got false`, fileName)
	}

	execCommand(t, "/bin/chmod", "000", fileName)
	execCommand(t, "/bin/chmod", "000", base)

	a2 := IsAccessible(fileName)
	if a2 {
		t.Errorf(`%s : expected not accessible but got true`, fileName)
	}
}

type mayberegs struct {
	path string
	reg  bool
}

var regData = []maybeexists{
	{"", false},
	{"   ", false},
	{"?|!", false},
	{".notfound", false},
	{".", false},
	{"    . ", false},
	{"testdata", false},
	{"testdata/", false},
	{"testdata/files", false},
	{"testdata/files/", false},
	{"testdata/files/01.txt", true},
	{"testdata/files/02.txt", true},
	{"testdata/files/linkto01", true},
	{"testdata/files/sub/03", false},
	{"testdata/files/sub/03.txt", true},
	{"./testdata/files/sub/03.txt", true},
	{"../files/testdata/files/sub/03.txt", true},
}

func TestIsRegular(t *testing.T) {
	for _, data := range regData {
		e := IsRegular(data.path)
		if e != data.exists {
			t.Errorf(`%s : expected regular "%t"`, data.path, data.exists)
		}
	}
}

func TestReadLines(t *testing.T) {
	path := "testdata/files/sub/03.txt"
	lines, err := ReadLines(path)
	if err != nil {
		t.Errorf("error reading lines from %s", path)
	}
	if len(lines) != 5 {
		t.Errorf("ReadLines(%s), expected %d lines but got %d", path, 5, len(lines))
	}
	filelines := []string{
		"Hi, my name is 03.",
		"",
		"I am multi...",
		"...",
		"lines!",
	}
	for index, actual := range lines {
		expected := filelines[index]
		if actual != expected {
			t.Errorf(`ReadLines(%s), line %d expected %q but got %q`, path, index, expected, actual)
		}
	}
}

func TestEachLine(t *testing.T) {
	path := "testdata/files/sub/03.txt"
	filelines := []string{}
	EachLine(path, func(line string) error {
		filelines = append(filelines, line)
		return nil
	})
	if len(filelines) != 5 {
		t.Errorf("EachLine(%s), expected %d lines but got %d", path, 5, len(filelines))
	}
	expectedlines := []string{
		"Hi, my name is 03.",
		"",
		"I am multi...",
		"...",
		"lines!",
	}
	for index, actual := range filelines {
		expected := expectedlines[index]
		if actual != expected {
			t.Errorf(`EachLine(%s), line %d expected %q but got %q`, path, index, expected, actual)
		}
	}
}

type maybeln struct {
	path  string
	isln bool
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
	{"testdata/files/linkto01", true},
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

func deleteFile(path string, t *testing.T) {
	if Exists(path) {
		err := os.Remove(path)
		if err != nil {
			t.Error("error deleting path", path)
		}
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error %v", err)
	}
}

func execCommand(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("error executing command %v\n%s", err, o)
	}
}
