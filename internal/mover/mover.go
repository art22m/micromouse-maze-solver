package mover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"jackson/internal/maze"
)

type Cell struct {
	Wall maze.Wall
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
	reqUrl := fmt.Sprint("http://", m.robotIP, "/move")

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

func (m VagifMover) getSensor() (*http.Response, error) {
	/* sensors POST:
	http://[robot_ip]/sensor
	{"id": "123456", "type": "all"}
	*/
	reqUrl := fmt.Sprint("http://", m.robotIP, "/sensor")

	reqBody, _ := json.Marshal(map[string]string{
		"id":   m.ID,
		"type": "all",
	})

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPost, reqUrl, requestBody)
	req.Header.Add("Content-Type", `application/json`)

	resp, err := http.DefaultClient.Do(req)
	return resp, err
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
	return Cell{
		Left:     0,
		Right:    0,
		Forward:  0,
		Backward: 0,
	}
}
