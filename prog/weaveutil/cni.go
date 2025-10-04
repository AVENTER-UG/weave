package main

import (
	"os"

	cni "github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"

	weaveapi "github.com/AVENTER-UG/weave/api"
	"github.com/AVENTER-UG/weave/common"
	ipamplugin "github.com/AVENTER-UG/weave/plugin/ipam"
	netplugin "github.com/AVENTER-UG/weave/plugin/net"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

func cniNet(args []string) error {
	weave := weaveapi.NewClient(os.Getenv("WEAVE_HTTP_ADDR"), common.Log)
	n := netplugin.NewCNIPlugin(weave)
	cni.PluginMain(n.CmdAdd, n.CmdCheck, n.CmdDel, version.All, bv.BuildString("weave net"))
	return nil
}

func cniIPAM(args []string) error {
	weave := weaveapi.NewClient(os.Getenv("WEAVE_HTTP_ADDR"), common.Log)
	i := ipamplugin.NewIpam(weave)
	cni.PluginMain(i.CmdAdd, i.CmdCheck, i.CmdDel, version.All, bv.BuildString("weave net"))
	return nil
}
