package wizardry

import (
	"encoding/binary"
)

// Spellbook contains a set of rules - at least one "" page, potentially others
type Spellbook map[string][]Rule

// AddRule appends a rule to the spellbook on the given page
func (sb Spellbook) AddRule(page string, rule Rule) {
	sb[page] = append(sb[page], rule)
}

// Rule is a single magic rule
type Rule struct {
	Level       int
	Offset      Offset
	Kind        Kind
	Description []byte
}

// Endianness describes the order in which a multi-byte number is stored
type Endianness int

// ByteOrder translates our in-house Endianness constant into a binary.ByteOrder decoder
func (en Endianness) ByteOrder() binary.ByteOrder {
	if en == BigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

// Swapped returns LittleEndian if you give it BigEndian, and vice versa
func (en Endianness) Swapped() Endianness {
	if en == BigEndian {
		return LittleEndian
	}
	return BigEndian
}

// MaybeSwapped returns swapped endianness if swap is true
func (en Endianness) MaybeSwapped(swap bool) Endianness {
	if !swap {
		return en
	}
	return en.Swapped()
}

const (
	// LittleEndian numbers are stored with the least significant byte first
	LittleEndian Endianness = iota
	// BigEndian numbers are stored with the most significant byte first
	BigEndian = iota
)

// Kind describes the type of tests a magic rule performs
type Kind struct {
	Family KindFamily
	Data   interface{}
}

// IntegerKind describes how to perform a test on an integer
type IntegerKind struct {
	ByteWidth   int
	Endianness  Endianness
	Signed      bool
	DoAnd       bool
	AndValue    uint64
	IntegerTest IntegerTest
	Value       int64
	MatchAny    bool
}

// IntegerTest describes which comparison to perform on an integer
type IntegerTest int

const (
	// IntegerTestEqual tests that two integers have the same value
	IntegerTestEqual IntegerTest = iota
	// IntegerTestNotEqual tests that two integers have different values
	IntegerTestNotEqual = iota
	// IntegerTestLessThan tests that one integer is less than the other
	IntegerTestLessThan = iota
	// IntegerTestGreaterThan tests that one integer is greater than the other
	IntegerTestGreaterThan = iota
)

// StringKind describes how to match a string pattern
type StringKind struct {
	Value  []byte
	Negate bool
	Flags  stringTestFlags
}

// SearchKind describes how to look for a fixed pattern
type SearchKind struct {
	Value  []byte
	MaxLen int
}

// KindFamily groups tests in families (all integer tests, for example)
type KindFamily int

const (
	// KindFamilyInteger tests numbers for equality, inequality, etc.
	KindFamilyInteger KindFamily = iota
	// KindFamilyString looks for a string, with casing and whitespace rules
	KindFamilyString = iota
	// KindFamilySearch looks for a precise string in a slice of the target
	KindFamilySearch = iota
	// KindFamilyDefault succeeds if no tests succeeded before on that level
	KindFamilyDefault = iota
	// KindFamilyClear resets the matched flag for that level
	KindFamilyClear = iota
)

// Offset describes where to look to compare something
type Offset struct {
	OffsetType OffsetType
	IsRelative bool
	Direct     int64
	Indirect   *IndirectOffset
}

// OffsetType describes whether an offset is direct or indirect
type OffsetType int

const (
	// OffsetTypeIndirect is an offset read from somewhere in a file
	OffsetTypeIndirect OffsetType = iota
	// OffsetTypeDirect is an offset directly specified by the magic
	OffsetTypeDirect = iota
)

// IndirectOffset indicates where to look in a file to find the real offset
type IndirectOffset struct {
	IsRelative                 bool
	ByteWidth                  int
	Endianness                 Endianness
	OffsetAddress              int64
	OffsetAdjustmentType       OffsetAdjustment
	OffsetAdjustmentIsRelative bool
	OffsetAdjustmentValue      int64
}

// OffsetAdjustment describes which operation to apply to an offset
type OffsetAdjustment int

const (
	// OffsetAdjustmentNone is a no-op
	OffsetAdjustmentNone OffsetAdjustment = iota
	// OffsetAdjustmentAdd adds a value
	OffsetAdjustmentAdd = iota
	// OffsetAdjustmentSub subtracts a value
	OffsetAdjustmentSub = iota
	// OffsetAdjustmentMul multiplies by a value
	OffsetAdjustmentMul = iota
	// OffsetAdjustmentDiv divides by a value
	OffsetAdjustmentDiv = iota
)