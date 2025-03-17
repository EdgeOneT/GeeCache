package geecache

// 根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	// 用于从对应group查找缓存值
	Get(group string, key string) ([]byte, error)
}
