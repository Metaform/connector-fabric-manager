package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelTopologicalSort_EmptyGraph(t *testing.T) {
	g := NewGraph[string]()

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in empty graph")
	assert.Empty(t, sorted.SortedOrder, "Expected empty sorted order")
	assert.Empty(t, sorted.ParallelLevels, "Expected no parallel levels")
}

func TestParallelTopologicalSort_SingleVertex(t *testing.T) {
	g := NewGraph[string]()
	value := "A"
	g.AddVertex("A", &value)

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in single vertex graph")
	assert.Len(t, sorted.SortedOrder, 1, "Expected 1 vertex in sorted order")
	assert.Equal(t, "A", sorted.SortedOrder[0].ID, "Expected vertex A")
	assert.Len(t, sorted.ParallelLevels, 1, "Expected 1 parallel level")
	assert.Equal(t, 0, sorted.ParallelLevels[0].Level, "Expected level 0")
	assert.Len(t, sorted.ParallelLevels[0].Vertices, 1, "Expected 1 vertex in level 0")
}

func TestParallelTopologicalSort_LinearChain(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D"}

	// Create vertices A -> B -> C -> D
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "D")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in linear chain")

	// Check sorted order
	expectedOrder := []string{"A", "B", "C", "D"}
	assert.Len(t, sorted.SortedOrder, len(expectedOrder), "Expected correct number of vertices in sorted order")
	for i, vertex := range sorted.SortedOrder {
		assert.Equal(t, expectedOrder[i], vertex.ID, "Expected vertex %s at position %d", expectedOrder[i], i)
	}

	// Check parallel levels - should be 4 levels with 1 vertex each
	assert.Len(t, sorted.ParallelLevels, 4, "Expected 4 parallel levels")
	for i, level := range sorted.ParallelLevels {
		assert.Equal(t, i, level.Level, "Expected level %d", i)
		assert.Len(t, level.Vertices, 1, "Expected 1 vertex in level %d", i)
		assert.Equal(t, expectedOrder[i], level.Vertices[0].ID, "Expected vertex %s in level %d", expectedOrder[i], i)
	}
}

func TestParallelTopologicalSort_ParallelVertices(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D"}

	// Create vertices where A, B, C can execute in parallel
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "D")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in parallel vertices graph")

	// Check that we have 2 levels
	assert.Len(t, sorted.ParallelLevels, 2, "Expected 2 parallel levels")

	// Level 0 should have A, B, C (sorted alphabetically)
	level0 := sorted.ParallelLevels[0]
	assert.Equal(t, 0, level0.Level, "Expected level 0")
	assert.Len(t, level0.Vertices, 3, "Expected 3 vertices in level 0")
	expectedLevel0 := []string{"A", "B", "C"}
	for i, vertex := range level0.Vertices {
		assert.Equal(t, expectedLevel0[i], vertex.ID, "Expected vertex %s at position %d in level 0", expectedLevel0[i], i)
	}

	// Level 1 should have D
	level1 := sorted.ParallelLevels[1]
	assert.Equal(t, 1, level1.Level, "Expected level 1")
	assert.Len(t, level1.Vertices, 1, "Expected 1 vertex in level 1")
	assert.Equal(t, "D", level1.Vertices[0].ID, "Expected vertex D in level 1")
}

func TestParallelTopologicalSort_ComplexDAG(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D", "E", "F"}

	// Create a complex DAG:
	// A -> C -> E
	// B -> C -> F
	// B -> D -> F
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "E")
	g.AddEdge("C", "F")
	g.AddEdge("D", "F")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in complex DAG")

	// Check that we have 4 levels
	assert.Len(t, sorted.ParallelLevels, 3, "Expected 3 parallel levels")

	// Level 0: A, B (no dependencies)
	level0 := sorted.ParallelLevels[0]
	assert.Len(t, level0.Vertices, 2, "Expected 2 vertices in level 0")
	level0IDs := []string{level0.Vertices[0].ID, level0.Vertices[1].ID}
	expectedLevel0 := []string{"A", "B"}
	assert.ElementsMatch(t, expectedLevel0, level0IDs, "Expected level 0 to have correct vertices")

	// Level 1: C, D (depend on A and B respectively)
	level1 := sorted.ParallelLevels[1]
	assert.Len(t, level1.Vertices, 2, "Expected 2 vertices in level 1")
	level1IDs := []string{level1.Vertices[0].ID, level1.Vertices[1].ID}
	expectedLevel1 := []string{"C", "D"}
	assert.ElementsMatch(t, expectedLevel1, level1IDs, "Expected level 1 to have correct vertices")

	// Level 2: E (depends on C), F (depends on C and D)
	level2 := sorted.ParallelLevels[2]
	assert.Len(t, level2.Vertices, 2, "Expected 2 vertex in level 2")
	assert.Equal(t, "E", level2.Vertices[0].ID, "Expected vertex E in level 2")
}

