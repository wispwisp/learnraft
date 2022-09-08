package node

import "testing"

func TestNodeAdd(t *testing.T) {
	var nodesInfo NodesInfo

	nodesInfo.Add(&NodeInfo{Uri: "127.0.0.1:8090"})
	nodesInfo.Add(&NodeInfo{Uri: "127.0.0.1:8091"})

	nodes := nodesInfo.Get()
	expectedLen := 2

	if got := len(nodes); got != expectedLen {
		t.Errorf("Len of nodes = %q, want %q", got, expectedLen)
	}
}

func TestNodeAddNonUnique(t *testing.T) {
	var nodesInfo NodesInfo

	nodesInfo.Add(&NodeInfo{Uri: "127.0.0.1:8090"})
	nodesInfo.Add(&NodeInfo{Uri: "127.0.0.1:8090"})

	nodes := nodesInfo.Get()
	expectedLen := 1

	if got := len(nodes); got != expectedLen {
		t.Errorf("Len of nodes = %q, want %q", got, expectedLen)
	}
}
