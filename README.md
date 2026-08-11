# Weave Net

This repository contains a fork of Weave Net, the first product developed by Weaveworks. Since Weaveworks has shut down, this repo aims to continue maintaining Weave Net, and to publish releases regularly.

[![Go Report Card](https://goreportcard.com/badge/github.com/AVENTER-UG/weave)](https://goreportcard.com/report/github.com/AVENTER-UG/weave)
[![Docker Pulls](https://img.shields.io/docker/pulls/avhost/weave-kube "Number of times the weave-kube image was pulled from the Docker Hub")](https://hub.docker.com/r/rajchaudhuri/weave-kube)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/AVENTER-UG/weave?include_prereleases)](https://github.com/AVENTER-UG/weave/releases)
[![Unique CVE count in all images](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2FAVENTER-UG%2Fweave%2Fmaster%2Freweave%2Fscans%2Fbadge.json&label=CVE%20count "The number of unique CVEs reported by scanning all images")](reweave/scans/report.md)

## Requirements

### Weave Net 4.x

Weave Net 4.x uses **native nftables** for its firewall, NAT, IPsec, and Network
Policy rules. It no longer calls `iptables`, `iptables-nft`, `iptables-legacy`,
or `ipset`, and the 4.x container images do not install those tools. Weave owns
an isolated `inet weave` table and does not modify unrelated nftables tables.

The host must provide Linux `nf_tables` kernel support. The `nft` userspace tool
is included in the Weave images, so it does not have to be installed separately
on the host.

Docker may use its own native nftables firewall backend alongside Weave. To run
the entire host without active iptables rules, use a Docker release that supports
the nftables backend and configure `/etc/docker/daemon.json` with:

```json
{
  "firewall-backend": "nftables"
}
```

Validate the Docker configuration before restarting the daemon. Existing Docker
daemon settings, such as custom runtimes, must remain in the same JSON object.

The release images are tagged `4.0.0`, for example:

```text
avhost/weave:4.0.0
avhost/weaveexec:4.0.0
avhost/weave-kube:4.0.0
avhost/weave-npc:4.0.0
```

### Weave Net 3.x and earlier

Weave Net 3.x and earlier use the iptables/ipset implementation. They require
`iptables`, `iptables-nft`, or `iptables-legacy` and do not provide the native
nftables backend introduced in 4.x.

When upgrading an existing installation, stop the old Weave container before
launching 4.x. Keep the previous image available until peer connectivity, Weave
DNS, container egress, and the `inet weave` ruleset have been verified.

## Funding

[![](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=H553XE4QJ9GJ8)

## Using Weave Net in other ways

Please refer to the [documentation](https://AVENTER-UG.github.io/weave).

## Building Weave Net

Details can be found [here](reweave/BUILDING.md). 

## Documentation status

The public documentation that used to exist in the `site` directory has been moved to the `original/site` directory. A new `website` directory has been created, and populated with the content of the `original/site` directory, rearranged and reformatted for being built with Jekyll and published to the GitHub pages site [https://aventer-ug.github.io/weave](https://aventer-ug.github.io/weave).

The documentation will now be maintained and published from the `website` directory exclusively.
