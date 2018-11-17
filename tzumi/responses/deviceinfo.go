package responses

type DeviceInfo struct {
	DeviceName string `xml:"device_name,attr"`
	DeviceType string `xml:"device_type,attr"`
	TvType     string `xml:"tv_type,attr"`
	HwVersion  string `xml:"hwversion,attr"`
	SwVersion  string `xml:"softwareversion,attr"`
}
