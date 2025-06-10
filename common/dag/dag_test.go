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
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopologicalSort_EmptyGraph(t *testing.T) {
	g := NewGraph[string]()

	sorted := g.TopologicalSort()

	assert.False(t, sorted.HasCycle, "Empty graph shouldn't have cycles")
	assert.Empty(t, sorted.SortedOrder, "Expected empty result")
	assert.Empty(t, sorted.CyclePath, "Expected empty cycle path")
}

func TestTopologicalSort_SingleVertex(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", stringPtr("Node A"))

	sorted := g.TopologicalSort()

	assert.False(t, sorted.HasCycle, "Single vertex graph shouldn't have cycles")
	assert.Len(t, sorted.SortedOrder, 1, "Expected 1 vertex")
	assert.Equal(t, "A", sorted.SortedOrder[0].ID, "Expected vertex A")
	assert.Empty(t, sorted.CyclePath, "Expected empty cycle path")
}

func TestTopologicalSort_SimpleChain(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", stringPtr("Node A"))
	g.AddVertex("B", stringPtr("Node B"))
	g.AddVertex("C", stringPtr("Node C"))
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	sorted := g.TopologicalSort()

	assert.False(t, sorted.HasCycle, "Simple chain shouldn't have cycles")

	expectedOrder := []string{"A", "B", "C"}
	assert.Len(t, sorted.SortedOrder, len(expectedOrder), "Expected correct number of vertices")

	for i, expected := range expectedOrder {
		assert.Equal(t, expected, sorted.SortedOrder[i].ID, "Expected correct vertex at position %d", i)
	}
	assert.Empty(t, sorted.CyclePath, "Expected empty cycle path")
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
	g.AddVertex("E", stringPtr("Node E"))
	g.AddVertex("D", stringPtr("Node D"))
	g.AddVertex("C", stringPtr("Node C"))
	g.AddVertex("B", stringPtr("Node B"))
	g.AddVertex("A", stringPtr("Node A"))
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")

	sorted := g.TopologicalSort()

	assert.False(t, sorted.HasCycle, "DAG shouldn't have cycles")

	// Verify topological ordering constraints
	visited := make(map[string]bool)
	for _, v := range sorted.SortedOrder {
		visited[v.ID] = true
		for _, edge := range v.Edges {
			assert.False(t, visited[edge.ID], "Invalid topological order: %s appears before its dependency %s", edge.ID, v.ID)
		}
	}

	assert.Len(t, sorted.SortedOrder, 5, "Expected 5 vertices")

	assert.True(t, isBeforeInSlice("A", "B", sorted.SortedOrder), "A should come before B")
	assert.True(t, isBeforeInSlice("A", "C", sorted.SortedOrder), "A should come before C")
	assert.True(t, isBeforeInSlice("B", "D", sorted.SortedOrder), "B should come before D")
	assert.True(t, isBeforeInSlice("C", "D", sorted.SortedOrder), "C should come before D")
	assert.True(t, isBeforeInSlice("D", "E", sorted.SortedOrder), "D should come before E")
	assert.Empty(t, sorted.CyclePath, "Expected empty cycle path")
}

func TestTopologicalSort_CycleDetection(t *testing.T) {
	g := NewGraph[string]()
	g.AddVertex("A", stringPtr("Node A"))
	g.AddVertex("B", stringPtr("Node B"))
	g.AddVertex("C", stringPtr("Node C"))
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A") // creates a cycle

	sorted := g.TopologicalSort()

	assert.True(t, sorted.HasCycle, "Expected cycle to be detected")
	assert.NotEmpty(t, sorted.CyclePath, "Expected cycle path to be provided")
	assert.Empty(t, sorted.SortedOrder, "Cyclic graph should not produce complete topological ordering")
}

