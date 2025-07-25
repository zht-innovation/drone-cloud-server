package shared

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils/logger"

	ws "github.com/gorilla/websocket"
)

// mqtt info
const (
	MQTT_BROKER = "http://110.42.101.86:18083"
	API_PREFIX  = "/api/v5"
)

// etc channels
const (
	COORS = "coors"
)

const (
	DRONE_INFO = iota
	RUNNING_STATUS
)

// running_status
const (
	PICKING   = iota // Now picking up the goods
	PICKED           // Have picked the goods
	DELIVERED        // Have delivered done
)

// drone flying type
const (
	STABLIZE = iota
	ACRO
	ALT_HOLD
	AUTO
	GUIDED
	LOITER
	RTL
	CIRCLE
	POSITION
	LAND
	OF_LOITER
	DRIFT
	SPORT
	FLIP
	AUTOTUNE
	POSHOLD
	BRAKE
	THROW
	AVOID_ADSB
	GUIDED_NOGPS
	SMART_RTL
	FLOWHOLD
	FOLLOW
	ZIGZAG
	SYSTEMID
	AUTOROTATE
	AUTO_RTL
	MANUAL
)

// drone flying state
const (
	MAV_STATE_UNINIT = iota
	MAV_STATE_BOOT
	MAV_STATE_CALIBRATING
	MAV_STATE_STANDBY
	MAV_STATE_ACTIVE
	MAV_STATE_CRITICAL
	MAV_STATE_EMERGENCY
	MAX_STATE_POWEROFF
	MAX_STATE_FLIGHT_TERMINATION
)

// Result is the standard response structure
type Result struct {
	Code int16        `json:"code"`
	Msg  string       `json:"msg"`
	Data *interface{} `json:"data,omitempty"`
}

// This structure is used to represent the global position of a drone
type GlobalPosition struct {
	TimeBootMs  uint32 `json:"time_boot_ms"`
	Lat         int32  `json:"lat"`
	Lon         int32  `json:"lon"`
	Alt         int32  `json:"alt"`
	RelativeAlt int32  `json:"relative_alt"`
	Vx          int16  `json:"vx"`
	Vy          int16  `json:"vy"`
	Vz          int16  `json:"vz"`
	Hdg         uint16 `json:"hdg"`
}

// drone attitude information
type Attitude struct {
	Roll       float64 `json:"roll"`
	Pitch      float64 `json:"pitch"`
	Yaw        float64 `json:"yaw"`
	RollSpeed  float64 `json:"rollspeed"`
	PitchSpeed float64 `json:"pitchspeed"`
	YawSpeed   float64 `json:"yawspeed"`
}

// system status information
type SysStatus struct {
	OnboardControlSensorsPresent uint32 `json:"onboard_control_sensors_present"`
	OnboardControlSensorsEnabled uint32 `json:"onboard_control_sensors_enabled"`
	OnboardControlSensorsHealth  uint32 `json:"onboard_control_sensors_health"`
	Load                         uint16 `json:"load"`
	VoltageBattery               uint16 `json:"voltage_battery"`
	CurrentBattery               int16  `json:"current_battery"`
	BatteryRemaining             int8   `json:"battery_remaining"`
	DropRateComm                 uint16 `json:"drop_rate_comm"`
	ErrorsComm                   uint16 `json:"errors_comm"`
	ErrorsCount1                 uint16 `json:"errors_count1"`
	ErrorsCount2                 uint16 `json:"errors_count2"`
	ErrorsCount3                 uint16 `json:"errors_count3"`
	ErrorsCount4                 uint16 `json:"errors_count4"`
}

// Motor information
type Motor struct {
	Current     uint16 `json:"current"`
	Voltage     uint16 `json:"voltage"`
	Speed       uint16 `json:"speed"`
	Temperature uint16 `json:"temperature"`
}

type DroneData struct {
	DID string `json:"did"` // Drone ID, mac address of the drone

	GLOBAL_POSITION_INT *GlobalPosition `json:"GLOBAL_POSITION_INT"`
	ATTITUDE            *Attitude       `json:"ATTITUDE"`
	SYS_STATUS          *SysStatus      `json:"SYS_STATUS"`
	MOTOR               *Motor          `json:"MOTOR"`

	MODE                      uint8 `json:"MODE"`
	STATUS                    uint8 `json:"STATUS"`
	TYPE                      uint8 `json:"TYPE"`
	GPS_NUM                   uint8 `json:"GPS_NUM"`
	REMOTE_CONTROL_CONNECTION bool  `json:"REMOTE_CONTROL_CONNECTION"`
	FLIGHT_CONTROLER_UNLOCK   bool  `json:"FLIGHT_CONTROLER_UNLOCK"`
}

type RunningStatus struct {
	TYPE           uint8 `json:"TYPE"`
	RUNNING_STATUS uint8 `json:"RUNNING_STATUS"`
}

// Coordinate waypoint
type Coordinate [2]float64

// Coordinates start point and end point
// Coords[0] is the start point, Coords[1] is the end
type Coordinates struct {
	Coords []Coordinate `json:"coords"`
}

// HandleErrorReqMethod hanles the error when the request method is not supported
func HandleErrorReqMethod(rs *Result) {
	rs.Code = rsp.INVALID_PARAMS
	rs.Msg = rsp.CodeToMsgMap[rsp.INVALID_PARAMS]
}

// HandleReqBodyDecode handles decoding the request body
// r is the request body, v is the target struct to decode into, rs is the response result
func HandleReqBodyDecode(r io.ReadCloser, v any, rs *Result) bool {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		_ = r.Close()
		logger.Error("request decode error: %v", err)
		rs.Code = rsp.INVALID_PARAMS
		rs.Msg = rsp.CodeToMsgMap[rsp.INVALID_PARAMS]
		return true
	}
	return false
}

// HandleResBodyEncode handles encoding the response body
// w is the response writer, rs is the response result
func HandleResBodyEncode(w io.Writer, rs *Result) {
	err := json.NewEncoder(w).Encode(rs)
	if err != nil {
		logger.Error("JSON encode error: %v", err)
		rs.Code = rsp.SERVER_ERROR
		rs.Msg = rsp.CodeToMsgMap[rsp.SERVER_ERROR]
		if rs.Data != nil {
			rs.Data = nil
		}
	}
}

var Upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Ctx = context.Background()
