package hdhomerun

import (
	"context"
	"fmt"
	"github.com/yeyus/go-tzumi/hdhomerun/protocol"
	"log"
	"net"
)

func (hd *HDHomerunEmulator) controlServer(ctx context.Context) (err error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Printf("[controlServer] error while opening socket: %s", err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("[controlServer] error while accepting connection: %s", err)
			break
		}
		go handleConnection(c)
	}

	return
}

func handleConnection(c net.Conn) {
	log.Printf("[controlServer] serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	defer log.Printf("[controlServer] connection closed")

	buf := make([]byte, 3072)
	for {
		size, err := c.Read(buf)
		if err != nil {
			log.Printf("[controlServer] error while reading: %s", err)
			continue
		}

		pkt := protocol.NewPacket(buf[:size])
		fmt.Printf("[controlServer] packet: %v", pkt)
	}
}
