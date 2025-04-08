package grammar

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/scim2/filter-parser/v2/internal/types"
)

func ValueFilter(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ValueLogExpOr,
	)
}

func ValueFilterAll(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		ValueFilter,
		ValueFilterNot,
	})
}

func ValueFilterNot(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.ValueFilterNot,
		TypeStrings: typ.Stringer,
		Value: op.And{
			parser.CheckStringCI("not"),
			op.MinZero(SP),
			'(',
			op.MinZero(SP),
			ValueFilter,
			op.MinZero(SP),
			')',
		},
	})
}

func ValueLogExpOr(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.ValueLogExpOr,
		TypeStrings: typ.Stringer,
		Value: op.And{
			ValueLogExpAnd,
			op.MinZero(op.And{
				op.MinZero(SP),
				parser.CheckStringCI("or"),
				op.MinZero(SP),
				ValueLogExpAnd,
			}),
		},
	})
}

func ValueLogExpAnd(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.ValueLogExpAnd,
		TypeStrings: typ.Stringer,
		Value: op.And{
			FilterValue,
			op.MinZero(op.And{
				op.MinZero(SP),
				parser.CheckStringCI("and"),
				op.MinZero(SP),
				FilterValue,
			}),
		},
	})
}

func ValuePath(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.ValuePath,
		TypeStrings: typ.Stringer,
		Value: op.And{
			AttrPath,
			op.MinZero(SP),
			'[',
			op.MinZero(SP),
			ValueFilterAll,
			op.MinZero(SP),
			']',
		},
	})
}
