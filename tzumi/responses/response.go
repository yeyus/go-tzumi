package responses

import "fmt"

type ResponseType string

const SignalStatusAck ResponseType = "signal_status_ack"
const LoginAck ResponseType = "login_ack"
const CheckLockAck ResponseType = "check_lock_ack"
const TunerAck ResponseType = "tuner_ack"

type Response struct {
	Type         ResponseType `xml:"type,attr"`
	AckInfo      AckInfo      `xml:"ack_info"`
	DeviceInfo   DeviceInfo   `xml:"deviceinfo,omitempty"`
	SignalStatus SignalStatus `xml:"signal_status,omitempty"`
}

func (r Response) String() string {
	switch r.Type {
	case LoginAck:
		return fmt.Sprintf("Login ACK (%d, %s) -> [Name: %s, Type: %s, System: %s, Hw: %s, Sw: %s]", r.AckInfo.ReturnValue, r.AckInfo.Message, r.DeviceInfo.DeviceName, r.DeviceInfo.DeviceType, r.DeviceInfo.TvType, r.DeviceInfo.HwVersion, r.DeviceInfo.SwVersion)
	case CheckLockAck:
		return fmt.Sprintf("Check Lock ACK (%d, %s)", r.AckInfo.ReturnValue, r.AckInfo.Message)
	case TunerAck:
		return fmt.Sprintf("Tuner ACK (%d, %s)", r.AckInfo.ReturnValue, r.AckInfo.Message)
	case SignalStatusAck:
		return fmt.Sprintf("Signal Status ACK (%d, %s) -> [Strength: %d, Quality: %d, QAM: %d]", r.AckInfo.ReturnValue, r.AckInfo.Message, r.SignalStatus.Strength, r.SignalStatus.Quality, r.SignalStatus.Qam)
	default:
		return "Unknown ACK message"
	}
}
