package mover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"jackson/internal/maze"
)

type Cell struct {
	Wall maze.Wall
}

type CellResp struct {
	Laser struct {
		Back    int `json:"1"`
		Left    int `json:"2"`
		Right45 int `json:"3"`
		Front   int `json:"4"`
		Right   int `json:"5"`
		Left45  int `json:"6"`
	} `json:"laser"`
	Imu struct {
		Roll  int `json:"roll"`
		Pitch int `json:"pitch"`
		Yaw   int `json:"yaw"`
	} `json:"imu"`
}

type Mover interface {
	Forward(int)
	Backward(int)

	Left()
	Right()

	CellState() Cell
}

type VagifMover struct {
	angle   int // from 0 to 360 degrees
	front   int
	back    int
	left    int
	right   int
	robotIP string
	ID      string
}

func NewVagifMover(robotIP, ID string) *VagifMover {
	return &VagifMover{
		angle:   0,
		robotIP: robotIP,
		ID:      ID,
	}
}

func (m VagifMover) move(direction string, value int) (*http.Response, error) {
	/* move PUT:
	http://[robot_ip]/move
	{"id": "123456", "direction":"forward", "len": 100}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.robotIP, "move")

	reqBody, _ := json.Marshal(struct {
		Id        string `json:"id"`
		Direction string `json:"direction"`
		Len       int    `json:"len"`
	}{
		Id:        m.ID,
		Direction: direction,
		Len:       value,
	})
	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPut, reqUrl, requestBody)
	req.Header.Add("Content-Type", `application/json`)

	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func (m VagifMover) getSensor() (*CellResp, error) {
	/* sensors POST:
	http://[robot_ip]/sensor
	{"id": "123456", "type": "all"}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.robotIP, "sensor")

	reqBody, _ := json.Marshal(map[string]string{
		"id":   m.ID,
		"type": "all",
	})

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPost, reqUrl, requestBody)
	req.Header.Add("Content-Type", `application/json`)

	resp, err := http.DefaultClient.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the struct
	var cellResp CellResp
	err = json.Unmarshal(body, &cellResp)
	if err != nil {
		return nil, err
	}
	return &cellResp, err
}

const (
	Front     = 0
	Right     = 90
	Down      = 180
	Left      = 270
	Tolerance = 5 // допустимая погрешность
)

func (m VagifMover) closestDirectionAndAngle() (string, int) {
	directions := map[string]int{
		"Front": Front,
		"Right": Right,
		"Down":  Down,
		"Left":  Left,
	}

	minDiff := 360.0
	closest := ""
	var angleDiff int

	for dir, angle := range directions {
		diff := angle - m.angle
		if diff > 180 {
			diff -= 360
		} else if diff < -180 {
			diff += 360
		}

		if math.Abs(float64(diff)) < minDiff {
			minDiff = math.Abs(float64(diff))
			closest = dir
			angleDiff = diff
		}
	}

	return closest, angleDiff
}

func (m VagifMover) isNotAimedAtCenter() bool {
	_, angleDiff := m.closestDirectionAndAngle()
	return math.Abs(float64(angleDiff)) > Tolerance
}

// Метод для центрирования робота к ближайшей оси
func (m VagifMover) centering() {
	_, angleDiff := m.closestDirectionAndAngle()

	if angleDiff > 0 {
		m.RotateRight(int(angleDiff))
	} else if angleDiff < 0 {
		m.RotateRight(int(-angleDiff))
	} else {
		fmt.Println("Робот уже отцентрован.")
	}
}

func (m VagifMover) Forward(cell int) {
	if m.isNotAimedAtCenter() {
		m.centering()
	}

	// transform cell parameter to mm
	// send command to mouse
	// check position and angle
	// save angle
}

func (m VagifMover) Backward(cell int) {
	// same as forward
}

func (m VagifMover) RotateLeft(degrees int) {
	_, err := m.move("left", degrees)
	if err != nil {
		return
	}
}

func (m VagifMover) Left() {
	// get current angle from memory
}

func (m VagifMover) RotateRight(degrees int) {
	_, err := m.move("right", degrees)
	if err != nil {
		return
	}
}
func (m VagifMover) Right() {
	// same as Left
}

func (m VagifMover) CellState() Cell {
	return Cell{}
}
