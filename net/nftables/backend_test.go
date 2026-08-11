package nftables

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	knftables "sigs.k8s.io/knftables"
)

func newFakeBackend() (*NFTables, *knftables.Fake) {
	fake := knftables.NewFake(family, tableName)
	backend := newWithInterface(fake)
	if err := backend.ensureTable(); err != nil {
		panic(err)
	}
	return backend, fake
}

func TestEnsureTableIsIdempotent(t *testing.T) {
	nft, _ := newFakeBackend()
	require.NoError(t, nft.ensureTable())
}

func TestBaseChainsUseNativeNFTHooksAndPriorities(t *testing.T) {
	forward := chainObject("filter", "FORWARD")
	require.Equal(t, "filter_FORWARD", forward.Name)
	require.Equal(t, knftables.FilterType, *forward.Type)
	require.Equal(t, knftables.ForwardHook, *forward.Hook)
	require.Equal(t, knftables.FilterPriority, *forward.Priority)

	postrouting := chainObject("nat", "POSTROUTING")
	require.Equal(t, "nat_POSTROUTING", postrouting.Name)
	require.Equal(t, knftables.NATType, *postrouting.Type)
	require.Equal(t, knftables.PostroutingHook, *postrouting.Hook)
	require.Equal(t, knftables.SNATPriority, *postrouting.Priority)
}

func TestChainExistsSanitizesLegacyChainName(t *testing.T) {
	nft, _ := newFakeBackend()
	require.NoError(t, nft.NewChain("mangle", "WEAVE-CANARY"))

	exists, err := nft.ChainExists("mangle", "WEAVE-CANARY")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestInsertAtFirstPositionWorksForEmptyChain(t *testing.T) {
	nft, fake := newFakeBackend()
	require.NoError(t, nft.Insert("filter", "FORWARD", 1, "-j", "DROP"))
	require.Contains(t, fake.Dump(), "drop")
}

func TestInsertImmediatelyAfterLastRuleAppends(t *testing.T) {
	nft, fake := newFakeBackend()
	require.NoError(t, nft.Insert("filter", "FORWARD", 1, "-j", "DROP"))
	require.NoError(t, nft.Insert("filter", "FORWARD", 2, "-j", "ACCEPT"))
	dump := fake.Dump()
	require.Less(t, strings.Index(dump, "drop"), strings.Index(dump, "accept"))
}

func TestAppendUniqueUsesNativeRuleAndStableIdentity(t *testing.T) {
	nft, fake := newFakeBackend()

	err := nft.AppendUnique("filter", "INPUT", "-i", "weave", "-p", "udp", "--dport", "53", "-j", "ACCEPT")
	require.NoError(t, err)
	dump := fake.Dump()
	require.Contains(t, dump, `iifname "weave" meta l4proto udp udp dport 53 accept`)
	require.Contains(t, dump, `comment "weave:v1:`)
}

func TestAppendUniqueDoesNotDuplicateRuleWithMatchingComment(t *testing.T) {
	nft, fake := newFakeBackend()
	rulespec := []string{"-d", "10.32.0.0/12", "-j", "ACCEPT"}

	require.NoError(t, nft.AppendUnique("filter", "WEAVE-EXPOSE", rulespec...))
	require.NoError(t, nft.AppendUnique("filter", "WEAVE-EXPOSE", rulespec...))
	require.Equal(t, 1, strings.Count(fake.Dump(), "ip daddr 10.32.0.0/12 accept"))
}

func TestDeleteFindsRuleByIdentityAndDeletesItsHandle(t *testing.T) {
	nft, fake := newFakeBackend()
	rulespec := []string{"-s", "10.32.0.0/12", "-j", "MASQUERADE"}
	require.NoError(t, nft.Append("nat", "WEAVE", rulespec...))

	err := nft.Delete("nat", "WEAVE", rulespec...)
	require.NoError(t, err)
	require.NotContains(t, fake.Dump(), "ip saddr 10.32.0.0/12 masquerade")
}

func TestClearChainCreatesAndFlushesScopedChain(t *testing.T) {
	nft, fake := newFakeBackend()
	require.NoError(t, nft.Append("mangle", "WEAVE-IPSEC-IN", "-j", "DROP"))

	err := nft.ClearChain("mangle", "WEAVE-IPSEC-IN")
	require.NoError(t, err)
	dump := fake.Dump()
	require.Contains(t, dump, "chain inet weave mangle_WEAVE_IPSEC_IN")
	require.NotContains(t, dump, "drop")
}

func TestDestroyRemovesOnlyWeaveTable(t *testing.T) {
	nft, fake := newFakeBackend()
	require.NoError(t, nft.Append("filter", "INPUT", "-j", "DROP"))
	require.NoError(t, nft.Destroy())
	require.Empty(t, fake.Dump())
}
