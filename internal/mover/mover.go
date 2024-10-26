package mover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type float float32

type VagifMover struct {
	angle   float
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

func (m VagifMover) isNotAimedAtCenter() bool {
	return true
}

func (m VagifMover) centering() {

}

func (m VagifMover) Forward(cell int) {
	if m.isNotAimedAtCenter() {
		m.centering()
	}
	// get current angle from memory
	// if angle incorrect -> small rotation
	// transform cell parameter to mm
	// send command to mouse
	// check position and angle
	// save angle
}

func (m VagifMover) Backward(cell int) {
	// same as forward
}

func (m VagifMover) Left() {
	// get current angle from memory
}

func (m VagifMover) Right() {
	// same as Left
}

func (m VagifMover) CellState() Cell {
	return Cell{}
}
