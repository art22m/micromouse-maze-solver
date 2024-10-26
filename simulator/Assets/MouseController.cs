using System;
using Unity.Mathematics;
using UnityEngine;

public class MouseController : MonoBehaviour
{
    public GameObject mousePrefab;

    public float mouseWidth;
    public float mouseHeight;

    public float mouseSpeed;
    public float mouseRotationSpeed;

    public GameObject trail;

    public GameObject mouse;
    public Rigidbody2D rb;

    private State state = State.Idle;
    private Action onFinish = null;
    private float moveCycles = 0;
    private float maxMoveCycles = 0;
    private Vector2 positionTarget;
    private float rotationTarget;

    public void InitMouse(float cellSize, int x, int y) {
        mouse = Instantiate(mousePrefab, new Vector3(cellSize * x + cellSize / 2, -(cellSize * y + cellSize / 2)), Quaternion.identity);
        mouse.transform.localScale = new Vector3(mouseWidth, mouseHeight);
        rb = mouse.GetComponent<Rigidbody2D>();
    }

    public bool MoveForward(float distance, Action onFinish) {
        return move(mouse.transform.up * distance, onFinish);
    }

    public bool MoveBackward(float distance, Action onFinish) {
        return move(-mouse.transform.up * distance, onFinish);
    }

    public bool RotateRight(float angle, Action onFinish) {
        return rotate(-angle, onFinish);
    }

    public bool RotateLeft(float angle, Action onFinish) {
        return rotate(angle, onFinish);
    }

    private bool move(Vector2 diff, Action onFinish) {
        if (state != State.Idle) {
            return false;
        }

        GameObject t = Instantiate(trail, rb.position, Quaternion.Euler(0, 0, rb.rotation));
        t.transform.localScale = new Vector3(mouseHeight / 2, mouseHeight / 2, 1);

        float distance = diff.magnitude;

        moveCycles = 0;
        maxMoveCycles = distance / mouseSpeed + 1;
        positionTarget = rb.position + diff;
        state = State.Moving;

        this.onFinish = onFinish;
        return true; 
    }

    private bool rotate(float angleDiff, Action onFinish) {
        if (state != State.Idle) {
            return false;
        }

        moveCycles = 0;
        maxMoveCycles = math.abs(angleDiff) / mouseRotationSpeed + 1;
        rotationTarget = rb.rotation + angleDiff;
        state = State.Rotating;

        this.onFinish = onFinish;
        return true;  
    }

    void FixedUpdate() {
        if (state == State.Idle) {
            if (onFinish != null) {
                onFinish();
                onFinish = null;
            }
            return;
        }

        if (state == State.Moving) {
            Vector2 diff = positionTarget - rb.position;
            if (diff.magnitude < mouseSpeed) {
                rb.MovePosition(positionTarget);
                state = State.Idle;
            } else {
                rb.MovePosition(rb.position + diff.normalized * mouseSpeed);
                moveCycles++;
                if (moveCycles >= maxMoveCycles) {
                    state = State.Idle; 
                }
            }
        }

        if (state == State.Rotating) {
            float diff = rotationTarget - rb.rotation;
            if (math.abs(diff) < mouseRotationSpeed) {
                rb.MoveRotation(rotationTarget);
                state = State.Idle;
            } else {
                if (diff > 0) {
                    rb.MoveRotation(rb.rotation + mouseRotationSpeed);
                } else {
                    rb.MoveRotation(rb.rotation - mouseRotationSpeed);
                }
                moveCycles++;
                if (moveCycles >= maxMoveCycles) {
                    state = State.Idle; 
                }
            }
        }
    }
}

enum State {
    Idle,
    Moving,
    Rotating
}
