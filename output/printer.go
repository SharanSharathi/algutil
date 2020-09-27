package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Printer struct {
	writer io.Writer
	buf    *bufio.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		writer: w,
		buf:    bufio.NewWriter(w),
	}
}

func (p *Printer) Line(values ...interface{}) {
	fmt.Fprintln(p.buf, values...)
}

func (p *Printer) Format(format string, values ...interface{}) {
	fmt.Fprintf(p.buf, format, values...)
}

func (p *Printer) Slice(sl interface{}) {
	slice := fmt.Sprint(sl)
	slice = strings.TrimLeft(slice, " [")
	slice = strings.TrimRight(slice, " ]")

	fmt.Fprint(p.buf, slice)
}

func (p *Printer) Flush() {
	p.buf.Flush()
}

func (p *Printer) Tee(w io.Writer) {
	p.Flush()
	p.writer = io.MultiWriter(p.writer, w)
	p.buf.Reset(p.writer)
}

type FilePrinter struct {
	file *os.File
	*Printer
}

func NewFilePrinter(file string, append bool) *FilePrinter {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND

	if !append {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(file, flags, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &FilePrinter{
		file:    f,
		Printer: NewPrinter(f),
	}
}

func (fp *FilePrinter) Close() {
	fp.Flush()
	fp.file.Close()
	fp = nil
}
