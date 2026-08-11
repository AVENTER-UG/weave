package nftables

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	knftables "sigs.k8s.io/knftables"
)

const (
	family    = knftables.InetFamily
	tableName = "weave"
)

type NFTables struct {
	nft knftables.Interface
}

func New() (*NFTables, error) {
	nft, err := knftables.New(family, tableName)
	if err != nil {
		return nil, fmt.Errorf("initialize nftables: %w", err)
	}
	backend := newWithInterface(nft)
	if err := backend.ensureTable(); err != nil {
		return nil, err
	}
	return backend, nil
}

func newWithInterface(nft knftables.Interface) *NFTables {
	return &NFTables{nft: nft}
}

func stringPtr(value string) *string { return &value }
func intPtr(value int) *int          { return &value }

func (n *NFTables) ensureTable() error {
	tx := n.nft.NewTransaction()
	tx.Add(&knftables.Table{})
	if err := n.nft.Run(context.Background(), tx); err != nil {
		return fmt.Errorf("create nftables table %s: %w", tableName, err)
	}
	return nil
}

// Destroy removes the table owned by Weave without touching any other nftables state.
func (n *NFTables) Destroy() error {
	tx := n.nft.NewTransaction()
	tx.Destroy(&knftables.Table{})
	if err := n.nft.Run(context.Background(), tx); err != nil && !knftables.IsNotFound(err) {
		return fmt.Errorf("destroy nftables table %s: %w", tableName, err)
	}
	return nil
}

type baseChainSpec struct {
	chainType knftables.BaseChainType
	hook      knftables.BaseChainHook
	priority  knftables.BaseChainPriority
}

var baseChains = map[string]baseChainSpec{
	"filter_INPUT":       {knftables.FilterType, knftables.InputHook, knftables.FilterPriority},
	"filter_FORWARD":     {knftables.FilterType, knftables.ForwardHook, knftables.FilterPriority},
	"filter_OUTPUT":      {knftables.FilterType, knftables.OutputHook, knftables.FilterPriority},
	"mangle_PREROUTING":  {knftables.FilterType, knftables.PreroutingHook, knftables.ManglePriority},
	"mangle_INPUT":       {knftables.FilterType, knftables.InputHook, knftables.ManglePriority},
	"mangle_FORWARD":     {knftables.FilterType, knftables.ForwardHook, knftables.ManglePriority},
	"mangle_OUTPUT":      {knftables.FilterType, knftables.OutputHook, knftables.ManglePriority},
	"mangle_POSTROUTING": {knftables.FilterType, knftables.PostroutingHook, knftables.ManglePriority},
	"nat_PREROUTING":     {knftables.NATType, knftables.PreroutingHook, knftables.DNATPriority},
	"nat_INPUT":          {knftables.NATType, knftables.InputHook, knftables.DNATPriority},
	"nat_OUTPUT":         {knftables.NATType, knftables.OutputHook, knftables.DNATPriority},
	"nat_POSTROUTING":    {knftables.NATType, knftables.PostroutingHook, knftables.SNATPriority},
}

func chainObject(table, chain string) *knftables.Chain {
	name := chainName(table, chain)
	object := &knftables.Chain{Name: name}
	if spec, found := baseChains[name]; found {
		object.Type = &spec.chainType
		object.Hook = &spec.hook
		object.Priority = &spec.priority
	}
	return object
}

func (n *NFTables) createChain(table, chain string) error {
	name := chainName(table, chain)
	chains, err := n.nft.List(context.Background(), "chain")
	if err != nil && !knftables.IsNotFound(err) {
		return fmt.Errorf("list nftables chains: %w", err)
	}
	for _, existing := range chains {
		if existing == name {
			return nil
		}
	}
	tx := n.nft.NewTransaction()
	tx.Create(chainObject(table, chain))
	if err := n.nft.Run(context.Background(), tx); err != nil && !knftables.IsAlreadyExists(err) {
		return fmt.Errorf("create nftables chain %s/%s: %w", table, chain, err)
	}
	return nil
}

func (n *NFTables) NewChain(table, chain string) error {
	return n.createChain(table, chain)
}

func (n *NFTables) ChainExists(table, chain string) (bool, error) {
	wanted := chainName(table, chain)
	chains, err := n.nft.List(context.Background(), "chain")
	if err != nil {
		if knftables.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("list nftables chains: %w", err)
	}
	for _, existing := range chains {
		if existing == wanted {
			return true, nil
		}
	}
	return false, nil
}

func (n *NFTables) ClearChain(table, chain string) error {
	if err := n.createChain(table, chain); err != nil {
		return err
	}
	tx := n.nft.NewTransaction()
	tx.Flush(&knftables.Chain{Name: chainName(table, chain)})
	if err := n.nft.Run(context.Background(), tx); err != nil {
		return fmt.Errorf("clear nftables chain %s/%s: %w", table, chain, err)
	}
	return nil
}

func (n *NFTables) DeleteChain(table, chain string) error {
	tx := n.nft.NewTransaction()
	tx.Destroy(&knftables.Chain{Name: chainName(table, chain)})
	if err := n.nft.Run(context.Background(), tx); err != nil {
		return fmt.Errorf("delete nftables chain %s/%s: %w", table, chain, err)
	}
	return nil
}

