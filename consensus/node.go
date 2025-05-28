package consensus

import (
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ID        int
	State     NodeState // Follower/Candidate/Leader
	CurrentTerm int
	Log       []LogEntry
	CommitIndex int
	peers     []*Node
	
	// Channels
	voteCh    chan RequestVote
	appendCh  chan AppendEntries
	crashCh   chan struct{}
}

type LogEntry struct {
	Term    int
	Command string
}

func (n *Node) Run() {
	rand.Seed(time.Now().UnixNano())
	electionTimeout := time.Duration(rand.Intn(150)+150) * time.Millisecond
	
	for {
		select {
		case <-time.After(electionTimeout):
			if n.State == Follower {
				n.startElection()
			}
			
		case ae := <-n.appendCh:
			if ae.Term >= n.CurrentTerm {
				n.handleAppendEntries(ae)
				electionTimeout = resetTimer()
			}
			
		case vote := <-n.voteCh:
			n.handleVoteRequest(vote)
			
		case <-n.crashCh:
			return
		}
	}
}

func (n *Node) startElection() {
	n.State = Candidate
	n.CurrentTerm++
	
	// Request votes from peers
	votes := 1 // self-vote
	for _, peer := range n.peers {
		if peer.sendVoteRequest(n.CurrentTerm) {
			votes++
		}
	}
	
	if votes > len(n.peers)/2 {
		n.becomeLeader()
	}
}

func (n *Node) handleAppendEntries(ae AppendEntries) bool {
	// Log replication logic
	if ae.LeaderCommit > n.CommitIndex {
		n.CommitIndex = min(ae.LeaderCommit, len(n.Log)-1)
	}
	return true
}
