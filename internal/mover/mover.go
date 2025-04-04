package mover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"jackson/internal/maze"
)

type Mover interface {
	Forward(int)
	Backward(int)

	Left()
	Right()
	Rotate()

	CellState(d maze.Direction) Cell

	Reset()
}

type Cell struct {
	Wall maze.Wall
}

type CellResp struct {
	Laser struct {
		Back    float64 `json:"1"`
		Left    float64 `json:"2"`
		Right45 float64 `json:"3"`
		Front   float64 `json:"4"`
		Right   float64 `json:"5"`
		Left45  float64 `json:"6"`
	} `json:"laser"`
	Imu struct {
		Roll  float64 `json:"roll"`
		Pitch float64 `json:"pitch"`
		Yaw   float64 `json:"yaw"`
	} `json:"imu"`
}

const wallThreshold float64 = 140

func (c CellResp) ToCell(robotDir maze.Direction) Cell {
	var w maze.Wall
	if c.Laser.Back < wallThreshold {
		w.Add(maze.Wall(maze.Down.GlobalFrom(robotDir)))
	}
	if c.Laser.Front < wallThreshold {
		w.Add(maze.Wall(maze.Up.GlobalFrom(robotDir)))
	}
	if c.Laser.Left < wallThreshold {
		w.Add(maze.Wall(maze.Left.GlobalFrom(robotDir)))
	}
	if c.Laser.Right < wallThreshold {
		w.Add(maze.Wall(maze.Right.GlobalFrom(robotDir)))
	}
	return Cell{
		Wall: w,
	}
}

type baseMover struct {
	motorsIP  string
	sensorsIP string
	id        string

	logger *logrus.Entry
}

func (m *baseMover) move(direction string, value int) {
	/* move PUT:
	http://[robot_ip]/move
	{"id": "123456", "direction":"forward", "len": 100}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.motorsIP, "move")
	m.logger.Infof("send /move to %s, dir=%s, val=%v\n", reqUrl, direction, value)

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
		m.logger.Fatal(err)
	}

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPut, reqUrl, requestBody)
	if err != nil {
		m.logger.Fatal(err)
	}
	req.Header.Add("Content-Type", `application/json`)

	//client := retryablehttp.NewClient()
	//client.Logger = m.logger
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.logger.Fatal(err)
	}
	m.logger.Info("get /move", resp.Body)
}

const sensorValueRetryThreshold = 25000.0

func (m *baseMover) getSensor() *CellResp {
	/* sensors POST:
	http://[robot_ip]/sensor
	{"id": "123456", "type": "all"}
	*/
	reqUrl := fmt.Sprintf("http://%s/%s", m.sensorsIP, "sensor")
	m.logger.Infof("send /sensor to %s\n", reqUrl)

	reqBody, err := json.Marshal(map[string]string{
		"id":   m.id,
		"type": "all",
	})
	if err != nil {
		m.logger.Fatal(err)
	}

	requestBody := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(http.MethodPost, reqUrl, requestBody)
	if err != nil {
		m.logger.Fatal(err)
	}
	req.Header.Add("Content-Type", `application/json`)

	//client := retryablehttp.NewClient()
	//client.Logger = m.logger
	//resp, err := client.Do(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.logger.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.Fatal(err)
	}
	m.logger.Info("get /sensor", string(body))

	var cellResp CellResp
	err = json.Unmarshal(body, &cellResp)
	if err != nil {
		m.logger.Fatal(err)
	}

	l := cellResp.Laser
	for _, v := range []float64{
		l.Left, l.Right, l.Front, l.Back, l.Left45, l.Right45,
	} {
		if v > sensorValueRetryThreshold {
			m.logger.Error("values is more then sensorValueRetryThreshold")
			time.Sleep(50 * time.Millisecond)
			return m.getSensor()
		}
	}
	return &cellResp
}
