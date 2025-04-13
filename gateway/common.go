package gateway

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	rsp "zhtcloud/pkg/response"
	"zhtcloud/utils/logger"

	ws "github.com/gorilla/websocket"
)

const (
	COORS = "coors"
)

const (
	DRONE_INFO = iota
	RUNNING_STATUS
)

/* running_status */
const (
	PICKING   = iota // Now picking up the goods
	PICKED           // Have picked the goods
	DELIVERED        // Have delivered done
)

// 飞行模式
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

// 飞行状态
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

// Result 响应基本结构体
type Result struct {
	Code int16        `json:"code"`
	Msg  string       `json:"msg"`
	Data *interface{} `json:"data,omitempty"`
}

// 无人机坐标定位
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

// 无人机旋转或俯仰角度
type Attitude struct {
	Roll       float64 `json:"roll"`
	Pitch      float64 `json:"pitch"`
	Yaw        float64 `json:"yaw"`
	RollSpeed  float64 `json:"rollspeed"`
	PitchSpeed float64 `json:"pitchspeed"`
	YawSpeed   float64 `json:"yawspeed"`
}

// 系统状态信息
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

type DroneData struct {
	GLOBAL_POSITION_INT *GlobalPosition `json:"GLOBAL_POSITION_INT"`
	ATTITUDE            *Attitude       `json:"ATTITUDE"`
	SYS_STATUS          *SysStatus      `json:"SYS_STATUS"`
	MODE                uint8           `json:"MODE"`
	STATUS              uint8           `json:"STATUS"`
	TYPE                uint8           `json:"TYPE"`
}

type RunningStatus struct {
	TYPE           uint8 `json:"TYPE"`
	RUNNING_STATUS uint8 `json:"RUNNING_STATUS"`
}

// Coordinate 坐标点
type Coordinate [2]float64

// Coordinates 起始坐标和终点坐标
type Coordinates struct {
	Coords []Coordinate `json:"coords"`
}

// HandleErrorReqMethod 处理错误的请求方法
func HandleErrorReqMethod(rs *Result) {
	rs.Code = rsp.INVALID_PARAMS
	rs.Msg = rsp.CodeToMsgMap[rsp.INVALID_PARAMS]
}

// HandleReqBodyDecode 处理解析请求体，如果返回true，则说明需要调用函数直接return
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

// HandleResBodyEncode 处理编码响应体
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
var ctx = context.Background()