func TestParallelTopologicalSort_CyclicGraph(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C"}

	// Create a cycle: A -> B -> C -> A
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")

	sorted := g.ParallelTopologicalSort()

	assert.True(t, sorted.HasCycle, "Expected cycle detection in cyclic graph")
	assert.Empty(t, sorted.SortedOrder, "Expected empty sorted order for cyclic graph")
	assert.Empty(t, sorted.ParallelLevels, "Expected no parallel levels for cyclic graph")
}

func TestParallelTopologicalSort_SelfLoop(t *testing.T) {
	g := NewGraph[string]()
	value := "A"
	g.AddVertex("A", &value)
	g.AddEdge("A", "A") // Self loop

	sorted := g.ParallelTopologicalSort()

	assert.True(t, sorted.HasCycle, "Expected cycle detection for self loop")
}

func TestParallelTopologicalSort_DiamondPattern(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D"}

	// Create diamond pattern: A -> B,C -> D
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in diamond pattern")

	// Should have 3 levels
	assert.Len(t, sorted.ParallelLevels, 3, "Expected 3 parallel levels")

	// Level 0: A
	assert.Len(t, sorted.ParallelLevels[0].Vertices, 1, "Expected 1 vertex in level 0")
	assert.Equal(t, "A", sorted.ParallelLevels[0].Vertices[0].ID, "Expected vertex A in level 0")

	// Level 1: B, C
	level1 := sorted.ParallelLevels[1]
	assert.Len(t, level1.Vertices, 2, "Expected 2 vertices in level 1")
	level1IDs := []string{level1.Vertices[0].ID, level1.Vertices[1].ID}
	expectedLevel1 := []string{"B", "C"}
	assert.ElementsMatch(t, expectedLevel1, level1IDs, "Expected level 1 to have correct vertices")

	// Level 2: D
	assert.Len(t, sorted.ParallelLevels[2].Vertices, 1, "Expected 1 vertex in level 2")
	assert.Equal(t, "D", sorted.ParallelLevels[2].Vertices[0].ID, "Expected vertex D in level 2")
}

func TestParallelTopologicalSort_DisconnectedComponents(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D"}

	// Create two disconnected components: A -> B and C -> D
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "B")
	g.AddEdge("C", "D")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in disconnected components")

	// Should have 2 levels
	assert.Len(t, sorted.ParallelLevels, 2, "Expected 2 parallel levels")

	// Level 0: A, C (no dependencies)
	level0 := sorted.ParallelLevels[0]
	assert.Len(t, level0.Vertices, 2, "Expected 2 vertices in level 0")
	level0IDs := []string{level0.Vertices[0].ID, level0.Vertices[1].ID}
	expectedLevel0 := []string{"A", "C"}
	assert.ElementsMatch(t, expectedLevel0, level0IDs, "Expected level 0 to have correct vertices")

	// Level 1: B, D
	level1 := sorted.ParallelLevels[1]
	assert.Len(t, level1.Vertices, 2, "Expected 2 vertices in level 1")
	level1IDs := []string{level1.Vertices[0].ID, level1.Vertices[1].ID}
	expectedLevel1 := []string{"B", "D"}
	assert.ElementsMatch(t, expectedLevel1, level1IDs, "Expected level 1 to have correct vertices")
}

func TestParallelTopologicalSort_VertexValues(t *testing.T) {
	g := NewGraph[int]()

	// Test with integer values
	values := []int{10, 20, 30}
	g.AddVertex("A", &values[0])
	g.AddVertex("B", &values[1])
	g.AddVertex("C", &values[2])
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle")

	// Check that values are preserved
	assert.Equal(t, 10, sorted.SortedOrder[0].Value, "Expected value 10 for vertex A")
	assert.Equal(t, 20, sorted.SortedOrder[1].Value, "Expected value 20 for vertex B")
	assert.Equal(t, 30, sorted.SortedOrder[2].Value, "Expected value 30 for vertex C")
}

func TestParallelTopologicalSort_LargeParallelSet(t *testing.T) {
	g := NewGraph[string]()

	// Create a large set of parallel vertices all pointing to one final vertex
	numParallel := 10
	for i := 0; i < numParallel; i++ {
		id := string(rune('A' + i))
		g.AddVertex(id, &id)
	}
	final := "Z"
	g.AddVertex(final, &final)

	// All parallel vertices point to Z
	for i := 0; i < numParallel; i++ {
		id := string(rune('A' + i))
		g.AddEdge(id, final)
	}

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle in large parallel set")

	// Should have 2 levels
	assert.Len(t, sorted.ParallelLevels, 2, "Expected 2 parallel levels")

	// Level 0 should have all parallel vertices
	assert.Len(t, sorted.ParallelLevels[0].Vertices, numParallel, "Expected %d vertices in level 0", numParallel)

	// Level 1 should have the final vertex
	assert.Len(t, sorted.ParallelLevels[1].Vertices, 1, "Expected 1 vertex in level 1")
	assert.Equal(t, final, sorted.ParallelLevels[1].Vertices[0].ID, "Expected vertex %s in level 1", final)
}

