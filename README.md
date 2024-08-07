#Implementation: 
Given an unweighted, undirected graph of N nodes and E edges, a source node S, and a destination node D, we need to find the shortest path from node S to node D in the graph.
  
For this BFS algorithm is used and implemented in Golang.

#Requirements:
1. Post a graph, returning an ID to be used in subsequent operations
2. Get the shortest path between two vertices in a previously posted graph
3. Delete a graph from the server

#Running & Testing
1.	Clone the repository to your local using
git clone https://github.com/Deepaks86/GraphTask.git

Files: 
graphActivity.go : BFS implementation
graphActivity_test.go: Test cases
run the graphActivity.go using 
    go run graphActivity.go

Requirement 1: Create a Graph
Using Terminal: 
curl --location 'http://localhost:9000/graph/create' \
--header 'Content-Type: application/json' \
--data '{
    "edges": 
    [[0, 1],[1, 2], [0, 3], [3, 4],[4, 7], [3,7], [6, 7], [4, 5], [4, 6], [5, 6]]
}'
Using Postman:
URL: http://localhost:9000/graph/create
Method: Post
Request Body(Application/JSON): 
{
    "edges": 
    [[0, 1],[1, 2], [0, 3], [3, 4],[4, 7], [3,7], [6, 7], [4, 5], [4, 6], [5, 6]]
}
Response:
{"id":"1"}

Requirement 2: Get the Shortest Path
Using Terminal:
curl --location 'http://localhost:9000/graph/shortest-path' \
--header 'Content-Type: application/json' \
--data '{
    "id": "1",
    "start": 1,
    "end": 7 
}'

Using Postman:
URL: http://localhost:9000/graph/shortest-path
Method: Post
Request Body(Application Json): 
{
    "id": "1",
    "start": 1,
    "end": 7 
}
Response: {"Shortest path":[1,0,3,7]}

Requirement 3: Delete a Graph
Using Terminal:
curl --location 'http://localhost:9000/graph/delete' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1
}'
Using Postman:
URL: http://localhost:9000/graph/delete
Method: Post
Request Body(Application/JSON):
{
    "id": 1
}
Response: {"deleted Graph ID":[1]}

#Performing the Test Cases:

Unit test & Functional Test: go test -v
Performance Test: go test -bench=.

 








