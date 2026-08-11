package main

import (
	"fmt"

	"github.com/AVENTER-UG/weave/net/nftables"
)

func nftAdd(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: weaveutil nft-add <table> <chain> <rule>...")
	}
	nft, err := nftables.New()
	if err != nil {
		return err
	}
	return nft.AppendUnique(args[0], args[1], args[2:]...)
}

func nftDelete(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: weaveutil nft-delete <table> <chain> <rule>...")
	}
	nft, err := nftables.New()
	if err != nil {
		return err
	}
	return nft.Delete(args[0], args[1], args[2:]...)
}

func nftDestroy(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: weaveutil nft-destroy")
	}
	nft, err := nftables.New()
	if err != nil {
		return err
	}
	return nft.Destroy()
}
