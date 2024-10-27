using UnityEngine;
using System;
using System.Collections.Concurrent;
using System.Net;
using System.Text;
using System.IO;

public class Runner : MonoBehaviour
{
    public MapBuilder builder;
    public MouseController mouseController;
    
    public int mapSize;
    public float poleSize;
    public float wallLength;
    
    [SerializeField] private int _Port = 8080;
    private HttpListener listener;

    private static readonly ConcurrentQueue<Action> _executionQueue = new ConcurrentQueue<Action>();

    void Start() {
        builder.build(mapSize, poleSize, wallLength);
        mouseController.InitMouse(poleSize + wallLength, 0, mapSize - 1);
        Serve();
    }

    void LateUpdate() {
        if (Input.GetKeyUp(KeyCode.UpArrow)) {
            bool status = mouseController.MoveForward(poleSize + wallLength, () => {});
        }
        if (Input.GetKeyUp(KeyCode.DownArrow)) {
            bool status = mouseController.MoveBackward(poleSize + wallLength, () => {});
        }
        if (Input.GetKeyUp(KeyCode.RightArrow)) {
            bool status = mouseController.RotateRight(90, () => {});
        }
        if (Input.GetKeyUp(KeyCode.LeftArrow)) {
            bool status = mouseController.RotateLeft(90, () => {});
        }

        while (_executionQueue.TryDequeue(out var action)) {
            action?.Invoke();
        }
    }

    private void Serve() {
        listener = new HttpListener();
        listener.Prefixes.Add($"http://*:{_Port}/");
        listener.Start();
        listener.BeginGetContext(new AsyncCallback(OnRequestCallback), null);
    }

    private void OnRequestCallback (IAsyncResult result) {
        HttpListenerContext context = listener.EndGetContext(result);
        HttpListenerRequest req = context.Request;
        HttpListenerResponse resp = context.Response;

        string path = req.RawUrl;

        StreamReader r = new StreamReader(req.InputStream);
        string body = r.ReadToEnd();

        if (path == "/move") {
            _executionQueue.Enqueue(() => {
                Debug.Log(path);

                Action onFinish = () => {
                    resp.SendChunked = false;
                    resp.StatusCode = 200;
                    resp.Close();

                    if (listener.IsListening) {
                        listener.BeginGetContext(new AsyncCallback(OnRequestCallback), null);
                    }
                };

                bool status = OnMoveRequest(body, onFinish);

                if (!status) {
                    Debug.LogError("Move request not executed");
                    resp.SendChunked = false;
                    resp.StatusCode = 200;
                    resp.Close();

                    if (listener.IsListening) {
                        listener.BeginGetContext(new AsyncCallback(OnRequestCallback), null);
                    }
                } 
            });
        } else if (path == "/sensor") {
            _executionQueue.Enqueue(() => {
                Debug.Log(path);

                SensorResponse sensorResponse = OnSensorRequest();

                string res = JsonUtility.ToJson(sensorResponse);

                res = res.Replace("_", "");

                resp.Headers.Set("Content-Type", "application/json");
                byte[] buffer = Encoding.UTF8.GetBytes(res);
                resp.SendChunked = false;
                resp.StatusCode = 200;
                resp.ContentLength64 = buffer.Length;
                using Stream ros = resp.OutputStream;
                ros.Write(buffer, 0, buffer.Length);
                resp.Close();

                if (listener.IsListening) {
                    listener.BeginGetContext(new AsyncCallback(OnRequestCallback), null);
                } 
            });
        } else {
            Debug.LogError("Not found - " + path);

            resp.SendChunked = false;
		    resp.StatusCode = 404;
            resp.Close();

            if (listener.IsListening) {
                listener.BeginGetContext(new AsyncCallback(OnRequestCallback), null);
            }
        }
	}

    private bool OnMoveRequest(string body, Action onFinish) {
        MoveRequest moveRequest;
        try {
            moveRequest = JsonUtility.FromJson<MoveRequest>(body);
        } catch(Exception e) {
            Debug.LogError(e);
            return false;
        }

        if (moveRequest == null) {
            Debug.LogError("Invalid json: " + body);
            return false;
        }

        if (moveRequest.direction == "forward") {
            return mouseController.MoveForward(moveRequest.len, onFinish);
        }
        if (moveRequest.direction == "backward") {
            return mouseController.MoveBackward(moveRequest.len, onFinish);
        }
        if (moveRequest.direction == "left") {
            return mouseController.RotateLeft(moveRequest.len, onFinish);
        }
        if (moveRequest.direction == "right") {
            return mouseController.RotateRight(moveRequest.len, onFinish);
        }
        
        Debug.LogError("Invalid direction: " + moveRequest.direction); 

        return false;
    }

    private SensorResponse OnSensorRequest() {
        SensorResponse resp = new SensorResponse();

        SensorReading reading = mouseController.mouse.GetComponent<Mouse>().readSensors();

        resp.imu = new Imu();
        resp.imu.yaw = (reading.yaw % 360 + 360) % 360;
        resp.laser = new Laser();
        resp.laser._1 = reading.backward;
        resp.laser._2 = reading.left;
        resp.laser._3 = reading.rightDiagonal;
        resp.laser._4 = reading.forward;
        resp.laser._5 = reading.right;
        resp.laser._6 = reading.left;

        return resp;
    }
}

[Serializable]
class MoveRequest {
    public float len;
    public string direction;
    public string id;
}

[Serializable]
class Laser {
    public float _1;
    public float _2;
    public float _3;
    public float _4;
    public float _5;
    public float _6;
}

[Serializable]
class Imu {
    public float roll;
    public float pitch;
    public float yaw;
}

[Serializable]
class SensorResponse {
    public Laser laser;
    public Imu imu;
}

