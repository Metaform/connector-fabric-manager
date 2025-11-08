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

package dag

import (
	"fmt"
	"sort"
)

const (
	unvisited = 0
	visiting  = 1
	visited   = 2
)

// Vertex represents a node in the graph
type Vertex[T any] struct {
	ID    string
	Value T
	Edges []*Vertex[T]
}

func (v *Vertex[T]) String() string {
	return fmt.Sprintf("Vertex{ID: %s, Edges: %d}", v.ID, len(v.Edges))
}

// Graph represents a directed graph using Vertices
type Graph[T any] struct {
	Vertices map[string]*Vertex[T]
}

// ParallelLevel represents a group of vertices that can be executed in parallel
type ParallelLevel[T any] struct {
	Level    int          // The execution level (0 = first to execute)
	Vertices []*Vertex[T] // Vertices that can be executed in parallel at this level
}

// SortResult contains both the sorted order and any detected cycles
type SortResult[T any] struct {
	SortedOrder []*Vertex[T] // Topological sort order
	HasCycle    bool         // Whether the graph contains a cycle
	CyclePath   []string     // The path of vertex IDs that form the cycle (empty if no cycle)
}

// ParallelSortResult contains both the sorted order and parallel execution groups
type ParallelSortResult[T any] struct {
	SortResult[T]
	ParallelLevels []ParallelLevel[T] // Vertices grouped by parallel execution levels
}

// NewVertex creates a new vertex with the given ID and value
func NewVertex[T any](id string, value *T) *Vertex[T] {
	return &Vertex[T]{
		ID:    id,
		Value: *value,
		Edges: make([]*Vertex[T], 0),
	}
}

// NewGraph creates a new empty graph
func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		Vertices: make(map[string]*Vertex[T]),
	}
}

// AddVertex adds a vertex with given ID and value to the graph
func (g *Graph[T]) AddVertex(id string, value *T) {
	if _, exists := g.Vertices[id]; !exists {
		g.Vertices[id] = NewVertex(id, value)
	}
}

