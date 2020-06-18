package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type FileParser interface {
	Next()
	HasNext() bool
	Token() Token
	Fatal(string, ...interface{})
}

type Token interface {
	Is(...string) bool
	Last() uint8
	Split(string) []string
	String() string
}

type fileParser struct {
	scanner  *bufio.Scanner
	tokens   []string
	token    string
	hasNext  bool
	row, col int
}

type token string

func (f *fileParser) Token() Token {
	return token(f.token)
}

func (f *fileParser) HasNext() bool {
	return f.hasNext
}

func (f *fileParser) Next() {
	f.col++
	if f.col >= len(f.tokens) {
		f.nextLine()
	}
	if f.hasNext {
		f.nextWord()
	}
}

func (f *fileParser) Fatal(format string, props ...interface{}) {
	message := fmt.Sprintf(format, props...)
	err := fmt.Errorf("at %d.%d: %s", f.row, f.col, message)
	log.Fatal(err)
}

func (f *fileParser) nextWord() {
	f.token = f.tokens[f.col]
	if f.token == "" {
		f.Next()
	}
}

func (f *fileParser) nextLine() {
	f.hasNext = f.scanner.Scan()
	if !f.hasNext {
		return
	}
	f.tokens = strings.Split(f.scanner.Text(), " ")
	f.row++
	f.col = 0
}

func (t token) Is(kinds ...string) bool {
	for _, kind := range kinds {
		if strings.HasPrefix(t.String(), kind) {
			return true
		}
	}
	return false
}

func (t token) Last() uint8 {
	return t[len(t)-1]
}

func (t token) Split(sep string) []string {
	return strings.Split(t.String(), sep)
}

func (t token) String() string {
	return string(t)
}
