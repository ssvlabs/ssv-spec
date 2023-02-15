package alea

import (
	"reflect"
	"sync"
)

type VCBCQueue struct {
	data     [][]*ProposalData
	priority []Priority
	mutex    sync.Mutex
}

func NewVCBCQueue() *VCBCQueue {
	return &VCBCQueue{
		data:     make([][]*ProposalData, 0),
		priority: make([]Priority, 0),
		mutex:    sync.Mutex{},
	}
}

func (queue *VCBCQueue) Enqueue(proposals []*ProposalData, priority Priority) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.data = append(queue.data, proposals)
	queue.priority = append(queue.priority, priority)
}

func (queue *VCBCQueue) Peek() ([]*ProposalData, Priority) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if len(queue.data) == 0 {
		return nil, Priority(0)
	}

	oldestProposals := queue.data[0]
	oldestPriority := queue.priority[0]

	return oldestProposals, oldestPriority
}

func (queue *VCBCQueue) PeekLast() ([]*ProposalData, Priority) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if len(queue.data) == 0 {
		return nil, Priority(0)
	}

	lastProposals := queue.data[len(queue.data)-1]
	lastPriority := queue.priority[len(queue.priority)-1]

	return lastProposals, lastPriority
}

func (queue *VCBCQueue) Dequeue() ([]*ProposalData, Priority) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if len(queue.data) == 0 {
		return nil, Priority(0)
	}

	lastProposals := queue.data[0]
	queue.data = queue.data[1:]

	lastPriority := queue.priority[0]
	queue.priority = queue.priority[1:]

	return lastProposals, lastPriority
}

func (queue *VCBCQueue) IsEmpty() bool {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	return len(queue.data) == 0
}

func (queue *VCBCQueue) Clear() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.data = nil
	queue.priority = nil
}

func (queue *VCBCQueue) GetValues() [][]*ProposalData {
	return queue.data
}

func (queue *VCBCQueue) GetPriorities() []Priority {
	return queue.priority
}

func (queue *VCBCQueue) HasProposal(proposalInstance *ProposalData) bool {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for _, proposals := range queue.data {
		for _, proposal := range proposals {
			if reflect.DeepEqual(proposal, proposalInstance) {
				return true
			}
		}
	}
	return false
}

func (queue *VCBCQueue) HasProposalList(proposalList []*ProposalData) bool {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for _, proposals := range queue.data {
		if reflect.DeepEqual(proposals, proposalList) {
			return true
		}
	}
	return false
}

func (queue *VCBCQueue) HasPriority(priority Priority) bool {
	for _, p := range queue.priority {
		if p == priority {
			return true
		}
	}
	return false
}
