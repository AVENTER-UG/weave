#!/bin/sh
# Initialisation of Weave Net pod: check Linux settings and install CNI plugin

set -e

[ -n "$WEAVE_DEBUG" ] && set -x

modprobe_safe() {
    modprobe $1 || echo "Ignore the error if \"$1\" is built-in in the kernel" >&2
}

nftables_exists() {
    command -v nft >/dev/null 2>&1 || {
        echo '"nft" executable does not exist' >&2
        return 1
    }
    nft --check 'add table inet weave_kube_test' >/dev/null
}

# Default for network policy
EXPECT_NPC=${EXPECT_NPC:-1}

# Ensure we have the required modules for NPC
if [ "${EXPECT_NPC}" != "0" ]; then
    modprobe_safe br_netfilter
    nftables_exists
fi

# kube-proxy requires that bridged traffic passes through netfilter
if ! BRIDGE_NF_ENABLED=$(cat /proc/sys/net/bridge/bridge-nf-call-iptables); then
    echo "Cannot detect bridge-nf support - network policy may not work reliably" >&2
else
    if [ "$BRIDGE_NF_ENABLED" != "1" ]; then
        echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables
    fi
fi

# This is where we expect the manifest to map host directories
HOST_ROOT=${HOST_ROOT:-/host}

# Install CNI plugin binary to typical CNI bin location
# with fall-back to CNI directory used by kube-up on GCI OS
if ! mkdir -p $HOST_ROOT/opt/cni/bin ; then
    if mkdir -p $HOST_ROOT/home/kubernetes/bin ; then
        export WEAVE_CNI_PLUGIN_DIR=$HOST_ROOT/home/kubernetes/bin
    else
        echo "Failed to install the Weave CNI plugin" >&2
        exit 1
    fi
fi
mkdir -p $HOST_ROOT/etc/cni/net.d
export HOST_ROOT
/home/weave/weave --local setup-cni