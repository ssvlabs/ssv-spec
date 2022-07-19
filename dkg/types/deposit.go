package types

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

func (s *SignDepositData) Output() ([]byte, error) {
	panic("implement me")
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
		Header: &MessageHeader{
			SessionId: s.Identifier[:],
			MsgType:   int32(DepositDataMsgType),
			Sender: uint64(s.config.Operator.OperatorID),
			Receiver:  0,
		},
		Data:       data,
	}
	return []Message{partialSigMsg}, nil
}

func (s *SignDepositData) ProcessMsg(msg *Message) ([]Message, error) {

	if msg.Header.MsgType != int32(DepositDataMsgType) {
		return nil, errors.New("invalid message type")
	}

	pMsg := &PartialDepositData{}
	err := pMsg.Decode(msg.Data)
	if err != nil {
		return nil, err
	}
	if err = s.validateDepositDataSig(pMsg); err != nil {
		return nil, errors.Wrap(err, "PartialDepositData invalid")
	}

	if found := s.PartialSignatures[pMsg.Signer]; found == nil {
		s.PartialSignatures[pMsg.Signer] = pMsg.Signature
	} else if bytes.Compare(found, pMsg.Signature) != 0 {
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
			Header: &MessageHeader{
				SessionId: s.Identifier[:],
				MsgType:   int32(PartialOutputMsgType),
				Sender: uint64(s.config.Operator.OperatorID),
				Receiver:  0,
			},
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

func (s *SignDepositData) partialSign() (*PartialDepositData, error) {
	share := bls.SecretKey{}
	err := share.Deserialize(s.key.SecretShare)
	if err != nil {
		return nil, err
	}

	root, _, err := types.GenerateETHDepositData(
		s.key.PublicKey,
		s.InitMsg.WithdrawalCredentials,
		s.InitMsg.Fork,
		types.DomainDeposit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate deposit data")
	}
	rawSig := share.SignByte(root[:])
	sigBytes := rawSig.Serialize()
	var sig spec.BLSSignature
	copy(sig[:], sigBytes)
	return &PartialDepositData{
		Signer:    s.config.Operator.OperatorID,
		Root:      root,
		Signature: sig[:],
	}, nil
}

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
