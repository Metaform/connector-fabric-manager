// Code generated from Query.g4 by ANTLR 4.13.2. DO NOT EDIT.

package query // Query
import "github.com/antlr4-go/antlr/v4"

// BaseQueryListener is a complete listener for a parse tree produced by QueryParser.
type BaseQueryListener struct{}

var _ QueryListener = &BaseQueryListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseQueryListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseQueryListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseQueryListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseQueryListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterPredicate is called when production predicate is entered.
func (s *BaseQueryListener) EnterPredicate(ctx *PredicateContext) {}

// ExitPredicate is called when production predicate is exited.
func (s *BaseQueryListener) ExitPredicate(ctx *PredicateContext) {}

// EnterCompoundPredicate is called when production compoundPredicate is entered.
func (s *BaseQueryListener) EnterCompoundPredicate(ctx *CompoundPredicateContext) {}

// ExitCompoundPredicate is called when production compoundPredicate is exited.
func (s *BaseQueryListener) ExitCompoundPredicate(ctx *CompoundPredicateContext) {}

// EnterAtomicPredicate is called when production atomicPredicate is entered.
func (s *BaseQueryListener) EnterAtomicPredicate(ctx *AtomicPredicateContext) {}

// ExitAtomicPredicate is called when production atomicPredicate is exited.
func (s *BaseQueryListener) ExitAtomicPredicate(ctx *AtomicPredicateContext) {}

// EnterFieldPath is called when production fieldPath is entered.
func (s *BaseQueryListener) EnterFieldPath(ctx *FieldPathContext) {}

// ExitFieldPath is called when production fieldPath is exited.
func (s *BaseQueryListener) ExitFieldPath(ctx *FieldPathContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BaseQueryListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BaseQueryListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterValue is called when production value is entered.
func (s *BaseQueryListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BaseQueryListener) ExitValue(ctx *ValueContext) {}

// EnterValueList is called when production valueList is entered.
func (s *BaseQueryListener) EnterValueList(ctx *ValueListContext) {}

// ExitValueList is called when production valueList is exited.
func (s *BaseQueryListener) ExitValueList(ctx *ValueListContext) {}

// EnterStringValue is called when production stringValue is entered.
func (s *BaseQueryListener) EnterStringValue(ctx *StringValueContext) {}

// ExitStringValue is called when production stringValue is exited.
func (s *BaseQueryListener) ExitStringValue(ctx *StringValueContext) {}

// EnterNumericValue is called when production numericValue is entered.
func (s *BaseQueryListener) EnterNumericValue(ctx *NumericValueContext) {}

// ExitNumericValue is called when production numericValue is exited.
func (s *BaseQueryListener) ExitNumericValue(ctx *NumericValueContext) {}

// EnterBooleanValue is called when production booleanValue is entered.
func (s *BaseQueryListener) EnterBooleanValue(ctx *BooleanValueContext) {}

// ExitBooleanValue is called when production booleanValue is exited.
func (s *BaseQueryListener) ExitBooleanValue(ctx *BooleanValueContext) {}
