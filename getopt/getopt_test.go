package getopt_test

import (
	"os"
	"testing"

	"github.com/sverrehu/goutils/getopt"
)

func TestMostOfIt(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "before", "-n-10", "-m", "12", "--number-two=13", "--help", "--text", "foo", "--", "after", "-even-after"}
	var n int
	var m uint
	var n2 uint
	var help bool
	var text string
	opts := []getopt.Option{
		{ShortName: 'h', LongName: "help", Type: getopt.Flag, Target: &help},
		{ShortName: 'n', LongName: "", Type: getopt.Integer, Target: &n},
		{ShortName: 'm', LongName: "", Type: getopt.UInteger, Target: &m},
		{ShortName: 0, LongName: "number-two", Type: getopt.UInteger, Target: &n2},
		{ShortName: 't', LongName: "text", Type: getopt.String, Target: &text},
	}
	getopt.Parse(&os.Args, opts, false)
	if got, want := n, -10; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := m, uint(12); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := n2, uint(13); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := help, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := text, "foo"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := len(os.Args), 4; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := os.Args[0], "cmd"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := os.Args[1], "before"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := os.Args[2], "after"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := os.Args[3], "-even-after"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
