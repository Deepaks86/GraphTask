package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit test for Graph methods
func TestAddEdge(t *testing.T) {
	graph := NewGraph()
	graph = graph.AddEdge(1, 2)

	assert.Contains(t, graph.AdjacencyList[1], 2, "Edge from 1 to 2 should be added")
	assert.Contains(t, graph.AdjacencyList[2], 1, "Edge from 2 to 1 should be added")
}

// Functional test for createGraph handler
func TestCreateGraph(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/graph/create", CreateGraph).Methods("POST")

	requestBody, _ := json.Marshal(GraphCreateRequest{
		Edges: [][]int{{1, 2}, {2, 3}},
	})
	req := httptest.NewRequest("POST", "/graph/create", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201 Created")

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err, "Failed to decode response body")
	assert.NotEmpty(t, response["id"], "Response should contain a graph ID")
}

// Performance test for getShortestPath handler
func BenchmarkGetShortestPath(b *testing.B) {
	router := mux.NewRouter()
	router.HandleFunc("/graph/shortest-path", GetShortestPath).Methods("POST")

	// Prepare a sample graph and request
	graph := NewGraph()
	for i := 0; i < 1000; i++ {
		graph = graph.AddEdge(i, (i+1)%1000)
	}
	graphID := "1"
	graphs[graphID] = graph

	requestBody, _ := json.Marshal(ShortestPathRequest{
		ID:    graphID,
		Start: 0,
		End:   999,
	})

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/graph/shortest-path", bytes.NewBuffer(requestBody))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
