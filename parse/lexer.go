package parse

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	d "github.com/JoergReinhardt/gatwd/data"
	f "github.com/JoergReinhardt/gatwd/functions"
)

//////////////////////////////////////////////////
type LineBuffer d.AsyncVal

func NewLineBuffer(callbacks ...func()) *LineBuffer {
	return (*LineBuffer)(d.NewAsync(callbacks...))
}
func (s LineBuffer) AsynVal() *d.AsyncVal   { return (*d.AsyncVal)(&s) }
func (s *LineBuffer) Subscribe(c ...func()) { (*d.AsyncVal)(s).Subscribe(c...) }
func (s LineBuffer) String() string {
	(&s).Lock()
	defer (&s).Unlock()
	return s.byteVec().String()
}
func (s LineBuffer) Lines() []string {
	return strings.Split(s.String(), "\n")
}
func (s LineBuffer) Fields() [][]string {
	var fields = [][]string{}
	for _, line := range s.Lines() {
		fields = append(fields, strings.Fields(line))
	}
	return fields
}
func (s *LineBuffer) setClean() {
	s.Clean = true
}
func (s *LineBuffer) SetClean() {
	s.Lock()
	defer s.Unlock()
	s.setClean()
}
func (s *LineBuffer) setDirty() {
	s.Clean = false
}
func (s *LineBuffer) callBack() {
	for _, call := range s.Calls {
		call()
	}
}
func (s *LineBuffer) byteVec() *d.ByteVec {
	return s.Native.(*d.ByteVec)
}
func (s *LineBuffer) bytes() []byte {
	return []byte(*s.byteVec())
}
func (s *LineBuffer) string() string {
	return string(s.bytes())
}
func (s *LineBuffer) Bytes() []byte {
	s.Lock()
	defer s.Unlock()
	return s.bytes()
}
func (s *LineBuffer) Runes() []rune {
	return []rune(s.String())
}
func (s *LineBuffer) len() int {
	return s.byteVec().Len()
}
func (s *LineBuffer) Len() int {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	return s.len()
}
func (s *LineBuffer) read(p *[]byte) (int, error) {
	var n = cap(*p)
	if n >= 0 && n < s.len() {
		*p = append(make([]byte, 0, n), s.rang(0, n)...)
		s.cut(0, n)
		return n, nil
	}

	return 0, fmt.Errorf(
		"could not read line from buffer\n"+
			"buffer index position: %d\n"+
			"buffer: %s\n", n, s.string())
}
func (s *LineBuffer) peak() byte {
	if len(s.bytes()) > 0 {
		return s.bytes()[0]
	}
	return byte(0)
}
func (s *LineBuffer) peakN(n int) []byte {
	if len(s.bytes()) > n {
		return s.bytes()[:n]
	}
	return nil
}
func (s *LineBuffer) Read(p *[]byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	s.setClean()

	return s.read(p)
}
func (s *LineBuffer) ReadString() (string, error) {
	var buf = make([]byte, 0, s.Len())
	n, err := s.Read(&buf)
	if err != nil {
		return string(buf), fmt.Errorf("error in lexer at position: %d"+
			" while trying to read string from line buffer\n"+
			"buffer content: ",
			strconv.Itoa(n),
			string(*s.Native.(*d.ByteVec)))
	}
	return string(buf), nil
}

// read line reads one line from buffer & either replaces p with the bytes read
// from buffer if length of p is zero, or appends bytes read from buffer to p,
// if it's length is greater than zero
func (s *LineBuffer) ReadLine(p *[]byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	s.setClean()

	var lines = strings.Split(s.string(), "\n")
	var length = len([]byte(lines[0]))
	s.cut(0, length)
	*p = []byte(lines[0])

	return length, nil
}

// writes the content of p to the underlying buffer
func (s *LineBuffer) Write(p []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	*(s.byteVec()) = append(s.bytes(), p...)
	return len(p), nil
}
func (s *LineBuffer) WriteRunes(r []rune) (int, error) {
	return s.WriteString(string(r))
}
func (s *LineBuffer) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}
func (s *LineBuffer) Insert(i, j int, b byte) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	(s.byteVec()).Insert(i, j, b)
}
func (s *LineBuffer) InsertSlice(i, j int, p []byte) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	(s.byteVec()).InsertSlice(i, j, p...)
}
func (s *LineBuffer) ReplaceSlice(i, j int, trail []byte) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	copy(([]byte(*s.Native.(*d.ByteVec)))[i:j], trail)
}
func (s *LineBuffer) Split(i int) (h, t []byte) {
	s.Lock()
	defer s.Unlock()
	var head, tail = d.SliceSplit(s.Native.(d.ByteVec).Slice(), i)
	for _, b := range head {
		h = append(h, byte(b.(d.ByteVal)))
	}
	for _, b := range tail {
		t = append(t, byte(b.(d.ByteVal)))
	}
	return h, t
}
func (s *LineBuffer) cut(i, j int) {
	(s.byteVec()).Cut(i, j)
}
func (s *LineBuffer) Cut(i, j int) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()

	s.cut(i, j)
}
func (s *LineBuffer) delete(i int) {
	(s.byteVec()).Delete(i)
}
func (s *LineBuffer) Delete(i int) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	s.delete(i)
}
func (s *LineBuffer) Get(i int) byte {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	return byte((s.byteVec()).Get(d.IntVal(i)).(d.ByteVal))
}
func (s *LineBuffer) rang(i, j int) []byte {
	return []byte((s.byteVec()).Range(i, j))
}
func (s *LineBuffer) Range(i, j int) []byte {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	return s.rang(i, j)
}
func (s *LineBuffer) UpdateTrailing(line []rune) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	var bytes = []byte(string(line))
	var buflen = len(s.bytes())
	var trailen = len(bytes)
	if buflen >= trailen {
		var end = buflen
		var start = end - trailen
		copy(([]byte(*s.Native.(*d.ByteVec)))[start:end], bytes)
	}
}

