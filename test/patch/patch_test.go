package patch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/test/mock"
)

type lengthPatcher struct{}

func (p *lengthPatcher) Visit(node *ast.Node) {
	switch n := (*node).(type) {
	case *ast.MemberNode:
		if prop, ok := n.Property.(*ast.StringNode); ok && prop.Value == "length" {
			ast.Patch(node, &ast.BuiltinNode{
				Name:      "len",
				Arguments: []ast.Node{n.Node},
			})
		}
	}
}

func TestPatch_length(t *testing.T) {
	program, err := expr.Compile(
		`String.length == 5`,
		expr.Env(mock.Env{}),
		expr.Patch(&lengthPatcher{}),
	)
	require.NoError(t, err)

	env := mock.Env{String: "hello"}
	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

type X struct{}

func (x *X) HelloCtx(ctx context.Context, text string) error {
	fmt.Println("hello:", text)
	return nil
}

func TestGoexrEngine(t *testing.T) {
	env := map[string]any{
		"_goctx_": context.TODO(),
		"_g_": map[string]*X{
			"rpc": &X{},
		},
		"text": "gonghuan",
	}

	exprStr := `let v = _g_.rpc.HelloCtx(_goctx_, text); v`

	_, err := expr.Compile(exprStr, expr.Env(env), expr.WithContext("_goctx_"))
	require.NoError(t, err)
}
