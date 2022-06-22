package stubdkg

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleDKG(t *testing.T) {
	types.InitBLS()

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}
	k := 3
	polyDegree := k - 1
	payloadToSign := "hello"

	// create polynomials for each operator
	poly := make(map[types.OperatorID][]bls.Fr)
	for _, id := range operators {
		coeff := make([]bls.Fr, 0)
		for i := 1; i <= polyDegree; i++ {
			c := bls.Fr{}
			c.SetByCSPRNG()
			coeff = append(coeff, c)
		}
		poly[id] = coeff
	}

	// create points for each operator
	points := make(map[types.OperatorID][]*bls.Fr)
	for _, id := range operators {
		for _, evalID := range operators {
			if points[evalID] == nil {
				points[evalID] = make([]*bls.Fr, 0)
			}

			res := &bls.Fr{}
			x := &bls.Fr{}
			x.SetInt64(int64(evalID))
			require.NoError(t, bls.FrEvaluatePolynomial(res, poly[id], x))

			points[evalID] = append(points[evalID], res)
		}
	}

	// calculate shares
	shares := make(map[types.OperatorID]*bls.SecretKey)
	pks := make(map[types.OperatorID]*bls.PublicKey)
	sigs := make(map[types.OperatorID]*bls.Sign)
	for id, ps := range points {
		var sum *bls.Fr
		for _, p := range ps {
			if sum == nil {
				sum = p
			} else {
				bls.FrAdd(sum, sum, p)
			}
		}
		shares[id] = bls.CastToSecretKey(sum)
		pks[id] = shares[id].GetPublicKey()
		sigs[id] = shares[id].Sign(payloadToSign)
	}

	// get validator pk
	validatorPK := bls.PublicKey{}
	idVec := make([]bls.ID, 0)
	pkVec := make([]bls.PublicKey, 0)
	for operatorID, pk := range pks {
		blsID := bls.ID{}
		err := blsID.SetDecString(fmt.Sprintf("%d", operatorID))
		require.NoError(t, err)
		idVec = append(idVec, blsID)

		pkVec = append(pkVec, *pk)
	}
	require.NoError(t, validatorPK.Recover(pkVec, idVec))
	fmt.Printf("validator pk: %vDKG\n", hex.EncodeToString(validatorPK.Serialize()))

	// reconstruct sig
	reconstructedSig := bls.Sign{}
	idVec = make([]bls.ID, 0)
	sigVec := make([]bls.Sign, 0)
	for operatorID, sig := range sigs {
		blsID := bls.ID{}
		err := blsID.SetDecString(fmt.Sprintf("%d", operatorID))
		require.NoError(t, err)
		idVec = append(idVec, blsID)

		sigVec = append(sigVec, *sig)

		if len(sigVec) >= k {
			break
		}
	}
	require.NoError(t, reconstructedSig.Recover(sigVec, idVec))
	fmt.Printf("reconstructed sig: %vDKG\n", hex.EncodeToString(reconstructedSig.Serialize()))

	// verify
	require.True(t, reconstructedSig.Verify(&validatorPK, payloadToSign))
}

type simulator struct {
	operators  []types.OperatorID
	identifier dkg.RequestID
	machines   []dkg.Protocol
}

func makeSimulator(n uint16) *simulator {
	var operators []types.OperatorID
	var identifier dkg.RequestID
	for i := 1; i < int(n)+1; i++ {
		operators = append(operators, types.OperatorID(i))
	}
	var machines []dkg.Protocol
	for _, operator := range operators {
		m := New(operator, identifier)
		machines = append(machines, m)
	}
	return &simulator{
		operators:  operators,
		identifier: identifier,
		machines:   machines,
	}
}

func (s *simulator) feedSuccess(t *testing.T, index int, msg *dkg.Message) []dkg.Message {
	out, err := s.machines[index-1].ProcessMsg(msg)
	assert.Nil(t, err)
	return out
}

func (s *simulator) initSuccess(t *testing.T, init *dkg.Init) []dkg.Message {
	var round1 []dkg.Message
	for _, machine := range s.machines {
		data, err := init.Encode()
		assert.Nil(t, err)
		out, err := machine.ProcessMsg(&dkg.Message{
			MsgType:    dkg.InitMsgType,
			Identifier: s.identifier,
			Data:       data,
		})
		assert.Nil(t, err)
		for _, message := range out {
			assert.Equal(t, dkg.ProtocolMsgType, message.MsgType)
			round1 = append(round1, message)
		}
	}

	return round1
}

