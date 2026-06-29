package a

// Thing is a Null-Object interface (it declares IsNull() bool).
type Thing interface {
	IsNull() bool
	Foo() string
}

type null struct{}

func (null) IsNull() bool { return true }
func (null) Foo() string  { return "" }

// NullThing is the canonical Null object.
func NullThing() Thing { return null{} }

// Embedded marker via an embedded interface is still detected.
type Base interface{ IsNull() bool }
type Derived interface {
	Base
	Bar() int
}

func okConstructor() Thing { return NullThing() } // ok
func okNamed() Thing {
	t := NullThing()
	return t // ok: not a bare nil
}

func badReturn() Thing         { return nil }      // want `return a Null object .* instead of nil for a\.Thing`
func badEmbedded() Derived     { return nil }      // want `return a Null object .* for a\.Derived`
func badMulti() (Thing, error) { return nil, nil } // want `return a Null object .* for a\.Thing`

func notNullObject() error { return nil } // ok: error has no IsNull()
func plainPointer() *int   { return nil } // ok: not an interface

func inLiteral() Thing {
	f := func() Thing { return nil } // want `return a Null object .* for a\.Thing`
	return f()
}

// --- construction-site checks (assignment, var, composite literals) ---

type Box struct {
	T Thing
	N int
}

func assignments() {
	var x Thing
	x = nil // want `assign a Null object .* for a\.Thing`
	_ = x

	var y Thing = nil // want `initialise a Null object .* for a\.Thing`
	_ = y

	var ok Thing = NullThing()
	ok = NullThing() // ok
	_ = ok
}

func construction() {
	_ = Box{T: nil}                // want `set field T to a Null object .* for a\.Thing`
	_ = Box{nil, 0}                // want `set field to a Null object .* for a\.Thing`
	_ = Box{NullThing(), 0}        // ok (positional, non-nil)
	_ = Box{N: 0}                  // ok (T field omitted; zero-value nil is not flagged)
	_ = Box{N: 1, T: nil}          // want `set field T to a Null object .* for a\.Thing`
	_ = []Thing{nil}               // want `set element to a Null object .* for a\.Thing`
	_ = [1]Thing{nil}              // want `set element to a Null object .* for a\.Thing`
	_ = []Thing{0: nil}            // want `set element to a Null object .* for a\.Thing`
	_ = [2]Thing{1: nil}           // want `set element to a Null object .* for a\.Thing`
	_ = map[string]Thing{"k": nil} // want `set map value to a Null object .* for a\.Thing`
	_ = map[Thing]int{nil: 1}      // ok (only map values are checked, not keys)
}

// --- keyed struct field set to a real (non-nil) field, and embedded struct
// keyed by its type name: exercises structFieldType's found-branch where the
// field is not a Null-Object interface, so nothing is reported.

type Inner struct{ T Thing }
type Outer struct {
	Inner
	N int
}

func keyedNonNull() {
	_ = Outer{Inner: Inner{T: NullThing()}, N: 1} // ok (Inner field is a struct; its T is non-nil)
}

// --- interfaces that declare an IsNull method but with the WRONG signature are
// NOT Null-Object interfaces, so a nil targeting them is left alone. These cover
// isNullObjectInterface's signature-shape rejection.

type WrongParams interface {
	IsNull(int) bool
}
type WrongResults interface {
	IsNull() (bool, error)
}
type WrongResultType interface {
	IsNull() int
}

func wrongSignatures() {
	var a WrongParams = nil      // ok: IsNull takes a parameter
	var b WrongResults = nil     // ok: IsNull returns two results
	var c WrongResultType = nil  // ok: IsNull returns a non-bool
	_, _, _ = a, b, c
}

func wrongSignatureReturn() WrongParams { return nil } // ok: not a Null-Object interface

// --- checkReturn mismatch branch: `return f()` feeds N results from a single
// call, so len(ret.Results) != results.Len() and the return is left alone.

func twoThings() (Thing, Thing) { return NullThing(), NullThing() }

func forwardCall() (Thing, Thing) {
	return twoThings() // ok: single call feeding multiple results (count mismatch, skipped)
}

func nakedReturn() (t Thing) {
	t = NullThing()
	return // ok: naked return, left alone
}
