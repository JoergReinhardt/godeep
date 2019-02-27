package parse

import (
	"strings"

	f "github.com/JoergReinhardt/gatwd/functions"
	l "github.com/JoergReinhardt/gatwd/lex"
)

//
//// lexer closes over a line buffer containing the input and a queue, to share
//// with the caller and emit it's tokens to.
type Lexer func() (f.IOBuffer, f.IOVector)

// main switch case expression matcher
func UtfPrefix(line string) bool {
	return strings.HasPrefix(line, l.UniCharString)
}
func AsciiPrefix(line string) bool {
	return strings.HasPrefix(line, l.AsciiString)
}
func DigitPrefix(line string) bool {
	return strings.HasPrefix(line, l.DigitString)
}
func KeyWordPrefix(line string) bool {
	return strings.HasPrefix(line, l.KeyWordString)
}

type StepFnc func(Lexer) f.StateFnc

// DO FUNCTIONS OF THE LEXER MONAD
func lexer(line string) f.StateFnc {
	var lex Lexer
	var do func(Lexer) f.StateFnc
	// return next state based on currently trailing characters
	switch {
	case UtfPrefix(line):
		do = consumeUtf
	case AsciiPrefix(line):
		do = consumeAscii
	case DigitPrefix(line):
		do = consumeDigits
	case KeyWordPrefix(line):
		do = consumeKeyword
	default:
		do = consumeLetters
	}
	return do(lex)
}

func consumeUtf(lex Lexer) f.StateFnc {
	var do StepFnc
	var curl string
	var queue f.IOVector

	for _, utf := range l.UniChars {
		if strings.HasPrefix(curl, utf) {
			item, _ := l.MatchUtf8(utf)
			queue.Append(NewSyntaxToken(item))
			curl = strings.TrimPrefix(curl, utf)
		}
	}

	do = lexer

	return lex.nextState(do)
}

func consumeAscii(lex Lexer) f.StateFnc {
	var do StepFnc
	var lbuf, queue, curl = lex()

	for _, asc := range l.Ascii {
		if strings.HasPrefix(curl, asc) {
			item, _ := l.MatchItem(asc)
			queue.Append(NewSyntaxToken(
				queue.CurrentPos(),
				item))
			curl = strings.TrimPrefix(curl, asc)
		}
	}

	do = lexer

	lex = newLexer(lbuf, queue, curl)
	return lex.nextState(do)
}

func consumeKeyword(lex Lexer) f.StateFnc {
	var do StepFnc
	var lbuf, queue, curl = lex()

	for _, keyword := range l.Keywords {
		if strings.HasPrefix(curl, keyword) {
			queue.Append(NewKeywordToken(
				queue.CurrentPos(),
				keyword))
			curl = strings.TrimPrefix(curl, keyword)
		}
	}

	do = lexer

	lex = newLexer(lbuf, queue, curl)
	return lex.nextState(do)
}

func consumeDigits(lex Lexer) f.StateFnc {
	var do StepFnc
	var lbuf, queue, curl = lex()
	var digits = []byte{}

	for DigitPrefix(curl) {
		digits = append(digits, curl[0])
	}

	queue.Append(NewDigitToken(
		queue.CurrentPos(),
		string(digits)))

	do = lexer

	lex = newLexer(lbuf, queue, curl)
	return lex.nextState(do)
}

func consumeLetters(lex Lexer) f.StateFnc {
	var do StepFnc
	var lbuf, queue, curl = lex()
	var letters = []rune{}
	var runes = []rune(curl)

	for strings.ContainsAny(string(runes[0]),
		l.LetterString+l.CapitalString) {

		letters = append(letters, runes[0])
		if len(runes) > 1 {
			curl = string(runes[1:])
			runes = []rune(curl)
		} else {
			curl = ""
			runes = runes[:0]
			break
		}
	}

	queue.Append(NewWordToken(
		queue.CurrentPos(),
		string(letters)))

	do = lexer

	lex = newLexer(lbuf, queue, curl)
	return lex.nextState(do)
}
