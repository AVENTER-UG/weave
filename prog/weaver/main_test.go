package main

import (
	"testing"

	weavetest "github.com/AVENTER-UG/weave/testing"
)

func TestMain(t *testing.T) {
	if weavetest.TrimTestArgs() {
		main()
	}
}