////////////////////////////////////////////////////
type TokenBuffer d.AsyncVal

func NewTokenBuffer(callbacks ...func()) *TokenBuffer {
	return (*TokenBuffer)(d.NewAsync(callbacks...))
}
func (s TokenBuffer) AsyncVal() *d.AsyncVal { return (*d.AsyncVal)(&s) }
func (s TokenBuffer) String() string {
	s.Lock()
	defer s.Unlock()
	var str = bytes.NewBuffer([]byte{})
	var l = len(s.dataSlice())
	for i, tok := range toks(s.slice()...) {
		str.WriteString(tok.String())
		if i < l-1 {
			str.WriteString("\n")
		}
	}

	return str.String()
}
func (s *TokenBuffer) dataSlice() d.DataSlice {
	return s.Native.(d.DataSlice)
}
func (s *TokenBuffer) slice() []d.Native {
	return s.dataSlice().Slice()
}
func (s *TokenBuffer) setDirty() {
	s.Clean = false
}
func (s *TokenBuffer) SetClean() {
	s.Lock()
	defer s.Unlock()
	s.Clean = true
}
func (s *TokenBuffer) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.dataSlice().Len()
}
func (s *TokenBuffer) Tokens() []Token {
	s.Lock()
	defer s.Unlock()
	return toks(s.slice()...)
}
func (s *TokenBuffer) Get(i int) Token {
	s.Lock()
	defer s.Unlock()
	return s.dataSlice().GetInt(i).(Token)
}
func (s *TokenBuffer) Range(i, j int) []Token {
	s.Lock()
	defer s.Unlock()
	var toks = []Token{}
	for _, dat := range s.dataSlice()[i:j] {
		toks = append(toks, dat.(Token))
	}
	return toks
}
func (s *TokenBuffer) Split(i int) (h, t []Token) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	var head, tail = d.SliceSplit(s.dataSlice(), i)
	return toks(head...), toks(tail...)
}
func (s *TokenBuffer) Set(i int, tok Token) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	s.dataSlice().SetInt(i, tok)
}
func (s *TokenBuffer) Append(toks ...Token) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	s.Native = d.SliceAppend(s.dataSlice().Slice(), nats(toks...)...)
}
func (s *TokenBuffer) Insert(i int, toks []Token) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	s.Native = d.SliceInsertVector(s.dataSlice(), i, nats(toks...)...)
}
func (s *TokenBuffer) Delete(i int) {
	s.Lock()
	defer s.Unlock()
	s.setDirty()
	s.Native = d.SliceDelete(s.dataSlice(), i)
}
func (t *TokenBuffer) Sort() {
	t.Lock()
	defer t.Unlock()
	var ts = tokSort(toks([]d.Native(t.Native.(d.DataSlice))...))
	sort.Sort(ts)
	t.Native = d.DataSlice(nats(ts...))
}
func (t *TokenBuffer) Search(pos int) int {
	return sort.Search(t.Len(), func(i int) bool {
		return pos < t.dataSlice().Slice()[i].(Token).Pos()
	})
}

//////
func nats(toks ...Token) []d.Native {
	var nats = []d.Native{}
	for _, nat := range toks {
		nats = append(nats, nat)
	}
	return nats
}
func toks(nats ...d.Native) []Token {
	var toks = []Token{}
	for _, nat := range nats {
		toks = append(toks, nat.(Token))
	}
	return toks
}

//////
type tokSort []Token

func (t tokSort) Len() int { return len(t) }
func (t tokSort) Less(i, j int) bool {
	return []Token(t)[i].Pos() <
		[]Token(t)[j].Pos()
}
func (t tokSort) Swap(i, j int) {
	[]Token(t)[i], []Token(t)[j] = []Token(t)[j], []Token(t)[i]
}

type doFnc func(Lexer) f.StateFnc
type Lexer func() (*LineBuffer, *TokenBuffer)

func (lex Lexer) Buffer() (*LineBuffer, *TokenBuffer) { l, t := lex(); return l, t }
func (lex Lexer) LineBuffer() *LineBuffer             { l, _ := lex(); return l }
func (lex Lexer) TokenBuffer() *TokenBuffer           { _, t := lex(); return t }
func (lex Lexer) nextState(do doFnc) f.StateFnc       { return do(lex) }
func newLexer(l *LineBuffer, t *TokenBuffer) Lexer {
	return func() (*LineBuffer, *TokenBuffer) { return l, t }
}

// lexer gets called, whenever linebuffer changes
func NewLexer(lbuf *LineBuffer) *TokenBuffer {
	// allocate new token buffer
	var tbuf = NewTokenBuffer()
	// enclose both buffers in a lexer instance
	var lex = newLexer(lbuf, tbuf)
	// retrieve state function from lexer closure over buffers
	var sf = lex.nextState(func(l Lexer) f.StateFnc {
		return lexer(lex)
	})
	// subscribe state functions run method to be called back once per
	// change in line buffer
	lbuf.Subscribe(sf.Run)
	// return token buffer reference
	return tbuf
}

func lexer(lex Lexer) f.StateFnc {
	var do doFnc

	return lex.nextState(do)
}
