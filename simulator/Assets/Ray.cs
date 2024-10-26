using UnityEngine;

public class Ray : MonoBehaviour
{
    public float distance;

    void FixedUpdate()
    {
        RaycastHit2D hit = Physics2D.Raycast(transform.position, transform.up, 100000, LayerMask.GetMask("Wall"));
        Debug.DrawRay(transform.position, transform.up * hit.distance, Color.green);
        distance = hit.distance;
    }
}
