package detector

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/wata727/tflint/issue"
)

type TerraformModulePinnedSourceDetector struct {
	source string
	line   int
	file   string
}

func NewTerraformModulePinnedSourceDetector(detector *Detector, file string, item *ast.ObjectItem) *TerraformModulePinnedSourceDetector {
	sourceToken, err := hclLiteralToken(item, "source")
	if err != nil {
		detector.Logger.Error(err)
		return nil
	}
	sourceText, err := detector.evalToString(sourceToken.Text)
	if err != nil {
		detector.Logger.Error(err)
		return nil
	}

	return &TerraformModulePinnedSourceDetector{
		source: sourceText,
		file:   file,
		line:   sourceToken.Pos.Line,
	}
}

func (d *TerraformModulePinnedSourceDetector) DetectPinnedModuleSource(issues *[]*issue.Issue) {
	lower := strings.ToLower(d.source)

	if strings.Contains(lower, "git") || strings.Contains(lower, "bitbucket") {
		if issue := d.detectGitSource(d.source); issue != nil {
			tmp := append(*issues, issue)
			*issues = tmp
		}
	} else if strings.HasPrefix(lower, "hg:") {
		if issue := d.detectMercurialSource(d.source); issue != nil {
			tmp := append(*issues, issue)
			*issues = tmp
		}
	}
}

func (d *TerraformModulePinnedSourceDetector) detectGitSource(source string) *issue.Issue {
	if strings.Contains(source, "ref=") {
		if strings.Contains(source, "ref=master") {
			return &issue.Issue{
				Type:    issue.WARNING,
				Message: fmt.Sprintf("Module source \"%s\" uses default ref \"master\"", source),
				Line:    d.line,
				File:    d.file,
			}
		}
	} else {
		return &issue.Issue{
			Type:    issue.WARNING,
			Message: fmt.Sprintf("Module source \"%s\" is not pinned", source),
			Line:    d.line,
			File:    d.file,
		}
	}

	return nil
}

func (d *TerraformModulePinnedSourceDetector) detectMercurialSource(source string) *issue.Issue {
	if strings.Contains(source, "rev=") {
		if strings.Contains(source, "rev=default") {
			return &issue.Issue{
				Type:    issue.WARNING,
				Message: fmt.Sprintf("Module source \"%s\" uses default rev \"default\"", source),
				Line:    d.line,
				File:    d.file,
			}
		}
	} else {
		return &issue.Issue{
			Type:    issue.WARNING,
			Message: fmt.Sprintf("Module source \"%s\" is not pinned", source),
			Line:    d.line,
			File:    d.file,
		}
	}

	return nil
}
