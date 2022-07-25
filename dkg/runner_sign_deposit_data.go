package dkg

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (r *Runner) startSigning() error {

	msg, err := r.partialSign()

	if err != nil {
		return err
	}

	r.PartialSignatures[r.Operator.OperatorID] = msg
	base, err := msg.ToBase()
	if err != nil {
		return err
	}
	err = r.broadcastMessages([]dkgtypes.Message{*base}, dkgtypes.PartialSingatureMsgType)
	return nil
}

func (r *Runner) partialSign() (*dkgtypes.ParsedPartialSigMessage, error) {
	share := bls.SecretKey{}
	err := share.Deserialize(r.keygenOutput.SecretShare)
	if err != nil {
		return nil, err
	}

	fork := spec.Version{}
	copy(fork[:], r.InitMsg.Fork)
	root, depData, err := types.GenerateETHDepositData(
		r.keygenOutput.PublicKey,
		r.InitMsg.WithdrawalCredentials,
		fork,
		types.DomainDeposit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate deposit data")
	}
	r.DepositDataRoot = make([]byte, len(root))
	copy(r.DepositDataRoot[:], root)

	rawSig := share.SignByte(root[:])
	sigBytes := rawSig.Serialize()
	var sig spec.BLSSignature
	copy(sig[:], sigBytes)

	copy(depData.DepositData.Signature[:], sigBytes)

	return &dkgtypes.ParsedPartialSigMessage{
		Header: &dkgtypes.MessageHeader{
			SessionId: r.Identifier[:],
			MsgType:   int32(dkgtypes.PartialSingatureMsgType),
			Sender:    uint64(r.Operator.OperatorID),
			Receiver:  0,
		},
		Body: &dkgtypes.PartialSigMsgBody{
			Signer:    uint64(r.Operator.OperatorID),
			Root:      root,
			Signature: sig[:],
		},
	}, nil
}

func (r *Runner) handlePartialSigMessage(baseMsg *dkgtypes.Message) error {

	msg := &dkgtypes.ParsedPartialSigMessage{}
	err := msg.FromBase(baseMsg)
	if err != nil {
		return err
	}

	if found := r.PartialSignatures[types.OperatorID(msg.Header.Sender)]; found == nil {
		r.PartialSignatures[types.OperatorID(msg.Header.Sender)] = msg
	} else if bytes.Compare(found.Body.Signature, msg.Body.Signature) != 0 {
		return errors.New("inconsistent partial signature received")
	}

	if len(r.PartialSignatures) > int(r.InitMsg.Threshold) {
		sigBytes := map[types.OperatorID][]byte{}
		for id, pSig := range r.PartialSignatures {
			if err := r.validateDepositDataSig(pSig.Body); err != nil {
				return errors.Wrap(err, "PartialSigMsgBody invalid")
			}
			sigBytes[id] = pSig.Body.Signature
		}

		sig, err := types.ReconstructSignatures(sigBytes)
		if err != nil {
			return err
		}

		// encrypt Operator's share
		encryptedShare, err := r.Config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.keygenOutput.SecretShare)
		if err != nil {
			return errors.Wrap(err, "could not encrypt share")
		}

		r.signOutput = &dkgtypes.SignedDepositDataMsgBody{
			RequestID:             r.Identifier[:],
			OperatorID:            uint64(r.Operator.OperatorID),
			EncryptedShare:        encryptedShare,
			Committee:             r.keygenOutput.Committee,
			Threshold:             r.InitMsg.Threshold,
			ValidatorPublicKey:    r.keygenOutput.PublicKey,
			WithdrawalCredentials: r.InitMsg.WithdrawalCredentials,
			DepositDataSignature:  sig.Serialize(),
		}
		return nil
	}
	return nil
}

func (r *Runner) validateDepositDataSig(msg *dkgtypes.PartialSigMsgBody) error {
	if !bytes.Equal(r.DepositDataRoot[:], msg.Root) {
		return errors.New("deposit data roots not equal")
	}
	sharePkBytes, err := r.findSignerPubKey(msg.Signer)
	if err != nil {
		return err
	}
	sharePk := &bls.PublicKey{} // TODO: cache this PubKey
	if err := sharePk.Deserialize(sharePkBytes); err != nil {
		return errors.Wrap(err, "could not deserialize public key")
	}

	sig := &bls.Sign{}
	if err := sig.Deserialize(msg.Signature); err != nil {
		return errors.Wrap(err, "could not deserialize partial sig")
	}

	root := make([]byte, 32)
	copy(root, r.DepositDataRoot[:])
	if !sig.VerifyByte(sharePk, root) {
		return errors.New("partial deposit data sig invalid")
	}

	return nil
}

func (r *Runner) findSignerPubKey(signer uint64) ([]byte, error) {

	index := -1
	for i, d := range r.InitMsg.OperatorIDs {
		if d == signer {
			index = i
		}
	}

	if index == -1 {
		return nil, errors.New("signer not part of committee")
	}

	// find operator and verify msg
	sharePkBytes := r.keygenOutput.SharePublicKeys[index]
	return sharePkBytes, nil
}
