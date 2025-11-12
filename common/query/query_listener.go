// Code generated from Query.g4 by ANTLR 4.13.2. DO NOT EDIT.

package query // Query
import "github.com/antlr4-go/antlr/v4"

// QueryListener is a complete listener for a parse tree produced by QueryParser.
type QueryListener interface {
	antlr.ParseTreeListener

	// EnterPredicate is called when entering the predicate production.
	EnterPredicate(c *PredicateContext)

	// EnterCompoundPredicate is called when entering the compoundPredicate production.
	EnterCompoundPredicate(c *CompoundPredicateContext)

	// EnterAtomicPredicate is called when entering the atomicPredicate production.
	EnterAtomicPredicate(c *AtomicPredicateContext)

	// EnterFieldPath is called when entering the fieldPath production.
	EnterFieldPath(c *FieldPathContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterValueList is called when entering the valueList production.
	EnterValueList(c *ValueListContext)

	// EnterStringValue is called when entering the stringValue production.
	EnterStringValue(c *StringValueContext)

	// EnterNumericValue is called when entering the numericValue production.
	EnterNumericValue(c *NumericValueContext)

	// EnterBooleanValue is called when entering the booleanValue production.
	EnterBooleanValue(c *BooleanValueContext)

	// ExitPredicate is called when exiting the predicate production.
	ExitPredicate(c *PredicateContext)

	// ExitCompoundPredicate is called when exiting the compoundPredicate production.
	ExitCompoundPredicate(c *CompoundPredicateContext)

	// ExitAtomicPredicate is called when exiting the atomicPredicate production.
	ExitAtomicPredicate(c *AtomicPredicateContext)

	// ExitFieldPath is called when exiting the fieldPath production.
	ExitFieldPath(c *FieldPathContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitValueList is called when exiting the valueList production.
	ExitValueList(c *ValueListContext)

	// ExitStringValue is called when exiting the stringValue production.
	ExitStringValue(c *StringValueContext)

	// ExitNumericValue is called when exiting the numericValue production.
	ExitNumericValue(c *NumericValueContext)

	// ExitBooleanValue is called when exiting the booleanValue production.
	ExitBooleanValue(c *BooleanValueContext)
}
