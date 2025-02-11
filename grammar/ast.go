package grammar

import (
	"fmt"
	"io"
	"strings"
)

// TODO - Probably should make most of what is in here un-exported

type Expression interface {
	ExpressionDump(w io.Writer, indent string, level int)
}

type UnaryOperator int

const (
	UnaryOpNot UnaryOperator = iota
)

func (op UnaryOperator) String() string {
	switch op {
	case UnaryOpNot:
		return "Not"
	default:
		return "UNKNOWN"
	}
}

type BinaryOperator int

const (
	BinaryOpAnd BinaryOperator = iota
	BinaryOpOr
)

func (op BinaryOperator) String() string {
	switch op {
	case BinaryOpAnd:
		return "And"
	case BinaryOpOr:
		return "Or"
	default:
		return "UNKNOWN"
	}
}

type MatchOperator int

const (
	MatchEqual MatchOperator = iota
	MatchNotEqual
	MatchIn
	MatchNotIn
	MatchIsEmpty
	MatchIsNotEmpty
	MatchMatches
	MatchNotMatches
	MatchGreaterOrEqualThan
	MatchGreaterThan
	MatchLesserOrEqualThan
	MatchLesserThan
)

func (op MatchOperator) String() string {
	switch op {
	case MatchEqual:
		return "Equal"
	case MatchNotEqual:
		return "Not Equal"
	case MatchIn:
		return "In"
	case MatchNotIn:
		return "Not In"
	case MatchIsEmpty:
		return "Is Empty"
	case MatchIsNotEmpty:
		return "Is Not Empty"
	case MatchMatches:
		return "Matches"
	case MatchNotMatches:
		return "Not Matches"
	case MatchGreaterThan:
		return "Greater Than"
	case MatchGreaterOrEqualThan:
		return "Greater or Equal Than"
	case MatchLesserThan:
		return "Lesser Than"
	case MatchLesserOrEqualThan:
		return "Lesser Than"
	default:
		return "UNKNOWN"
	}
}

type MatchValue struct {
	Raw       string
	Converted interface{}
}

type UnaryExpression struct {
	Operator UnaryOperator
	Operand  Expression
}

type BinaryExpression struct {
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

type SelectorType uint32

const (
	SelectorTypeUnknown = iota
	SelectorTypeBexpr
	SelectorTypeJsonPointer
)

type Selector struct {
	Type SelectorType
	Path []string
}

func (sel Selector) String() string {
	if len(sel.Path) == 0 {
		return ""
	}
	switch sel.Type {
	case SelectorTypeBexpr:
		return strings.Join(sel.Path, ".")
	case SelectorTypeJsonPointer:
		return strings.Join(sel.Path, "/")
	default:
		return ""
	}
}

type MatchExpression struct {
	Selector Selector
	Operator MatchOperator
	Value    *MatchValue
}
type CollectionNameBinding struct {
	Mode    CollectionBindMode
	Default string
	Index   string
	Value   string
}

func (b *CollectionNameBinding) String() string {
	switch b.Mode {
	case CollectionBindDefault:
		return fmt.Sprintf("Default ( %s )", b.Default)
	case CollectionBindIndex:
		return fmt.Sprintf("Index ( %s )", b.Index)
	case CollectionBindValue:
		return fmt.Sprintf("Value ( %s )", b.Value)
	case CollectionBindIndexAndValue:
		return fmt.Sprintf("Index & Value ( %s, %s )", b.Index, b.Value)
	default:
		return "UNKNOWN"
	}
}

type CollectionExpression struct {
	Operator    CollectionOperator
	Selector    Selector
	NameBinding CollectionNameBinding
	Expression  Expression
}

func (expr *CollectionExpression) GetSelector() Selector {
	return expr.Selector
}

func (expr *CollectionExpression) SetSelector(sel Selector) {
	expr.Selector = sel
}

type CollectionBindMode int

const (
	CollectionBindDefault CollectionBindMode = iota
	CollectionBindIndex
	CollectionBindValue
	CollectionBindIndexAndValue
)

type CollectionOperator int

const (
	CollectionOpAny CollectionOperator = iota
	CollectionOpAll
)

func (op CollectionOperator) String() string {
	switch op {
	case CollectionOpAny:
		return "any"
	case CollectionOpAll:
		return "all"
	default:
		return "UNKNOWN"
	}
}

func (expr *UnaryExpression) ExpressionDump(w io.Writer, indent string, level int) {
	localIndent := strings.Repeat(indent, level)
	fmt.Fprintf(w, "%s%s {\n", localIndent, expr.Operator.String())
	expr.Operand.ExpressionDump(w, indent, level+1)
	fmt.Fprintf(w, "%s}\n", localIndent)
}

func (expr *BinaryExpression) ExpressionDump(w io.Writer, indent string, level int) {
	localIndent := strings.Repeat(indent, level)
	fmt.Fprintf(w, "%s%s {\n", localIndent, expr.Operator.String())
	expr.Left.ExpressionDump(w, indent, level+1)
	expr.Right.ExpressionDump(w, indent, level+1)
	fmt.Fprintf(w, "%s}\n", localIndent)
}

func (expr *MatchExpression) ExpressionDump(w io.Writer, indent string, level int) {
	switch expr.Operator {
	case MatchEqual, MatchNotEqual, MatchIn, MatchNotIn, MatchLesserOrEqualThan, MatchGreaterThan, MatchGreaterOrEqualThan, MatchLesserThan:
		fmt.Fprintf(w, "%[1]s%[3]s {\n%[2]sSelector: %[4]v\n%[2]sValue: %[5]q\n%[1]s}\n", strings.Repeat(indent, level), strings.Repeat(indent, level+1), expr.Operator.String(), expr.Selector, expr.Value.Raw)
	default:
		fmt.Fprintf(w, "%[1]s%[3]s {\n%[2]sSelector: %[4]v\n%[1]s}\n", strings.Repeat(indent, level), strings.Repeat(indent, level+1), expr.Operator.String(), expr.Selector)
	}
}
func (expr *CollectionExpression) ExpressionDump(w io.Writer, indent string, level int) {
	localIndent := strings.Repeat(indent, level)
	fmt.Fprintf(w, "%[2]s%[3]s {\n%[1]s%[2]sSelector: %[4]v\n%[1]s%[2]sNameBinding: %[5]s\n", indent, localIndent, expr.Operator.String(), expr.Selector, expr.NameBinding.String())
	expr.Expression.ExpressionDump(w, indent, level+2)
	fmt.Fprintf(w, "%s}\n", localIndent)
}
