package filter

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/rra696/scim-filter-parser/v2/internal/grammar"
	typ "github.com/rra696/scim-filter-parser/v2/internal/types"
)

// ParseValuePath parses the given raw data as an ValuePath.
func ParseValuePath(raw []byte) (ValuePath, error) {
	return parseValuePath(raw, config{})
}

// ParseValuePathNumber parses the given raw data as an ValuePath with json.Number.
func ParseValuePathNumber(raw []byte) (ValuePath, error) {
	return parseValuePath(raw, config{useNumber: true})
}

func parseValuePath(raw []byte, c config) (ValuePath, error) {
	p, err := ast.New(raw)
	if err != nil {
		return ValuePath{}, err
	}
	node, err := grammar.ValuePath(p)
	if err != nil {
		return ValuePath{}, err
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		return ValuePath{}, err
	}
	return c.parseValuePath(node)
}

func (p config) parseValueFilter(node *ast.Node) (Expression, error) {
	switch t := node.Type; t {
	case typ.ValueLogExpOr:
		children := node.Children()
		if len(children) == 0 {
			return nil, invalidLengthError(typ.ValueLogExpOr, 1, 0)
		}

		if len(children) == 1 {
			return p.parseValueFilter(children[0])
		}

		var or LogicalExpression
		for _, node := range children {
			exp, err := p.parseValueFilter(node)
			if err != nil {
				return nil, err
			}
			switch {
			case or.Left == nil:
				or.Left = exp
			case or.Right == nil:
				or.Right = exp
				or.Operator = OR
			default:
				or = LogicalExpression{
					Left: &LogicalExpression{
						Left:     or.Left,
						Right:    or.Right,
						Operator: OR,
					},
					Right:    exp,
					Operator: OR,
				}
			}
		}
		return &or, nil
	case typ.ValueLogExpAnd:
		children := node.Children()
		if len(children) == 0 {
			return nil, invalidLengthError(typ.FilterAnd, 1, 0)
		}

		if len(children) == 1 {
			return p.parseFilterValue(children[0])
		}
		var and LogicalExpression
		for _, node := range children {
			exp, err := p.parseFilterValue(node)
			if err != nil {
				return nil, err
			}
			switch {
			case and.Left == nil:
				and.Left = exp
			case and.Right == nil:
				and.Right = exp
				and.Operator = AND
			default:
				and = LogicalExpression{
					Left: &LogicalExpression{
						Left:     and.Left,
						Right:    and.Right,
						Operator: AND,
					},
					Right:    exp,
					Operator: AND,
				}
			}
		}
		return &and, nil
	case typ.ValueFilterNot:
		children := node.Children()
		if l := len(children); l != 1 {
			return nil, invalidLengthError(typ.ValueFilterNot, 1, l)
		}

		valueFilter, err := p.parseValueFilter(children[0])
		if err != nil {
			return nil, err
		}
		return &NotExpression{
			Expression: valueFilter,
		}, nil
	default:
		return nil, invalidChildTypeError(typ.ValuePath, t)
	}
}

func (p config) parseValuePath(node *ast.Node) (ValuePath, error) {
	if node.Type != typ.ValuePath {
		return ValuePath{}, invalidTypeError(typ.ValuePath, node.Type)
	}

	children := node.Children()
	if l := len(children); l != 2 {
		return ValuePath{}, invalidLengthError(typ.ValuePath, 2, l)
	}

	attrPath, err := parseAttrPath(children[0])
	if err != nil {
		return ValuePath{}, err
	}

	valueFilter, err := p.parseValueFilter(children[1])
	if err != nil {
		return ValuePath{}, err
	}

	return ValuePath{
		AttributePath: attrPath,
		ValueFilter:   valueFilter,
	}, nil
}
