package hdhomerun

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/yeyus/go-tzumi/hdhomerun/protocol"
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
			} else if strings.Contains(addr.String(), "127.0.0.1") {
				// ignore localhost packages
				continue
			}

			log.Printf("packet-received: bytes=%d from=%s\n",
				n, addr.String())

			var rpy protocol.Packet
			pkt := protocol.NewPacket(buffer[:n])
			msg, err := protocol.ParseDiscoveryMsg(pkt)
			if err != nil {
				log.Printf("error parsing incoming discovery packet: %s", err)
				log.Printf("packet was: %v", pkt)
				continue
			}

			log.Printf("received packet: %v", pkt)
			if !msg.IsAddressedToSelf(hd.DeviceID) {
				log.Printf("packet wasn't addressed to this device!")
				continue
			}

			rpyMsg := protocol.DiscoveryMsg{
				Type:       protocol.HDHOMERUN_TYPE_DISCOVER_RPY,
				DeviceType: protocol.HDHOMERUN_DEVICE_TYPE_TUNER,
				DeviceId:   hd.DeviceID,
				TunerCount: 1,
			}
			rpy = rpyMsg.Packet()

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

			log.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
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
