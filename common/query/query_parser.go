// Code generated from Query.g4 by ANTLR 4.13.2. DO NOT EDIT.

package query // Query
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type QueryParser struct {
	*antlr.BaseParser
}

var QueryParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func queryParserInit() {
	staticData := &QueryParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "'='", "'!='",
		"'>'", "'>='", "'<'", "'<='", "'('", "')'", "'.'", "','",
	}
	staticData.SymbolicNames = []string{
		"", "AND", "OR", "IS", "NOT", "IN", "NULL_KEYWORD", "TRUE_KEYWORD",
		"FALSE_KEYWORD", "LIKE", "NOT_LIKE", "CONTAINS", "STARTS_WITH", "ENDS_WITH",
		"EQ", "NE", "GT", "GTE", "LT", "LTE", "LPAREN", "RPAREN", "DOT", "COMMA",
		"STRING", "NUMBER", "DIGIT", "IDENTIFIER", "WS",
	}
	staticData.RuleNames = []string{
		"predicate", "compoundPredicate", "atomicPredicate", "fieldPath", "comparisonOperator",
		"value", "valueList", "stringValue", "numericValue", "booleanValue",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 28, 100, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0, 27, 8, 0, 1, 1, 1, 1, 1, 1, 4, 1, 32,
		8, 1, 11, 1, 12, 1, 33, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 40, 8, 1, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 3, 2, 68, 8, 2, 1, 3, 1, 3, 1, 3, 5, 3, 73, 8, 3, 10, 3, 12, 3, 76,
		9, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 84, 8, 5, 1, 6, 1, 6, 1,
		6, 5, 6, 89, 8, 6, 10, 6, 12, 6, 92, 9, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9,
		1, 9, 1, 9, 0, 0, 10, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 0, 3, 1, 0, 1,
		2, 1, 0, 9, 19, 1, 0, 7, 8, 101, 0, 26, 1, 0, 0, 0, 2, 39, 1, 0, 0, 0,
		4, 67, 1, 0, 0, 0, 6, 69, 1, 0, 0, 0, 8, 77, 1, 0, 0, 0, 10, 83, 1, 0,
		0, 0, 12, 85, 1, 0, 0, 0, 14, 93, 1, 0, 0, 0, 16, 95, 1, 0, 0, 0, 18, 97,
		1, 0, 0, 0, 20, 21, 3, 2, 1, 0, 21, 22, 5, 0, 0, 1, 22, 27, 1, 0, 0, 0,
		23, 24, 3, 4, 2, 0, 24, 25, 5, 0, 0, 1, 25, 27, 1, 0, 0, 0, 26, 20, 1,
		0, 0, 0, 26, 23, 1, 0, 0, 0, 27, 1, 1, 0, 0, 0, 28, 31, 3, 4, 2, 0, 29,
		30, 7, 0, 0, 0, 30, 32, 3, 4, 2, 0, 31, 29, 1, 0, 0, 0, 32, 33, 1, 0, 0,
		0, 33, 31, 1, 0, 0, 0, 33, 34, 1, 0, 0, 0, 34, 40, 1, 0, 0, 0, 35, 36,
		5, 20, 0, 0, 36, 37, 3, 2, 1, 0, 37, 38, 5, 21, 0, 0, 38, 40, 1, 0, 0,
		0, 39, 28, 1, 0, 0, 0, 39, 35, 1, 0, 0, 0, 40, 3, 1, 0, 0, 0, 41, 42, 3,
		6, 3, 0, 42, 43, 3, 8, 4, 0, 43, 44, 3, 10, 5, 0, 44, 68, 1, 0, 0, 0, 45,
		46, 3, 6, 3, 0, 46, 47, 5, 3, 0, 0, 47, 48, 5, 6, 0, 0, 48, 68, 1, 0, 0,
		0, 49, 50, 3, 6, 3, 0, 50, 51, 5, 3, 0, 0, 51, 52, 5, 4, 0, 0, 52, 53,
		5, 6, 0, 0, 53, 68, 1, 0, 0, 0, 54, 55, 3, 6, 3, 0, 55, 56, 5, 5, 0, 0,
		56, 57, 5, 20, 0, 0, 57, 58, 3, 12, 6, 0, 58, 59, 5, 21, 0, 0, 59, 68,
		1, 0, 0, 0, 60, 61, 3, 6, 3, 0, 61, 62, 5, 4, 0, 0, 62, 63, 5, 5, 0, 0,
		63, 64, 5, 20, 0, 0, 64, 65, 3, 12, 6, 0, 65, 66, 5, 21, 0, 0, 66, 68,
		1, 0, 0, 0, 67, 41, 1, 0, 0, 0, 67, 45, 1, 0, 0, 0, 67, 49, 1, 0, 0, 0,
		67, 54, 1, 0, 0, 0, 67, 60, 1, 0, 0, 0, 68, 5, 1, 0, 0, 0, 69, 74, 5, 27,
		0, 0, 70, 71, 5, 22, 0, 0, 71, 73, 5, 27, 0, 0, 72, 70, 1, 0, 0, 0, 73,
		76, 1, 0, 0, 0, 74, 72, 1, 0, 0, 0, 74, 75, 1, 0, 0, 0, 75, 7, 1, 0, 0,
		0, 76, 74, 1, 0, 0, 0, 77, 78, 7, 1, 0, 0, 78, 9, 1, 0, 0, 0, 79, 84, 3,
		14, 7, 0, 80, 84, 3, 16, 8, 0, 81, 84, 3, 18, 9, 0, 82, 84, 5, 6, 0, 0,
		83, 79, 1, 0, 0, 0, 83, 80, 1, 0, 0, 0, 83, 81, 1, 0, 0, 0, 83, 82, 1,
		0, 0, 0, 84, 11, 1, 0, 0, 0, 85, 90, 3, 10, 5, 0, 86, 87, 5, 23, 0, 0,
		87, 89, 3, 10, 5, 0, 88, 86, 1, 0, 0, 0, 89, 92, 1, 0, 0, 0, 90, 88, 1,
		0, 0, 0, 90, 91, 1, 0, 0, 0, 91, 13, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 93,
		94, 5, 24, 0, 0, 94, 15, 1, 0, 0, 0, 95, 96, 5, 25, 0, 0, 96, 17, 1, 0,
		0, 0, 97, 98, 7, 2, 0, 0, 98, 19, 1, 0, 0, 0, 7, 26, 33, 39, 67, 74, 83,
		90,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// QueryParserInit initializes any static state used to implement QueryParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewQueryParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func QueryParserInit() {
	staticData := &QueryParserStaticData
	staticData.once.Do(queryParserInit)
}

