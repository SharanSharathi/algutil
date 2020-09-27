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

func (s *Scanner) readIfEmpty() {
	if len(s.line) == 0 {
		s.readNextLine()
	}
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

func (s *Scanner) SplitLineStrings() []string {
	splitted := bytes.Split(s.Line(), s.delim)
	result := make([]string, 0, len(splitted))

	for _, by := range splitted {
		str := string(by)

		if str != "" {
			result = append(result, str)
		}
	}

	return result
}

func (s *Scanner) SplitLineInts() []int {
	lineStrings := s.SplitLineStrings()
	lineInts := make([]int, 0, len(lineStrings))

	for _, str := range lineStrings {
		lineInts = append(lineInts, parseInt(str))
	}

	return lineInts
}

func (s *Scanner) SplitLineInt64s() []int64 {
	lineStrings := s.SplitLineStrings()
	lineInt64s := make([]int64, 0, len(lineStrings))

	for _, str := range lineStrings {
		lineInt64s = append(lineInt64s, parseInt64(str))
	}

	return lineInt64s
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

func parseInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (s *Scanner) skipPrefix(prefix []byte) {
	for bytes.HasPrefix(s.line, prefix) {
		s.line = s.line[len(prefix):]
		s.readIfEmpty()
	}
}

func (s *Scanner) NextString() string {
	s.readIfEmpty()
	s.skipPrefix(s.delim)

	ind := bytes.Index(s.line, s.delim)
	if ind == -1 {
		return s.LineString()
	}

	result := s.line[:ind]
	s.line = s.line[ind+len(s.delim):]

	return string(result)
}

func (s *Scanner) NextInt() int {
	return parseInt(s.NextString())
}

func (s *Scanner) NextInt64() int64 {
	return parseInt64(s.NextString())
}