func ruleComment(table, chain string, rulespec []string) string {
	raw := strings.Join(rulespec, "\x00")
	digest := sha256.Sum256([]byte(table + "\x00" + chain + "\x00" + raw))
	prefix := fmt.Sprintf("weave:v1:%x:", digest[:8])
	encoded := base64.RawURLEncoding.EncodeToString([]byte(raw))
	if len(prefix)+len(encoded) > 126 {
		return strings.TrimSuffix(prefix, ":")
	}
	return prefix + encoded
}

func decodeRuleComment(comment string) ([]string, bool) {
	parts := strings.SplitN(comment, ":", 4)
	if len(parts) != 4 || parts[0] != "weave" || parts[1] != "v1" || parts[3] == "" {
		return nil, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, false
	}
	return strings.Split(string(raw), "\x00"), true
}

func (n *NFTables) findRule(table, chain string, rulespec []string) (*knftables.Rule, error) {
	comment := ruleComment(table, chain, rulespec)
	rules, err := n.nft.ListRules(context.Background(), chainName(table, chain))
	if err != nil {
		if knftables.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list nftables rules in %s/%s: %w", table, chain, err)
	}
	for _, rule := range rules {
		if rule.Comment != nil && *rule.Comment == comment {
			return rule, nil
		}
	}
	return nil, nil
}

func (n *NFTables) Exists(table, chain string, rulespec ...string) (bool, error) {
	rule, err := n.findRule(table, chain, rulespec)
	return rule != nil, err
}

func (n *NFTables) add(table, chain string, index *int, unique, insert bool, rulespec []string) error {
	if err := n.createChain(table, chain); err != nil {
		return err
	}
	if unique {
		exists, err := n.Exists(table, chain, rulespec...)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
	}
	expression, err := translateRule(table, rulespec)
	if err != nil {
		return err
	}
	comment := ruleComment(table, chain, rulespec)
	tx := n.nft.NewTransaction()
	rule := &knftables.Rule{Chain: chainName(table, chain), Rule: expression, Comment: &comment, Index: index}
	if insert {
		tx.Insert(rule)
	} else {
		tx.Add(rule)
	}
	if err := n.nft.Run(context.Background(), tx); err != nil {
		return fmt.Errorf("add nftables rule to %s/%s (%v): %w", table, chain, rulespec, err)
	}
	return nil
}

func (n *NFTables) Append(table, chain string, rulespec ...string) error {
	return n.add(table, chain, nil, false, false, rulespec)
}

func (n *NFTables) AppendUnique(table, chain string, rulespec ...string) error {
	return n.add(table, chain, nil, true, false, rulespec)
}

func (n *NFTables) Insert(table, chain string, pos int, rulespec ...string) error {
	if pos < 1 {
		return fmt.Errorf("nftables rule position must be at least 1, got %d", pos)
	}
	if err := n.createChain(table, chain); err != nil {
		return err
	}
	rules, err := n.nft.ListRules(context.Background(), chainName(table, chain))
	if err != nil {
		return fmt.Errorf("list nftables rules in %s/%s: %w", table, chain, err)
	}
	if pos > len(rules) {
		return n.add(table, chain, nil, false, false, rulespec)
	}
	return n.add(table, chain, intPtr(pos-1), false, true, rulespec)
}

func (n *NFTables) Delete(table, chain string, rulespec ...string) error {
	rule, err := n.findRule(table, chain, rulespec)
	if err != nil {
		return err
	}
	if rule == nil {
		return nil
	}
	if rule.Handle == nil {
		return fmt.Errorf("nftables did not return a handle for rule in %s/%s", table, chain)
	}
	tx := n.nft.NewTransaction()
	tx.Delete(&knftables.Rule{Chain: chainName(table, chain), Handle: rule.Handle})
	if err := n.nft.Run(context.Background(), tx); err != nil {
		return fmt.Errorf("delete nftables rule from %s/%s (%v): %w", table, chain, rulespec, err)
	}
	return nil
}

func (n *NFTables) ListChains(table string) ([]string, error) {
	chains, err := n.nft.List(context.Background(), "chain")
	if err != nil {
		return nil, err
	}
	prefix := sanitizeName(table) + "_"
	result := make([]string, 0)
	for _, chain := range chains {
		if strings.HasPrefix(chain, prefix) {
			result = append(result, strings.TrimPrefix(chain, prefix))
		}
	}
	return result, nil
}

func (n *NFTables) List(table, chain string) ([]string, error) {
	rules, err := n.nft.ListRules(context.Background(), chainName(table, chain))
	if err != nil {
		return nil, err
	}
	result := []string{"-N " + chain}
	for _, rule := range rules {
		if rule.Comment == nil {
			result = append(result, "-A "+chain)
			continue
		}
		if spec, ok := decodeRuleComment(*rule.Comment); ok {
			result = append(result, "-A "+chain+" "+strings.Join(spec, " "))
		} else {
			result = append(result, "-A "+chain)
		}
	}
	return result, nil
}
