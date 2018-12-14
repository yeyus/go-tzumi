package hdhomerun

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/yeyus/go-tzumi/hdhomerun/protocol"
	"net"
	"time"
)

const maxBufferSize = 1024

func (hd *HDHomerunEmulator) discoveryServer(ctx context.Context) (err error) {
	// code from https://ops.tips/blog/udp-client-and-server-in-go/
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		return
	}

	defer pc.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n",
				n, addr.String())

			var rpy protocol.Packet
			pkt := protocol.NewPacket(buffer[:n])
			if pkt.Type == protocol.HDHOMERUN_TYPE_DISCOVER_REQ &&
				(binary.BigEndian.Uint32(pkt.Tags[0].Value) == HDHOMERUN_DEVICE_TYPE_WILDCARD || binary.BigEndian.Uint32(pkt.Tags[0].Value) == HDHOMERUN_DEVICE_TYPE_TUNER) &&
				(binary.BigEndian.Uint32(pkt.Tags[1].Value) == HDHOMERUN_DEVICE_ID_WILDCARD || binary.BigEndian.Uint32(pkt.Tags[1].Value) == hd.DeviceID) {
				rpyMsg := protocol.DiscoveryResponseMsg{
					DeviceType: HDHOMERUN_DEVICE_TYPE_TUNER,
					DeviceId:   hd.DeviceID,
				}
				rpy = rpyMsg.Packet()
			}
			fmt.Printf("Received packet: %v\n", pkt)

			deadline := time.Now().Add(5 * time.Second)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			// Write the packet's contents back to the client.
			n, err = pc.WriteTo(rpy.Serialize(), addr)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
