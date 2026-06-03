package files

import (
	"testing"
)

type samepathTestCase struct {
	p1     string
	p2     string
	equals bool
}

var testcases = []samepathTestCase{
	{"", " ", true},
	{".notfound", "../files/.notfound", true},
	{".notfound", `..\files\.notfound`, true},
	{".", "../files", true},
	{"testdata/", "./testdata", true},
	// negative cases
	{"testdata/files/01.txt", "testdata/files/02.txt", false},
	{"testdata", "testdata/files", false},
}

func TestSamePath(t *testing.T) {
	for _, data := range testcases {
		res := IsSamePath(data.p1, data.p2)
		if res != data.equals {
			t.Errorf(`Expected IsSamePath=%t for paths "%s" and "%s"`, data.equals, data.p1, data.p2)
		}
	}
}
