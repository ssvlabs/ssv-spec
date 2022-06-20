package dkg

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operator *Operator
	// InitMsg holds the init method which started this runner
	InitMsg *Init
	// Identifier unique for DKG session
	Identifier RequestID
	// ProtocolOutput holds the protocol output once it finishes
	ProtocolOutput *ProtocolOutput
	// DepositDataSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	I                 uint16

	protocol Protocol
	config   *Config
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, *SignedOutput, error) {
	// TODO - validate message

	switch msg.Message.MsgType {
	case ProtocolMsgType:
		pMsg := stubdkg.KeygenProtocolMsg{}
		err := pMsg.Decode(msg.Message.Data)
		if err != nil {
			return false, nil, err
		}
		finished, outgoing, err := r.protocol.ProcessMsg(&pMsg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}

		for _, message := range outgoing {
			r.config.Network.Broadcast(&message)
		}

		if finished {
			keygenOutput := r.protocol.Output()
			if keygenOutput != nil {
				share := bls.SecretKey{} // TODO: when to drop this object?
				share.Deserialize(keygenOutput.SecretShare[:])
				var pubkeys []bls.PublicKey // TODO: when to drop these objects?
				for _, key := range keygenOutput.SharePublicKeys {
					pk := bls.PublicKey{}
					pk.Deserialize(key[:])
					pubkeys = append(pubkeys, pk)
				}
				r.ProtocolOutput = &ProtocolOutput{
					&share,
					r.InitMsg.OperatorIDs,
					pubkeys,
					keygenOutput.PublicKey[:],
				}
				// TODO broadcast partial deposit data
				root, err := r.getDepositDataSigningRoot(keygenOutput.PublicKey)
				if err != nil {
					return false, nil, err
				}
				sig := share.SignByte(root[:])
				sig.Serialize()
				pSig := &stubdkg.PartialSignature{}
				pSig.I = r.I
				copy(pSig.SigmaI[:], sig.Serialize()[:])
				r.PartialSignatures[r.Operator.OperatorID] = pSig.SigmaI[:]
				r.config.Network.BroadcastPartialSignature(pSig)
			} else {
				return false, nil, errors.New("Unexpected state")
			}
		}
	case PartialSigType:
		pMsg := stubdkg.PartialSignature{}
		err := pMsg.Decode(msg.Message.Data)
		if err != nil {
			return false, nil, err
		}
		id := r.InitMsg.OperatorIDs[pMsg.I-1]
		if found := r.PartialSignatures[id]; found == nil {
			r.PartialSignatures[id] = pMsg.SigmaI[:]
		} else if bytes.Compare(found, pMsg.SigmaI[:]) != 0 {
			return false, nil, errors.New("inconsistent partial signature received")
		}
		if len(r.PartialSignatures) > int(r.InitMsg.Threshold) {

			sig, err := types.ReconstructSignatures(r.PartialSignatures)
			if err != nil {
				return false, nil, err
			}

			// encrypt Operator's share
			encryptedShare, err := r.config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.ProtocolOutput.Share.Serialize())
			if err != nil {
				return false, nil, errors.Wrap(err, "could not encrypt share")
			}

			ret, err := r.generateSignedOutput(&Output{
				Identifier:            r.Identifier,
				EncryptedShare:        encryptedShare,
				DKGSetSize:            uint16(len(r.InitMsg.OperatorIDs)),
				Threshold:             r.InitMsg.Threshold,
				ValidatorPubKey:       r.ProtocolOutput.ValidatorPK,
				WithdrawalCredentials: r.InitMsg.WithdrawalCredentials,
				SignedDepositData:     sig.Serialize(),
			})
			if err != nil {
				return false, nil, errors.Wrap(err, "could not generate dkg SignedOutput")
			}
			return true, ret, nil
		}
		// TODO: Do we need to aggregate the signed outputs.
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) reconstructDepositDataSignature() (types.Signature, error) {
	panic("implement")
}

func (r *Runner) generateSignedOutput(o *Output) (*SignedOutput, error) {
	sig, err := r.config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &SignedOutput{
		Data:      o,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	}, nil
}

func (r *Runner) getDepositDataSigningRoot(pubKey stubdkg.BlsPublicKey) (spec.Root, error) {
	var (
		domain   spec.Domain
		forkData spec.ForkData
	)
	message := spec.DepositMessage{
		PublicKey:             spec.BLSPubKey(pubKey),
		WithdrawalCredentials: r.InitMsg.WithdrawalCredentials,
		Amount:                32_000_000_000,
	}
	depositDomain := spec.DomainType{0x03, 0x00, 0x00, 0x00}

	msgRoot, err := message.HashTreeRoot()
	if err != nil {
		return [32]byte{}, err
	}

	copy(forkData.CurrentVersion[:], r.config.BeaconNetwork.ForkVersion())
	forkDataRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return [32]byte{}, err
	}

	copy(domain[:], depositDomain[:])
	copy(domain[4:], forkDataRoot[:])
	signingData := spec.SigningData{
		ObjectRoot: msgRoot,
		Domain:     domain,
	}
	root, err := signingData.HashTreeRoot()
	if err != nil {
		return [32]byte{}, err
	}
	return root, nil
}
