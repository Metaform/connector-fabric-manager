// Code generated from Query.g4 by ANTLR 4.13.2. DO NOT EDIT.

package query // Query
import "github.com/antlr4-go/antlr/v4"

type BaseQueryVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseQueryVisitor) VisitPredicate(ctx *PredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitCompoundPredicate(ctx *CompoundPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitAtomicPredicate(ctx *AtomicPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitFieldPath(ctx *FieldPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitValue(ctx *ValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitValueList(ctx *ValueListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitStringValue(ctx *StringValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitNumericValue(ctx *NumericValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseQueryVisitor) VisitBooleanValue(ctx *BooleanValueContext) interface{} {
	return v.VisitChildren(ctx)
}
