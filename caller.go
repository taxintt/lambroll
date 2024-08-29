package lambroll

import (
	"context"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type CallerIdentity struct {
	data     map[string]any
	Resolver func(ctx context.Context) (*sts.GetCallerIdentityOutput, error)
}

func newCallerIdentity(cfg aws.Config) *CallerIdentity {
	return &CallerIdentity{
		Resolver: func(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
			return sts.NewFromConfig(cfg).GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		},
	}
}

func (c *CallerIdentity) resolve(ctx context.Context) error {
	if c.data != nil {
		return nil
	}
	res, err := c.Resolver(ctx)
	if err != nil {
		return err
	}
	c.data = map[string]any{
		"Account": *res.Account,
		"Arn":     *res.Arn,
		"UserId":  *res.UserId,
	}
	return nil
}

func (c *CallerIdentity) Account(ctx context.Context) string {
	if err := c.resolve(ctx); err != nil {
		return ""
	}
	return c.data["Account"].(string)
}

func (c *CallerIdentity) JsonnetNativeFuncs(ctx context.Context) []*jsonnet.NativeFunction {
	return []*jsonnet.NativeFunction{
		{
			Name:   "caller_identity",
			Params: []ast.Identifier{},
			Func: func(params []any) (any, error) {
				if err := c.resolve(ctx); err != nil {
					return nil, err
				}
				return c.data, nil
			},
		},
	}
}

func (c *CallerIdentity) FuncMap(ctx context.Context) template.FuncMap {
	return template.FuncMap{
		"caller_identity": func() map[string]any {
			if err := c.resolve(ctx); err != nil {
				return nil
			}
			return c.data
		},
	}
}
