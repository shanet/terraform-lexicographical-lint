Terraform Lexicographical Sort
==============================

This is a small linter for Terraform that checks that all blocks have their attributes in lexicographical/alphabetical order. It uses the HCL parsing libraries so all HCL documents are supported with the same parser that Terraform and other HCL-related tools use.

## Example

Consider the following Terraform module call:

```hcl
module "an_unsorted_module" {
  foo    = "str1"
  bar    = "str2"
  source = "github.com/organization/repo"
}
```

Running `terraform-lexicographical-lint example.tf` yields the following output with the expected order of the block's arguments:

```bash
$ terraform-lexicographical-lint example.tf
test.tf:1 Block "module an_unsorted_module" expected order:
    source
    bar
    foo
```

In this case, `source` is considered a special argument which should always appear at the top of the block.

## Usage

```bash
go install github.com/shanet/terraform-lexicographical-lint@latest
terraform-lexicographical-lint [path to terraform files]
```

This assumes `$GOBIN` or `$GOPATH/bin` is in your `$PATH`. If not, the binary should be called directly from `$GOPATH/bin/terraform-lexicographical-lint`.

## Local Development

```bash
make
bin/terraform-lexicographical-lint [path to terraform files]
```

## Publishing new versions

Commit changes as normal then:

```bash
git tag vX.X.X
git push origin vX.X.X
go install github.com/shanet/terraform-lexicographical-lint@vX.X.X
```

## License

MIT
