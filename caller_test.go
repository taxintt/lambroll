package lambroll_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/fujiwara/lambroll"
)

func TestCallerIdentity(t *testing.T) {
	c := lambroll.NewCallerIdentity(aws.Config{})
	c.Resolver = func(_ context.Context) (*sts.GetCallerIdentityOutput, error) {
		return &sts.GetCallerIdentityOutput{
			Account: aws.String("123456789012"),
			Arn:     aws.String("arn:aws:iam::123456789012:user/test-user"),
			UserId:  aws.String("AIXXXXXXXXX"),
		}, nil
	}
	ctx := context.Background()
	if c.Account(ctx) != "123456789012" {
		t.Errorf("unexpected account id: %s", c.Account(ctx))
	}
}
