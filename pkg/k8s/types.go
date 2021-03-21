package k8s

type TargetPod struct {
	Addr  string
	Label string
}

func (targetPod *TargetPod) Equals(other *TargetPod) bool {
	isAddrEquals := targetPod.Addr == other.Addr
	isLabelEquals := targetPod.Label == other.Label
	return isAddrEquals && isLabelEquals
}