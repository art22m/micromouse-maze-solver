using UnityEngine;

public class Mouse : MonoBehaviour
{
    public bool isExact;

    public float forwardErrorMean;
    public float forwardErrorStd;
    public float forwardAbsErrorMean;
    public float forwardAbsErrorStd;
    public float backwardErrorMean;
    public float backwardErrorStd;
    public float backwardAbsErrorMean;
    public float backwardAbsErrorStd;
    public float rightErrorMean;
    public float rightErrorStd;
    public float rightAbsErrorMean;
    public float rightAbsErrorStd;
    public float leftErrorMean;
    public float leftErrorStd;
    public float leftAbsErrorMean;
    public float leftAbsErrorStd;
    public float yawAbsErrorStd;
    public float northYawReading;

    public uint forwardErrorSeed;
    public uint forwardAbsErrorSeed;
    public uint backwardErrorSeed;
    public uint backwardAbsErrorSeed;
    public uint rightErrorSeed;
    public uint rightAbsErrorSeed;
    public uint leftErrorSeed;
    public uint leftAbsErrorSeed;
    public uint yawAbsErrorSeed;

    public float laserUpdateIntervalSeconds;
    public float yawUpdateIntervalSeconds;
    public bool spreadSensorsUpdate;

    private Gaussian yawAbsError;

    public Ray forward;
    public Ray backward;
    public Ray right;
    public Ray left;
    public Ray diagonalRight;
    public Ray diagonalLeft;

    private SensorReading sensorReading;

    private float time;
    private float nextForwardUpdateAt;
    private float nextBackwardUpdateAt;
    private float nextLeftUpdateAt;
    private float nextRightUpdateAt;
    private float nextLeftDiagonalUpdateAt;
    private float nextRightDiagonalUpdateAt;
    private float nextYawUpdateAt;

    void Start() {
        yawAbsError = new Gaussian(0, yawAbsErrorStd, new Unity.Mathematics.Random(yawAbsErrorSeed));

        sensorReading = new SensorReading();
        sensorReading.forward = forward.getDistance();
        sensorReading.backward = backward.getDistance();
        sensorReading.left = left.getDistance();
        sensorReading.right = right.getDistance();
        sensorReading.leftDiagonal = diagonalLeft.getDistance();
        sensorReading.rightDiagonal = diagonalRight.getDistance();
        sensorReading.yaw = readYaw();

        time = 0;
        
        nextForwardUpdateAt = nextBackwardUpdateAt = nextLeftUpdateAt = nextRightDiagonalUpdateAt = nextRightUpdateAt = nextLeftDiagonalUpdateAt = laserUpdateIntervalSeconds;
        nextYawUpdateAt = yawUpdateIntervalSeconds;

        Unity.Mathematics.Random rand = new Unity.Mathematics.Random(1);
        if (spreadSensorsUpdate) {
            nextForwardUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);
            nextBackwardUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);
            nextLeftUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);
            nextRightUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);
            nextRightDiagonalUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);
            nextLeftDiagonalUpdateAt -= rand.NextFloat(laserUpdateIntervalSeconds);

            nextYawUpdateAt -= rand.NextFloat(yawUpdateIntervalSeconds);
        }
    }

    void OnCollisionExit2D(Collision2D col) {
        GetComponent<Rigidbody2D>().linearVelocity = Vector3.zero;
        GetComponent<Rigidbody2D>().angularVelocity = 0;
    }

    public SensorReading readSensors() {
        return sensorReading;
    }

    private float readYaw() {
        float rotation = GetComponent<Rigidbody2D>().rotation;

        rotation += northYawReading;

        if (!isExact) {
            rotation += yawAbsError.Next();
        }

        return rotation;
    }

    void Update() {
        if (isExact) {
            sensorReading.forward = forward.getDistance();
            sensorReading.backward = backward.getDistance();
            sensorReading.left = left.getDistance();
            sensorReading.right = right.getDistance();
            sensorReading.leftDiagonal = diagonalLeft.getDistance();
            sensorReading.rightDiagonal = diagonalRight.getDistance();
            sensorReading.yaw = readYaw();
        } else {
            time += Time.deltaTime;

            if (time > nextYawUpdateAt) {
                nextYawUpdateAt += yawUpdateIntervalSeconds;
                sensorReading.yaw = readYaw();
            }
            
            if (time > nextForwardUpdateAt) {
                nextForwardUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.forward = forward.getDistance();
            }

            if (time > nextBackwardUpdateAt) {
                nextBackwardUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.backward = backward.getDistance();
            }

            if (time > nextLeftUpdateAt) {
                nextLeftUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.left = left.getDistance();
            }
        
            if (time > nextRightUpdateAt) {
                nextRightUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.right = right.getDistance();
            }

             if (time > nextLeftDiagonalUpdateAt) {
                nextLeftDiagonalUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.leftDiagonal = diagonalLeft.getDistance();
            }
        
            if (time > nextRightDiagonalUpdateAt) {
                nextRightDiagonalUpdateAt += laserUpdateIntervalSeconds;
                sensorReading.rightDiagonal = diagonalRight.getDistance();
            }
        }
    }
}

public class SensorReading {
    public float forward;
    public float backward;
    public float left;
    public float right;
    public float leftDiagonal;
    public float rightDiagonal;

    public float pitch;
    public float roll;
    public float yaw;
}
