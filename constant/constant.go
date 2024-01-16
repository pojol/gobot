package constant

import (
	"sync"

	"github.com/pojol/braid-go/module/meta"
)

var statemu sync.Mutex
var serverState = meta.EWait

func SetServerState(state int32) {
	statemu.Lock()
	serverState = state
	statemu.Unlock()
}

func GetServerState() int32 {
	statemu.Lock()
	defer statemu.Unlock()

	return serverState
}

var clusterState = false
var clusterMu sync.Mutex

func SetClusterState(state bool) {
	clusterMu.Lock()
	clusterState = state
	clusterMu.Unlock()
}

func GetClusterState() bool {
	clusterMu.Lock()
	defer clusterMu.Unlock()

	return clusterState
}

var nodmap map[string]int = make(map[string]int)

func AddNode(ip string) {
	nodmap[ip] = 1
}

func RmvNode(ip string) {
	delete(nodmap, ip)
}

func GetNods() int {
	return len(nodmap)
}
