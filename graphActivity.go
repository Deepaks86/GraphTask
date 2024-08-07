package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Graph struct {
	AdjacencyList map[int][]int
}
type GraphCreateRequest struct {
	Edges [][]int `json:"edges"`
}

type GraphDeleteRequest struct {
	ID int `json:"id"`
}

type ShortestPathRequest struct {
	ID    string `json:"id"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

var (
	graphs = make(map[string]Graph)
	mutex  = &sync.Mutex{}
	nextID = 1
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/graph/create", CreateGraph).Methods("POST")
	router.HandleFunc("/graph/delete", DeleteGraph).Methods("POST")
	router.HandleFunc("/graph/shortest-path", GetShortestPath).Methods("POST")
	fmt.Println("Server is running on :9000")
	if err := http.ListenAndServe(":9000", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func CreateGraph(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var graphCreateRequest GraphCreateRequest
	if err := json.Unmarshal(body, &graphCreateRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	graph := NewGraph()
	for _, edge := range graphCreateRequest.Edges {
		if len(edge) == 2 {
			graph = graph.AddEdge(edge[0], edge[1])
		}
	}
	graphID := strconv.Itoa(nextID)
	graphs[graphID] = graph
	nextID++
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": graphID})
}

func NewGraph() Graph {
	return Graph{AdjacencyList: make(map[int][]int)}
}

func (g Graph) AddEdge(from, to int) Graph {
	g.AdjacencyList[from] = append(g.AdjacencyList[from], to)
	g.AdjacencyList[to] = append(g.AdjacencyList[to], from)
	return g
}

func DeleteGraph(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	var graphDeleteRequest GraphDeleteRequest
	if err := json.Unmarshal(body, &graphDeleteRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := graphs[strconv.Itoa(graphDeleteRequest.ID)]; exists {
		delete(graphs, strconv.Itoa(graphDeleteRequest.ID))
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Graph not found", http.StatusNotFound)
	}
}

func GetShortestPath(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var shortestPathRequest ShortestPathRequest
	if err := json.Unmarshal(body, &shortestPathRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	graph, exists := graphs[shortestPathRequest.ID]
	mutex.Unlock()
	fmt.Println(shortestPathRequest.ID, graph)
	if !exists {
		http.Error(w, "Graph not found", http.StatusNotFound)
		return
	}

	path, err := graph.ShortestPath(shortestPathRequest.Start, shortestPathRequest.End)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]int{"path": path})
	return
}

func (g Graph) ShortestPath(start, end int) ([]int, error) {
	_, startExists := g.AdjacencyList[start]
	_, endExists := g.AdjacencyList[end]
	if !startExists || !endExists {
		return nil, errors.New("one or both vertices do not exist")
	}

	visited := make(map[int]bool)
	queue := [][]int{{start}}
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node == end {
			return path, nil
		}

		for _, neighbor := range g.AdjacencyList[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]int{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}
	return nil, errors.New("no path found")
}
