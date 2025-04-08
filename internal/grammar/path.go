package grammar

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/rra696/scim-filter-parser/v2/internal/types"
)

func Path(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.Path,
		TypeStrings: typ.Stringer,
		Value: op.Or{
			op.And{
				ValuePath,
				op.Optional(SubAttr),
			},
			AttrPath,
		},
	})
}
