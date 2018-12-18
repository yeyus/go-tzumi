package hdhomerun

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yeyus/go-tzumi/hdhomerun/protocol"
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
		go hd.handleConnection(c)
	}

	return
}

func (hd *HDHomerunEmulator) handleConnection(c net.Conn) {
	log.Printf("[controlServer] serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	defer log.Printf("[controlServer] connection closed")

	buf := make([]byte, 3072)
	for {
		size, err := c.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("[controlServer] error while reading: %s", err)
			continue
		}

		pkt := protocol.NewPacket(buf[:size])
		log.Printf("[controlServer] packet: %+v", pkt)
		err = hd.commandLoop(pkt, c)
		if err != nil {
			log.Printf("[controlServer] error while processing packet: %s", err)
			continue
		}
	}
}

func (hd *HDHomerunEmulator) commandLoop(pkt *protocol.Packet, c net.Conn) error {
	switch pkt.Type {
	case protocol.HDHOMERUN_TYPE_GETSET_REQ:
		msg, err := protocol.ParseGetSetMsg(pkt)
		log.Printf("[commandLoop] GETSET REQ -> %s", msg.Property.GetValue())
		if err != nil {
			break
		}
		err = hd.processGetSetRequest(msg, c)
		if err != nil {
			log.Printf("[commandLoop] error while processing GETSET: %s", err)
			log.Printf("[commandLoop] unknown property %s", string(pkt.Tags[0].Value))
			return err
		}
	default:
		log.Printf("[controlServer] unrecognized command with type 0x%x", pkt.Type)
	}

	return nil
}

func (hd *HDHomerunEmulator) processGetSetRequest(msg *protocol.GetSetMsg, c net.Conn) error {
	switch msg.Property {
	case protocol.SysModel:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.SysModel,
			Value:    []byte("hdhomerun3_atsc\x00"),
		}
		log.Printf("[GetSetRequest] responding to %s with %s", getset.Property.GetValue(), string(getset.Value))
		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.SysHwModel:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.SysHwModel,
			Value:    []byte("HDHR3-US\x00"),
		}
		log.Printf("[GetSetRequest] responding to %s with %s", getset.Property.GetValue(), string(getset.Value))
		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.SysVersion:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.SysVersion,
			Value:    []byte("20170930\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.SysFeatures:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.SysFeatures,
			Value: []byte("channelmap: us-bcast us-cable us-hrc us-irc\n" +
				"modulation: 8vsb qam256 qam64\n" +
				"auto-modulation: auto auto6t auto6c qam\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.TunerStatus:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.TunerStatus,
			Value:    []byte("ch=qam:33 lock=qam256 ss=83 snq=90 seq=100 bps=38807712 pps=0\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.TunerStreamInfo:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.TunerStreamInfo,
			Value: []byte("3: 33.1 KBWB-HD\n" +
				"4: 33.4 AZTECA\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.TunerChannelMap:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.TunerChannelMap,
			Value:    []byte("us-bcast\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	case protocol.TunerProgram:
		getset := protocol.GetSetMsg{
			Type:     protocol.HDHOMERUN_TYPE_GETSET_RPY,
			Property: protocol.TunerProgram,
			Value:    []byte("0\x00"),
		}

		_, err := c.Write(getset.Packet().Serialize())
		return err
	default:
		return errors.New("unknown property")
	}
}
