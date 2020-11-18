package testerprogress

import (
	"fmt"
	"io"
	"strings"
)

const testerProgressBarDesc = "\r[%-50s]%3d%% %8d/%d"
const percentWidth = 50

type TesterProgress struct {
	dest    io.Writer
	total   int
	current int
	graph   rune
}

func New(dest io.Writer, total int) *TesterProgress {
	if total <= 0 {
		total = 1
	}
	return &TesterProgress{
		dest:    dest,
		total:   total,
		current: 0,
		graph:   '.',
	}
}

func (tp *TesterProgress) Step() {
	if tp.current == tp.total {
		return
	}

	tp.current = tp.current + 1
	percent := float64(tp.current) / float64(tp.total)
	barString := strings.Repeat(string(tp.graph), int(percent*percentWidth))

	fmt.Fprintf(tp.dest, testerProgressBarDesc, barString, int(percent*100), tp.current, tp.total)
	if tp.current >= tp.total {
		fmt.Fprintln(tp.dest)
	}
}
