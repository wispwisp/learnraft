package node

import (
	"sync"
)

type State int

const (
	LEADER    State = iota
	CANDIDATE State = iota
	FOLOWER   State = iota
)

type NodeState struct {
	mtx   sync.Mutex
	state State
	term  int
	uri   string
}

func NewNodeState(addr, port string) (nodeState *NodeState) {
	nodeState = &NodeState{uri: addr + ":" + port}

	return
}

func (ns *NodeState) GetUri() string {
	ns.mtx.Lock()
	defer ns.mtx.Unlock()
	return ns.uri
}

func (ns *NodeState) SetState(state State) {
	if state < LEADER || state > FOLOWER {
		return
	}

	ns.mtx.Lock()
	defer ns.mtx.Unlock()

	ns.state = state
}

func (ns *NodeState) GetState() State {
	ns.mtx.Lock()
	defer ns.mtx.Unlock()
	return ns.state
}
