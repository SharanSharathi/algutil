package debugger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

type Debugger struct {
	enabled bool
	output  io.Writer

	tracking  map[string]reflect.Value
	trackList []string
}

func New(enabled bool) *Debugger {
	return &Debugger{
		enabled:  enabled,
		output:   os.Stdout,
		tracking: make(map[string]reflect.Value),
	}
}

func (d *Debugger) Print(values ...interface{}) {
	if !d.enabled {
		return
	}

	fmt.Fprintln(d.output, values...)
	d.printTracked()
}

func (d *Debugger) Printf(format string, values ...interface{}) {
	if !d.enabled {
		return
	}

	fmt.Fprintf(d.output, format, values...)
	d.printTracked()
}

func (d *Debugger) PrintTrace() {
	if !d.enabled {
		return
	}

	fmt.Fprintln(d.output, getTrace())
}

func (d *Debugger) Track(name string, pointer interface{}) {
	ptr := reflect.ValueOf(pointer)
	if ptr.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("%T is not a pointer, expecting to pass pointer", pointer))
	}

	if _, yes := d.tracking[name]; !yes {
		d.trackList = append(d.trackList, name)
	}

	d.tracking[name] = ptr
}

func (d *Debugger) Unlink(name string) {
	delete(d.tracking, name)

	var pos int

	for i, N := range d.trackList {
		if N == name {
			pos = i
			break
		}
	}

	d.trackList = append(d.trackList[:pos], d.trackList[pos+1:]...)
}

func (d *Debugger) UnlinkAll() {
	d.trackList = d.trackList[:0]
	d.tracking = make(map[string]reflect.Value)
}

func (d *Debugger) SetOutput(output io.Writer) {
	d.output = output
}

func (d *Debugger) printTracked() {
	if len(d.trackList) == 0 {
		return
	}

	found := false
	maxLen := d.maxNameLength()

	for _, name := range d.trackList {
		ptr, ok := d.tracking[name]
		if !ok {
			continue
		}

		if ptr.IsNil() {
			continue
		}

		found = true

		fmt.Fprintf(d.output, "    %-[1]*[2]s : %v\n", maxLen, name,
			reflect.Indirect(ptr))
	}

	if found {
		fmt.Fprintln(d.output)
	}
}

func (d *Debugger) maxNameLength() (maxLen int) {
	for name, ptr := range d.tracking {
		if ptr.IsNil() {
			continue
		}

		if l := len(name); l > maxLen {
			maxLen = l
		}
	}

	return
}

const maxCallerDepth = 100

func getTrace() string {
	pc := make([]uintptr, maxCallerDepth+1)

	n := runtime.Callers(3, pc)
	if n > 0 {
		n--
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	buf := new(bytes.Buffer)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		buf.WriteByte('\n')

		buf.WriteString(frame.Function)
		buf.WriteString("\n    ")
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
	}

	buf.WriteByte('\n')

	return buf.String()
}
