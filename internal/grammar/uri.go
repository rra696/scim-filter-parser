package grammar

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/rra696/scim-filter-parser/v2/internal/types"
)

func URI(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.URI,
		TypeStrings: typ.Stringer,
		Value: op.MinOne(op.And{
			op.MinOne(op.Or{
				parser.CheckRuneRange('a', 'z'),
				parser.CheckRuneRange('A', 'Z'),
				parser.CheckRuneRange('0', '9'),
				'.',
			}),
			":",
		}),
	})
}
