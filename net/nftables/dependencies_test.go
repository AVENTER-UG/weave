package nftables

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectDependencyGraphDoesNotContainGoIPTables(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "./...")
	cmd.Dir = "../.."
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
	require.NotContains(t, strings.Split(string(output), "\n"), "github.com/coreos/go-iptables/iptables")
}
