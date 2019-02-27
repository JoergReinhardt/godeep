package parse

import d "github.com/JoergReinhardt/gatwd/data"
import f "github.com/JoergReinhardt/gatwd/functions"

// data to parse
type Token interface {
	f.Functional
	TypeTok() TyToken
	Data() d.Native
}

// Ident interface{}
//
// the ident interface is implemented by everything providing unique identification.
type Ident interface {
	f.Functional
	Ident() f.Functional // calls enclosed fnc, with enclosed parameters
}
