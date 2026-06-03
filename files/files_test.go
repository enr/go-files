package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
	{"testdata/files/01.txt/test", false},
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
	outputDir, err := ioutil.TempDir("", "filestest_copy")
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
	destinationDir, err := ioutil.TempDir("", "filestest_copyindir")
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
	{"testdata/files/01.txt/foo", false},
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
	path string
	isln bool
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

func TestBrokenLink(t *testing.T) {
	sourceFile := "testdata/files/brokenlink"
	ln := isSymlink(sourceFile)
	if !ln {
		t.Errorf("broken link %s expected to be link", sourceFile)
	}
}

// ---------------------------------------------------------------------------
// ExistsWithError
// ---------------------------------------------------------------------------

func TestExistsWithError(t *testing.T) {
	ok, err := ExistsWithError("testdata/files/01.txt")
	if !ok || err != nil {
		t.Errorf("ExistsWithError(existing): want (true, nil), got (%v, %v)", ok, err)
	}

	ok, err = ExistsWithError(".notfound")
	if ok || err == nil {
		t.Errorf("ExistsWithError(missing): want (false, non-nil), got (%v, %v)", ok, err)
	}

	// A path containing a null byte triggers a non-IsNotExist OS error, exercising
	// the fallthrough branch in existsWithError.
	ok, err = ExistsWithError("a\x00b")
	if ok || err == nil {
		t.Errorf("ExistsWithError(null-byte path): want (false, non-nil), got (%v, %v)", ok, err)
	}
}

// ---------------------------------------------------------------------------
// Sha1Sum error path
// ---------------------------------------------------------------------------

func TestSha1SumError(t *testing.T) {
	_, err := Sha1Sum(".notfound")
	if err == nil {
		t.Error("Sha1Sum(.notfound): expected error, got nil")
	}
}

// ---------------------------------------------------------------------------
// ReadLines error paths
// ---------------------------------------------------------------------------

func TestReadLinesError(t *testing.T) {
	_, err := ReadLines(".notfound")
	if err == nil {
		t.Error("ReadLines(.notfound): expected error, got nil")
	}
}

func TestReadLinesWhitespacePath(t *testing.T) {
	// cleanPath strips surrounding spaces — the trimmed path must resolve correctly.
	lines, err := ReadLines("  testdata/files/sub/03.txt  ")
	if err != nil {
		t.Fatalf("ReadLines with whitespace-padded path: unexpected error: %v", err)
	}
	if len(lines) != 5 {
		t.Errorf("ReadLines with whitespace-padded path: want 5 lines, got %d", len(lines))
	}
}

// ---------------------------------------------------------------------------
// EachLine error paths
// ---------------------------------------------------------------------------

func TestEachLineError(t *testing.T) {
	err := EachLine(".notfound", func(line string) error { return nil })
	if err == nil {
		t.Error("EachLine(.notfound): expected error, got nil")
	}
}

func TestEachLineCallbackError(t *testing.T) {
	sentinel := fmt.Errorf("stop")
	count := 0
	err := EachLine("testdata/files/sub/03.txt", func(line string) error {
		count++
		if count == 2 {
			return sentinel
		}
		return nil
	})
	if err != sentinel {
		t.Errorf("EachLine callback error: want sentinel error, got %v", err)
	}
	if count != 2 {
		t.Errorf("EachLine callback error: want 2 lines visited, got %d", count)
	}
}

func TestEachLineWhitespacePath(t *testing.T) {
	count := 0
	err := EachLine("  testdata/files/sub/03.txt  ", func(_ string) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("EachLine with whitespace-padded path: unexpected error: %v", err)
	}
	if count != 5 {
		t.Errorf("EachLine with whitespace-padded path: want 5 lines, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// IsSamePath — false case
// ---------------------------------------------------------------------------

func TestIsSamePathFalse(t *testing.T) {
	if IsSamePath("testdata/files/01.txt", "testdata/files/02.txt") {
		t.Error("IsSamePath: distinct files reported as same path")
	}
}

// ---------------------------------------------------------------------------
// CopyDir
// ---------------------------------------------------------------------------

func makeTempDir(t *testing.T, prefix string) string {
	t.Helper()
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("makeTempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		t.Fatalf("writeFile(%s): %v", path, err)
	}
	if _, err = fmt.Fprint(f, content); err != nil {
		f.Close()
		t.Fatalf("writeFile(%s): write: %v", path, err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("writeFile(%s): close: %v", path, err)
	}
}

func TestCopyDir(t *testing.T) {
	src := makeTempDir(t, "filestest_copydir_src")
	dst := makeTempDir(t, "filestest_copydir_dst")

	// src/
	//   a.txt
	//   sub/
	//     b.txt
	subDir := filepath.Join(src, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(src, "a.txt"), "hello", 0644)
	writeFile(t, filepath.Join(subDir, "b.txt"), "world", 0644)

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: unexpected error: %v", err)
	}

	for _, rel := range []string{"a.txt", filepath.Join("sub", "b.txt")} {
		if !Exists(filepath.Join(dst, rel)) {
			t.Errorf("CopyDir: expected %s to exist in dst", rel)
		}
	}
}

func TestCopyDirPreservesFilePermissions(t *testing.T) {
	src := makeTempDir(t, "filestest_copydir_perms_src")
	dst := makeTempDir(t, "filestest_copydir_perms_dst")

	writeFile(t, filepath.Join(src, "exec.sh"), "#!/bin/sh\n", 0755)

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	fi, err := os.Stat(filepath.Join(dst, "exec.sh"))
	if err != nil {
		t.Fatalf("stat dst/exec.sh: %v", err)
	}
	if fi.Mode()&0111 == 0 {
		t.Errorf("CopyDir: execute bit lost; mode = %v", fi.Mode())
	}
}

func TestCopyDirError(t *testing.T) {
	err := CopyDir(".notfound", "/tmp/doesnotmatter")
	if err == nil {
		t.Error("CopyDir(.notfound): expected error, got nil")
	}
}

func TestCopyDirFileError(t *testing.T) {
	// A broken symlink inside the source causes Copy to fail, exercising the
	// error-break path inside CopyDir's file-copy branch.
	src := makeTempDir(t, "filestest_copydir_ferr_src")
	dst := makeTempDir(t, "filestest_copydir_ferr_dst")

	if err := os.Symlink("/nonexistent/target", filepath.Join(src, "broken")); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	if err := CopyDir(src, dst); err == nil {
		t.Error("CopyDir with broken symlink: expected error, got nil")
	}
}

func TestCopyDirNested(t *testing.T) {
	src := makeTempDir(t, "filestest_copydir_nested_src")
	dst := makeTempDir(t, "filestest_copydir_nested_dst")

	depth := filepath.Join(src, "a", "b", "c")
	if err := os.MkdirAll(depth, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(depth, "deep.txt"), "deep", 0644)

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir nested: %v", err)
	}

	expected := filepath.Join(dst, "a", "b", "c", "deep.txt")
	if !Exists(expected) {
		t.Errorf("CopyDir nested: expected %s to exist", expected)
	}
}
