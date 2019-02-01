package lex

import (
	"strings"

	d "github.com/JoergReinhardt/godeep/data"
	"github.com/olekukonko/tablewriter"
)

///// SYNTAX DEFINITION /////
type SyntaxItemFlag d.BitFlag

func (t SyntaxItemFlag) Flag() d.BitFlag      { return d.BitFlag(t) }
func (t SyntaxItemFlag) Type() SyntaxItemFlag { return t }
func (t SyntaxItemFlag) Syntax() string       { return syntax[t] }

// all syntax items represented as string
func AllSyntax() string {
	str := &strings.Builder{}
	tab := tablewriter.NewWriter(str)
	for _, t := range AllItems() {
		row := []string{
			t.String(), syntax[t], matchSyntax[syntax[t]],
		}
		tab.Append(row)
	}
	tab.Render()
	return str.String()
}

// slice of all syntax items in there int constant form
func AllItems() []SyntaxItemFlag {
	var tt = []SyntaxItemFlag{}
	var i uint
	var t SyntaxItemFlag = 0
	for i < 63 {
		t = 1 << i
		i = i + 1
		tt = append(tt, SyntaxItemFlag(t))
	}
	return tt
}

//go:generate stringer -type=SyntaxItemFlag
const (
	None  SyntaxItemFlag = 1
	Blank SyntaxItemFlag = 1 << iota
	Underscore
	Asterisk
	Dot
	Comma
	Colon
	Semicolon
	Minus
	Plus
	Or
	Xor
	And
	Equal
	Lesser
	Greater
	Lesseq
	Greaterq
	LeftPar
	RightPar
	LeftBra
	RightBra
	LeftCur
	RightCur
	Slash
	Pipe
	Not
	Unequal
	Dec
	Inc
	DoubEqual
	TripEqual
	RightArrow
	LeftArrow
	FatLArrow
	FatRArrow
	DoubCol
	Sing_quote
	Doub_quote
	BackSlash
	Lambda
	Function
	Polymorph
	HeadWord
	TailWord
	InWord
	ConWord
	LetWord
	MutableWord
	WhereWord
	OtherwiseWord
	IfWord
	ThenWord
	ElseWord
	CaseWord
	OfWord
	DataWord
	TypeWord
	StringItem
	NumeralItem
	FuncIdent
	TypeIdent

	// parts of syntax that constitute indipendent items, but commonly end
	// up concatenated to other items without seperating space character to
	// split by, which results in the need for further processing.
	Appendices = Dot | Comma | Colon | Semicolon | LeftPar | RightPar | LeftBra |
		RightBra | LeftCur | RightCur
)

var keywords = map[string]SyntaxItemFlag{
	"in":        InWord,
	"con":       ConWord,
	"let":       LetWord,
	"mutable":   MutableWord,
	"where":     WhereWord,
	"otherwise": OtherwiseWord,
	"if":        IfWord,
	"then":      ThenWord,
	"else":      ElseWord,
	"case":      CaseWord,
	"of":        OfWord,
	"data":      DataWord,
	"type":      TypeWord,
}

