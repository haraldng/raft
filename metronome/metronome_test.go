package metronome

import (
	"testing"
)

// TestDistance checks the correctness of the distance function
func TestDistance(t *testing.T) {
	t1 := QuorumTuple{1, 2, 3}
	t2 := QuorumTuple{4, 5, 6}
	expected := 5

	result := distance(t1, t2)
	if result != expected {
		t.Errorf("distance(%v, %v) = %d; want %d", t1, t2, result, expected)
	}
}

// TestMaximizeDistanceOrdering checks if the maximize distance function works as expected
func TestMaximizeDistanceOrdering(t *testing.T) {
	tuples := []QuorumTuple{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	}
	expected := []QuorumTuple{
		{1, 2, 3},
		{10, 11, 12},
		{4, 5, 6},
		{7, 8, 9},
	}

	orderedTuples := maximizeDistanceOrdering(&tuples)
	if !equalQuorumTuples(orderedTuples, expected) {
		t.Errorf("maximizeDistanceOrdering() = %v; want %v", orderedTuples, expected)
	}
}

// TestNewMetronome checks the creation of a Metronome and ensures properties are correct
func TestNewMetronome(t *testing.T) {
	testCases := []int{3, 5, 7, 9, 11}

	for _, numNodes := range testCases {
		quorumSize := numNodes/2 + 1
		var allMetronomes []*Metronome

		for pid := 1; pid <= numNodes; pid++ {
			m := NewMetronome(NodeId(pid), numNodes, quorumSize)
			allMetronomes = append(allMetronomes, m)

			if pid == 1 {
				t.Logf("N=%d: ordering len: %d, critical len: %d", numNodes, len(m.MyOrdering), m.CriticalLen)
			}
			t.Log(m)
		}

		checkCriticalLen(t, allMetronomes)
	}
}

// Helper function to check if the critical lengths are the same and verify quorum assignments
func checkCriticalLen(t *testing.T, allMetronomes []*Metronome) {
	criticalLen := allMetronomes[0].CriticalLen

	// Ensure all metronomes have the same critical length
	for _, m := range allMetronomes {
		if m.CriticalLen != criticalLen {
			t.Errorf("Expected critical length %d, but got %d", criticalLen, m.CriticalLen)
		}
	}

	numNodes := len(allMetronomes)
	allOrderings := make([][]int, 0, numNodes)
	for _, m := range allMetronomes {
		allOrderings = append(allOrderings, m.MyOrdering)
	}

	quorumSize := numNodes/2 + 1
	numOps := len(allOrderings[0])
	h := make(map[int]int, numOps)
	for i := 0; i < numOps; i++ {
		h[i] = 0
	}

	// Check column by column
	for column := 0; column < numOps; column++ {
		for _, ordering := range allOrderings {
			opID := ordering[column]
			h[opID]++
		}

		if column == criticalLen-1 {
			// At critical length, all ops should have been assigned quorumSize times
			for _, count := range h {
				if count != quorumSize {
					t.Errorf("At column %d, expected quorumSize=%d assignments, but got %d", column, quorumSize, count)
				}
			}
		}
	}

	// Ensure all ops were assigned across all nodes
	for _, count := range h {
		if count != numNodes {
			t.Errorf("Expected %d ops assignment across nodes, but got %d", numNodes, count)
		}
	}
}

// Helper function to compare two slices of QuorumTuples
func equalQuorumTuples(a, b []QuorumTuple) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equal(a[i], b[i]) {
			return false
		}
	}
	return true
}
