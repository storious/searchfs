package cluster

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegistryHeartbeatRegistersNewNode(t *testing.T) {
	registry := NewRegistry()

	now := time.Now()

	err := registry.Heartbeat(Heartbeat{
		ID:       "node-1",
		Addr:     "http://localhost:9001",
		Capacity: 1024,
		Used:     128,
		Time:     now,
	})
	require.NoError(t, err)

	info, ok := registry.Get("node-1")
	require.True(t, ok)
	require.Equal(t, DataNodeID("node-1"), info.ID)
	require.Equal(t, "http://localhost:9001", info.Addr)
	require.Equal(t, uint64(1024), info.Capacity)
	require.Equal(t, uint64(128), info.Used)
	require.Equal(t, now, info.LastSeen)
	require.Equal(t, NodeAlive, info.State)
}

func TestRegistryHeartbeatUpdatesExistingNode(t *testing.T) {
	registry := NewRegistry()

	require.NoError(t, registry.Register(DataNodeInfo{
		ID:       "node-1",
		Addr:     "http://localhost:9001",
		Capacity: 1024,
		Used:     100,
		State:    NodeDead,
	}))

	now := time.Now()

	err := registry.Heartbeat(Heartbeat{
		ID:       "node-1",
		Addr:     "http://localhost:9001",
		Capacity: 2048,
		Used:     256,
		Time:     now,
	})
	require.NoError(t, err)

	info, ok := registry.Get("node-1")
	require.True(t, ok)
	require.Equal(t, uint64(2048), info.Capacity)
	require.Equal(t, uint64(256), info.Used)
	require.Equal(t, now, info.LastSeen)
	require.Equal(t, NodeAlive, info.State)
}

func TestRegistryHeartbeatRejectsInvalidInput(t *testing.T) {
	registry := NewRegistry()

	err := registry.Heartbeat(Heartbeat{
		Addr: "http://localhost:9001",
	})
	require.ErrorIs(t, err, ErrEmptyDataNodeID)

	err = registry.Heartbeat(Heartbeat{
		ID: "node-1",
	})
	require.ErrorIs(t, err, ErrEmptyDataNodeAddr)
}
