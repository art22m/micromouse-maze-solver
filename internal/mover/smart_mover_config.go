package mover

import (
	"encoding/json"
	"os"
	"time"
)

const (
	robotConfigPath       = "configs/robot_config.json"
	calibrationConfigPath = "configs/calibration_config.json"
)

type SmartMoverConfig struct {
	robot       RobotConfig
	calibration SmartMoverCalibrationConfig
}

func LoadSmartMoverConfig() SmartMoverConfig {
	return SmartMoverConfig{
		robot:       loadRobotConfig(),
		calibration: loadCalibrationConfig(),
	}
}

func loadRobotConfig() RobotConfig {
	var config RobotConfig
	configFile, err := os.Open(robotConfigPath)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}

func loadCalibrationConfig() SmartMoverCalibrationConfig {
	var config SmartMoverCalibrationConfig
	configFile, err := os.Open(calibrationConfigPath)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
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
