package k8s

// TargetPod represents the fetched pods from api-server to redirect incoming requests
type TargetPod struct {
	Addr  string
	Label string
}

// Equals method is for checking TargetPod equality
func (targetPod *TargetPod) Equals(other *TargetPod) bool {
	isAddrEquals := targetPod.Addr == other.Addr
	isLabelEquals := targetPod.Label == other.Label
	return isAddrEquals && isLabelEquals
}
