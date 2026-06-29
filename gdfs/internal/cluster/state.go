package cluster

import (
	"encoding/json"
	"fmt"
)

func (s NodeState) String() string {
	switch s {
	case NodeAlive:
		return "alive"
	case NodeSuspect:
		return "suspect"
	case NodeDead:
		return "dead"
	default:
		return "unknown"
	}
}

func (s NodeState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *NodeState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value {
	case "unknown":
		*s = NodeUnknown
	case "alive":
		*s = NodeAlive
	case "suspect":
		*s = NodeSuspect
	case "dead":
		*s = NodeDead
	default:
		return fmt.Errorf("unknown node state: %s", value)
	}

	return nil
}
