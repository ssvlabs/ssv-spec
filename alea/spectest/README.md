# Tests for Alea- [X]BFT Module


## List

### Flow
- [X] Normal case 4 operators
- [X] Normal case 7 operators
- [X] Normal case 10 operators
- [X] Normal case 13 operators

### Messages
- [X] Proposal message creation
- [X] FillGap message creation
- [X] Filler message creation
- [X] VCBCSend message creation
- [X] VCBCReady message creation
- [X] VCBCFinal message creation
- [X] VCBCRequest message creation
- [X] VCBCAnswer message creation
- [X] ABAInit message creation
- [X] ABAAux message creation
- [X] ABAConf message creation
- [X] ABAFinish message creation

- [X] Message with nil identifier
- [X] Message with empty identifier
- [X] Message with type unknown
- [X] Encoding proposal data
- [X] Encoding proposal data and root
- [X] Get root of a message
- [X] Message with nil data
- [X] Message with empty data
- [X] Message with short signature
- [X] Message with long signature
- [X] Message with no signers
- [X] Message with duplicate signatures
- [X] Message with multi signers
- [X] Message with invalid data
- [X] Message with invalid signer

### Proposal
- [X] Proposal with multi signer
- [X] Proposal with unknown signer
- [X] Proposal with wrong height
- [X] Proposal with wrong signature
- [X] Proposal with empty data
- [X] Proposal with post decided proposal

### VCBCSend
- [X] VCBCSend with multi signer
- [X] VCBCSend with unknown signer
- [X] VCBCSend with wrong height
- [X] VCBCSend with wrong signature
- [X] VCBCSend with empty data
- [X] VCBCSend after batch size of proposals
- [X] 2 VCBCSend after 2 batch size of proposals
- [X] VCBCSend receive from other node
- [X] VCBCSend receive multiple msgs from other nodes
- [X] VCBCSend wrong author
- [X] VCBCSend duplicated
- [X] VCBCSend (duplicated) same author and priority with different data

### VCBCReady
- [X] VCBCReady with multi signer
- [X] VCBCReady with unknown signer
- [X] VCBCReady with wrong height
- [X] VCBCReady with wrong signature
- [X] VCBCReady with empty hash
- [X] VCBCReady with wrong hash
- [X] VCBCReady duplicated
- [X] VCBCReady receive msg
- [X] VCBCReady receive quorum
- [X] VCBCReady with unexpected author

### VCBCFinal
- [X] VCBCFinal with multi signer
- [X] VCBCFinal with unknown signer
- [X] VCBCFinal with wrong height
- [X] VCBCFinal with wrong signature
- [X] VCBCFinal with empty aggregated message bytes (proof)
- [X] VCBCFinal with empty hash
- [X] VCBCFinal with wrong aggregated signature
- [X] VCBCFinal duplicated
- [X] VCBCFinal receive message
- [X] VCBCFinal receive message than receive request

### VCBCRequest
- [X] VCBCRequest with multi signer
- [X] VCBCRequest with unknown signer
- [X] VCBCRequest with wrong height
- [X] VCBCRequest with wrong signature
- [X] VCBCRequest with wrong author
- [X] VCBCRequest with wrong priority
- [X] VCBCRequest receive message

### VCBCAnswer
- [X] VCBCAnswer with multi signer
- [X] VCBCAnswer with unknown signer
- [X] VCBCAnswer with wrong height
- [X] VCBCAnswer with wrong signature
- [X] VCBCAnswer with empty data
- [X] VCBCAnswer with wrong data
- [X] VCBCAnswer with wrong priority
- [X] VCBCAnswer with wrong answer

### ABAInit
- [X] ABAInit with multi signer
- [X] ABAInit with unknown signer
- [X] ABAInit with wrong height
- [X] ABAInit with wrong signature
- [X] ABAInit with invalid vote
- [X] ABAInit duplicated
- [X] ABAInit receive
- [X] ABAInit receive quorum
- [X] ABAInit receive two quorums
- [X] ABAInit abastart


### ABAAux
- [X] ABAAux with multi signer
- [X] ABAAux with unknown signer
- [X] ABAAux with wrong height
- [X] ABAAux with wrong signature
- [X] ABAAux with invalid vote
- [X] ABAAux duplicated
- [X] ABAAux receive
- [X] ABAAux receive quorum
- [X] ABAAux receive two quorum
- [X] ABAAux not reaching quorum


### ABAConf
- [X] ABAConf with multi signer
- [X] ABAConf with unknown signer
- [X] ABAConf with wrong height
- [X] ABAConf with wrong signature
- [ ] ABAConf with empty data
- [ ] ABAConf duplicated
- [ ] ABAConf with invalid votes
- [ ] ABAConf reaching quorum
- [ ] ABAConf not reaching quorum
- [ ] ABAConf reaching quorum, len(values)== 1, coin same value
- [ ] ABAConf reaching quorum, len(values)== 1, coin different value
- [ ] ABAConf reaching quorum, len(values)== 2
- [ ] ABAConf reaching quorum, len(values)== 1, coin same value but alredy sent Finish with same value
- [ ] ABAConf reaching quorum, len(values)== 1, coin same value but alredy sent Finish with different value

### ABAFinish
- [X] ABAFinish with multi signer
- [X] ABAFinish with unknown signer
- [X] ABAFinish with wrong height
- [X] ABAFinish with wrong signature
- [ ] ABAFinish with empty data
- [ ] ABAFinish with invalid vote
- [ ] ABAFinish duplicated
- [ ] ABAFinish weak support
- [ ] ABAFinish weak support but already sent Finish
- [ ] ABAFinish strong support

### FillGap
- [X] FillGap with multi signer
- [X] FillGap with unknown signer
- [X] FillGap with wrong height
- [X] FillGap with wrong signature
- [X] FillGap with receive
- [X] FillGap with wrong priority
- [X] FillGap with wrong answer

### Filler
- [X] Filler with multi signer
- [X] Filler with unknown signer
- [X] Filler with wrong height
- [X] Filler with wrong signature
- [X] Filler with empty data
- [X] Filler with wrong answer
- [X] Filler with wrong data
- [X] Filler with wrong priority