package internal

import (
	"io"
	"regexp"

	"github.com/maruel/panicparse/v2/stack"
)

const myPathFormat = relPath // fullPath //

func WriteBucketsToConsole(out io.Writer, p *Palette, a *stack.Aggregated, needsEnv bool, filter, match *regexp.Regexp) error {
	pf := myPathFormat

	if needsEnv {
		_, _ = io.WriteString(out, "\nTo see all goroutines, visit https://github.com/maruel/panicparse#gotraceback\n\n")
	}
	srcLen, pkgLen := calcBucketsLengths(a, pf)
	multi := len(a.Buckets) > 1
	for _, e := range a.Buckets {
		header := p.BucketHeader(e, pf, multi)
		if filter != nil && filter.MatchString(header) {
			continue
		}
		if match != nil && !match.MatchString(header) {
			continue
		}
		_, _ = io.WriteString(out, header)
		_, _ = io.WriteString(out, p.StackLines(&e.Signature, srcLen, pkgLen, pf))
	}
	return nil
}

func CreatedByString(s *stack.Signature) string {
	return myPathFormat.createdByString(s)
}