// NewQueryParser produces a new parser instance for the optional input antlr.TokenStream.
func NewQueryParser(input antlr.TokenStream) *QueryParser {
	QueryParserInit()
	this := new(QueryParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &QueryParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Query.g4"

	return this
}

// QueryParser tokens.
const (
	QueryParserEOF           = antlr.TokenEOF
	QueryParserAND           = 1
	QueryParserOR            = 2
	QueryParserIS            = 3
	QueryParserNOT           = 4
	QueryParserIN            = 5
	QueryParserNULL_KEYWORD  = 6
	QueryParserTRUE_KEYWORD  = 7
	QueryParserFALSE_KEYWORD = 8
	QueryParserLIKE          = 9
	QueryParserNOT_LIKE      = 10
	QueryParserCONTAINS      = 11
	QueryParserSTARTS_WITH   = 12
	QueryParserENDS_WITH     = 13
	QueryParserEQ            = 14
	QueryParserNE            = 15
	QueryParserGT            = 16
	QueryParserGTE           = 17
	QueryParserLT            = 18
	QueryParserLTE           = 19
	QueryParserLPAREN        = 20
	QueryParserRPAREN        = 21
	QueryParserDOT           = 22
	QueryParserCOMMA         = 23
	QueryParserSTRING        = 24
	QueryParserNUMBER        = 25
	QueryParserDIGIT         = 26
	QueryParserIDENTIFIER    = 27
	QueryParserWS            = 28
)

// QueryParser rules.
const (
	QueryParserRULE_predicate          = 0
	QueryParserRULE_compoundPredicate  = 1
	QueryParserRULE_atomicPredicate    = 2
	QueryParserRULE_fieldPath          = 3
	QueryParserRULE_comparisonOperator = 4
	QueryParserRULE_value              = 5
	QueryParserRULE_valueList          = 6
	QueryParserRULE_stringValue        = 7
	QueryParserRULE_numericValue       = 8
	QueryParserRULE_booleanValue       = 9
)

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompoundPredicate() ICompoundPredicateContext
	EOF() antlr.TerminalNode
	AtomicPredicate() IAtomicPredicateContext

	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_predicate
	return p
}

