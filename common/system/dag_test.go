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

import (
	"testing"
)

func TestTopologicalSort_EmptyGraph(t *testing.T) {
	g := NewGraph[string]()

	result, hasCycle := g.TopologicalSort()

	if hasCycle {
		t.Error("Empty graph shouldn't have cycles")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d vertices", len(result))
	}
}

func TestTopologicalSort_SingleVertex(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", StringPtr("Node A"))

	result, hasCycle := g.TopologicalSort()

	if hasCycle {
		t.Error("Single vertex graph shouldn't have cycles")
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 vertex, got %d", len(result))
	}
	if result[0].ID != "A" {
		t.Errorf("Expected vertex A, got %s", result[0].ID)
	}
}

func TestTopologicalSort_SimpleChain(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", StringPtr("Node A"))
	g.AddVertex("B", StringPtr("Node B"))
	g.AddVertex("C", StringPtr("Node C"))
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	result, hasCycle := g.TopologicalSort()

	if hasCycle {
		t.Error("Simple chain shouldn't have cycles")
	}

	expectedOrder := []string{"A", "B", "C"}
	if len(result) != len(expectedOrder) {
		t.Errorf("Expected %d vertices, got %d", len(expectedOrder), len(result))
	}

	for i, expected := range expectedOrder {
		if result[i].ID != expected {
			t.Errorf("At position %d: expected %s, got %s", i, expected, result[i].ID)
		}
	}
}

func TestTopologicalSort_ComplexDAG(t *testing.T) {
	g := NewGraph[string]()
	// Create a more complex DAG:
	//     A
	//    / \
	//   B   C
	//    \ /
	//     D
	//     |
	//     E
	g.AddVertex("E", StringPtr("Node E"))
	g.AddVertex("D", StringPtr("Node D"))
	g.AddVertex("C", StringPtr("Node C"))
	g.AddVertex("B", StringPtr("Node B"))
	g.AddVertex("A", StringPtr("Node A"))
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")

	result, hasCycle := g.TopologicalSort()

	if hasCycle {
		t.Error("DAG shouldn't have cycles")
	}

	// Verify topological ordering constraints
	visited := make(map[string]bool)
	for _, v := range result {
		visited[v.ID] = true
		for _, edge := range v.Edges {
			if visited[edge.ID] {
				t.Errorf("Invalid topological order: %s appears before its dependency %s",
					edge.ID, v.ID)
			}
		}
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 vertices, got %d", len(result))
	}

	if !isBeforeInSlice("A", "B", result) {
		t.Error("A should come before B")
	}
	if !isBeforeInSlice("A", "C", result) {
		t.Error("A should come before C")
	}
	if !isBeforeInSlice("B", "D", result) {
		t.Error("B should come before D")
	}
	if !isBeforeInSlice("C", "D", result) {
		t.Error("C should come before D")
	}
	if !isBeforeInSlice("D", "E", result) {
		t.Error("D should come before E")
	}
}

func TestTopologicalSort_CycleDetection(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", StringPtr("Node A"))
	g.AddVertex("B", StringPtr("Node B"))
	g.AddVertex("C", StringPtr("Node C"))
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A") // creates a cycle

	result, hasCycle := g.TopologicalSort()

	if !hasCycle {
		t.Error("Expected cycle to be detected")
	}

	if len(result) == len(g.vertices) {
		t.Error("Cyclic graph should not produce complete topological ordering")
	}
}

// Checks if one vertex comes before another in the result
func isBeforeInSlice(first, second string, vertices []*Vertex[string]) bool {
	firstIndex := -1
	secondIndex := -1

	for i, v := range vertices {
		if v.ID == first {
			firstIndex = i
		}
		if v.ID == second {
			secondIndex = i
		}
	}

	return firstIndex != -1 && secondIndex != -1 && firstIndex < secondIndex
}

func StringPtr(s string) *string {
	return &s
}
