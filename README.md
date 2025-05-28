# raft-consensus-go

A minimal Raft-style consensus implementation featuring leader election, log replication, and fault tolerance.

## Features
- Timeout-based leader election
- Log replication with quorum commitment
- Node crash recovery
- Network message dropout handling
- CLI demo with 5-node cluster

### Components:
- **Node States**: Follower, Candidate, Leader
- **Election Timeout**: 150-300ms randomized
- **Heartbeat Interval**: 50ms
- **Log Structure**: Term-indexed entries with LSM-style persistence

## Quick Start
