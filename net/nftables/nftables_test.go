package nftables

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranslateLegacyRuleSpecToNativeNFTExpression(t *testing.T) {
	tests := []struct {
		name  string
		table string
		rule  []string
		want  string
	}{
		{
			name:  "interface protocol address port and verdict",
			table: "filter",
			rule:  []string{"-i", "docker0", "-p", "tcp", "--dst", "172.17.0.1", "--dport", "6783", "-j", "DROP"},
			want:  `iifname "docker0" meta l4proto tcp ip daddr 172.17.0.1 tcp dport 6783 drop`,
		},
		{
			name:  "negated output and conntrack state",
			table: "filter",
			rule:  []string{"-i", "weave", "!", "-o", "weave", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT"},
			want:  `iifname "weave" oifname != "weave" ct state { related, established } accept`,
		},
		{
			name:  "native set lookup and return",
			table: "nat",
			rule:  []string{"-m", "set", "--match-set", "weaver-no-masq-local", "dst", "-m", "comment", "--comment", "Prevent SNAT", "-j", "RETURN"},
			want:  `ip daddr @weaver_no_masq_local return`,
		},
		{
			name:  "negated set lookup",
			table: "filter",
			rule:  []string{"-s", "192.168.48.0/24", "-m", "set", "!", "--match-set", "weave-except", "src", "-j", "ACCEPT"},
			want:  `ip saddr 192.168.48.0/24 ip saddr != @weave_except accept`,
		},
		{
			name:  "packet mark test and set",
			table: "filter",
			rule:  []string{"-m", "mark", "!", "--mark", "0x40000/0x40000", "-j", "MARK", "--set-xmark", "0x20000/0x20000"},
			want:  `meta mark & 0x40000 != 0x40000 meta mark set meta mark & 0xfffdffff | 0x20000`,
		},
		{
			name:  "nflog",
			table: "filter",
			rule:  []string{"-m", "state", "--state", "NEW", "-j", "NFLOG", "--nflog-group", "86"},
			want:  `ct state new log group 86`,
		},
		{
			name:  "esp spi jump",
			table: "mangle",
			rule:  []string{"-s", "10.0.0.2", "-d", "10.0.0.1", "-p", "esp", "-m", "esp", "--espspi", "0xbeef", "-j", "WEAVE-IPSEC-IN-MARK"},
			want:  `ip saddr 10.0.0.2 ip daddr 10.0.0.1 meta l4proto esp esp spi 0xbeef jump mangle_WEAVE_IPSEC_IN_MARK`,
		},
		{
			name:  "outbound policy missing",
			table: "filter",
			rule:  []string{"!", "-p", "esp", "-m", "policy", "--dir", "out", "--pol", "none", "-m", "mark", "--mark", "0x20000/0x20000", "-j", "DROP"},
			want:  `meta l4proto != esp rt ipsec missing meta mark & 0x20000 == 0x20000 drop`,
		},
		{
			name:  "bridge physical output",
			table: "filter",
			rule:  []string{"-m", "physdev", "--physdev-is-bridged", "--physdev-out=vethwe-bridge", "-j", "ACCEPT"},
			want:  `meta bri_oifname "vethwe-bridge" accept`,
		},
		{
			name:  "local address type",
			table: "filter",
			rule:  []string{"-m", "addrtype", "!", "--src-type", "LOCAL", "-j", "DROP"},
			want:  `fib saddr type != local drop`,
		},
		{
			name:  "masquerade",
			table: "nat",
			rule:  []string{"-s", "10.32.0.0/12", "!", "-d", "10.32.0.0/12", "-j", "MASQUERADE"},
			want:  `ip saddr 10.32.0.0/12 ip daddr != 10.32.0.0/12 masquerade`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translateRule(tt.table, tt.rule)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNativeNFTTranslationRejectsUnknownLegacyOption(t *testing.T) {
	_, err := translateRule("filter", []string{"--unknown", "value", "-j", "ACCEPT"})
	require.ErrorContains(t, err, "unsupported iptables argument")
}

func TestObjectNamesAreScopedByLogicalTable(t *testing.T) {
	require.Equal(t, "filter_INPUT", chainName("filter", "INPUT"))
	require.Equal(t, "mangle_WEAVE_IPSEC_IN", chainName("mangle", "WEAVE-IPSEC-IN"))
	require.Equal(t, "weave_namespace_selector", setName("weave-namespace-selector"))
	require.Equal(t, "weave_I239Zp_sCvoVt_D6u_A_2_YEk", setName("weave-I239Zp%sCvoVt*D6u=A!2]YEk"))
}
