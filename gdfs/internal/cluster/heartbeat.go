package cluster

import "time"

type Heartbeat struct {
	ID       DataNodeID
	Addr     string
	Capacity uint64
	Used     uint64
	Time     time.Time
}

func (r *Registry) Heartbeat(hb Heartbeat) error {
	if hb.ID == "" {
		return ErrEmptyDataNodeID
	}
	if hb.Addr == "" {
		return ErrEmptyDataNodeAddr
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := hb.Time
	if now.IsZero() {
		now = time.Now()
	}

	info := DataNodeInfo{
		ID:       hb.ID,
		Addr:     hb.Addr,
		Capacity: hb.Capacity,
		Used:     hb.Used,
		LastSeen: now,
		State:    NodeAlive,
	}

	r.nodes[hb.ID] = info
	return nil
}
