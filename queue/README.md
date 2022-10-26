
# VALIDATOR QUEUE

## Introduction

The queue objective is to receive messages and pop the right msg by to specify priority.
The priority is based on what will achieve the quickest consensus and keep the node in sync with other peers.

### Queue

Receive msg and add to stack. queue provide interface with ADD, POP, DELETE functions. 

### Index Points

- Height
- Slot
- Consensus
  - Propose
  - Prepare
  - Commit
  - Commit (decided)
  - Change Round
- Post Consensus
  - PostConsensusPartialSig
  - RandaoPartialSig
  - SelectionProofPartialSig
  - ContributionProofs

### Priority

1. Higher height (all above the current height)
2. Post Consensus (by the current slot)
3. All Decided
4. Commit by current height
5. Consensus by state (explained below)
6. Lower Height (all below current height)


**Why Need Current Height?**
> Assuming queue added the following msg's by this order
> 1. propose (height 10)
> 2. change round*2 (height 9)
> 3. Prepare*3 (height 10)
> 4. Commit*3 (height 10)
> 5. Post Consensus (slot 100)
>
>if queue don't know the current height and goes only by highest height, 
it will pop all height msg's and only then the post consensus (which don't hold height) so there is no way to know when to 
look for post consensus. if queue knows what is the expected height, once no more msgs for height it will look for PC msg type.  

**Why need current slot?**
> Assuming queue added the following msg's by this order
> 1. Post Consensus (slot 200)
> 2. Post Consensus (slot 200)
> 3. Post Consensus (slot 100)
> 
> In this case without the current slot the queue will consume the PC with 200 slot first and only after that the 100 slot. this will result in validation fail for the 200 slot cause the runner already runs with the 100 slot.

**Consensus By State**
> Priority for consensus msg's should be by the state if the instance.
> for example, in case where there are Proposal and Prepare msg's for the same height, we need to process Proposal before Prepare. if not, the Prepare msg's will fail in validation and be wasted 

**Lower height**
> We need to pop msg's lower than the current height up to X below. 

