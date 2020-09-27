package debug

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type Debugger struct {
	enabled bool
	output  io.Writer

	tracking  map[string]reflect.Value
	trackList []string
}

func NewDebugger(enabled bool) *Debugger {
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

func (d *Debugger) Trace(name string, pointer interface{}) {
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
