package mover

import "time"

type SmartMoverConfig struct {
	robot       RobotConfig
	calibration SmartMoverCalibrationConfig
}

func loadSmartMoverConfig() SmartMoverConfig {
	// TODO
	return SmartMoverConfig{}
}

type RobotConfig struct {
	// min delay between /move call and subsequent api calls
	SensorsDelay  time.Duration `json:"sensors_delay"`
	MovementDelay time.Duration `json:"movement_delay"`

	// robot parameters
	Width          int `json:"width"`
	Length         int `json:"length"`
	FrontLaserZero int `json:"front_laser_zero"`
	BackLaserZero  int `json:"back_laser_zero"`
	SideLasersZero int `json:"side_lasers_zero"`
}

type SmartMoverCalibrationConfig struct {
	ForwardRatio   float32 `json:"forward_ratio"`
	BackwardRatio  float32 `json:"backward_ratio"`
	TurnRightRatio float32 `json:"right_turn_ratio"`
	TurnLeftRatio  float32 `json:"turn_left_ratio"`
	Turn180Ratio   float32 `json:"turn_180_ratio"`
	// can't turn less than this
	MinTurn int `json:"min_turn"`
}
