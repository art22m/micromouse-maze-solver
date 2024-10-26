using Unity.Mathematics;
using UnityEngine;
using System;
using System.Collections.Generic;
using System.Linq;

public class MapBuilder : MonoBehaviour
{
    public GameObject square;
    public Camera sceneCamera;

    public int seed;
    public int loops;
    public bool removeDanglingPoles;

    public void build(int size, float poleSize, float wallLength) {
        System.Random rand = new System.Random(seed);

        Map map = new Map(size, loops, rand);

        float cellSize = poleSize + wallLength;

        GameObject[,] poles = new GameObject[size + 1, size + 1];

        for (int x = 0; x <= size; x++) {
            for (int y = 0; y <= size; y++) {
                GameObject pole = Instantiate(square, new Vector3(x * cellSize, -y * cellSize), Quaternion.identity, transform);
                pole.transform.localScale = new Vector3(poleSize, poleSize, 1);
                poles[x,y] = pole;
            }
        }

        sceneCamera.transform.position = new Vector3(cellSize * size / 2, -cellSize * size / 2, -10);
        sceneCamera.orthographicSize = cellSize * size / 2 + 20;

        for (int x = 0; x < size; x++) {
            for (int y = 0; y < size; y++) {
                Cell cell = map.cells[x,y];
                if (!cell.conUp) {
                    GameObject wall = Instantiate(square, new Vector3(x * cellSize + cellSize / 2, -(y * cellSize)), quaternion.identity, transform);
                    wall.transform.localScale = new Vector3(wallLength, poleSize);
                } 
                if (!cell.conLeft) {
                    GameObject wall = Instantiate(square, new Vector3(x * cellSize, -(y * cellSize + cellSize / 2)), quaternion.identity, transform);
                    wall.transform.localScale = new Vector3(poleSize, wallLength);
                } 
            }
        }

        for (int x = 0; x < size; x++) {
            int y = size - 1;
            Cell cell = map.cells[x,y];
            if (!cell.conDown) {
                GameObject wall = Instantiate(square, new Vector3(x * cellSize + cellSize / 2, -((y + 1) * cellSize)), quaternion.identity, transform);
                wall.transform.localScale = new Vector3(wallLength, poleSize);
            } 
        } 

        for (int y = 0; y < size; y++) {
            int x = size - 1;
            Cell cell = map.cells[x,y];
            if (!cell.conRight) {
                GameObject wall = Instantiate(square, new Vector3((x + 1) * cellSize, -(y * cellSize + cellSize / 2)), quaternion.identity, transform); 
                wall.transform.localScale = new Vector3(poleSize, wallLength);
            } 
        } 

        if (removeDanglingPoles) {
            for (int x = 0; x < size - 1; x++) {
                for (int y = 0; y < size - 1; y++) {
                    Cell cell = map.cells[x,y];
                    if (cell.conRight && cell.conDown && cell.down.conRight && cell.right.conDown) {
                        poles[x + 1, y + 1].SetActive(false);
                    }
                }
            }
        }
    }
}

class Map {
    public Cell[,] cells;

    public Map(int size, int loops, System.Random rand) {
        cells = new Cell[size, size];
        for (int x = 0; x < size; x++) {
            for (int y = 0; y < size; y++) {
                cells[x,y] = new Cell(false, false, false, false);
            }
        }

        for (int x = 0; x < size; x++) {
            for (int y = 0; y < size; y++) {
                Cell cell = cells[x, y];
                if (y > 0) cell.up = cells[x, y - 1];
                if (y < size - 1) cell.down = cells[x, y + 1];
                if (x > 0) cell.left = cells[x - 1, y];
                if (x < size - 1) cell.right = cells[x + 1, y];
            }
        }   

        BuildNoLoop(size, rand);

        while (loops > 0) {
            int x = rand.Next(size);
            int y = rand.Next(size);
        
            Cell cell = cells[x, y];
        
            int dir = rand.Next(4);
            if (dir == 0) {
                if (cell.left != null && !cell.conLeft) {
                    cell.conLeft = true;
                    cell.left.conRight = true;
                    loops--;
                }
            }
            if (dir == 1) {
                if (cell.right != null && !cell.conRight) {
                    cell.conRight = true;
                    cell.right.conLeft = true;
                    loops--;
                }
            }
            if (dir == 2) {
                if (cell.up != null && !cell.conUp) {
                    cell.conUp = true;
                    cell.up.conDown = true;
                    loops--;
                }
            }
            if (dir == 3) {
                if (cell.down != null && !cell.conDown) {
                    cell.conDown = true;
                    cell.down.conUp = true;
                    loops--;
                }
            }
        }
    }

    private void BuildNoLoop(int size, System.Random rand) {
        CarvePassagesFrom(cells[0, 0], rand);
    }

    // Generated by gpt-o1-preview
    private void CarvePassagesFrom(Cell currentCell, System.Random rand)
    {
        currentCell.visited = true;

        var neighbors = new List<Tuple<Cell, string>>();

        if (currentCell.up != null && !currentCell.up.visited)
            neighbors.Add(Tuple.Create(currentCell.up, "up"));
        if (currentCell.down != null && !currentCell.down.visited)
            neighbors.Add(Tuple.Create(currentCell.down, "down"));
        if (currentCell.left != null && !currentCell.left.visited)
            neighbors.Add(Tuple.Create(currentCell.left, "left"));
        if (currentCell.right != null && !currentCell.right.visited)
            neighbors.Add(Tuple.Create(currentCell.right, "right"));

        neighbors = neighbors.OrderBy(n => rand.Next()).ToList();

        foreach (var neighborTuple in neighbors)
        {
            Cell neighborCell = neighborTuple.Item1;
            string direction = neighborTuple.Item2;

            if (!neighborCell.visited)
            {
                if (direction == "up")
                {
                    currentCell.conUp = true;
                    neighborCell.conDown = true;
                }
                else if (direction == "down")
                {
                    currentCell.conDown = true;
                    neighborCell.conUp = true;
                }
                else if (direction == "left")
                {
                    currentCell.conLeft = true;
                    neighborCell.conRight = true;
                }
                else if (direction == "right")
                {
                    currentCell.conRight = true;
                    neighborCell.conLeft = true;
                }
                CarvePassagesFrom(neighborCell, rand);
            }
        }
    }
}

class Cell {
    public Cell up, down, left, right;
    public bool conUp, conDown, conLeft, conRight;

    // used by CarvePassagesFrom
    public bool visited = false;

    public Cell(bool conUp, bool conDown, bool conLeft, bool conRight) {
        this.conUp = conUp;
        this.conRight = conRight;
        this.conDown = conDown;
        this.conLeft = conLeft;
    }   
}
