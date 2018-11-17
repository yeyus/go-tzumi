package tzumi

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func (t *TzumiMagicTV) tsLoop() {
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

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleTSRequest(conn, tsconn)
	}
}

func handleTSRequest(conn net.Conn, tsconn net.Conn) {
	log.Println("Accepted new TS connection.")
	defer conn.Close()
	defer log.Println("Closed TS connection.")

	for {
		buf := make([]byte, 1024)
		size, err := tsconn.Read(buf)
		if err != nil {
			log.Printf("Error reading connection, %s", err)
			return
		}
		data := buf[:size]
		log.Printf("Read new data from connection, bytes %d", len(data))
		conn.Write(data)
	}
}
