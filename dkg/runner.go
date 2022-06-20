package dkg

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Config struct {
	// Protocol the DKG protocol implementation
	Protocol      Protocol
	BeaconNetwork ssv.BeaconNetwork
	Network       Network
	// OperatorID the node's operator ID
	OperatorID types.OperatorID
	// Identifier unique for DKG session
	Identifier types.MessageID
	Signer     types.DKGSigner
	// PubKey signing key for all message
	PubKey *ecdsa.PublicKey
	// EncryptionPubKey encryption pubkey for shares
	EncryptionPubKey *rsa.PublicKey
}

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operators             []types.OperatorID
	I                     uint16
	Threshold             uint16
	WithdrawalCredentials []byte

	ProtocolOutput        *ProtocolOutput
	PartialSignatures     map[types.OperatorID][]byte
	DepositDataSignatures map[types.OperatorID]*PartialDepositData

	protocol Protocol
	config   *Config
}

func NewRunner(initMsg *Init, config *Config) (*Runner, error) {
	var i uint16
	for i0, id := range initMsg.OperatorIDs {
		if id == config.OperatorID {
			i = uint16(i0) + 1
		}
	}
	if i == 0 {
		return nil, errors.New("invalid request")
	}
	runner := &Runner{
		I:                     i,
		Operators:             initMsg.OperatorIDs,
		Threshold:             initMsg.Threshold,
		WithdrawalCredentials: initMsg.WithdrawalCredentials,
		PartialSignatures:     map[types.OperatorID][]byte{},
		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},

		protocol: config.Protocol,
		config:   config,
	}

	if err := runner.protocol.Start(initMsg); err != nil {
		return nil, errors.Wrap(err, "could not start dkg protocol")
	}

	return runner, nil
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
					r.Operators,
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
				r.PartialSignatures[r.config.OperatorID] = pSig.SigmaI[:]
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
		id := r.Operators[pMsg.I-1]
		if found := r.PartialSignatures[id]; found == nil {
			r.PartialSignatures[id] = pMsg.SigmaI[:]
		} else if bytes.Compare(found, pMsg.SigmaI[:]) != 0 {
			return false, nil, errors.New("inconsistent partial signature received")
		}
		if len(r.PartialSignatures) > int(r.Threshold) {

			sig, err := types.ReconstructSignatures(r.PartialSignatures)
			if err != nil {
				return false, nil, err
			}
			// encrypt operator's share
			encryptedShare, err := r.config.Signer.Encrypt(r.config.EncryptionPubKey, r.ProtocolOutput.Share.Serialize())
			if err != nil {
				return false, nil, errors.Wrap(err, "could not encrypt share")
			}

			ret, err := r.generateSignedOutput(&Output{
				Identifier:            r.config.Identifier,
				EncryptedShare:        encryptedShare,
				DKGSetSize:            uint16(len(r.Operators)),
				Threshold:             r.Threshold,
				ValidatorPubKey:       r.ProtocolOutput.ValidatorPK,
				WithdrawalCredentials: r.WithdrawalCredentials,
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
	sig, err := r.config.Signer.SignDKGOutput(o, r.config.PubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &SignedOutput{
		Data:      o,
		Signer:    r.config.OperatorID,
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
		WithdrawalCredentials: r.WithdrawalCredentials,
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
