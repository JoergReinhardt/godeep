package functions

type ()

///////////////////////////////////////////////////////////////////////////////
//// COMPOSITION OF DEFINED FUNCTIONS
///
// define the curryed function, so that it accepts the argument types of the g
// function passed as second argument to the constructor and the return type of
// the g function passed as its second argument.
func Curry(f, g FuncDef) FuncDef {
	if f.TypeArgs().Match(g.TypeRet()) {
		return Define(Lambda(

			func(args ...Expression) Expression {

				// call f with the result of calling g applying
				// the arguments if any are given
				if len(args) > 0 {
					return f.Call(g.Call(args...))
				}
				return f.Call(g.Call())
			}),

			// define a function by composing both type ids with
			// the argument type of g passed as second argument and
			// the return type of f passed as first argument
			Def(g.TypeId(), Def(f.TypeId())),
			f.TypeRet(),  // ‥.return type of g &
			g.TypeArgs(), //‥.argument type of f
		)
	}
	return Define(NewNone(), None, None)
}

///////////////////////////////////////////////////////////////////////////////
//// CONTINUATION COMPOSITION
///
// flatten flattens sequences of sequences to one dimension
func Flatten(con Continuation) VecVal {
	return VecVal(func(args ...Expression) []Expression {
		if len(args) > 0 {
			con = Flatten(con.Call(args...).(Sequential))
		}
		var head, tail = con.Continue()
		if head.Type().Match(Sequences) {
			return append(
				Flatten(head.(Sequential))(),
				Flatten(tail.(Sequential))()...)
		}
		return append([]Expression{head},
			Flatten(tail.(Sequential))()...)
	})
}

// map returns a continuation calling the map function for every element
func Map(
	con Continuation,
	mapf func(Expression) Expression,
) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = Map(
					NewSequence(args...), mapf,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	return SeqVal(func(args ...Expression) (Expression, SeqVal) {
		if len(args) > 0 {
			con = con.Call(args...).(Continuation)
		}
		var head, tail = con.Continue()
		// skip none instances, when tail has further elements
		if head.Type().Match(None) && !tail.Empty() {
			return Map(tail, mapf)()
		}
		return mapf(head), Map(tail, mapf)
	})
}

// apply returns a continuation called on every element of the continuation.
// when continuation is called pssing arguments those, are passed to apply
// alongside the current element
func Apply(
	con Continuation,
	apply func(Expression, ...Expression) Expression,
) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = Apply(
					NewSequence(args...), apply,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	var head, tail = con.Continue()
	return SeqVal(func(args ...Expression) (Expression, SeqVal) {
		if len(args) > 0 {
			var head, tail = apply(head, args...),
				Apply(tail, apply)
				// skip none
			if head.Type().Match(None) && !tail.Empty() {
				return Apply(tail, apply)()
			}

		}
		var head, tail = apply(head), Apply(tail, apply)
		// skip none
		if head.Type().Match(None) && !tail.Empty() {
			return Apply(tail, apply)()
		}
		return head, tail
	})
}

// fold takes a continuation, an initial expression and a fold function. the
// fold function is called for every element of the continuation and passed the
// current element and init expression and returns a possbly altered init
// element to pass to next call
func Fold(
	con Continuation,
	init Expression,
	fold func(init, head Expression) Expression,
) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = Fold(
					NewSequence(args...), init, fold,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	var (
		head, tail = con.Continue()
		result     = fold(init, head)
	)
	// skip none instances, when tail has further elements
	if result.Type().Match(None) && !tail.Empty() {
		for result.Type().Match(None) && !tail.Empty() {
			head, tail = tail.Continue()
			result = fold(init, head)
		}
	}
	init = result
	return SeqVal(func(args ...Expression) (Expression, SeqVal) {
		if len(args) > 0 {
			return init.Call(args...), Fold(tail, init, fold)
		}
		return init, Fold(tail, init, fold)
	})
}

