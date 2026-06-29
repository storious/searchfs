package cluster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeastUsedPlacementAllocatesAliveNodesWithMostFreeSpace(t *testing.T) {
	policy := NewLeastUsedPlacement()

	nodes := []DataNodeInfo{
		{ID: "node-1", Addr: "http://localhost:9001", Capacity: 1000, Used: 900, State: NodeAlive},
		{ID: "node-2", Addr: "http://localhost:9002", Capacity: 1000, Used: 100, State: NodeAlive},
		{ID: "node-3", Addr: "http://localhost:9003", Capacity: 1000, Used: 500, State: NodeAlive},
	}

	selected, err := policy.Allocate(100, 2, nodes)
	require.NoError(t, err)
	require.Len(t, selected, 2)

	require.Equal(t, DataNodeID("node-2"), selected[0].ID)
	require.Equal(t, DataNodeID("node-3"), selected[1].ID)
}

func TestLeastUsedPlacementSkipsNonAliveNodes(t *testing.T) {
	policy := NewLeastUsedPlacement()

	nodes := []DataNodeInfo{
		{ID: "dead", Addr: "http://localhost:9001", Capacity: 1000, Used: 0, State: NodeDead},
		{ID: "suspect", Addr: "http://localhost:9002", Capacity: 1000, Used: 0, State: NodeSuspect},
		{ID: "alive", Addr: "http://localhost:9003", Capacity: 1000, Used: 100, State: NodeAlive},
	}

	selected, err := policy.Allocate(100, 1, nodes)
	require.NoError(t, err)
	require.Len(t, selected, 1)
	require.Equal(t, DataNodeID("alive"), selected[0].ID)
}

func TestLeastUsedPlacementRejectsInsufficientNodes(t *testing.T) {
	policy := NewLeastUsedPlacement()

	nodes := []DataNodeInfo{
		{ID: "node-1", Addr: "http://localhost:9001", Capacity: 1000, Used: 100, State: NodeAlive},
	}

	_, err := policy.Allocate(100, 2, nodes)
	require.ErrorIs(t, err, ErrNotEnoughDataNodes)
}

func TestLeastUsedPlacementRejectsInsufficientCapacity(t *testing.T) {
	policy := NewLeastUsedPlacement()

	nodes := []DataNodeInfo{
		{ID: "node-1", Addr: "http://localhost:9001", Capacity: 1000, Used: 950, State: NodeAlive},
	}

	_, err := policy.Allocate(100, 1, nodes)
	require.ErrorIs(t, err, ErrNoAliveDataNodes)
}