func TestPrintTopologicalSort(t *testing.T) {
	tests := []struct {
		name           string
		setupGraph     func() *Graph[string]
		expectedOutput string
	}{
		{
			name: "Simple linear graph",
			setupGraph: func() *Graph[string] {
				g := NewGraph[string]()
				val1, val2, val3 := "Task A", "Task B", "Task C"
				g.AddVertex("A", &val1)
				g.AddVertex("B", &val2)
				g.AddVertex("C", &val3)
				g.AddEdge("A", "B")
				g.AddEdge("B", "C")
				return g
			},
			expectedOutput: `Topological Sort Results:
Linear Order: [&{A Task A [0x`,
		},
		{
			name: "Graph with parallel execution levels",
			setupGraph: func() *Graph[string] {
				g := NewGraph[string]()
				val1, val2, val3, val4 := "Task A", "Task B", "Task C", "Task D"
				g.AddVertex("A", &val1)
				g.AddVertex("B", &val2)
				g.AddVertex("C", &val3)
				g.AddVertex("D", &val4)
				g.AddEdge("A", "C")
				g.AddEdge("B", "D")
				return g
			},
			expectedOutput: `Topological Sort Results:
Linear Order: [&{A Task A [0x`,
		},
		{
			name: "Graph with cycle",
			setupGraph: func() *Graph[string] {
				g := NewGraph[string]()
				val1, val2, val3 := "Task A", "Task B", "Task C"
				g.AddVertex("A", &val1)
				g.AddVertex("B", &val2)
				g.AddVertex("C", &val3)
				g.AddEdge("A", "B")
				g.AddEdge("B", "C")
				g.AddEdge("C", "A") // Creates a cycle
				return g
			},
			expectedOutput: "Error: Graph contains a cycle",
		},
		{
			name: "Empty graph",
			setupGraph: func() *Graph[string] {
				return NewGraph[string]()
			},
			expectedOutput: `Topological Sort Results:
Linear Order: []

Parallel Execution Levels:`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Setup graph and get result
			graph := tt.setupGraph()
			result := graph.ParallelTopologicalSort()

			// Call the function we're testing
			PrintTopologicalSort(result)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify the output contains expected content
			if tt.name == "Graph with cycle" {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, but got %q", tt.expectedOutput, output)
				}
			} else if tt.name == "Empty graph" {
				if !strings.Contains(output, "Linear Order: []") || !strings.Contains(output, "Parallel Execution Levels:") {
					t.Errorf("Expected output to contain empty graph structure, but got %q", output)
				}
			} else if tt.name == "Single vertex graph" {
				if !strings.Contains(output, "Single Task") || !strings.Contains(output, "Level 0") {
					t.Errorf("Expected output to contain single vertex information, but got %q", output)
				}
			} else {
				// For other tests, check that output contains key elements
				if !strings.Contains(output, "Topological Sort Results:") {
					t.Errorf("Expected output to contain 'Topological Sort Results:', but got %q", output)
				}
				if !strings.Contains(output, "Linear Order:") {
					t.Errorf("Expected output to contain 'Linear Order:', but got %q", output)
				}
				if !strings.Contains(output, "Parallel Execution Levels:") {
					t.Errorf("Expected output to contain 'Parallel Execution Levels:', but got %q", output)
				}
			}
		})
	}
}

func TestPrintTopologicalSort_DetailedOutput(t *testing.T) {
	// Test for detailed output format verification
	g := NewGraph[int]()
	val1, val2, val3, val4 := 1, 2, 3, 4
	g.AddVertex("A", &val1)
	g.AddVertex("B", &val2)
	g.AddVertex("C", &val3)
	g.AddVertex("D", &val4)

	// Create a diamond-shaped dependency: A -> B, A -> C, B -> D, C -> D
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := g.ParallelTopologicalSort()
	PrintTopologicalSort(result)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify specific format elements
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have at least 5 lines: title, linear order, empty line, parallel title, level info
	if len(lines) < 5 {
		t.Errorf("Expected at least 5 lines of output, got %d", len(lines))
	}

	// First line should be the title
	if !strings.Contains(lines[0], "Topological Sort Results:") {
		t.Errorf("First line should contain title, got: %s", lines[0])
	}

	// Should contain parallel execution levels
	foundParallelSection := false
	for _, line := range lines {
		if strings.Contains(line, "Parallel Execution Levels:") {
			foundParallelSection = true
			break
		}
	}
	if !foundParallelSection {
		t.Error("Output should contain 'Parallel Execution Levels:' section")
	}

	// Should contain level information
	foundLevelInfo := false
	for _, line := range lines {
		if strings.Contains(line, "Level") && strings.Contains(line, "can execute in parallel") {
			foundLevelInfo = true
			break
		}
	}
	if !foundLevelInfo {
		t.Error("Output should contain level execution information")
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

func stringPtr(s string) *string {
	return &s
}
