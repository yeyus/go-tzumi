package hdhomerun

import (
	"context"
	"testing"
)

func TestDiscoveryServer(t *testing.T) {
	hd := HDHomerunEmulator{
		DeviceID: 0xDEADBEEF,
	}

	hd.discoveryServer(context.TODO())

	for {
	}
}
