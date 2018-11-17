package responses

type SignalStatus struct {
	Strength int `xml:"strength,attr"`
	Quality  int `xml:"quality,attr"`
	Qam      int `xml:"qam,attr"`
}
