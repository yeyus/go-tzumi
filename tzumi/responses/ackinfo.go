package responses

type AckInfo struct {
	ReturnValue int    `xml:"ret,attr"`
	Message     string `xml:"errmsg,attr"`
}
