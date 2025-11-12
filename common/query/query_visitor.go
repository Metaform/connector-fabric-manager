// Code generated from Query.g4 by ANTLR 4.13.2. DO NOT EDIT.

package query // Query
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by QueryParser.
type QueryVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by QueryParser#predicate.
	VisitPredicate(ctx *PredicateContext) interface{}

	// Visit a parse tree produced by QueryParser#compoundPredicate.
	VisitCompoundPredicate(ctx *CompoundPredicateContext) interface{}

	// Visit a parse tree produced by QueryParser#atomicPredicate.
	VisitAtomicPredicate(ctx *AtomicPredicateContext) interface{}

	// Visit a parse tree produced by QueryParser#fieldPath.
	VisitFieldPath(ctx *FieldPathContext) interface{}

	// Visit a parse tree produced by QueryParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by QueryParser#value.
	VisitValue(ctx *ValueContext) interface{}

	// Visit a parse tree produced by QueryParser#valueList.
	VisitValueList(ctx *ValueListContext) interface{}

	// Visit a parse tree produced by QueryParser#stringValue.
	VisitStringValue(ctx *StringValueContext) interface{}

	// Visit a parse tree produced by QueryParser#numericValue.
	VisitNumericValue(ctx *NumericValueContext) interface{}

	// Visit a parse tree produced by QueryParser#booleanValue.
	VisitBooleanValue(ctx *BooleanValueContext) interface{}
}