func InitEmptyPredicateContext(p *PredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_predicate
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) CompoundPredicate() ICompoundPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompoundPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompoundPredicateContext)
}

func (s *PredicateContext) EOF() antlr.TerminalNode {
	return s.GetToken(QueryParserEOF, 0)
}

func (s *PredicateContext) AtomicPredicate() IAtomicPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAtomicPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAtomicPredicateContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterPredicate(s)
	}
}

func (s *PredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitPredicate(s)
	}
}

func (s *PredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) Predicate() (localctx IPredicateContext) {
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, QueryParserRULE_predicate)
	p.SetState(26)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(20)
			p.CompoundPredicate()
		}
		{
			p.SetState(21)
			p.Match(QueryParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(23)
			p.AtomicPredicate()
		}
		{
			p.SetState(24)
			p.Match(QueryParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompoundPredicateContext is an interface to support dynamic dispatch.
type ICompoundPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAtomicPredicate() []IAtomicPredicateContext
	AtomicPredicate(i int) IAtomicPredicateContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	CompoundPredicate() ICompoundPredicateContext
	RPAREN() antlr.TerminalNode

	// IsCompoundPredicateContext differentiates from other interfaces.
	IsCompoundPredicateContext()
}

type CompoundPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompoundPredicateContext() *CompoundPredicateContext {
	var p = new(CompoundPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_compoundPredicate
	return p
}

func InitEmptyCompoundPredicateContext(p *CompoundPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_compoundPredicate
}

func (*CompoundPredicateContext) IsCompoundPredicateContext() {}

func NewCompoundPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundPredicateContext {
	var p = new(CompoundPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_compoundPredicate

	return p
}

func (s *CompoundPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *CompoundPredicateContext) AllAtomicPredicate() []IAtomicPredicateContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAtomicPredicateContext); ok {
			len++
		}
	}

	tst := make([]IAtomicPredicateContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAtomicPredicateContext); ok {
			tst[i] = t.(IAtomicPredicateContext)
			i++
		}
	}

	return tst
}

func (s *CompoundPredicateContext) AtomicPredicate(i int) IAtomicPredicateContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAtomicPredicateContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAtomicPredicateContext)
}

func (s *CompoundPredicateContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(QueryParserAND)
}

func (s *CompoundPredicateContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(QueryParserAND, i)
}

func (s *CompoundPredicateContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(QueryParserOR)
}

func (s *CompoundPredicateContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(QueryParserOR, i)
}

func (s *CompoundPredicateContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(QueryParserLPAREN, 0)
}

func (s *CompoundPredicateContext) CompoundPredicate() ICompoundPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompoundPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompoundPredicateContext)
}

func (s *CompoundPredicateContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(QueryParserRPAREN, 0)
}

func (s *CompoundPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompoundPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterCompoundPredicate(s)
	}
}

func (s *CompoundPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitCompoundPredicate(s)
	}
}

func (s *CompoundPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitCompoundPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) CompoundPredicate() (localctx ICompoundPredicateContext) {
	localctx = NewCompoundPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, QueryParserRULE_compoundPredicate)
	var _la int

	p.SetState(39)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case QueryParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(28)
			p.AtomicPredicate()
		}
		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == QueryParserAND || _la == QueryParserOR {
			{
				p.SetState(29)
				_la = p.GetTokenStream().LA(1)

				if !(_la == QueryParserAND || _la == QueryParserOR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(30)
				p.AtomicPredicate()
			}

			p.SetState(33)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case QueryParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(35)
			p.Match(QueryParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(36)
			p.CompoundPredicate()
		}
		{
			p.SetState(37)
			p.Match(QueryParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAtomicPredicateContext is an interface to support dynamic dispatch.
type IAtomicPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FieldPath() IFieldPathContext
	ComparisonOperator() IComparisonOperatorContext
	Value() IValueContext
	IS() antlr.TerminalNode
	NULL_KEYWORD() antlr.TerminalNode
	NOT() antlr.TerminalNode
	IN() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ValueList() IValueListContext
	RPAREN() antlr.TerminalNode

	// IsAtomicPredicateContext differentiates from other interfaces.
	IsAtomicPredicateContext()
}

type AtomicPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtomicPredicateContext() *AtomicPredicateContext {
	var p = new(AtomicPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_atomicPredicate
	return p
}

func InitEmptyAtomicPredicateContext(p *AtomicPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_atomicPredicate
}

func (*AtomicPredicateContext) IsAtomicPredicateContext() {}

func NewAtomicPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomicPredicateContext {
	var p = new(AtomicPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_atomicPredicate

	return p
}

func (s *AtomicPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomicPredicateContext) FieldPath() IFieldPathContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldPathContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldPathContext)
}

func (s *AtomicPredicateContext) ComparisonOperator() IComparisonOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonOperatorContext)
}

func (s *AtomicPredicateContext) Value() IValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *AtomicPredicateContext) IS() antlr.TerminalNode {
	return s.GetToken(QueryParserIS, 0)
}

