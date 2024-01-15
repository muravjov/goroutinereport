package report

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"sort"
	"text/tabwriter"

	"github.com/maruel/panicparse/v2/stack"
	"github.com/muravjov/goroutinereport/internal"
)

type createItem struct {
	sig   *stack.Signature
	count int
	IDs   []int
}

type createList []createItem

func (a createList) Len() int      { return len(a) }
func (a createList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a createList) Less(i, j int) bool {
	return a[i].count < a[j].count
}

func Report(reader io.Reader, writer io.Writer) error {
	opts := stack.DefaultOpts()

	var snapshot *stack.Snapshot
	var snapshotSuffix []byte

	accumulate := func(s *stack.Snapshot) {
		if s == nil {
			return
		}

		if snapshot == nil {
			snapshot = s
		} else {
			snapshot.Goroutines = append(snapshot.Goroutines, s.Goroutines...)
		}
	}

	for {
		s, suffix, err := stack.ScanSnapshot(reader, io.Discard, opts)
		if err != nil {
			if err == io.EOF {
				accumulate(s)
				snapshotSuffix = suffix
				break
			}
			return err
		}

		reader = io.MultiReader(bytes.NewReader(suffix), reader)
		accumulate(s)
	}

	if snapshot == nil {
		return errors.New("no go stacks found")
	}

	if snapshot.IsRace() {
		return errors.New("not for race detector stacks")
	}

	aggregated := snapshot.Aggregate(stack.AnyValue)
	buckets := aggregated.Buckets

	cl := createList{}
	for _, b := range buckets {
		// fmtPrintln(sig.CreatedBy.FullSourceLine(), len(gList))

		cl = append(cl, createItem{
			&b.Signature,
			len(b.IDs),
			b.IDs,
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
	//fmtPrintf := func(format string, a ...interface{}) (n int, err error) {
	//	return fmt.Fprintf(writer, format, a...)
	//}

	for i := range cl {
		ci := cl[len(cl)-i-1]

		name := "UnknownFunc:UnknownLine"
		if n := internal.CreatedByString(ci.sig); n != "" {
			name = n
		}

		writeColumn("%s\t%d:%v", name, ci.count, ci.IDs)
	}
	w.Flush()

	if len(snapshotSuffix) > 0 {
		fmtPrintln()
		fmtPrintln("==== Suffix ================")
		fmtPrintln()

		fmtPrintln(string(snapshotSuffix))
	}

	fmtPrintln()
	fmtPrintln("=======================================")
	fmtPrintln()

	p := &internal.Palette{}
	var filter *regexp.Regexp
	var match *regexp.Regexp

	return internal.WriteBucketsToConsole(writer, p, aggregated, false, filter, match)
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
