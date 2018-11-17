package commands

import "fmt"

type TuneRequest struct {
	Frequency int
	Bandwidth int
	Plp       int
	Program   int
}

func NewTuneRequest(frequency int, program int) TuneRequest {
	return TuneRequest{
		Frequency: frequency,
		Bandwidth: 8000,
		Plp:       0,
		Program:   program,
	}
}

func (c TuneRequest) Serialize() string {
	return fmt.Sprintf("<msg type=\"tune_req\"><params tv_type=\"ATSC\" freq=\"%d\" bandwidth=\"%d\" plp=\"%d\" programNumber=\"%d\"/></msg>", c.Frequency, c.Bandwidth, c.Plp, c.Program)
}

func (c TuneRequest) ToHuman() string {
	return fmt.Sprintf("TuneRequest -> Freq: %dKhz Bandwidth: %dKhz Program: %d", c.Frequency, c.Bandwidth, c.Program)
}
