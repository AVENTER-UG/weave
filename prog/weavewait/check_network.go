//go:build !iface && !mcast
// +build !iface,!mcast

package main

func checkNetwork() error {
	return nil
}
