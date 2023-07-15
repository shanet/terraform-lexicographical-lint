package main

type SortedNode struct {
	Name string
	Line int
	Type int
}

func (this SortedNode) IsSpecial() bool {
	return (this.Name == "count" ||
		this.Name == "source" ||
		this.Name == "providers" ||
		this.Name == "lifecycle" ||
		this.Name == "metadata" ||
		this.Name == "for_each")
}
