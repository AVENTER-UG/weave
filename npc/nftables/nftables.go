package nftables

// Interface is the subset of nftables rule operations used by the policy controller.
type Interface interface {
	Append(table, chain string, rulespec ...string) error
	Delete(table, chain string, rulespec ...string) error
	Insert(table, chain string, pos int, rulespec ...string) error
}
