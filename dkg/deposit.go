package dkg

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type SignDepositData struct {
	// Identifier unique for DKG session
	Identifier RequestID
	// InitMsg holds the details of this instance
	InitMsg *Init
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	DepositDataRoot   spec.Root
	key               *KeygenOutput
	config            ProtocolConfig
}

func NewSignDepositData(init *Init, key *KeygenOutput, config ProtocolConfig) *SignDepositData {

	return &SignDepositData{
		InitMsg:           init,
		PartialSignatures: map[types.OperatorID][]byte{},
		key:               key,
		config:            config,
	}
}

func (s *SignDepositData) Start() ([]Message, error) {
	pSig, err := s.partialSign()
	if err != nil {
		return nil, err
	}
	data, err := pSig.Encode()
	if err != nil {
		return nil, err
	}
	partialSigMsg := Message{
		MsgType:    PartialSigType,
		Identifier: s.Identifier,
		Data:       data,
	}
	return []Message{partialSigMsg}, nil
}

func (s *SignDepositData) ProcessMsg(msg *Message) ([]Message, error) {

	if msg.MsgType != PartialSigType {
		return nil, errors.New("invalid message type")
	}

	pMsg := PartialSignature{}
	err := pMsg.Decode(msg.Data)
	if err != nil {
		return nil, err
	}
	id := s.InitMsg.OperatorIDs[pMsg.I-1]
	if found := s.PartialSignatures[id]; found == nil {
		s.PartialSignatures[id] = pMsg.SigmaI[:]
	} else if bytes.Compare(found, pMsg.SigmaI[:]) != 0 {
		return nil, errors.New("inconsistent partial signature received")
	}
	if len(s.PartialSignatures) > int(s.InitMsg.Threshold) {

		sig, err := types.ReconstructSignatures(s.PartialSignatures)
		if err != nil {
			return nil, err
		}

		// encrypt Operator's share
		encryptedShare, err := s.config.Signer.Encrypt(s.config.Operator.EncryptionPubKey, s.key.SecretShare)
		if err != nil {
			return nil, errors.Wrap(err, "could not encrypt share")
		}

		signedOut, err := s.generateSignedOutput(&Output{
			RequestID:             s.Identifier,
			ShareIndex:            s.key.Index,
			EncryptedShare:        encryptedShare,
			DKGSetSize:            uint16(len(s.InitMsg.OperatorIDs)),
			Threshold:             s.InitMsg.Threshold,
			SharePubKeys:          s.key.SharePublicKeys,
			ValidatorPubKey:       s.key.PublicKey,
			WithdrawalCredentials: s.InitMsg.WithdrawalCredentials,
			DepositDataSignature:  sig.Serialize(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not generate dkg SignedOutput")
		}
		data, err := signedOut.Encode()
		if err != nil {
			return nil, err
		}
		outMsg := Message{
			MsgType:    PartialOutputMsgType,
			Identifier: s.Identifier,
			Data:       data,
		}
		return []Message{outMsg}, nil
	}
	return nil, nil
}

func (s *SignDepositData) generateSignedOutput(o *Output) (*SignedOutput, error) {
	sig, err := s.config.Signer.SignDKGOutput(o, s.config.Operator.ETHAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &SignedOutput{
		Data:      o,
		Signer:    s.config.Operator.OperatorID,
		Signature: sig,
	}, nil
}

func (s *SignDepositData) partialSign() (*PartialSignature, error) {
	share := bls.SecretKey{}
	err := share.Deserialize(s.key.SecretShare)
	if err != nil {
		return nil, err
	}

	root, err := s.getDepositDataSigningRoot(s.key.PublicKey)
	if err != nil {
		return nil, err
	}
	rawSig := share.SignByte(root[:])
	sigBytes := rawSig.Serialize()
	var sig spec.BLSSignature
	copy(sig[:], sigBytes)
	return &PartialSignature{
		I:      s.key.Index,
		SigmaI: sig,
	}, nil
}

func (s *SignDepositData) getDepositDataSigningRoot(pubKey []byte) (spec.Root, error) {
	var (
		domain   spec.Domain
		forkData spec.ForkData
		pk48     spec.BLSPubKey
	)
	copy(pk48[:], pubKey)
	message := spec.DepositMessage{
		PublicKey:             pk48,
		WithdrawalCredentials: s.InitMsg.WithdrawalCredentials,
		Amount:                32_000_000_000,
	}
	depositDomain := spec.DomainType{0x03, 0x00, 0x00, 0x00}

	msgRoot, err := message.HashTreeRoot()
	if err != nil {
		return [32]byte{}, err
	}

	forkVersion := s.config.BeaconNetwork.ForkVersion()
	copy(forkData.CurrentVersion[:], forkVersion[:])
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

// TODO: Standardize PartialDepositData and PartialSignature
func (s *SignDepositData) validateDepositDataSig(msg *PartialDepositData) error {
	if !bytes.Equal(s.DepositDataRoot[:], msg.Root) {
		return errors.New("deposit data roots not equal")
	}

	index := -1
	for i, d := range s.InitMsg.OperatorIDs {
		if d == msg.Signer {
			index = i
		}
	}

	if index == -1 {
		return errors.New("signer not part of committee")
	}

	// find operator and verify msg
	sharePkBytes := s.key.SharePublicKeys[index]
	sharePk := &bls.PublicKey{} // TODO: cache this PubKey
	if err := sharePk.Deserialize(sharePkBytes); err != nil {
		return errors.Wrap(err, "could not deserialize public key")
	}

	sig := &bls.Sign{}
	if err := sig.Deserialize(msg.Signature); err != nil {
		return errors.Wrap(err, "could not deserialize partial sig")
	}
	if !sig.VerifyByte(sharePk, s.DepositDataRoot[:]) {
		return errors.New("partial deposit data sig invalid")
	}

	return nil
}
