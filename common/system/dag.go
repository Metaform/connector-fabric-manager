//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package system

// Vertex represents a node in the graph
type Vertex[T any] struct {
	ID    string
	Value T
	Edges []*Vertex[T]
}

// NewVertex creates a new vertex with the given ID and value
func NewVertex[T any](id string, value *T) *Vertex[T] {
	return &Vertex[T]{
		ID:    id,
		Value: *value,
		Edges: make([]*Vertex[T], 0),
	}
}

// Graph represents a directed graph using vertices
type Graph[T any] struct {
	vertices map[string]*Vertex[T]
}

// NewGraph creates a new empty graph
func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		vertices: make(map[string]*Vertex[T]),
	}
}

// AddVertex adds a vertex with given ID and value to the graph
func (g *Graph[T]) AddVertex(id string, value *T) {
	if _, exists := g.vertices[id]; !exists {
		g.vertices[id] = NewVertex(id, value)
	}
}

// AddEdge adds a directed edge from vertex with ID u to vertex with ID v
func (g *Graph[T]) AddEdge(fromID, toID string) {
	// Get vertices, return if either doesn't exist
	uVertex, uExists := g.vertices[fromID]
	vVertex, vExists := g.vertices[toID]
	if !uExists || !vExists {
		return
	}

	// Check if edge already exists
	for _, edge := range uVertex.Edges {
		if edge.ID == toID {
			return // Edge already exists
		}
	}

	uVertex.Edges = append(uVertex.Edges, vVertex)
}

// TopologicalSort performs topological sorting of the graph
// Returns the sorted order of IDs and whether the graph has a cycle
func (g *Graph[T]) TopologicalSort() ([]*Vertex[T], bool) {
	// Calculate in-degree for all vertices
	inDegree := make(map[string]int)

	// Initialize in-degree for all vertices
	for id := range g.vertices {
		inDegree[id] = 0
	}

	// Calculate in-degrees
	for _, vertex := range g.vertices {
		for _, edge := range vertex.Edges {
			inDegree[edge.ID]++
		}
	}

	// Create a queue and enqueue vertices with in-degree 0
	var queue []*Vertex[T]
	for id, vertex := range g.vertices {
		if inDegree[id] == 0 {
			queue = append(queue, vertex)
		}
	}

	// Process vertices
	var result []*Vertex[T]
	visited := 0
	totalVertices := len(g.vertices)

	for len(queue) > 0 {
		// Dequeue a vertex
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)
		visited++

		// Reduce in-degree of adjacent vertices
		for _, v := range u.Edges {
			inDegree[v.ID]--
			if inDegree[v.ID] == 0 {
				queue = append(queue, v)
			}
		}
	}

	// If visited count doesn't match vertices count, there's a cycle
	hasCycle := visited != totalVertices
	return result, hasCycle
}

// GetVertex returns the vertex with the given ID and a boolean indicating if it exists
func (g *Graph[T]) GetVertex(id string) (*Vertex[T], bool) {
	vertex, exists := g.vertices[id]
	return vertex, exists
}

// GetValue returns the value associated with the given vertex ID and a boolean indicating if it exists
func (g *Graph[T]) GetValue(id string) (T, bool) {
	if vertex, exists := g.vertices[id]; exists {
		return vertex.Value, true
	}
	var zero T
	return zero, false
}
