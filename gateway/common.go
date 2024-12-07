package gateway

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

// 无人机坐标定位
type GlobalPosition struct {
	TimeBootMs uint32 `json:"time_boot_ms"`
	Lat	int32 `json:"lat"`
	Lon int32 `json:"lon"`
	Alt int32 `json:"alt"`
	RelativeAlt int32 `json:"relative_alt"`
	Vx int16 `json:"vx"`
	Vy int16 `json:"vy"`
	Vz int16 `json:"vz"`
	Hdg uint16 `json:"hdg"`
}

// 无人机旋转或俯仰角度
type Attitude struct {
	Roll float64 `json:"roll"`
	Pitch float64 `json:"pitch"`
	Yaw float64 `json:"yaw"`
	RollSpeed float64 `json:"rollspeed"`
	PitchSpeed float64 `json:"pitchspeed"`
	YawSpeed float64 `json:"yawspeed"`
}

// 系统状态信息
type SysStatus struct {
	OnboardControlSensorsPresent uint32 `json:"onboard_control_sensors_present"`
	OnboardControlSensorsEnabled uint32 `json:"onboard_control_sensors_enabled"`
	OnboardControlSensorsHealth	uint32 `json:"onboard_control_sensors_health"`
	Load uint16 `json:"load"`
	VoltageBattery uint16 `json:"voltage_battery"`
	CurrentBattery int16 `json:"current_battery"`
	BatteryRemaining int8 `json:"battery_remaining"`
	DropRateComm uint16 `json:"drop_rate_comm"`
	ErrorsComm uint16 `json:"errors_comm"`
	ErrorsCount1 uint16 `json:"errors_count1"`
	ErrorsCount2 uint16 `json:"errors_count2"`
	ErrorsCount3 uint16 `json:"errors_count3"`
	ErrorsCount4 uint16 `json:"errors_count4"`
}

type DroneData struct {
	GLOBAL_POSITION_INT *GlobalPosition `json:"GLOBAL_POSITION_INT"`
	ATTITUDE *Attitude `json:"ATTITUDE"`
	SYS_STATUS *SysStatus `json:"SYS_STATUS"`
	MODE uint8 `json:"MODE"`
	STATUS uint8 `json:"STATUS"`
}