// AddEdge adds a directed edge from vertex with ID u to vertex with ID v
func (g *Graph[T]) AddEdge(fromID, toID string) {
	// Get Vertices, return if either doesn't exist
	uVertex, uExists := g.Vertices[fromID]
	vVertex, vExists := g.Vertices[toID]
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

// GetVertex returns the vertex with the given ID and a boolean indicating if it exists
func (g *Graph[T]) GetVertex(id string) (*Vertex[T], bool) {
	vertex, exists := g.Vertices[id]
	return vertex, exists
}

// TopologicalSort performs topological sorting of the graph
// Returns a SortResult containing the sorted order and cycle information
func (g *Graph[T]) TopologicalSort() *SortResult[T] {
	result := &SortResult[T]{
		SortedOrder: make([]*Vertex[T], 0),
		HasCycle:    false,
		CyclePath:   []string{},
	}

	// First check for cycles
	hasCycle, cyclePath := g.detectCycleWithPath()
	result.HasCycle = hasCycle
	result.CyclePath = cyclePath

	if hasCycle {
		return result
	}

	// Calculate in-degree for all Vertices
	inDegree := make(map[string]int)

	// Initialize in-degree for all Vertices
	for id := range g.Vertices {
		inDegree[id] = 0
	}

	// Calculate in-degrees
	for _, vertex := range g.Vertices {
		for _, edge := range vertex.Edges {
			inDegree[edge.ID]++
		}
	}

	// Create a queue and enqueue Vertices with in-degree 0
	var queue []*Vertex[T]
	for id, vertex := range g.Vertices {
		if inDegree[id] == 0 {
			queue = append(queue, vertex)
		}
	}

	// Process Vertices
	visited := 0
	totalVertices := len(g.Vertices)

	for len(queue) > 0 {
		// Dequeue a vertex
		u := queue[0]
		queue = queue[1:]
		result.SortedOrder = append(result.SortedOrder, u)
		visited++

		// Reduce in-degree of adjacent Vertices
		for _, v := range u.Edges {
			inDegree[v.ID]--
			if inDegree[v.ID] == 0 {
				queue = append(queue, v)
			}
		}
	}

	// If visited count doesn't match Vertices count, there's a cycle
	if visited != totalVertices {
		result.HasCycle = true
		// Clear sorted order since the result is invalid
		result.SortedOrder = make([]*Vertex[T], 0)
	}

	return result
}

// ParallelTopologicalSort performs topological sorting and groups vertices into parallel execution levels.
// Returns a result containing the sorted order, parallel levels, and cycle information.
func (g *Graph[T]) ParallelTopologicalSort() *ParallelSortResult[T] {
	result := &ParallelSortResult[T]{
		SortResult: SortResult[T]{
			SortedOrder: make([]*Vertex[T], 0),
			HasCycle:    false,
			CyclePath:   []string{},
		},
		ParallelLevels: make([]ParallelLevel[T], 0),
	}

	// First detect cycles with path information
	hasCycle, cyclePath := g.detectCycleWithPath()
	result.HasCycle = hasCycle
	result.CyclePath = cyclePath

	if hasCycle {
		return result
	}

	// Get the topological sort result
	sortResult := g.TopologicalSort()
	result.SortResult = *sortResult

	if sortResult.HasCycle {
		return result
	}

	// Calculate in-degrees for parallel level grouping
	inDegree := make(map[string]int)

	// Initialize in-degrees
	for _, vertex := range sortResult.SortedOrder {
		inDegree[vertex.ID] = 0
	}

	// Count incoming edges (dependencies)
	for _, vertex := range sortResult.SortedOrder {
		for _, edge := range vertex.Edges {
			inDegree[edge.ID]++
		}
	}

	// Group vertices by parallel execution levels
	remaining := make(map[string]*Vertex[T])
	for _, vertex := range sortResult.SortedOrder {
		remaining[vertex.ID] = vertex
	}

	levelIndex := 0

	for len(remaining) > 0 {
		// Find vertices with no dependencies at this level
		currentLevel := make([]*Vertex[T], 0)
		for vertexID, vertex := range remaining {
			if inDegree[vertexID] == 0 {
				currentLevel = append(currentLevel, vertex)
			}
		}

		if len(currentLevel) == 0 {
			// This shouldn't happen if there's no cycle, but safety check
			result.HasCycle = true
			return result
		}

		// Sort current level by ID for consistent output
		sort.Slice(currentLevel, func(i, j int) bool {
			return currentLevel[i].ID < currentLevel[j].ID
		})

		// Add current level to parallel levels
		result.ParallelLevels = append(result.ParallelLevels, ParallelLevel[T]{
			Level:    levelIndex,
			Vertices: make([]*Vertex[T], len(currentLevel)),
		})
		copy(result.ParallelLevels[levelIndex].Vertices, currentLevel)

		// Remove processed vertices and update in-degrees
		for _, vertex := range currentLevel {
			delete(remaining, vertex.ID)

			// Reduce in-degree of dependent vertices
			for _, dependent := range vertex.Edges {
				if _, exists := remaining[dependent.ID]; exists {
					inDegree[dependent.ID]--
				}
			}
		}

		levelIndex++
	}
	return result
}

// GetValue returns the value associated with the given vertex ID and a boolean indicating if it exists
func (g *Graph[T]) GetValue(id string) (T, bool) {
	if vertex, exists := g.Vertices[id]; exists {
		return vertex.Value, true
	}
	var zero T
	return zero, false
}

// GetDependents returns the dependents of a given vertex
func (g *Graph[T]) GetDependents(vertex Vertex[T]) []string {
	dependents := make([]string, 0, len(vertex.Edges))
	for _, edge := range vertex.Edges {
		dependents = append(dependents, edge.ID)
	}

	sort.Strings(dependents)
	return dependents
}

// GetDependencies returns the dependencies of a given vertex
func (g *Graph[T]) GetDependencies(vertexID string) []string {
	dependencies := make([]string, 0)

	// Since we reversed the edge direction, we need to find vertices that point to this one
	sortResult := g.TopologicalSort()
	for _, v := range sortResult.SortedOrder {
		for _, edge := range v.Edges {
			if edge.ID == vertexID {
				dependencies = append(dependencies, v.ID)
			}
		}
	}

	sort.Strings(dependencies)
	return dependencies
}

// detectCycleWithPath performs cycle detection using DFS and returns the cycle path if found
func (g *Graph[T]) detectCycleWithPath() (bool, []string) {
	// States of vertices unvisited (0), visiting (1), visited (2)
	visitState := make(map[string]int)
	parent := make(map[string]string)

	// Initialize all vertices as unvisited
	for id := range g.Vertices {
		visitState[id] = unvisited
	}

	// Try DFS from each unvisited vertex
	for id := range g.Vertices {
		if visitState[id] == unvisited {
			if hasCycle, cyclePath := g.detectCycle(id, visitState, parent); hasCycle {
				return true, cyclePath
			}
		}
	}

	return false, []string{}
}

// detectCycle performs DFS traversal to detect cycles and returns the cycle path if found
func (g *Graph[T]) detectCycle(vertexID string, visitState map[string]int, parent map[string]string) (bool, []string) {
	visitState[vertexID] = visiting // Mark as visiting

	vertex := g.Vertices[vertexID]
	for _, neighbor := range vertex.Edges {
		neighborID := neighbor.ID

		if visitState[neighborID] == visiting {
			// Back-edge found - cycle detected
			// Reconstruct the cycle path
			cycle := []string{neighborID}
			current := vertexID

			// Trace back from current vertex to the start of the cycle
			for current != neighborID {
				cycle = append(cycle, current)
				current = parent[current]
			}

			// Add the starting vertex again to complete the cycle
			cycle = append(cycle, neighborID)

			// Reverse to get the correct order
			for i, j := 0, len(cycle)-1; i < j; i, j = i+1, j-1 {
				cycle[i], cycle[j] = cycle[j], cycle[i]
			}

			return true, cycle
		} else if visitState[neighborID] == unvisited {
			parent[neighborID] = vertexID
			if hasCycle, cyclePath := g.detectCycle(neighborID, visitState, parent); hasCycle {
				return true, cyclePath
			}
		}
	}

	visitState[vertexID] = visited // Mark as visited
	return false, nil
}

// PrintTopologicalSort prints the results of topological sorting
func PrintTopologicalSort[T any](result *ParallelSortResult[T]) {
	if result.HasCycle {
		fmt.Println("Error: Graph contains a cycle")
		return
	}

	fmt.Println("Topological Sort Results:")
	fmt.Printf("Linear Order: %v\n\n", result.SortedOrder)

	fmt.Println("Parallel Execution Levels:")
	for _, level := range result.ParallelLevels {
		fmt.Printf("  Level %d (can execute in parallel): %v\n", level.Level, level.Vertices)
	}
	fmt.Println()
}

// PrintGraph prints the graph structure
func PrintGraph[T any](graph *Graph[T]) {
	fmt.Println("Graph Structure:")
	sortResult := graph.TopologicalSort()
	for _, vertex := range sortResult.SortedOrder {
		dependents := graph.GetDependents(*vertex)
		dependencies := graph.GetDependencies(vertex.ID)

		fmt.Printf("  %s:\n", vertex.ID)
		if len(dependencies) > 0 {
			fmt.Printf("    Dependencies: %v\n", dependencies)
		}
		if len(dependents) > 0 {
			fmt.Printf("    Dependents: %v\n", dependents)
		}
	}
	fmt.Println()
}
