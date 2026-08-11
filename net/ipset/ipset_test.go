package ipset

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	knftables "sigs.k8s.io/knftables"
)

func newFakeSetBackend(t *testing.T) (Interface, *knftables.Fake) {
	t.Helper()
	fake := knftables.NewFake(knftables.InetFamily, nftTableName)
	tx := fake.NewTransaction()
	tx.Create(&knftables.Table{})
	require.NoError(t, fake.Run(t.Context(), tx))
	return newWithInterface(log.Default(), 0, fake), fake
}

func TestHashIPUsesNativeNFTSetAndElements(t *testing.T) {
	sets, fake := newFakeSetBackend(t)
	require.NoError(t, sets.Create("weave-pods", HashIP))
	require.NoError(t, sets.AddEntry("pod-a", "weave-pods", "10.32.0.4", "pod a"))

	dump := fake.Dump()
	require.Contains(t, dump, "set inet weave weave_pods { type ipv4_addr")
	require.Contains(t, dump, "element inet weave weave_pods { 10.32.0.4 comment \"pod a\" }")
}

func TestHashNetUsesIntervalNFTSet(t *testing.T) {
	sets, fake := newFakeSetBackend(t)
	require.NoError(t, sets.Create("weaver-no-masq-local", HashNet))
	require.NoError(t, sets.AddEntry("owner", "weaver-no-masq-local", "10.32.0.0/12", ""))

	dump := fake.Dump()
	require.Contains(t, dump, "set inet weave weaver_no_masq_local { type ipv4_addr ; flags interval")
	require.Contains(t, dump, "element inet weave weaver_no_masq_local { 10.32.0.0/12 }")
}

func TestEntryRemainsUntilLastOwnerDeletesIt(t *testing.T) {
	sets, fake := newFakeSetBackend(t)
	require.NoError(t, sets.Create("weave-pods", HashIP))
	require.NoError(t, sets.AddEntry("policy-a", "weave-pods", "10.32.0.4", ""))
	require.NoError(t, sets.AddEntry("policy-b", "weave-pods", "10.32.0.4", ""))
	require.NoError(t, sets.DelEntry("policy-a", "weave-pods", "10.32.0.4"))
	require.Equal(t, 1, strings.Count(fake.Dump(), "10.32.0.4"))
	require.NoError(t, sets.DelEntry("policy-b", "weave-pods", "10.32.0.4"))
	require.NotContains(t, fake.Dump(), "10.32.0.4")
}

func TestListSetIsRejectedBecauseNFTDoesNotNestSets(t *testing.T) {
	sets, _ := newFakeSetBackend(t)
	err := sets.Create("weave-namespaces", ListSet)
	require.ErrorContains(t, err, "list:set")
}

func TestNFTSetNameSanitizesSelectorHashes(t *testing.T) {
	require.Equal(t, "weave_I239Zp_sCvoVt_D6u_A_2_YEk", nftSetName("weave-I239Zp%sCvoVt*D6u=A!2]YEk"))
}
