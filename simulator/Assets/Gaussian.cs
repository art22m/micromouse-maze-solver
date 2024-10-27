using Unity.Mathematics;

public class Gaussian {
    private float mean;
    private float std;
    private Random rand;

    public Gaussian(float mean, float std, Random rand) {
        this.mean = mean;
        this.std = std;
        this.rand = rand;
    }

    public float Next() {
        float x1 = rand.NextFloat(1);
        float x2 = rand.NextFloat(1);
        float y = math.sqrt(-2 * math.log(x1)) * math.cos(2 * math.PI * x2);
        float res = y * std + mean;
        if (math.abs(res - mean) > std * 1.5) {
            return Next();
        }
        return res;
    }
}