// continuation of elements not matched by test
func Filter(
	con Continuation,
	filter func(Expression) bool,
) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = Filter(
					NewSequence(args...), filter,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	var (
		init = NewSequence()
		fold = func(init, head Expression) Expression {
			if filter(head) {
				return NewNone()
			}
			return init.(SeqVal).Cons(head)
		}
	)
	return Fold(con, init, fold)
}

// continuation of elements matched by test
func Pass(
	con Continuation,
	pass func(Expression) bool,
) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = Pass(
					NewSequence(args...), pass,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	var (
		init = NewSequence()
		fold = func(init, head Expression) Expression {
			if pass(head) {
				return init.(SeqVal).Cons(head)
			}
			return NewNone()
		}
	)
	return Fold(con, init, fold)
}

// take-n is a variation of fold that takes an initial continuation cuts and
// returns it as continuation of vector instances of length n
func TakeN(con Continuation, n int) SeqVal {
	if con.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			if len(args) > 0 {
				var head, tail = TakeN(
					NewSequence(NewVector(args...)), n,
				).Continue()
				return head, tail.(SeqVal)
			}
			return NewNone(), nil
		})
	}
	var (
		init  = NewVector()
		takeN = func(init Expression, arg Expression) Expression {
			var (
				vector = init.(VecVal)
			)
			if vector.Len() == n {
				return NewVector(arg)
			}
			return vector.Put(arg)
		}
	)
	return Filter(Fold(con, init, takeN),
		func(arg Expression) bool {
			return arg.(VecVal).Len() < n
		})
}

// split is a variation of fold that splits either a continuation of pairs, or
// takes two arguments at a time and splits those into continuation of left and
// right values and returns those as elements of a pair
func Split(
	con Continuation,
	left, right func(arg Expression) bool,
) SeqVal {
	var (
		pair  = NewPair(NewSequence(), NewSequence())
		split = func(init, arg Expression) Expression {
			var (
				sleft  = init.(ValPair).Left().(SeqVal)
				sright = init.(ValPair).Right().(SeqVal)
			)
			var matchL, matchR = left(arg), right(arg)
			if !matchL && !matchR {
				return NewNone()
			}
			if matchL {
				sleft = sleft.Cons(arg).(SeqVal)
			}
			if matchR {
				sright = sright.Cons(arg).(SeqVal)
			}
			return NewPair(sleft, sright)
		}
	)
	return Fold(con, pair, split)
}

// zip expects two continuations and a function to create a list of resulting
// elements each created from the two current continuation heads, using the
// passed zip function. if arguments are passed calling the lists call method,
// the results call method is called after heads have been zipped, passing on
// those arguemnts.
func Zip(
	left, right Continuation,
	zip func(l, r Expression) Expression,
) SeqVal {
	if left.Empty() && right.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			return NewNone(), nil
		})
	}
	return SeqVal(func(args ...Expression) (Expression, SeqVal) {
		var (
			gh, gt = left.Continue()
			fh, ft = right.Continue()
		)
		if len(args) > 0 {
			// pass arguments to results call method
			return zip(gh, fh).Call(args...),
				Zip(gt, ft, zip)
		}
		return zip(gh, fh), Zip(gt, ft, zip)
	})
}

// bind works similar to zip, but the bind function takes additional arguments
// passed to the bound list when calling the call method with arguments
// directly (instead of passing on to results call method, analog to map/apply)
func Bind(
	f, g Continuation,
	bind func(l, r Expression, args ...Expression) Expression,
) SeqVal {
	if f.Empty() && g.Empty() {
		return SeqVal(func(args ...Expression) (Expression, SeqVal) {
			return NewNone(), nil
		})
	}
	return SeqVal(func(args ...Expression) (Expression, SeqVal) {
		var (
			gh, gt = f.Continue()
			fh, ft = g.Continue()
		)
		if len(args) > 0 {
			// pass arguments on to bind
			return bind(gh, fh, args...),
				Bind(gt, ft, bind)
		}
		return bind(gh, fh), Bind(gt, ft, bind)
	})
}
