package node

type Vote struct {
	NodeName string `json:"node_name"`
}

type VoteResponse struct {
	NewLeader string `json:"new_leader"`
}
