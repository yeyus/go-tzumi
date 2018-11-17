package tzumi

import (
	"testing"

	"github.com/yeyus/go-tzumi/tzumi/responses"
)

func TestParseLoginMsg(t *testing.T) {
	tz, _ := NewTzumiMagicTV("127.0.0.1")

	tz.Debug = true

	var a responses.Response = tz.parse("<msg type=\"login_ack\"><ack_info ret=\"0\" errmsg=\"OK\"<deviceinfo device_name=\"Mobile Tv Box\"device_type=\"Beta2.3\" tv_type=\"ATSC\" hwversion=\"2.5\" softwareversion=\"2.8\"/></msg>")
	if a.Type != "login_ack" {
		t.Errorf("msg type should be login_ack but got %s", a.Type)
	}

	if a.AckInfo.ReturnValue != 0 {
		t.Errorf("ret should be 0 but got %d", a.AckInfo.ReturnValue)
	}

	if a.AckInfo.Message != "OK" {
		t.Errorf("errmsg should be OK but got %s", a.AckInfo.Message)
	}

	if a.DeviceInfo.DeviceName != "Mobile Tv Box" {
		t.Errorf("device_name should be Mobile Tv Box but got %s", a.DeviceInfo.DeviceName)
	}

	if a.DeviceInfo.DeviceType != "Beta2.3" {
		t.Errorf("devicetype should be Beta2.3 but got %s", a.DeviceInfo.DeviceType)
	}

	if a.DeviceInfo.TvType != "ATSC" {
		t.Errorf("tvtype should be ATSC but got %s", a.DeviceInfo.TvType)
	}

	if a.DeviceInfo.HwVersion != "2.5" {
		t.Errorf("hwversion should be 2.5 but got %s", a.DeviceInfo.HwVersion)
	}

	if a.DeviceInfo.SwVersion != "2.8" {
		t.Errorf("swversion should be 2.8 but got %s", a.DeviceInfo.SwVersion)
	}
}

func TestParseSignalStatusMsg(t *testing.T) {
	tz, _ := NewTzumiMagicTV("127.0.0.1")

	tz.Debug = true

	a := tz.parse("<msg type=\"signal_status_ack\"><ack_info ret=\"0\" errmsg=\"OK\" /><signal_status strength=\"-128\" quality=\"0\" qam=\"0\"/></msg>")
	if a.Type != "signal_status_ack" {
		t.Errorf("msg type should be signal_status_ack but got %s", a.Type)
	}

	if a.AckInfo.ReturnValue != 0 {
		t.Errorf("ret should be 0 but got %d", a.AckInfo.ReturnValue)
	}

	if a.AckInfo.Message != "OK" {
		t.Errorf("errmsg should be OK but got %s", a.AckInfo.Message)
	}

	if a.SignalStatus.Strength != -128 {
		t.Errorf("strength should be -128 but got %d", a.SignalStatus.Strength)
	}

	if a.SignalStatus.Quality != 0 {
		t.Errorf("quality should be 0 but got %d", a.SignalStatus.Quality)
	}

	if a.SignalStatus.Qam != 0 {
		t.Errorf("qam should be 0 but got %d", a.SignalStatus.Qam)
	}

}
