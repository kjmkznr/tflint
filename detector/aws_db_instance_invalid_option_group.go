package detector

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/wata727/tflint/issue"
)

type AwsDBInstanceInvalidOptionGroupDetector struct {
	*Detector
	IssueType    string
	Target       string
	DeepCheck    bool
	optionGroups map[string]bool
}

func (d *Detector) CreateAwsDBInstanceInvalidOptionGroupDetector() *AwsDBInstanceInvalidOptionGroupDetector {
	return &AwsDBInstanceInvalidOptionGroupDetector{
		Detector:     d,
		IssueType:    issue.ERROR,
		Target:       "aws_db_instance",
		DeepCheck:    true,
		optionGroups: map[string]bool{},
	}
}

func (d *AwsDBInstanceInvalidOptionGroupDetector) PreProcess() {
	resp, err := d.AwsClient.DescribeOptionGroups()
	if err != nil {
		d.Logger.Error(err)
		d.Error = true
		return
	}

	for _, optionGroup := range resp.OptionGroupsList {
		d.optionGroups[*optionGroup.OptionGroupName] = true
	}
}

func (d *AwsDBInstanceInvalidOptionGroupDetector) Detect(file string, item *ast.ObjectItem, issues *[]*issue.Issue) {
	optionGroupToken, err := hclLiteralToken(item, "option_group_name")
	if err != nil {
		d.Logger.Error(err)
		return
	}
	optionGroup, err := d.evalToString(optionGroupToken.Text)
	if err != nil {
		d.Logger.Error(err)
		return
	}

	if !d.optionGroups[optionGroup] {
		issue := &issue.Issue{
			Type:    d.IssueType,
			Message: fmt.Sprintf("\"%s\" is invalid option group name.", optionGroup),
			Line:    optionGroupToken.Pos.Line,
			File:    file,
		}
		*issues = append(*issues, issue)
	}
}
