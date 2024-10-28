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
	sensorsDelay  time.Duration `json:"sensors_delay"`
	movementDelay time.Duration `json:"movement_delay"`

	// robot parameters
	width          int
	length         int
	frontLaserZero int
	backLaserZero  int
	sideLasersZero int
}

type SmartMoverCalibrationConfig struct {
	forwardRatio   float32 `json:"forward_ratio"`
	backwardRatio  float32 `json:"backward_ratio"`
	turnRightRatio float32 `json:"right_turn_ratio"`
	turnLeftRatio  float32 `json:"turn_left_ratio"`
	turn180Ratio   float32 `json:"turn_180_ratio"`
	// can't turn less than this
	minTurn int `json:"min_turn"`
}
