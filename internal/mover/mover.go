package mover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"jackson/internal/maze"
)

type Mover interface {
	Forward(int)
	Backward(int)

	Left()
	Right()
	Rotate()

	CellState(d maze.Direction) Cell
}

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

const wallThreshold int = 100

func (c CellResp) ToCell(robotDir maze.Direction) Cell {
	var w maze.Wall
	if c.Laser.Back < wallThreshold {
		switch robotDir {
		case maze.Up:
			w.Add(maze.D)
		case maze.Right:
			w.Add(maze.L)
		case maze.Down:
			w.Add(maze.U)
		case maze.Left:
			w.Add(maze.R)
		}
	}
	if c.Laser.Front < wallThreshold {
		switch robotDir {
		case maze.Up:
			w.Add(maze.U)
		case maze.Right:
			w.Add(maze.R)
		case maze.Down:
			w.Add(maze.D)
		case maze.Left:
			w.Add(maze.L)
		}
	}
	if c.Laser.Left < wallThreshold {
		switch robotDir {
		case maze.Up:
			w.Add(maze.L)
		case maze.Right:
			w.Add(maze.U)
		case maze.Down:
			w.Add(maze.R)
		case maze.Left:
			w.Add(maze.D)
		}
	}
	if c.Laser.Right < wallThreshold {
		switch robotDir {
		case maze.Up:
			w.Add(maze.R)
		case maze.Right:
			w.Add(maze.D)
		case maze.Down:
			w.Add(maze.L)
		case maze.Left:
			w.Add(maze.U)
		}
	}

	log.Printf("wall: %s, dir: %s", w, robotDir)

	return Cell{
		Wall: w,
	}
}

type baseMover struct {
	motorsIP  string
	sensorsIP string
	id        string
}

func (m baseMover) move(direction string, value int) (*http.Response, error) {
	/* move PUT:
	http://[robot_ip]/move
	{"id": "123456", "direction":"forward", "len": 100}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.motorsIP, "move")
	log.Printf("send /move to %s, dir=%s, val=%v\n", reqUrl, direction, value)

	reqBody, err := json.Marshal(struct {
		Id        string `json:"id"`
		Direction string `json:"direction"`
		Len       int    `json:"len"`
	}{
		Id:        m.id,
		Direction: direction,
		Len:       value,
	})
	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPut, reqUrl, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", `application/json`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("get /move", resp)
	return resp, err
}

func (m baseMover) getSensor() (*CellResp, error) {
	/* sensors POST:
	http://[robot_ip]/sensor
	{"id": "123456", "type": "all"}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.sensorsIP, "sensor")
	log.Printf("send /sensor to %s\n", reqUrl)

	reqBody, err := json.Marshal(map[string]string{
		"id":   m.id,
		"type": "all",
	})
	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPost, reqUrl, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", `application/json`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("get /sensor", resp)

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
