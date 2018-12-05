package tzumi

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/yeyus/go-tzumi/tzumi/buffer"
)

func (t *TzumiMagicTV) tsLoop() {
	// create buffer
	t.TSBuffer = buffer.NewBuffer(64)

	// open local port
	l, err := net.Listen("tcp", ":"+strconv.Itoa(TZUMI_TS_PORT))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to TS connections on port", strconv.Itoa(TZUMI_TS_PORT))
	defer l.Close()

	// call remote end
	tsconn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, TZUMI_TS_PORT))
	if err != nil {
		log.Panicln(err)
	}
	defer tsconn.Close()

	go t.readLoop(tsconn)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go t.handleTSRequest(conn)
	}
}

func (t *TzumiMagicTV) readLoop(tsconn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := tsconn.Read(buf)
		if err == io.EOF {
			log.Printf("[readLoop] Tzumi hung connection: %s", err)
			t.Close()
			panic(err)
		} else if err != nil {
			log.Printf("[readLoop] unexpected error reading TS: %s", err)
			continue
		}

		t.TSBuffer.Write(buf[:n])
	}
}

func (t *TzumiMagicTV) handleTSRequest(conn net.Conn) {
	log.Println("Accepted new TS connection.")
	reader := t.TSBuffer.NewReader()
	defer conn.Close()
	defer log.Println("Closed TS connection.")

	for {
		buf := make([]byte, 1024)
		size, err := reader.Read(buf)
		if err != nil {
			log.Printf("Error reading connection, %s", err)
			return
		} else if size == 0 {
			log.Printf("read 0")
			continue
		}
		data := buf[:size]
		log.Printf("Read new data from connection, bytes %d", len(data))
		conn.Write(data)

		// detect if connection is closed
		one := []byte{}
		conn.SetReadDeadline(time.Now())
		if _, err := conn.Read(one); err == io.EOF {
			log.Printf("%s detected client left", conn.RemoteAddr)
			break
		} else {
			conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
	}
}
