package hdhomerun

import (
	"context"
	"testing"
)

func TestControlServer(t *testing.T) {
	hd := HDHomerunEmulator{
		DeviceID: 0xDEADBEEF,
	}

	ctx := context.TODO()
	hd.discoveryServer(ctx)
	hd.controlServer(ctx)
}
