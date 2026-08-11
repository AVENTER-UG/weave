package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructuredLoggingAdapterIsInitialized(t *testing.T) {
	require.NotNil(t, Logging)
	require.NotPanics(t, func() {
		Logging.Infof("signal handler logging test")
	})
}
