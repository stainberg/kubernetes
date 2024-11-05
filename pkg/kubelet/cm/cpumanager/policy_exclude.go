package cpumanager

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/topology"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/utils/cpuset"
)

const (
	PolicyExclude policyName = "exclude"
)

var _ Policy = &excludePolicy{}

type excludePolicy struct {
	topology *topology.CPUTopology
	reserved cpuset.CPUSet
}

func NewExcludePolicy(topology *topology.CPUTopology, reservedCPUs cpuset.CPUSet) (Policy, error) {
	policy := &excludePolicy{
		topology: topology,
		reserved: reservedCPUs,
	}
	klog.InfoS("Reserved CPUs not available for exclusive assignment", "reserved = ", policy.reserved)
	return policy, nil
}

func (p *excludePolicy) Name() string {
	return string(PolicyExclude)
}

func (p *excludePolicy) Start(s state.State) error {
	if err := p.validateState(s); err != nil {
		klog.ErrorS(err, "Exclude policy invalid state, please drain node and remove policy state file")
		return err
	}
	return nil
}

func (p *excludePolicy) validateState(s state.State) error {
	allCPUs := p.topology.CPUDetails.CPUs()
	if p.reserved.Size() == 0 {
		s.SetDefaultCPUSet(allCPUs)
	} else {
		s.SetDefaultCPUSet(allCPUs.Difference(p.reserved))
	}
	return nil
}

func (p *excludePolicy) GetAllocatableCPUs(_ state.State) cpuset.CPUSet {
	return p.topology.CPUDetails.CPUs().Difference(p.reserved)
}

func (p *excludePolicy) Allocate(_ state.State, _ *v1.Pod, _ *v1.Container) error {
	return nil
}

func (p *excludePolicy) RemoveContainer(_ state.State, _ string, _ string) error {
	return nil
}

func (p *excludePolicy) GetTopologyHints(_ state.State, _ *v1.Pod, _ *v1.Container) map[string][]topologymanager.TopologyHint {
	return nil
}

func (p *excludePolicy) GetPodTopologyHints(_ state.State, _ *v1.Pod) map[string][]topologymanager.TopologyHint {
	return nil
}