var matchSyntax = map[string]string{
	"⊥":  "",
	" ":  " ",
	"_":  "_",
	"∗":  "*",
	".":  ".",
	",":  ",",
	":":  ":",
	";":  ";",
	"-":  "-",
	"+":  "+",
	"∨":  "OR",
	"※":  "XOR",
	"∧":  "AND",
	"=":  "=",
	"≪":  "<<",
	"≫":  ">>",
	"≤":  "=<",
	"≥":  ">=",
	"(":  "(",
	")":  ")",
	"[":  "[",
	"]":  "]",
	"{":  "{",
	"}":  "}",
	"/":  "/",
	"¬":  "!",
	"≠":  "!=",
	"--": "--",
	"++": "++",
	"==": "==",
	"≡":  "===",
	"→":  "->",
	"←":  "<-",
	"⇐":  "<=",
	"⇒":  "=>",
	"∷":  "::",
	`'`:  `'`,
	`"`:  `"`,
	`\`:  `\`,
	`ϝ`:  `\f`,
	"λ":  `\x`,
	"x":  "x",
	"xs": "xs",
}

var syntax = map[SyntaxItemFlag]string{
	None:       "⊥",
	Blank:      " ",
	Underscore: "_",
	Asterisk:   "∗",
	Dot:        ".",
	Comma:      ",",
	Colon:      ":",
	Semicolon:  ";",
	Minus:      "-",
	Plus:       "+",
	Or:         "∨",
	Xor:        "※",
	And:        "∧",
	Equal:      "=",
	Lesser:     "≪",
	Greater:    "≫",
	Lesseq:     "≤",
	Greaterq:   "≥",
	LeftPar:    "(",
	RightPar:   ")",
	LeftBra:    "[",
	RightBra:   "]",
	LeftCur:    "{",
	RightCur:   "}",
	Slash:      "/",
	Pipe:       "|",
	Not:        "¬",
	Unequal:    "≠",
	Dec:        "--",
	Inc:        "++",
	DoubEqual:  "==",
	TripEqual:  "≡",
	RightArrow: "→",
	LeftArrow:  "←",
	FatLArrow:  "⇐",
	FatRArrow:  "⇒",
	DoubCol:    "∷",
	Sing_quote: `'`,
	Doub_quote: `"`,
	BackSlash:  `\`,
	Function:   "ϝ",
	Polymorph:  "Ф",
	Lambda:     "λ",
	HeadWord:   "x",
	TailWord:   "xs",
}

var match = map[string]SyntaxItemFlag{
	"":    None,
	" ":   Blank,
	"_":   Underscore,
	"*":   Asterisk,
	".":   Dot,
	",":   Comma,
	":":   Colon,
	";":   Semicolon,
	"-":   Minus,
	"+":   Plus,
	"OR":  Or,
	"XOR": Xor,
	"AND": And,
	"=":   Equal,
	"<<":  Lesser,
	">>":  Greater,
	"=<":  Lesseq,
	">=":  Greaterq,
	"(":   LeftPar,
	")":   RightPar,
	"[":   LeftBra,
	"]":   RightBra,
	"{":   LeftCur,
	"}":   RightCur,
	"/":   Slash,
	"|":   Pipe,
	"!":   Not,
	"!=":  Unequal,
	"--":  Dec,
	"++":  Inc,
	"==":  DoubEqual,
	"===": TripEqual,
	"->":  RightArrow,
	"<-":  LeftArrow,
	"<=":  FatLArrow,
	"=>":  FatRArrow,
	"::":  DoubCol,
	`'`:   Sing_quote,
	`"`:   Doub_quote,
	`\`:   BackSlash,
	`\f`:  Function,
	`\F`:  Polymorph,
	`\x`:  Lambda,
	"x":   HeadWord,
	"xs":  TailWord,
}

// checks if a string matches any representation of a syntax item and yields
// it, if that's the case
func MatchString(tos string) Item { return syntaxItem(match[tos]) }

// convert item string representation from editable to pretty
func ASCIIToUtf8(tos ...string) []SyntaxItemFlag {
	var ti = []SyntaxItemFlag{}
	for _, s := range tos {
		ti = append(ti, match[s])
	}
	return ti
}

// convert item string representation from pretty to editable
func Utf8ToASCII(tos ...string) string {
	var sto string
	for _, s := range tos {
		sto = sto + matchSyntax[s]
	}
	return sto
}

// item is a bitflag of course
type Item interface {
	Flag() d.BitFlag
	Type() SyntaxItemFlag
	String() string
}

type syntaxItem d.BitFlag
type stringItem struct {
	SyntaxItemFlag
	string
}

func (t syntaxItem) Type() SyntaxItemFlag { return SyntaxItemFlag(t.Flag()) }
func (t stringItem) Type() SyntaxItemFlag { return SyntaxItemFlag(t.Flag()) }

// pretty utf-8 version of syntax item
func (t syntaxItem) String() string { return SyntaxItemFlag(t.Flag()).Syntax() }
func (t stringItem) String() string { return t.string }

// provides an alternative string representation that can be edited without
// having to produce utf-8 digraphs
func (t syntaxItem) StringAlt() string { return matchSyntax[syntax[SyntaxItemFlag(t.Flag())]] }
func (t stringItem) StringAlt() string { return matchSyntax[syntax[SyntaxItemFlag(t.Flag())]] }
func (t syntaxItem) Flag() d.BitFlag   { return d.Flag.Flag() }
func (t stringItem) Flag() d.BitFlag   { return d.Flag.Flag() }
