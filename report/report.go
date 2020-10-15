package report

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"text/tabwriter"

	"github.com/maruel/panicparse/stack"
)

type createItem struct {
	sig   *stack.Signature
	count int
}

type createList []createItem

func (a createList) Len() int      { return len(a) }
func (a createList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a createList) Less(i, j int) bool {
	return a[i].count < a[j].count
}

func Report(reader io.Reader, writer io.Writer) error {
	goroutines, err := stack.ParseDump(reader, ioutil.Discard, false)
	if err != nil {
		return err
	}

	buckets := stack.Aggregate(goroutines.Goroutines, stack.AnyValue)

	cl := createList{}
	for _, b := range buckets {
		// fmtPrintln(sig.CreatedBy.FullSourceLine(), len(gList))

		cl = append(cl, createItem{
			&b.Signature,
			len(b.IDs),
		})
	}

	sort.Sort(cl)

	const padding = 3
	var twFlags uint = 0 // tabwriter.Debug
	w := tabwriter.NewWriter(writer, 0, 0, padding, ' ', twFlags)
	writeColumn := func(format string, a ...interface{}) {
		fmt.Fprintln(w, fmt.Sprintf(format, a...))
	}

	writeColumn("Created By\tGoroutine Count")
	// writeColumn("")
	// w.Write([]byte("\n"))
	writeColumn("\t")

	fmtPrintln := func(a ...interface{}) (n int, err error) {
		return fmt.Fprintln(writer, a...)
	}
	fmtPrintf := func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(writer, format, a...)
	}

	for i, _ := range cl {
		ci := cl[len(cl)-i-1]
		// fmtPrintf("%s:%s\t\t%d\n", ci.sig.CreatedBy.Func.Name(), ci.sig.CreatedBy.FullSrcLine(), ci.count)
		writeColumn("%s:%s\t%d", ci.sig.CreatedBy.Func.Name(), ci.sig.CreatedBy.FullSrcLine(), ci.count)
	}
	w.Flush()

	fmtPrintln()
	fmtPrintln("=======================================")
	fmtPrintln()

	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.FullSrcLine()); l > srcLen {
				srcLen = l
			}
			if l := len(line.Func.PkgName()); l > pkgLen {
				pkgLen = l
			}
		}
	}

	for i, _ := range cl {
		ci := cl[len(cl)-i-1]

		// Print the goroutine header.
		extra := ""
		if s := ci.sig.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if ci.sig.Locked {
			extra += " [locked]"
		}
		if c := ci.sig.CreatedByString(false); c != "" {
			extra += " [Created by " + c + "]"
		}
		fmtPrintf("%d: %s%s\n", ci.count, ci.sig.State, extra)

		// Print the stack lines.
		for _, line := range ci.sig.Stack.Calls {
			fmtPrintf(
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.FullSrcLine(),
				line.Func.Name(), &line.Args)
		}
		if ci.sig.Stack.Elided {
			io.WriteString(os.Stdout, "    (...)\n")
		}
	}
	return nil
}

func ReportSelf(writer io.Writer) error {
	// :COPY_N_PASTE: runtime/pprof/pprof.go:writeGoroutineStacks

	// We don't know how big the buffer needs to be to collect
	// all the goroutines. Start with 1 MB and try a few times, doubling each time.
	// Give up and use a truncated trace if 64 MB is not enough.
	buf := make([]byte, 1<<20)
	for i := 0; ; i++ {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		if len(buf) >= 64<<20 {
			// Filled 64 MB - stop there.
			break
		}
		buf = make([]byte, 2*len(buf))
	}

	return Report(bytes.NewReader(buf), writer)
}
