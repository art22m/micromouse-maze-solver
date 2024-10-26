using UnityEngine;

public class Mouse : MonoBehaviour
{
    public Ray forward;
    public Ray backward;
    public Ray right;
    public Ray left;
    public Ray diagonalRight;
    public Ray diagonalLeft;

    void OnCollisionExit2D(Collision2D col) {
        GetComponent<Rigidbody2D>().linearVelocity = Vector3.zero;
        GetComponent<Rigidbody2D>().angularVelocity = 0;
    }
}