func TestParallelTopologicalSort_LevelConsistency(t *testing.T) {
	g := NewGraph[string]()
	values := []string{"A", "B", "C", "D", "E"}

	// Create: A -> B -> D, A -> C -> D, D -> E
	for _, v := range values {
		g.AddVertex(v, &v)
	}
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")

	sorted := g.ParallelTopologicalSort()

	assert.False(t, sorted.HasCycle, "Expected no cycle")

	// Verify level numbering is sequential
	for i, level := range sorted.ParallelLevels {
		assert.Equal(t, i, level.Level, "Expected level %d", i)
	}

	// Verify total vertices match
	totalVertices := 0
	for _, level := range sorted.ParallelLevels {
		totalVertices += len(level.Vertices)
	}
	assert.Equal(t, len(values), totalVertices, "Expected total vertices to match across levels")

	// Verify sorted order length matches
	assert.Len(t, sorted.SortedOrder, len(values), "Expected sorted order length to match")
}

func TestParallelTopologicalSort_CycleDetection(t *testing.T) {
	g := NewGraph[string]()

	t.Run("SimpleCycle", func(t *testing.T) {
		// Clear the graph
		g = NewGraph[string]()

		// Create a simple cycle: A -> B -> C -> A
		values := []string{"A", "B", "C"}
		for _, v := range values {
			g.AddVertex(v, &v)
		}
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "A")

		sorted := g.ParallelTopologicalSort()

		// Assert cycle detection
		assert.True(t, sorted.HasCycle, "Expected cycle to be detected")

		// Assert cycle path is not empty
		assert.NotEmpty(t, sorted.CyclePath, "Expected non-empty cycle path")

		// Assert cycle path length (should be 4: start -> A -> B -> C -> A)
		assert.Len(t, sorted.CyclePath, 4, "Expected cycle path length of 4")

		// Assert cycle path starts and ends with the same vertex
		assert.Equal(t, sorted.CyclePath[0], sorted.CyclePath[len(sorted.CyclePath)-1],
			"Expected cycle path to start and end with the same vertex")

		// Assert all vertices in the cycle are present
		cycleVertices := make(map[string]bool)
		for _, vertex := range sorted.CyclePath[:len(sorted.CyclePath)-1] { // Exclude the duplicate end vertex
			cycleVertices[vertex] = true
		}
		assert.Contains(t, cycleVertices, "A", "Expected vertex A in cycle path")
		assert.Contains(t, cycleVertices, "B", "Expected vertex B in cycle path")
		assert.Contains(t, cycleVertices, "C", "Expected vertex C in cycle path")
		assert.Len(t, cycleVertices, 3, "Expected exactly 3 unique vertices in cycle")

		// Assert that consecutive vertices in the path form valid edges
		for i := 0; i < len(sorted.CyclePath)-1; i++ {
			from := sorted.CyclePath[i]
			to := sorted.CyclePath[i+1]
			// Verify edge exists
			fromVertex, exists := g.GetVertex(from)
			assert.True(t, exists, "Expected vertex %s to exist", from)

			hasEdge := false
			for _, edge := range fromVertex.Edges {
				if edge.ID == to {
					hasEdge = true
					break
				}
			}
			assert.True(t, hasEdge, "Expected edge from %s to %s", from, to)
		}
	})

	t.Run("SelfLoop", func(t *testing.T) {
		// Clear the graph
		g = NewGraph[string]()

		// Create a self-loop: A -> A
		value := "A"
		g.AddVertex("A", &value)
		g.AddEdge("A", "A")

		sorted := g.ParallelTopologicalSort()

		// Assert cycle detection
		assert.True(t, sorted.HasCycle, "Expected self-loop cycle to be detected")

		// Assert cycle path
		assert.NotEmpty(t, sorted.CyclePath, "Expected non-empty cycle path for self-loop")
		assert.Len(t, sorted.CyclePath, 2, "Expected cycle path length of 2 for self-loop")
		assert.Equal(t, "A", sorted.CyclePath[0], "Expected cycle to start with A")
		assert.Equal(t, "A", sorted.CyclePath[1], "Expected cycle to end with A")
	})

	t.Run("ComplexCycle", func(t *testing.T) {
		// Clear the graph
		g = NewGraph[string]()

		// Create a more complex cycle: A -> B -> C -> D -> B (cycle is B -> C -> D -> B)
		values := []string{"A", "B", "C", "D"}
		for _, v := range values {
			g.AddVertex(v, &v)
		}
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "D")
		g.AddEdge("D", "B") // Creates cycle B -> C -> D -> B

		sorted := g.ParallelTopologicalSort()

		// Assert cycle detection
		assert.True(t, sorted.HasCycle, "Expected cycle to be detected in complex graph")

		// Assert cycle path
		assert.NotEmpty(t, sorted.CyclePath, "Expected non-empty cycle path")
		assert.Len(t, sorted.CyclePath, 4, "Expected cycle path length of 4")

		// The cycle should be B -> C -> D -> B
		expectedCycleVertices := map[string]bool{"B": true, "C": true, "D": true}
		actualCycleVertices := make(map[string]bool)
		for _, vertex := range sorted.CyclePath[:len(sorted.CyclePath)-1] {
			actualCycleVertices[vertex] = true
		}

		assert.Equal(t, expectedCycleVertices, actualCycleVertices,
			"Expected cycle to contain vertices B, C, D")
	})

}
