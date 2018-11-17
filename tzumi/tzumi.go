package tzumi

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/yeyus/go-tzumi/tzumi/commands"
	"github.com/yeyus/go-tzumi/tzumi/responses"
)

const TZUMI_COMMAND_PORT = 6000
const TZUMI_TS_PORT = 8000

type TzumiMagicTV struct {
	Host            string
	CommandConn     net.Conn
	CommandChannel  chan string
	responseChannel chan responses.Response
	Debug           bool
	State           State
}

type Callback func(response responses.Response)

func NewTzumiMagicTV(host string) (*TzumiMagicTV, error) {
	t := &TzumiMagicTV{}

	t.State = DISCONNECTED

	t.CommandChannel = make(chan string)
	t.responseChannel = make(chan responses.Response)

	t.Host = host
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, TZUMI_COMMAND_PORT))
	if err != nil {
		log.Printf("Error opening socket to %s on command port: %s", host, err)
		return t, err
	}
	t.CommandConn = conn

	go t.tcpLoop()

	return t, err
}

func (t *TzumiMagicTV) tcpLoop() {
	for {
		// non blocking channel read
		// if there's command send it, if not skip to read
		select {
		case cmd, more := <-t.CommandChannel:
			if !more {
				break
			}
			fmt.Fprintf(t.CommandConn, cmd)
		default:
			// no command received
			// do read
			err := t.CommandConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				log.Println("Nothing has been read:", err)
				// do something else, for example create new conn
				return
			}

			recvBuf := make([]byte, 1024)

			_, err = t.CommandConn.Read(recvBuf[:]) // recv data
			if err != nil {
				continue
			}
			log.Printf("[Tzumi] calling process response with buffer %s", recvBuf)
			t.processResponse(recvBuf)
		}
	}
}

func (t *TzumiMagicTV) processResponse(recvbuf []byte) {
	responseBody := strings.Trim(string(recvbuf), " \t\n\r")

	if t.Debug {
		log.Printf("[DEBUG] raw received: %s", responseBody)
	}

	response := t.parse(responseBody)

	if t.Debug {
		log.Printf("[DEBUG] parsed response: %s", response.String())
	}

	switch response.Type {
	case responses.LoginAck:
		t.setState(CONNECTED)
	case responses.TunerAck:
		if response.AckInfo.ReturnValue == 0 {
			t.setState(TUNED)
		} else {
			t.setState(CONNECTED)
		}
	}
	t.responseChannel <- response
}

func (t *TzumiMagicTV) SendCommand(command commands.Command) {
	log.Printf("Sending command: %s", command.ToHuman())
	if t.Debug {
		log.Printf("[DEBUG] raw send: %s", command.Serialize())
	}
	t.CommandChannel <- command.Serialize()
}

func (t *TzumiMagicTV) Login(cb Callback) {
	t.SendCommand(commands.NewLoginRequest())
	response := <-t.responseChannel
	cb(response)
}

func (t *TzumiMagicTV) Tune(frequency int, program int, cb Callback) {
	t.SendCommand(commands.NewTuneRequest(frequency, program))
	t.setState(TUNING)
	response := <-t.responseChannel
	cb(response)
}

func (t *TzumiMagicTV) CheckLock(cb Callback) {
	t.SendCommand(commands.NewCheckLockRequest())
	response := <-t.responseChannel
	cb(response)
}

func (t *TzumiMagicTV) GetSignalStatus(cb Callback) {
	t.SendCommand(commands.NewSignalStatusRequest())
	response := <-t.responseChannel
	cb(response)
}

func (t *TzumiMagicTV) Close() {
	log.Printf("Clossing channel and tcp processing...")
	t.setState(DISCONNECTED)
}
