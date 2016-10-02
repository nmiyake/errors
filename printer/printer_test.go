package printer_test

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/pkg/errors/printer"
)

func TestPrintSingleStack(t *testing.T) {
	cases := []struct{
		err error
		want[] string
	}{
		// regular error prints directly
		//{io.EOF, []string{"EOF"}},
		// multiple errors wrapped with "Wrap" coalesce into a single stack
		{loadConfig(), []string {
			"EOF",
			"failed to open foo",
			"failed to parse file",
			"failed to load config",
			"github.com/pkg/errors/printer_test.openFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfig",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintSingleStack",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error contains multiple stacks that don't coalesce, prints "%+v" output
		{loadConfigInChannel(), []string{
			"EOF",
			"failed to open foo in channel",
			"github.com/pkg/errors/printer_test.openFileInChannel.func1",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/pkg/errors/printer_test.parseFileInChannel",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigInChannel",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintSingleStack",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t/usr/local/go/src/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t/usr/local/go/src/runtime/asm_amd64.s:[0-9]+",
			"failed to load file",
			"github.com/pkg/errors/printer_test.loadConfigInChannel",
			"\t/Volumes/git/go/src/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintSingleStack",
			"\t/Volumes/git/go/src/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t/usr/local/go/src/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t/usr/local/go/src/runtime/asm_amd64.s:[0-9]+",
		}},
	}

	for i, currCase := range cases {
		got := strings.Split(printer.PrintSingleStack(currCase.err), "\n")

		if len(got) != len(currCase.want) {
			t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want, got)
			continue
		}

		for j := range got {
			if !regexp.MustCompile("^" + currCase.want[j] + "$").MatchString(got[j]) {
				t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want[j], got[j])
				break
			}
		}
	}
}

func TestPrintStackWithMessages(t *testing.T) {
	cases := []struct{
		err error
		want[] string
	}{
		// regular error prints directly
		{io.EOF, []string{"EOF"}},
		// multiple errors wrapped with "Wrap" annotates stack with messages in correct spots
		{loadConfig(), []string {
			"EOF",
			"failed to open foo",
			"github.com/pkg/errors/printer_test.openFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"failed to parse file",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"failed to load config",
			"github.com/pkg/errors/printer_test.loadConfig",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error causes do not alternate between stackTracer and non-stackTracer, prints "%+v" output
		{loadConfigWithMsg(), []string {
			"EOF",
			"failed to open foo",
			"github.com/pkg/errors/printer_test.openFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigWithMsg",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigWithMsg",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to load config",
		}},
		// if error causes do not alternate between stackTracer and non-stackTracer, prints "%+v" output
		{loadConfigWithStack(), []string {
			"EOF",
			"failed to open foo",
			"github.com/pkg/errors/printer_test.openFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigWithStack",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/pkg/errors/printer_test.parseFile",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigWithStack",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"github.com/pkg/errors/printer_test.loadConfigWithStack",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error contains multiple stacks, all stacks are included and annotated at correct location
		{loadConfigInChannel(), []string{
			"EOF",
			"failed to open foo in channel",
			"github.com/pkg/errors/printer_test.openFileInChannel.func1",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/pkg/errors/printer_test.parseFileInChannel",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"failed to load file",
			"github.com/pkg/errors/printer_test.loadConfigInChannel",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"github.com/pkg/errors/printer_test.TestPrintStackWithMessages",
			"\t.+/github.com/pkg/errors/printer/printer_test.go:[0-9]+",
			"testing.tRunner",
			"\t/usr/local/go/src/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t/usr/local/go/src/runtime/asm_amd64.s:[0-9]+",
		}},
	}

	for i, currCase := range cases {
		got := strings.Split(printer.PrintStackWithMessages(currCase.err), "\n")

		if len(got) != len(currCase.want) {
			t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want, got)
			continue
		}

		for j := range got {
			if !regexp.MustCompile("^" + currCase.want[j] + "$").MatchString(got[j]) {
				t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want[j], got[j])
				break
			}
		}
	}
}

func loadConfig() error {
	return errors.Wrap(parseFile(), "failed to load config")
}

func loadConfigWithMsg() error {
	return errors.WithMessage(parseFile(), "failed to load config")
}

func loadConfigWithStack() error {
	return errors.WithStack(parseFile())
}

func parseFile() error {
	return errors.Wrap(openFile(), "failed to parse file")
}

func openFile() error {
	return errors.Wrap(io.EOF, "failed to open foo")
}

func loadConfigInChannel() error {
	return errors.Wrapf(parseFileInChannel(), "failed to load file")
}

func parseFileInChannel() error {
	return errors.Wrap(openFileInChannel(), "failed to parse file")
}

func openFileInChannel() error {
	errs := make(chan error)
	go func() { errs <- errors.Wrap(io.EOF, "failed to open foo in channel") }()
	return <-errs
}