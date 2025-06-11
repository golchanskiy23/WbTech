package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type PType = string

const size = 5

var (
	WALKING    = PType("WALKING")
	ROAD       = PType("ROAD")
	PUBLIC     = PType("PUBLIC")
	pointTypes = []PType{WALKING, ROAD, PUBLIC}
)

type Point struct {
	x, y      int
	pointType PType
}

type RouteStrategy interface {
	buildRoad(Point, Point, [][]Point) ([]Point, error)
}

func bfs(start, end Point, graph [][]Point, allowed PType) ([]Point, error) {

	type pair struct {
		pt   Point
		path []Point
	}

	visited := make([][]bool, size)
	for i := range visited {
		visited[i] = make([]bool, size)
	}

	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}

	queue := []pair{{start, []Point{start}}}
	visited[start.x][start.y] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.pt.x == end.x && curr.pt.y == end.y {
			return curr.path, nil
		}

		for i := 0; i < 4; i++ {
			nx := curr.pt.x + dx[i]
			ny := curr.pt.y + dy[i]

			if nx >= 0 && ny >= 0 && nx < size && ny < size {
				neighbor := graph[nx][ny]
				if !visited[nx][ny] && neighbor.pointType == allowed {
					visited[nx][ny] = true
					newPath := append([]Point{}, curr.path...)
					newPath = append(newPath, neighbor)
					queue = append(queue, pair{neighbor, newPath})
				}
			}
		}
	}
	return nil, fmt.Errorf("no path found")
}

type RoadStrategy struct {
}

func (br RoadStrategy) buildRoad(a, b Point, graph [][]Point) ([]Point, error) {
	if a.pointType != b.pointType {
		return nil, errors.New("Point types are not equal\n")
	}
	if b.pointType != ROAD {
		return nil, errors.New("invalid point types are not equal\n")
	}
	return bfs(a, b, graph, ROAD)
}

type WalkingStrategy struct {
}

func (br WalkingStrategy) buildRoad(a, b Point, graph [][]Point) ([]Point, error) {
	if a.pointType != b.pointType {
		return nil, errors.New("Point types are not equal\n")
	}
	if b.pointType != WALKING {
		return nil, errors.New("invalid point types are not equal\n")
	}
	return bfs(a, b, graph, WALKING)
}

type PublicTransportStrategy struct {
}

func (br PublicTransportStrategy) buildRoad(a, b Point, graph [][]Point) ([]Point, error) {
	if a.pointType != b.pointType {
		return nil, errors.New("Point types are not equal\n")
	}
	if b.pointType != PUBLIC {
		return nil, errors.New("invalid point types are not equal\n")
	}
	return bfs(a, b, graph, PUBLIC)
}

type Navigator struct {
	buildStrategy RouteStrategy
	graph         [][]Point
}

func generateGraph() [][]Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	graph := make([][]Point, size)
	for i := range graph {
		graph[i] = make([]Point, size)
		for j := range graph[i] {
			ptype := pointTypes[r.Intn(len(pointTypes))]
			graph[i][j] = Point{i, j, ptype}
		}
	}
	return graph
}

func main() {
	graph := generateGraph()
	navigator := Navigator{buildStrategy: WalkingStrategy{}}

	src := graph[0][0]
	destination := graph[size-1][size-1]
	src.pointType = WALKING
	destination.pointType = WALKING
	graph[0][0] = src
	graph[size-1][size-1] = destination

	ans, err := navigator.buildStrategy.buildRoad(src, destination, graph)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(ans)-1; i++ {
		fmt.Printf("[%d %d] ->", ans[i].x, ans[i].y)
	}
	fmt.Printf("[%d %d]\n", ans[len(ans)-1].x, ans[len(ans)-1].y)
}
