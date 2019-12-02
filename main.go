package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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

func report(f *os.File) {
	goroutines, err := stack.ParseDump(f, ioutil.Discard, false)
	if err != nil {
		log.Println(err)
		return
	}

	buckets := stack.Aggregate(goroutines.Goroutines, stack.AnyValue)

	cl := createList{}
	for _, b := range buckets {
		// fmt.Println(sig.CreatedBy.FullSourceLine(), len(gList))

		cl = append(cl, createItem{
			&b.Signature,
			len(b.IDs),
		})
	}

	sort.Sort(cl)

	const padding = 3
	var twFlags uint = 0 // tabwriter.Debug
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', twFlags)
	writeColumn := func(format string, a ...interface{}) {
		fmt.Fprintln(w, fmt.Sprintf(format, a...))
	}

	writeColumn("Created By\tGoroutine Count")
	// writeColumn("")
	// w.Write([]byte("\n"))
	writeColumn("\t")

	for i, _ := range cl {
		ci := cl[len(cl)-i-1]
		// fmt.Printf("%s:%s\t\t%d\n", ci.sig.CreatedBy.Func.Name(), ci.sig.CreatedBy.FullSrcLine(), ci.count)
		writeColumn("%s:%s\t%d", ci.sig.CreatedBy.Func.Name(), ci.sig.CreatedBy.FullSrcLine(), ci.count)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("=======================================")
	fmt.Println()

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
		fmt.Printf("%d: %s%s\n", ci.count, ci.sig.State, extra)

		// Print the stack lines.
		for _, line := range ci.sig.Stack.Calls {
			fmt.Printf(
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.FullSrcLine(),
				line.Func.Name(), &line.Args)
		}
		if ci.sig.Stack.Elided {
			io.WriteString(os.Stdout, "    (...)\n")
		}
	}

}

func main() {
	file := os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		file = f
	}
	report(file)
}