func (s *AtomicPredicateContext) NULL_KEYWORD() antlr.TerminalNode {
	return s.GetToken(QueryParserNULL_KEYWORD, 0)
}

func (s *AtomicPredicateContext) NOT() antlr.TerminalNode {
	return s.GetToken(QueryParserNOT, 0)
}

func (s *AtomicPredicateContext) IN() antlr.TerminalNode {
	return s.GetToken(QueryParserIN, 0)
}

func (s *AtomicPredicateContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(QueryParserLPAREN, 0)
}

func (s *AtomicPredicateContext) ValueList() IValueListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueListContext)
}

func (s *AtomicPredicateContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(QueryParserRPAREN, 0)
}

func (s *AtomicPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomicPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtomicPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterAtomicPredicate(s)
	}
}

func (s *AtomicPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitAtomicPredicate(s)
	}
}

func (s *AtomicPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitAtomicPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) AtomicPredicate() (localctx IAtomicPredicateContext) {
	localctx = NewAtomicPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, QueryParserRULE_atomicPredicate)
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(41)
			p.FieldPath()
		}
		{
			p.SetState(42)
			p.ComparisonOperator()
		}
		{
			p.SetState(43)
			p.Value()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(45)
			p.FieldPath()
		}
		{
			p.SetState(46)
			p.Match(QueryParserIS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(47)
			p.Match(QueryParserNULL_KEYWORD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(49)
			p.FieldPath()
		}
		{
			p.SetState(50)
			p.Match(QueryParserIS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(51)
			p.Match(QueryParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(52)
			p.Match(QueryParserNULL_KEYWORD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(54)
			p.FieldPath()
		}
		{
			p.SetState(55)
			p.Match(QueryParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(56)
			p.Match(QueryParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(57)
			p.ValueList()
		}
		{
			p.SetState(58)
			p.Match(QueryParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(60)
			p.FieldPath()
		}
		{
			p.SetState(61)
			p.Match(QueryParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(62)
			p.Match(QueryParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(63)
			p.Match(QueryParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(64)
			p.ValueList()
		}
		{
			p.SetState(65)
			p.Match(QueryParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldPathContext is an interface to support dynamic dispatch.
type IFieldPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsFieldPathContext differentiates from other interfaces.
	IsFieldPathContext()
}

type FieldPathContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldPathContext() *FieldPathContext {
	var p = new(FieldPathContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_fieldPath
	return p
}

func InitEmptyFieldPathContext(p *FieldPathContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_fieldPath
}

func (*FieldPathContext) IsFieldPathContext() {}

func NewFieldPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldPathContext {
	var p = new(FieldPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_fieldPath

	return p
}

func (s *FieldPathContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldPathContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(QueryParserIDENTIFIER)
}

func (s *FieldPathContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(QueryParserIDENTIFIER, i)
}

func (s *FieldPathContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(QueryParserDOT)
}

func (s *FieldPathContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(QueryParserDOT, i)
}

func (s *FieldPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterFieldPath(s)
	}
}

func (s *FieldPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitFieldPath(s)
	}
}

func (s *FieldPathContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitFieldPath(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) FieldPath() (localctx IFieldPathContext) {
	localctx = NewFieldPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, QueryParserRULE_fieldPath)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(QueryParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(74)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == QueryParserDOT {
		{
			p.SetState(70)
			p.Match(QueryParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(71)
			p.Match(QueryParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(76)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonOperatorContext is an interface to support dynamic dispatch.
type IComparisonOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EQ() antlr.TerminalNode
	NE() antlr.TerminalNode
	GT() antlr.TerminalNode
	GTE() antlr.TerminalNode
	LT() antlr.TerminalNode
	LTE() antlr.TerminalNode
	LIKE() antlr.TerminalNode
	NOT_LIKE() antlr.TerminalNode
	CONTAINS() antlr.TerminalNode
	STARTS_WITH() antlr.TerminalNode
	ENDS_WITH() antlr.TerminalNode

	// IsComparisonOperatorContext differentiates from other interfaces.
	IsComparisonOperatorContext()
}

type ComparisonOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonOperatorContext() *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_comparisonOperator
	return p
}

func InitEmptyComparisonOperatorContext(p *ComparisonOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_comparisonOperator
}

func (*ComparisonOperatorContext) IsComparisonOperatorContext() {}

func NewComparisonOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_comparisonOperator

	return p
}

func (s *ComparisonOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonOperatorContext) EQ() antlr.TerminalNode {
	return s.GetToken(QueryParserEQ, 0)
}

func (s *ComparisonOperatorContext) NE() antlr.TerminalNode {
	return s.GetToken(QueryParserNE, 0)
}

func (s *ComparisonOperatorContext) GT() antlr.TerminalNode {
	return s.GetToken(QueryParserGT, 0)
}

func (s *ComparisonOperatorContext) GTE() antlr.TerminalNode {
	return s.GetToken(QueryParserGTE, 0)
}

func (s *ComparisonOperatorContext) LT() antlr.TerminalNode {
	return s.GetToken(QueryParserLT, 0)
}

func (s *ComparisonOperatorContext) LTE() antlr.TerminalNode {
	return s.GetToken(QueryParserLTE, 0)
}

func (s *ComparisonOperatorContext) LIKE() antlr.TerminalNode {
	return s.GetToken(QueryParserLIKE, 0)
}

func (s *ComparisonOperatorContext) NOT_LIKE() antlr.TerminalNode {
	return s.GetToken(QueryParserNOT_LIKE, 0)
}

func (s *ComparisonOperatorContext) CONTAINS() antlr.TerminalNode {
	return s.GetToken(QueryParserCONTAINS, 0)
}

func (s *ComparisonOperatorContext) STARTS_WITH() antlr.TerminalNode {
	return s.GetToken(QueryParserSTARTS_WITH, 0)
}

func (s *ComparisonOperatorContext) ENDS_WITH() antlr.TerminalNode {
	return s.GetToken(QueryParserENDS_WITH, 0)
}

func (s *ComparisonOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitComparisonOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) ComparisonOperator() (localctx IComparisonOperatorContext) {
	localctx = NewComparisonOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, QueryParserRULE_comparisonOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(77)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1048064) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StringValue() IStringValueContext
	NumericValue() INumericValueContext
	BooleanValue() IBooleanValueContext
	NULL_KEYWORD() antlr.TerminalNode

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_value
	return p
}

func InitEmptyValueContext(p *ValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_value
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) StringValue() IStringValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringValueContext)
}

func (s *ValueContext) NumericValue() INumericValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumericValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumericValueContext)
}

func (s *ValueContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *ValueContext) NULL_KEYWORD() antlr.TerminalNode {
	return s.GetToken(QueryParserNULL_KEYWORD, 0)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitValue(s)
	}
}

func (s *ValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, QueryParserRULE_value)
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case QueryParserSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(79)
			p.StringValue()
		}

	case QueryParserNUMBER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(80)
			p.NumericValue()
		}

	case QueryParserTRUE_KEYWORD, QueryParserFALSE_KEYWORD:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(81)
			p.BooleanValue()
		}

	case QueryParserNULL_KEYWORD:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(82)
			p.Match(QueryParserNULL_KEYWORD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueListContext is an interface to support dynamic dispatch.
type IValueListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllValue() []IValueContext
	Value(i int) IValueContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsValueListContext differentiates from other interfaces.
	IsValueListContext()
}

type ValueListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueListContext() *ValueListContext {
	var p = new(ValueListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_valueList
	return p
}

func InitEmptyValueListContext(p *ValueListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_valueList
}

func (*ValueListContext) IsValueListContext() {}

func NewValueListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueListContext {
	var p = new(ValueListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_valueList

	return p
}

func (s *ValueListContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueListContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *ValueListContext) Value(i int) IValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ValueListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(QueryParserCOMMA)
}

func (s *ValueListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(QueryParserCOMMA, i)
}

func (s *ValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterValueList(s)
	}
}

func (s *ValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitValueList(s)
	}
}

func (s *ValueListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitValueList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) ValueList() (localctx IValueListContext) {
	localctx = NewValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, QueryParserRULE_valueList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		p.Value()
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == QueryParserCOMMA {
		{
			p.SetState(86)
			p.Match(QueryParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(87)
			p.Value()
		}

		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringValueContext is an interface to support dynamic dispatch.
type IStringValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING() antlr.TerminalNode

	// IsStringValueContext differentiates from other interfaces.
	IsStringValueContext()
}

type StringValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringValueContext() *StringValueContext {
	var p = new(StringValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_stringValue
	return p
}

func InitEmptyStringValueContext(p *StringValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_stringValue
}

func (*StringValueContext) IsStringValueContext() {}

func NewStringValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringValueContext {
	var p = new(StringValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_stringValue

	return p
}

func (s *StringValueContext) GetParser() antlr.Parser { return s.parser }

func (s *StringValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(QueryParserSTRING, 0)
}

func (s *StringValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterStringValue(s)
	}
}

func (s *StringValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitStringValue(s)
	}
}

func (s *StringValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitStringValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) StringValue() (localctx IStringValueContext) {
	localctx = NewStringValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, QueryParserRULE_stringValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		p.Match(QueryParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INumericValueContext is an interface to support dynamic dispatch.
type INumericValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode

	// IsNumericValueContext differentiates from other interfaces.
	IsNumericValueContext()
}

type NumericValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumericValueContext() *NumericValueContext {
	var p = new(NumericValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_numericValue
	return p
}

func InitEmptyNumericValueContext(p *NumericValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_numericValue
}

func (*NumericValueContext) IsNumericValueContext() {}

func NewNumericValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumericValueContext {
	var p = new(NumericValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_numericValue

	return p
}

func (s *NumericValueContext) GetParser() antlr.Parser { return s.parser }

func (s *NumericValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(QueryParserNUMBER, 0)
}

func (s *NumericValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumericValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumericValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterNumericValue(s)
	}
}

func (s *NumericValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitNumericValue(s)
	}
}

func (s *NumericValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitNumericValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) NumericValue() (localctx INumericValueContext) {
	localctx = NewNumericValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, QueryParserRULE_numericValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		p.Match(QueryParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanValueContext is an interface to support dynamic dispatch.
type IBooleanValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUE_KEYWORD() antlr.TerminalNode
	FALSE_KEYWORD() antlr.TerminalNode

	// IsBooleanValueContext differentiates from other interfaces.
	IsBooleanValueContext()
}

type BooleanValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanValueContext() *BooleanValueContext {
	var p = new(BooleanValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_booleanValue
	return p
}

func InitEmptyBooleanValueContext(p *BooleanValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = QueryParserRULE_booleanValue
}

func (*BooleanValueContext) IsBooleanValueContext() {}

func NewBooleanValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanValueContext {
	var p = new(BooleanValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = QueryParserRULE_booleanValue

	return p
}

func (s *BooleanValueContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanValueContext) TRUE_KEYWORD() antlr.TerminalNode {
	return s.GetToken(QueryParserTRUE_KEYWORD, 0)
}

func (s *BooleanValueContext) FALSE_KEYWORD() antlr.TerminalNode {
	return s.GetToken(QueryParserFALSE_KEYWORD, 0)
}

func (s *BooleanValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.EnterBooleanValue(s)
	}
}

func (s *BooleanValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(QueryListener); ok {
		listenerT.ExitBooleanValue(s)
	}
}

func (s *BooleanValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case QueryVisitor:
		return t.VisitBooleanValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *QueryParser) BooleanValue() (localctx IBooleanValueContext) {
	localctx = NewBooleanValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, QueryParserRULE_booleanValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(97)
		_la = p.GetTokenStream().LA(1)

		if !(_la == QueryParserTRUE_KEYWORD || _la == QueryParserFALSE_KEYWORD) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
