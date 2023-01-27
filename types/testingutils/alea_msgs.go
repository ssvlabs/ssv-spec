package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var MultiSignAleaMsg = func(sks []*bls.SecretKey, ids []types.OperatorID, msg *alea.Message) *alea.SignedMessage {
	if len(sks) == 0 || len(ids) != len(sks) {
		panic("sks != ids")
	}
	var signed *alea.SignedMessage
	for i, sk := range sks {
		if signed == nil {
			signed = SignAleaMsg(sk, ids[i], msg)
		} else {
			if err := signed.Aggregate(SignAleaMsg(sk, ids[i], msg)); err != nil {
				panic(err.Error())
			}
		}
	}

	return signed
}

var SignAleaMsg = func(sk *bls.SecretKey, id types.OperatorID, msg *alea.Message) *alea.SignedMessage {
	domain := types.PrimusTestnet
	sigType := types.QBFTSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig := sk.SignByte(r)

	return &alea.SignedMessage{
		Message:   msg,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}
var ProposalDataBytesAlea = func(data []byte) []byte {
	d := &alea.ProposalData{
		Data: data,
	}
	ret, _ := d.Encode()
	return ret
}
var FillGapDataBytes = func(operatorID types.OperatorID, priority alea.Priority) []byte {
	d := &alea.FillGapData{
		OperatorID: operatorID,
		Priority:   priority,
	}
	ret, _ := d.Encode()
	return ret
}
var FillerDataBytes = func(entries [][]*alea.ProposalData, priorities []alea.Priority, operatorID types.OperatorID) []byte {
	d := &alea.FillerData{
		Entries:    entries,
		Priorities: priorities,
		OperatorID: operatorID,
	}
	ret, _ := d.Encode()
	return ret
}
var ABAInitDataBytes = func(vote byte, round alea.Round) []byte {
	d := &alea.ABAInitData{
		Vote:  vote,
		Round: round,
	}
	ret, _ := d.Encode()
	return ret
}
var ABAAuxDataBytes = func(vote byte, round alea.Round) []byte {
	d := &alea.ABAAuxData{
		Vote:  vote,
		Round: round,
	}
	ret, _ := d.Encode()
	return ret
}
var ABAConfDataBytes = func(votes []byte, round alea.Round) []byte {
	d := &alea.ABAConfData{
		Votes: votes,
		Round: round,
	}
	ret, _ := d.Encode()
	return ret
}
var ABAFinishDataBytes = func(vote byte) []byte {
	d := &alea.ABAFinishData{
		Vote: vote,
	}
	ret, _ := d.Encode()
	return ret
}
var VCBCSendDataBytes = func(proposals []*alea.ProposalData, priority alea.Priority, author types.OperatorID) []byte {
	d := &alea.VCBCSendData{
		Proposals: proposals,
		Priority:  priority,
		Author:    author,
	}
	ret, _ := d.Encode()
	return ret
}
var VCBCReadyDataBytes = func(hash []byte, priority alea.Priority, author types.OperatorID) []byte {
	d := &alea.VCBCReadyData{
		Hash:     hash,
		Priority: priority,
		Author:   author,
	}
	ret, _ := d.Encode()
	return ret
}
var VCBCFinalDataBytes = func(hash []byte, priority alea.Priority, proof types.Signature, author types.OperatorID) []byte {
	d := &alea.VCBCFinalData{
		Hash:     hash,
		Priority: priority,
		Proof:    proof,
		Author:   author,
	}
	ret, _ := d.Encode()
	return ret
}
var VCBCRequestDataBytes = func(priority alea.Priority, author types.OperatorID) []byte {
	d := &alea.VCBCRequestData{
		Priority: priority,
		Author:   author,
	}
	ret, _ := d.Encode()
	return ret
}
var VCBCAnswerDataBytes = func(proposals []*alea.ProposalData, priority alea.Priority, proof types.Signature, author types.OperatorID) []byte {
	d := &alea.VCBCAnswerData{
		Proposals: proposals,
		Priority:  priority,
		Proof:     proof,
		Author:    author,
	}
	ret, _ := d.Encode()
	return ret
}
