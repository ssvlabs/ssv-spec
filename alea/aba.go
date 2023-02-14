package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) StartAgreementComponent() error {

	if i.verbose {
		fmt.Println("StartingAgreementComponent")
	}
	for {
		if i.verbose {
			fmt.Println("\tAC: Round:", i.State.ACState.ACRound)
		}

		// check if it should stop performing agreement
		if i.State.StopAgreement {
			break
		}

		// calculate the round leader (to get value to be decided on)
		leader := i.config.GetProposerF()(i.State, Round(i.State.ACState.ACRound))
		if i.verbose {
			fmt.Println("\tLeader:", leader)
		}

		// get the local queue associated with the leader's id (create if there isn't one)
		if _, exists := i.State.VCBCState.Queues[leader]; !exists {
			i.State.VCBCState.Queues[leader] = NewVCBCQueue()
		}
		queue := i.State.VCBCState.Queues[leader]

		// get the value of the queue with the lowest priority value
		value, priority := queue.Peek()
		if i.verbose {
			fmt.Println("\tPeek:", value, priority)
		}

		// decide own vote
		vote := byte(0)
		if value == nil {
			vote = byte(0)
		} else {
			vote = byte(1)
		}
		if i.verbose {
			fmt.Println("\tvote:", vote)
		}

		// start ABA protocol
		result, err := i.StartABA(vote)
		if err != nil {
			return errors.Wrap(err, "failed to start ABA and get result")
		}
		if i.verbose {
			fmt.Println("\tABA result:", result)
		}
		if i.State.StopAgreement {
			break
		}

		if result == 1 {
			// if the protocol agreed on the value of the leader replica, deliver it

			// if ABA decided 1 but own vote was 0, start recover mechanism to get VCBC messages not received from leader
			if vote == 0 {
				if i.verbose {
					fmt.Println("\tresult 1 but voted 0")
				}
				if !i.State.VCBCState.HasM(leader, priority) {
					// create FILLGAP message
					fillerContLen := i.State.FillerContainer.Len(i.State.AleaDefaultRound)
					fillGapMsg, err := CreateFillGap(i.State, i.config, leader, priority)
					if err != nil {
						return errors.Wrap(err, "StartAgreementComponent: failed to create FillGap message")
					}
					if i.verbose {
						fmt.Println("\tBroadcasting fill gap")
					}
					i.Broadcast(fillGapMsg)
					// wait for the value to be received
					i.WaitFillGapResponse(leader, priority, fillerContLen)
					if i.verbose {
						fmt.Println("\tgot filler response")
					}
				}
			}

			// get decided value
			value, priority = queue.Peek()
			if i.verbose {
				fmt.Println("\tpeek before delivering:", value, priority)
			}

			// remove the value from the queue and add it to S
			queue.Dequeue()
			if i.verbose {
				fmt.Println("\tqueue dequeued. New queue:", queue)
			}
			i.State.Delivered.Enqueue(value, priority)
			if i.verbose {
				fmt.Println("\tdelivered enqueued values:", i.State.Delivered)
			}
			// return the value to the client
			i.Deliver(value)
		}
		// increment the round number
		i.State.ACState.IncrementRound()
	}
	return nil
}

func (i *Instance) WaitFillGapResponse(leader types.OperatorID, priority Priority, fillerContLen int) {
	// gets the leader queue
	queue := i.State.VCBCState.Queues[leader]
	currentFillerNum := fillerContLen
	for {
		// if has the desired priority, returns
		_, localPriority := queue.Peek()
		if localPriority >= priority {
			return
		}

		// waits until a FILLER signal is received (actived on the uponFiller function)
		for {
			newLen := i.State.FillerContainer.Len(i.State.AleaDefaultRound)
			if newLen > currentFillerNum {
				currentFillerNum = newLen
				break
			}
		}
	}
}

func (i *Instance) StartABA(vote byte) (byte, error) {
	// set ABA's input value
	i.State.ACState.GetCurrentABAState().SetVInput(i.State.ACState.GetCurrentABAState().Round, vote)

	// broadcast INIT message with input vote
	initMsg, err := CreateABAInit(i.State, i.config, vote, i.State.ACState.GetCurrentABAState().Round, i.State.ACState.ACRound)
	if err != nil {
		return byte(2), errors.Wrap(err, "StartABA: failed to create ABA Init message")
	}
	i.Broadcast(initMsg)

	// update sent flag
	i.State.ACState.GetCurrentABAState().SetSentInit(i.State.ACState.GetCurrentABAState().Round, vote, true)

	// process own init msg
	i.uponABAInit(initMsg)

	// wait until channel Terminate receives a signal
	for {
		if i.State.ACState.GetCurrentABAState().Terminate || i.State.StopAgreement {
			break
		}
	}

	// i.State.ACState.GetCurrentABAState().Terminate = false

	// returns the decided value
	return i.State.ACState.GetCurrentABAState().Vdecided, nil
}
