package lang

import (
	d "github.com/JoergReinhardt/godeep/data"
)

///// SYNTAX DEFINITION /////
type TokType d.BitFlag

func (t TokType) Flag() d.BitFlag { return d.BitFlag(t) }
func (t TokType) Syntax() string  { return syntax[t] }

//go:generate stringer -type=TokType
const (
	None  TokType = 1
	Blank TokType = 1 << iota
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
	Not
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
	BckSla
	Lambda
	Number
	Letter
	Capital
	GenType
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
	TypeIdent
	FuncIdent
)

var syntax = map[TokType]string{
	None:          "",
	Blank:         " ",
	Underscore:    "_",
	Asterisk:      "∗",
	Dot:           ".",
	Comma:         ",",
	Colon:         ":",
	Semicolon:     ";",
	Minus:         "-",
	Plus:          "+",
	Or:            "∨",
	Xor:           "※",
	And:           "∧",
	Equal:         "=",
	Lesser:        "≪",
	Greater:       "≫",
	Lesseq:        "≤",
	Greaterq:      "≥",
	LeftPar:       "(",
	RightPar:      ")",
	LeftBra:       "[",
	RightBra:      "]",
	LeftCur:       "{",
	RightCur:      "}",
	Slash:         "/",
	Not:           "≠",
	Dec:           "--",
	Inc:           "++",
	DoubEqual:     "==",
	TripEqual:     "≡",
	RightArrow:    "←",
	LeftArrow:     "→",
	FatLArrow:     "⇐",
	FatRArrow:     "⇒",
	DoubCol:       "∷",
	Sing_quote:    `'`,
	Doub_quote:    `"`,
	BckSla:        `\`,
	Lambda:        "λ",
	Number:        "[0-9]",
	Letter:        "[a-z]",
	Capital:       "[A-Z]",
	GenType:       "[[a-w]|y|z]",
	HeadWord:      "x",
	TailWord:      "xs",
	InWord:        "in",
	ConWord:       "con",
	LetWord:       "let",
	MutableWord:   "mutable",
	WhereWord:     "where",
	OtherwiseWord: "otherwise",
	IfWord:        "if",
	ThenWord:      "then",
	ElseWord:      "else",
	CaseWord:      "case",
	OfWord:        "of",
	DataWord:      "data",
	TypeWord:      "type",
	TypeIdent:     "[A-z][a-z]*",
	FuncIdent:     "([a-w|y|z][a-z])|(x[a-r|t-z])",
}

//// Token type according to text, scanner tokenizer.
type Token d.BitFlag

func Con(t d.Typed) Token {
	if tok, ok := t.(d.Type); ok {
		return Token(tok.Flag())
	}
	return Token(t.Flag())
}
func newTypeToken(typ d.Type) Token {
	return Token(typ.Flag())
}
func newSyntaxToken(typ TokType) Token {
	return Token(typ.Flag())
}
func (t Token) Flag() d.BitFlag { return d.BitFlag(t) }
func (t Token) Text() string {
	return syntax[TokType(t.Flag())]
}
func (t Token) String() string {
	return syntax[TokType(t.Flag())]
}
func (t Token) Type() d.BitFlag { return d.BitFlag(t) }