package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)



func (i *Instance) StartAgreementComponent() error {

	// init round number to 0
	i.State.ACRound = 0
	for {

		// check if it should stop performing agreement
		if i.State.StopAgreement {
			break
		}

		// calculate the round leader (to get value to be decided on)
        leader := i.config.GetProposerF()(i.State,i.State.ACRound)

        // get the local queue associated with the leader's id (create if there isn't one)
		if _, exists := i.State.queues[leader]; !exists {
			i.State.queues[leader] = NewVCBCQueue()
		}
		queue := i.State.queues[leader]

        // get the value of the queue with the lowest priority value
        value, priority := queue.Peek()

		// decide own vote
		vote := byte(0)
        if value == nil {
			vote = byte(0)
        } else {
			vote = byte(1)
		}

		// start ABA protocol
		result := i.StartABA(vote)

		if result == 1 { 
			// if the protocol agreed on the value of the leader replica, deliver it
			
			// if ABA decided 1 but own vote was 0, start recover mechanism to get VCBC messages not received from leader
			if vote == 0 {
				// create FILLGAP message
				fillGapMsg, err := CreateFillGap(i.State, i.config, leader, priority)
				if err != nil {
					errors.Wrap(err,"failed to create FillGap message")
				}
				i.Broadcast(fillGapMsg)
				// wait for the value to be received
				i.WaitFillGapResponse(leader,priority)
			}

			// get decided value
			value, priority = queue.Peek()

			// remove the value from the queue and add it to S
			queue.Dequeue()
			i.State.S.Enqueue(value, priority)
			// return the value to the client
			i.Deliver(value)
		}
		// increment the round number
		i.State.ACRound++
	}
	return nil
}

func (i *Instance) WaitFillGapResponse(leader types.OperatorID, priority Priority) {
	// gets the leader queue
	queue := i.State.queues[leader]
	for {
		// if has the desired priority, returns
		_, localPriority := queue.Peek()
		if localPriority >= priority {
			return
		}

		// waits until a FILLER signal is received (actived on the uponFiller function)
		<- i.State.FillerMsgReceived
	}
}

func (i *Instance) StartABA(vote byte) byte {
	// set ABA's input value
	i.State.ABAState.Vin = vote

	// broadcast INIT message with input vote
	initMsg, err := CreateABAInit(i.State, i.config, vote, i.State.ABAState.Round)
	if err != nil {
		errors.Wrap(err,"failed to create ABA Init message")
	}
	i.Broadcast(initMsg)

	// update sent flag
	if vote == 1 {
		i.State.ABAState.SentInit1 = true
	} else {
		i.State.ABAState.SentInit0 = true
	}

	// wait until channel Terminate receives a signal
	<- i.State.ABAState.Terminate

	// returns the decided value
	return i.State.ABAState.Vdecided
}


func CreateABA(state *State, config IConfig, vote byte, round Round) (*SignedMessage, error) {
	abaData := &ABAData{
		Vote:			vote,
		Round:			round,					
	}
	dataByts, err := abaData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode aba data")
	}
	msg := &Message{
		MsgType:    ABAMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing aba msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
