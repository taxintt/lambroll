package lambroll

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/kylelemons/godebug/diff"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// DiffOption represents options for Diff()
type DiffOption struct {
	Src        string `help:"function zip archive or src dir" default:"."`
	CodeSha256 bool   `help:"diff of code sha256" default:"false"`

	ExcludeFileOption
}

// Diff prints diff of function.json compared with latest function
func (app *App) Diff(ctx context.Context, opt *DiffOption) error {
	if err := opt.Expand(); err != nil {
		return err
	}

	newFunc, err := app.loadFunction(app.functionFilePath)
	if err != nil {
		return fmt.Errorf("failed to load function: %w", err)
	}
	fillDefaultValues(newFunc)
	name := *newFunc.FunctionName

	var latest *types.FunctionConfiguration
	var code *types.FunctionCodeLocation

	var tags Tags
	var currentCodeSha256 string
	var packageType types.PackageType
	if res, err := app.lambda.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: &name,
	}); err != nil {
		return fmt.Errorf("failed to GetFunction %s: %w", name, err)
	} else {
		latest = res.Configuration
		code = res.Code
		tags = res.Tags
		currentCodeSha256 = *res.Configuration.CodeSha256
		packageType = res.Configuration.PackageType
	}
	latestFunc := newFunctionFrom(latest, code, tags)

	latestJSON, _ := marshalJSON(latestFunc)
	newJSON, _ := marshalJSON(newFunc)

	if ds := diff.Diff(string(latestJSON), string(newJSON)); ds != "" {
		fmt.Println(color.RedString("---" + app.functionArn(ctx, name)))
		fmt.Println(color.GreenString("+++" + app.functionFilePath))
		fmt.Println(coloredDiff(ds))
	}

	if err := validateUpdateFunction(latest, code, newFunc); err != nil {
		return err
	}

	if opt.CodeSha256 {
		if packageType != types.PackageTypeZip {
			return fmt.Errorf("code-sha256 is only supported for Zip package type")
		}
		zipfile, _, err := prepareZipfile(opt.Src, opt.excludes)
		if err != nil {
			return err
		}
		h := sha256.New()
		if _, err := io.Copy(h, zipfile); err != nil {
			return err
		}
		newCodeSha256 := base64.StdEncoding.EncodeToString(h.Sum(nil))
		prefix := "CodeSha256: "
		if ds := diff.Diff(prefix+currentCodeSha256, prefix+newCodeSha256); ds != "" {
			fmt.Println(color.RedString("---" + app.functionArn(ctx, name)))
			fmt.Println(color.GreenString("+++" + "--src=" + opt.Src))
			fmt.Println(coloredDiff(ds))
		}
	}

	return nil
}

func coloredDiff(src string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(line, "-") {
			b.WriteString(color.RedString(line) + "\n")
		} else if strings.HasPrefix(line, "+") {
			b.WriteString(color.GreenString(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}
