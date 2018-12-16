package hdhomerun

import (
	"context"
	"testing"
)

func TestDiscoveryServer(t *testing.T) {
	t.Skip("skipping manual test")

	hd := HDHomerunEmulator{
		DeviceID: 0xDEADBEEF,
	}

	hd.discoveryServer(context.TODO())

	for {
	}
}