func (s *simulator) checkAndFeedRoundMessages(t *testing.T, round Round, messages []dkg.Message, finished bool) []dkg.Message {
	groupSize := len(s.operators)
	var nextRound []dkg.Message

	for _, message := range messages {
		msg := KeygenProtocolMsg{}
		err := msg.Decode(message.Data)
		assert.Nil(t, err)
		assert.Equal(t, round, msg.RoundNumber)

		for i := 1; i < int(groupSize)+1; i++ {
			if i == int(msg.Receiver) || (msg.Receiver == uint16(0) && msg.Sender != uint16(i)) {
				out := s.feedSuccess(t, i, &message)
				nextRound = append(nextRound, out...)
			}
		}
	}
	if finished {
		assert.Equal(t, dkg.KeygenOutputType, nextRound[len(nextRound)-1].MsgType)
	}

	return nextRound
}

func checkTwoOfThree(t *testing.T, pubKey []byte, indices []bls.Fr, shares []bls.Fr, participating []int) {
	var quorumIndices, quorumShares []bls.Fr
	for _, ind := range participating {
		quorumIndices = append(quorumIndices, indices[ind-1])
		quorumShares = append(quorumShares, shares[ind-1])
	}
	var secret bls.Fr
	bls.FrLagrangeInterpolation(&secret, indices[:2], shares[:2])
	secretBytes := secret.Serialize()
	var sk bls.SecretKey
	sk.Deserialize(secretBytes)
	assert.Equal(t, pubKey, sk.GetPublicKey().Serialize())
}

func TestStub(t *testing.T) {
	types.InitBLS()
	threshold := uint16(1)
	groupSize := uint16(3)
	simulator := makeSimulator(groupSize)

	round1 := simulator.initSuccess(t, &dkg.Init{
		Nonce:                 0,
		OperatorIDs:           simulator.operators,
		Threshold:             threshold,
		WithdrawalCredentials: []byte(""),
	})

	assert.Equal(t, 3, len(round1))
	round2 := simulator.checkAndFeedRoundMessages(t, Round(1), round1, false)
	assert.Equal(t, 3, len(round2))
	round3 := simulator.checkAndFeedRoundMessages(t, Round(2), round2, false)
	assert.Equal(t, 6, len(round3))
	round4 := simulator.checkAndFeedRoundMessages(t, Round(3), round3, false)
	assert.Equal(t, 3, len(round4))
	outputMsgs := simulator.checkAndFeedRoundMessages(t, Round(4), round4, true)
	assert.Equal(t, 3, len(outputMsgs))
	var lastOut *dkg.KeygenOutput
	indices := make([]bls.Fr, groupSize)
	shares := make([]bls.Fr, groupSize)
	for _, outputMsg := range outputMsgs {
		assert.Equal(t, dkg.KeygenOutputType, outputMsg.MsgType)
		output := dkg.KeygenOutput{}
		err := output.Decode(outputMsg.Data)
		assert.Nil(t, err)
		assert.Equal(t, threshold, output.Threshold)
		assert.Equal(t, groupSize, output.ShareCount)
		if lastOut != nil {
			assert.Equal(t, lastOut.PublicKey, output.PublicKey)
			assert.Equal(t, lastOut.SharePublicKeys, output.SharePublicKeys)
		}
		skBytes := make([]byte, len(output.SecretShare))
		copy(skBytes, output.SecretShare)
		i := output.Index - 1
		err = shares[i].SetBigEndianMod(skBytes)
		assert.Nil(t, err)
		indices[i].SetInt64(int64(output.Index))

		lastOut = &output
	}

	var pubKeys [][]byte

	for _, share := range shares {
		pubKeys = append(pubKeys, bls.CastToSecretKey(&share).GetPublicKey().Serialize())
	}
	assert.Equal(t, lastOut.SharePublicKeys, pubKeys)

	for i := 1; i < int(groupSize)+1; i++ {
		fr := bls.Fr{}
		fr.SetInt64(int64(i))
		indices = append(indices, fr)
	}

	checkTwoOfThree(t, lastOut.PublicKey, indices, shares, []int{1, 2})
	checkTwoOfThree(t, lastOut.PublicKey, indices, shares, []int{1, 3})
	checkTwoOfThree(t, lastOut.PublicKey, indices, shares, []int{2, 3})
}
