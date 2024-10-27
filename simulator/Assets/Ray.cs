using Unity.Mathematics;
using UnityEditor.SearchService;
using UnityEngine;

public class Ray : MonoBehaviour
{
    public float maxReading;

    public float errorMean;
    public float errorStd;
    public float absErrorMean;
    public float absErrorStd;

    public uint errorSeed;
    public uint absErrorSeed;

    private Gaussian error;
    private Gaussian absError;

    private bool isExact;

    private float distance;

    void Awake() {
        error = new Gaussian(errorMean, errorStd, new Unity.Mathematics.Random(errorSeed));
        absError = new Gaussian(absErrorMean, absErrorStd, new Unity.Mathematics.Random(absErrorSeed));
    }

    void Start() {
        isExact = transform.parent.GetComponent<Mouse>().isExact;
    }

    void Update() {
        RaycastHit2D hit = Physics2D.Raycast(transform.position, transform.up, 100000, LayerMask.GetMask("Wall"));
        Debug.DrawRay(transform.position, transform.up * hit.distance, Color.green);
        distance = math.min(hit.distance, maxReading);
    }

    public float getDistance() {
        if (isExact) {
            return distance;
        }
        return distance + distance * error.Next() + absError.Next();
    }
}
