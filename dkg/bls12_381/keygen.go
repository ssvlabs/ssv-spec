package bls12_381

import (
	"encoding/json"
	bls_tss "github.com/RockX-SG/bls-tss"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
)

const (
	blsCurve1 = "bls12_381_1"
)

type IKeygen interface {
	Init() []stubdkg.KeygenProtocolMsg
	Output() *stubdkg.LocalKeyShare
	HandleMessage(msg *stubdkg.KeygenProtocolMsg) (bool, []stubdkg.KeygenProtocolMsg, error)
}

type KeygenWrapper struct {
	inner *bls_tss.KeygenSimple
}

func New(myIndex, threshold, groupSize int) *KeygenWrapper {
	inner := bls_tss.NewKeygenSimple(myIndex, threshold, groupSize)
	return &KeygenWrapper{
		inner: inner,
	}
}

func (k *KeygenWrapper) Init() ([]stubdkg.KeygenProtocolMsg, error) {
	innerOutgoing := k.inner.Init()
	outgoing, err := decodeOutgoing(innerOutgoing)
	if err != nil {
		return nil, err
	}
	return outgoing, nil
}

func (k *KeygenWrapper) Output() *stubdkg.LocalKeyShare {
	if k.inner.Output() == nil {
		return nil
	}
	jsonStr := *k.inner.Output()
	innerOutput := bls_tss.LocalKey{}
	err := json.Unmarshal([]byte(jsonStr), innerOutput)
	if err != nil {
		return nil
	}

	var sharePubKeys []stubdkg.BlsPublicKey
	for _, pk := range innerOutput.VkVec {
		var pk48 stubdkg.BlsPublicKey
		copy(pk48[:], i2b(pk.Point))
		sharePubKeys = append(sharePubKeys, pk48)
	}

	share := stubdkg.LocalKeyShare{
		Index:           uint16(innerOutput.SharedKey.I),
		Threshold:       uint16(innerOutput.SharedKey.T),
		ShareCount:      uint16(innerOutput.SharedKey.N),
		SharePublicKeys: sharePubKeys,
	}

	copy(share.PublicKey[:], i2b(innerOutput.SharedKey.Vk.Point))
	copy(share.SecretShare[:], i2b(innerOutput.SharedKey.SkI.Scalar))
	return &share
}

func (k *KeygenWrapper) HandleMessage(msg *stubdkg.KeygenProtocolMsg) (bool, []stubdkg.KeygenProtocolMsg, error) {
	innerMsg := bls_tss.KeygenRoundMsg{
		Sender:   int(msg.Sender),
		Receiver: msg.Receiver,
	}

	switch msg.RoundNumber {
	case stubdkg.KG_R1:
		roundMsg, err := msg.GetRound1Data()
		if err != nil {
			return false, nil, err
		}
		innerMsg.Body.Round1 = &bls_tss.KeygenRound1{Com: b2i(roundMsg.Commitment[:])}
	case stubdkg.KG_R2:
		roundMsg, err := msg.GetRound2Data()
		if err != nil {
			return false, nil, err
		}
		innerMsg.Body.Round2 = &bls_tss.KeygenRound2{
			BlindFactor: b2i(roundMsg.BlindFactor[:]),
			YI: struct {
				Curve string `json:"curve"`
				Point []int  `json:"point"`
			}{
				Curve: blsCurve1,
				Point: b2i(roundMsg.YI[:]),
			},
		}
	case stubdkg.KG_R3:
		roundMsg, err := msg.GetRound3Data()
		if err != nil {
			return false, nil, err
		}
		var coms []struct {
			Curve string `json:"curve"`
			Point []int  `json:"point"`
		}
		for _, commitment := range roundMsg.Commitments {
			var com struct {
				Curve string `json:"curve"`
				Point []int  `json:"point"`
			}
			com.Curve = blsCurve1
			com.Point = b2i(commitment[:])
			coms = append(coms, com)
		}
		innerMsg.Body.Round3 = &bls_tss.KeygenRound3{
			I:           int(msg.Sender),
			T:           roundMsg.Parameters.Threshold,
			N:           roundMsg.Parameters.ShareCount,
			J:           int(msg.Receiver),
			Commitments: coms,
			Share: struct {
				Curve  string `json:"curve"`
				Scalar []int  `json:"scalar"`
			}{
				Curve:  blsCurve1,
				Scalar: b2i(roundMsg.ShareIJ[:]),
			},
		}
	case stubdkg.KG_R4:
		roundMsg, err := msg.GetRound4Data()
		if err != nil {
			return false, nil, err
		}
		innerMsg.Body.Round4 = &bls_tss.KeygenRound4{
			Pk: struct {
				Curve string `json:"curve"`
				Point []int  `json:"point"`
			}{
				Curve: blsCurve1,
				Point: b2i(roundMsg.Pk[:]),
			},
			PkTRandCommitment: struct {
				Curve string `json:"curve"`
				Point []int  `json:"point"`
			}{
				Curve: blsCurve1,
				Point: b2i(roundMsg.PkTRandCommitment[:]),
			},
			ChallengeResponse: struct {
				Curve  string `json:"curve"`
				Scalar []int  `json:"scalar"`
			}{
				Curve:  blsCurve1,
				Scalar: b2i(roundMsg.ChallengeResponse[:]),
			},
		}
	}
	bytes, err := json.Marshal(innerMsg)
	finished, innerOutgoing, err := k.inner.Handle(string(bytes))
	outgoing, err := decodeOutgoing(innerOutgoing)
	if err != nil {
		return false, nil, err
	}
	return finished, outgoing, err
}

