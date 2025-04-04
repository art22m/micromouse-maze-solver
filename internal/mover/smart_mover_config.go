package mover

import (
	"encoding/json"
	"os"
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

func saveCalibrationConfig(config SmartMoverCalibrationConfig) {
	configFile, err := os.OpenFile(calibrationConfigPath, os.O_WRONLY, 0x666)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(config)
	if err != nil {
		panic(err)
	}
}

type RobotConfig struct {
	// min delay between /move call and subsequent api calls
	SensorsDelayMs  int `json:"sensors_delay_ms"`
	MovementDelayMs int `json:"movement_delay_ms"`

	// robot parameters
	Width                 int `json:"width"`
	Length                int `json:"length"`
	FrontLaserZero        int `json:"front_laser_zero"`
	BackLaserZero         int `json:"back_laser_zero"`
	SideLasersZero        int `json:"side_lasers_zero"`
	MaxDerivationFromAxis int `json:"max_derivation_from_axis"`

	// map parameters
	CellSize int `json:"cell_size"`
}

type SmartMoverCalibrationConfig struct {
	ForwardRatio   float32 `json:"forward_ratio"`
	BackwardRatio  float32 `json:"backward_ratio"`
	TurnRightRatio float32 `json:"turn_right_ratio"`
	TurnLeftRatio  float32 `json:"turn_left_ratio"`
	Turn180Ratio   float32 `json:"turn_180_ratio"`
	// can't turn less than this
	MinTurn int `json:"min_turn"`
}
