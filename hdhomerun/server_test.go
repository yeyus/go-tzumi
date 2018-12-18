package hdhomerun

import (
	"context"
	"testing"
)

func TestControlServer(t *testing.T) {
	hd := HDHomerunEmulator{
		DeviceID: 0x1231E599,
	}

	ctx := context.TODO()
	go hd.discoveryServer(ctx)
	hd.controlServer(ctx)
}