func decodeOutgoing(innerOutgoing []string) ([]stubdkg.KeygenProtocolMsg, error) {
	var outgoing []stubdkg.KeygenProtocolMsg
	for _, msg0 := range innerOutgoing {
		out0 := bls_tss.KeygenRoundMsg{}
		err := json.Unmarshal([]byte(msg0), out0)
		if err != nil {
			return nil, err
		}

		out := stubdkg.KeygenProtocolMsg{}

		if out0.Body.Round4 != nil {
			data := &stubdkg.KeygenRound4Data{
				Pk:                stubdkg.BlsPublicKey{},
				PkTRandCommitment: stubdkg.BlsPublicKey{},
				ChallengeResponse: stubdkg.BlsScalar{},
			}
			copy(data.Pk[:], i2b(out0.Body.Round4.Pk.Point))
			copy(data.PkTRandCommitment[:], i2b(out0.Body.Round4.PkTRandCommitment.Point))
			copy(data.ChallengeResponse[:], i2b(out0.Body.Round4.ChallengeResponse.Scalar))
			err := out.SetRound4Data(data)
			if err != nil {
				return nil, err
			}
		} else if out0.Body.Round3 != nil {
			data := &stubdkg.KeygenRound3Data{
				Parameters: struct {
					Threshold  int `json:"threshold"`
					ShareCount int `json:"share_count"`
				}{
					Threshold:  out0.Body.Round3.T,
					ShareCount: out0.Body.Round3.N,
				},
			}
			var commitments []stubdkg.BlsPublicKey
			for _, commitment := range out0.Body.Round3.Commitments {
				var pt stubdkg.BlsPublicKey
				copy(pt[:], i2b(commitment.Point))
				commitments = append(commitments, pt)
			}
			data.Commitments = commitments
			copy(data.ShareIJ[:], i2b(out0.Body.Round3.Share.Scalar))

			err := out.SetRound3Data(data)
			if err != nil {
				return nil, err
			}
		} else if out0.Body.Round2 != nil {

			data := &stubdkg.KeygenRound2Data{}
			data.BlindFactor = i2b(out0.Body.Round2.BlindFactor)
			copy(data.YI[:], i2b(out0.Body.Round2.YI.Point))
			err := out.SetRound2Data(data)
			if err != nil {
				return nil, err
			}
		} else if out0.Body.Round1 != nil {
			var com []byte
			for _, i := range out0.Body.Round1.Com {
				com = append(com, byte(i))
			}
			data := &stubdkg.KeygenRound1Data{}
			copy(data.Commitment[:], com)
			err := out.SetRound1Data(data)
			if err != nil {
				return nil, err
			}
		}

		out.Sender = uint16(out0.Sender)

		if receiver, ok := out0.Receiver.(uint16); ok {
			out.Receiver = receiver
		}
		outgoing = append(outgoing, out)
	}
	return outgoing, nil
}

func i2b(input []int) []byte {
	var out []byte
	for _, i := range input {
		out = append(out, byte(i))
	}
	return out
}

func b2i(input []byte) []int {
	var out []int
	for _, b := range input {
		out = append(out, int(b))
	}
	return out
}
