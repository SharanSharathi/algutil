package input

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type Scanner struct {
	line   []byte
	delim  []byte
	reader *bufio.Reader

	reachedEOF bool
}

func NewScanner(r io.Reader, delim string) *Scanner {
	return &Scanner{
		delim:  []byte(delim),
		reader: bufio.NewReader(r),
	}
}

func (s *Scanner) readIfEmpty() {
	if len(s.line) == 0 {
		s.readNextLine()
	}
}

func (s *Scanner) readNextLine() {
	line, err := s.reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			panic(err)
		}

		s.reachedEOF = true
		return
	}

	length := len(line)
	drop := 0

	if length > 0 && line[length-1] == '\n' {
		drop = 1

		if length > 1 && line[length-2] == '\r' {
			drop = 2
		}
	}

	s.line = line[:length-drop]
}

func (s *Scanner) skipPrefix(prefix []byte) {
	if len(prefix) == 0 {
		return
	}

	for bytes.HasPrefix(s.line, prefix) {
		s.line = s.line[len(prefix):]
		s.readIfEmpty()
	}
}

func (s *Scanner) Next() []byte {
	s.readIfEmpty()
	s.skipPrefix(s.delim)

	if len(s.line) == 0 {
		return nil
	}

	if len(s.delim) == 0 { // no delim; return one byte
		result := []byte{s.line[0]}
		s.line = s.line[1:]

		return result
	}

	ind := bytes.Index(s.line, s.delim)
	if ind == -1 {
		return s.Line()
	}

	result := s.line[:ind]
	s.line = s.line[ind+len(s.delim):]

	return result
}

func (s *Scanner) NextString() string {
	return string(s.Next())
}

func (s *Scanner) NextInt() int {
	return parseInt(s.NextString())
}

func (s *Scanner) NextUint() uint {
	return parseUint(s.NextString())
}

func (s *Scanner) NextInt64() int64 {
	return parseInt64(s.NextString())
}

func (s *Scanner) NextUint64() uint64 {
	return parseUint64(s.NextString())
}

func (s *Scanner) Line() []byte {
	s.readIfEmpty()

	line := s.line
	s.line = nil

	return line
}

func (s *Scanner) LineString() string {
	return string(s.Line())
}

func (s *Scanner) SliceOfBytes() [][]byte {
	splitted := bytes.Split(s.Line(), s.delim)
	bytesSlice := make([][]byte, 0, len(splitted))

	for _, bytes := range splitted {
		if len(bytes) > 0 {
			bytesSlice = append(bytesSlice, bytes)
		}
	}

	return bytesSlice
}

func (s *Scanner) SliceOfStrings() []string {
	splitted := bytes.Split(s.Line(), s.delim)
	stringSlice := make([]string, 0, len(splitted))

	for _, bytes := range splitted {
		if len(bytes) > 0 {
			stringSlice = append(stringSlice, string(bytes))
		}
	}

	return stringSlice
}

func (s *Scanner) SliceOfInts() []int {
	stringSlice := s.SliceOfStrings()
	intSlice := make([]int, len(stringSlice))

	for i, str := range stringSlice {
		intSlice[i] = parseInt(str)
	}

	return intSlice
}

func (s *Scanner) SliceOfUints() []uint {
	stringSlice := s.SliceOfStrings()
	uintSlice := make([]uint, len(stringSlice))

	for i, str := range stringSlice {
		uintSlice[i] = parseUint(str)
	}

	return uintSlice
}

func (s *Scanner) SliceOfInt64s() []int64 {
	stringSlice := s.SliceOfStrings()
	int64Slice := make([]int64, len(stringSlice))

	for i, str := range stringSlice {
		int64Slice[i] = parseInt64(str)
	}

	return int64Slice
}

func (s *Scanner) SliceOfUint64s() []uint64 {
	stringSlice := s.SliceOfStrings()
	uint64Slice := make([]uint64, len(stringSlice))

	for i, str := range stringSlice {
		uint64Slice[i] = parseUint64(str)
	}

	return uint64Slice
}

func (s *Scanner) ReachedEOF() bool {
	return s.reachedEOF
}

func parseInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return i
}

func parseUint(str string) uint {
	i, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		panic(err)
	}

	return uint(i)
}

func parseInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func parseUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}
