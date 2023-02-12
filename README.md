Terraform Lexicographical Sort
==============================

This is a small linter for Terraform that checks that all blocks have their attributes in lexicographical/alphabetical order.

## Usage

```
go install github.com/shanet/terraform-lexicographical-lint@latest
$GOPATH/bin/terraform-lexicographical-lint [path to terraform files]
```

## Local Development

```
make
bin/terraform-lexicographical-lint [path to terraform files]
```

## License

MIT
