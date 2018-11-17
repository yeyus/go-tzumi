package tzumi

import (
	"encoding/xml"
	"log"
	"strings"

	"github.com/yeyus/go-tzumi/tzumi/responses"
)

func (t *TzumiMagicTV) parse(response string) responses.Response {

	// fix login_ack where xml is not well formatted
	response = strings.Replace(response, "\"<deviceinfo", "\"/><deviceinfo", -1)

	if t.Debug {
		log.Printf("[DEBUG] parsing command: %s", response)
	}

	var reply responses.Response
	xml.Unmarshal([]byte(response), &reply)

	return reply
}
