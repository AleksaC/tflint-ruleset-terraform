# terraform_module_pinned_source

Disallow using `depends_on` on modules.

## Example

```hcl
module "foo" {
  source = "./bar"

  depends_on = [resource.baz]
}
```

```
$ tflint
1 issue(s) found:

Warning: depends_on set on module "foo" (terraform_module_depends_on)

  on template.tf line 4:
   4:   depends_on = [resource.baz]

Reference: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.14.0/docs/rules/terraform_module_depends_on.md
```

## Why

Putting `depends_on` on module can cause resources inside the module to be unnecessarily
recreated. Notable example of this is when a module contains data sources. `depends_on`
causes reading of data sources to be deffered to apply phase, causing the dependent
resources to get recreated.

`depends_on` is often used where a direct dependency can be established instead, by
passing an output of one resource as input of another. Using `depends_on` on modules
often suggests poor separation of concerns or bloated workspaces where components
that have different lifecycle are provisioned together.

## How To Fix

Establish a direct dependency between resources that depend on each other. If that's
not possible, reconsider your module boundaries, maybe the dependent resources
should be created inside the same module. In some cases you should do the opposite,
and consider provisioning resources in separate workspaces. Finally, as a last
resort, in some cases it's possible to put a depends_on on some other resource
and establish a transitive dependency between resources.
