package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
)

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

	program, err := expr.Compile(exprStr, expr.Env(env), expr.WithContext("_goctx_"))
	fmt.Println(program, err)
}
