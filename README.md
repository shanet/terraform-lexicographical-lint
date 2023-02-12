Terraform Lexicographical Sort
==============================

This is a small linter for Terraform that checks that all blocks have their attributes in lexicographical/alphabetical order.

## Usage

```
go install github.com/shanet/terraform-lexicographical-lint@latest
terraform-lexicographical-lint [path to terraform files]
```

This assumes `$GOBIN` or `$GOPATH/bin` is in your `$PATH`. If not, the binary should be called directly from `$GOPATH/bin/terraform-lexicographical-lint`.

## Local Development

```
make
bin/terraform-lexicographical-lint [path to terraform files]
```

## License

MIT
