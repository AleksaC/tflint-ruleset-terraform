package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-terraform/project"
	"github.com/terraform-linters/tflint-ruleset-terraform/terraform"
)

// TerraformModuleDependsOnRule checks whether `depends_on“ attribure is set on a module
type TerraformModuleDependsOnRule struct {
	tflint.DefaultRule
}

// NewTerraformModuleDependsOnRule returns the new rule
func NewTerraformModuleDependsOnRule() *TerraformModuleDependsOnRule {
	return &TerraformModuleDependsOnRule{}
}

// Name returns the rule name
func (r *TerraformModuleDependsOnRule) Name() string {
	return "terraform_module_depends_on"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModuleDependsOnRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformModuleDependsOnRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModuleDependsOnRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether `depends_on“ attribure is set on a module
func (r *TerraformModuleDependsOnRule) Check(rr tflint.Runner) error {
	runner := rr.(*terraform.Runner)

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "depends_on"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "failed to call GetModuleContent()",
				Detail:   err.Error(),
			},
		}
	}

	for _, block := range body.Blocks {
		if attr, exists := block.Body.Attributes["depends_on"]; exists {
			return runner.EmitIssue(
				r,
				fmt.Sprintf(`depends_on set on module "%s"`, block.Labels[0]),
				attr.Expr.Range(),
			)
		}
	}

	return nil
}
