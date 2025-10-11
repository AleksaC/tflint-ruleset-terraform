package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformModuleDependsOn(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "Has depends_on",
			Content: `
module "invalid" {
  source = "./irrelevant"

  depends_on = [resource.not_relevant]
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleDependsOnRule(),
					Message: `depends_on set for module "invalid"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 16},
						End:      hcl.Pos{Line: 5, Column: 39},
					},
				},
			},
		},
		{
			Name: "Does not have depends_on",
			Content: `
module "valid" {
  source = "./irrelevant"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformModuleDependsOnRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := testRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Runner.(*helper.Runner).Issues)
		})
	}
}
