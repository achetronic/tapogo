package types

// RequestParamsSpec represents TODO
type RequestParamsSpec struct {
	DeviceOn bool `json:"device_on"`
}

// RequestSpec represents TODO
type RequestSpec struct {
	Method          string             `json:"method"`
	RequestTimeMils int                `json:"requestTimeMils"`
	Params          *RequestParamsSpec `json:"params,omitempty"`
}

// ResponseSpec represents TODO
type ResponseSpec struct {
	Result *struct {

		// GetDeviceInfo related fields
		DeviceId           string `json:"device_id,omitempty"`
		FwVer              string `json:"fw_ver,omitempty"`
		HwVer              string `json:"hw_ver,omitempty"`
		Type               string `json:"type,omitempty"`
		Model              string `json:"model,omitempty"`
		Mac                string `json:"mac,omitempty"`
		HwId               string `json:"hw_id,omitempty"`
		FwId               string `json:"fw_id,omitempty"`
		OemId              string `json:"oem_id,omitempty"`
		Ip                 string `json:"ip,omitempty"`
		TimeDiff           int    `json:"time_diff,omitempty"`
		Ssid               string `json:"ssid,omitempty"`
		Rssi               int    `json:"rssi,omitempty"`
		SignalLevel        int    `json:"signal_level,omitempty"`
		AutoOffStatus      string `json:"auto_off_status,omitempty"`
		AutoOffRemainTime  int    `json:"auto_off_remain_time,omitempty"`
		Latitude           int    `json:"latitude,omitempty"`
		Longitude          int    `json:"longitude,omitempty"`
		Lang               string `json:"lang,omitempty"`
		Avatar             string `json:"avatar,omitempty"`
		Region             string `json:"region,omitempty"`
		Specs              string `json:"specs,omitempty"`
		Nickname           string `json:"nickname,omitempty"`
		HasSetLocationInfo bool   `json:"has_set_location_info,omitempty"`
		DeviceOn           bool   `json:"device_on,omitempty"`
		OnTime             int    `json:"on_time,omitempty"`
		DefaultStates      *struct {
			Type  string `json:"type,omitempty"`
			State *struct {
			} `json:"state,omitempty"`
		} `json:"default_states,omitempty"`
		Overheated            bool   `json:"overheated,omitempty"`
		PowerProtectionStatus string `json:"power_protection_status,omitempty"`
		OvercurrentStatus     string `json:"overcurrent_status,omitempty"`

		// GetEnergyUsage related fields
		TodayRuntime      int    `json:"today_runtime,omitempty"`
		MonthRuntime      int    `json:"month_runtime,omitempty"`
		TodayEnergy       int    `json:"today_energy,omitempty"`
		MonthEnergy       int    `json:"month_energy,omitempty"`
		LocalTime         string `json:"local_time,omitempty"`
		ElectricityCharge []int  `json:"electricity_charge,omitempty"`
		CurrentPower      int    `json:"current_power,omitempty"`

		//
	} `json:"result,omitempty"`
	ErrorCode int `json:"error_code"`
}
