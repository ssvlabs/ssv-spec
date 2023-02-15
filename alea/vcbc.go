package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) StartVCBC(priority Priority) error {

	author := i.State.Share.OperatorID
	proposals := i.State.VCBCState.GetM(author, priority)

	// create VCBCSend message and broadcasts
	msgToBroadcast, err := CreateVCBCSend(i.State, i.config, proposals, priority, author)
	if err != nil {
		return errors.Wrap(err, "StartVCBC: failed to create VCBCSend message")
	}
	if i.verbose {
		fmt.Println("\tbroadcasting VCBCSend")
	}
	i.Broadcast(msgToBroadcast)

	if err = i.AddOwnVCBCReady(proposals, priority); err != nil {
		return errors.Wrap(err, "StartVCBC: could not perform own VCBCReady")
	}
	if i.verbose {
		fmt.Println("\tCreated own VCBCReady")
	}

	return nil
}

func (i *Instance) AddOwnVCBCReady(proposals []*ProposalData, priorioty Priority) error {

	hash, err := GetProposalsHash(proposals)
	if err != nil {
		return errors.Wrap(err, "AddOwnVCBCReady: could not get hash of proposals")
	}
	// create VCBCReady message with proof
	vcbcReadyMsg, err := CreateVCBCReady(i.State, i.config, hash, priorioty, i.State.Share.OperatorID)
	if err != nil {
		return errors.Wrap(err, "AddOwnVCBCReady: failed to create VCBCReady message with proof")
	}
	i.uponVCBCReady(vcbcReadyMsg)
	return nil
}

func (i *Instance) AddVCBCOutput(proposals []*ProposalData, priority Priority, author types.OperatorID) {

	// initializes queue of the author if it doesn't exist
	if _, exists := i.State.VCBCState.Queues[author]; !exists {
		i.State.VCBCState.Queues[author] = NewVCBCQueue()
	}

	// gets the sender's associated queue
	queue := i.State.VCBCState.Queues[author]

	// check if it was already delivered
	if i.State.Delivered.HasProposalList(proposals) {
		return
	}

	// check if queue alreasy has proposals and priority
	if queue.HasPriority(priority) {
		return
	}

	// store proposals and priorioty value
	queue.Enqueue(proposals, priority)
}
