package sign

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type SignDepositData struct {
	// Identifier unique for DKG session
	Identifier dkgtypes.RequestID
	// InitMsg holds the details of this instance
	InitMsg *dkgtypes.Init
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	DepositDataRoot   spec.Root
	output            *dkgtypes.SignedDepositDataMsgBody
	key               *dkgtypes.LocalKeyShare
	config            dkgtypes.ProtocolConfig
}

func (s *SignDepositData) Output() ([]byte, error) {
	if s.output == nil {
		return nil, nil
	}
	return s.output.Encode()
}

func NewSignDepositData(init *dkgtypes.Init, key *dkgtypes.LocalKeyShare, config dkgtypes.ProtocolConfig) *SignDepositData {

	return &SignDepositData{
		Identifier:        dkgtypes.RequestID{},
		InitMsg:           init,
		PartialSignatures: map[types.OperatorID][]byte{},
		DepositDataRoot:   spec.Root{},
		output:            nil,
		key:               key,
		config:            config,
	}
}

func (s *SignDepositData) Start() ([]dkgtypes.Message, error) {
	pSig, err := s.partialSign()
	if err != nil {
		return nil, err
	}
	//data, err := pSig.Encode()
	//if err != nil {
	//	return nil, err
	//}
	s.PartialSignatures[s.config.Operator.OperatorID] = pSig.Signature
	partialSigMsg := dkgtypes.ParsedPartialSigMessage{
		Header: &dkgtypes.MessageHeader{
			SessionId: s.Identifier[:],
			MsgType:   int32(dkgtypes.DepositDataMsgType),
			Sender:    uint64(s.config.Operator.OperatorID),
			Receiver:  0,
		},
		Body: pSig,
	}
	base, err := partialSigMsg.ToBase()
	if err != nil {
		return nil, err
	}
	return []dkgtypes.Message{*base}, nil
}

func (s *SignDepositData) ProcessMsg(msg *dkgtypes.Message) ([]dkgtypes.Message, error) {

	if msg.Header.MsgType != int32(dkgtypes.DepositDataMsgType) {
		return nil, errors.New("invalid message type")
	}

	pMsg := &dkgtypes.ParsedPartialSigMessage{}
	err := pMsg.FromBase(msg)
	if err != nil {
		return nil, err
	}
	if err = s.validateDepositDataSig(pMsg.Body); err != nil {
		return nil, errors.Wrap(err, "PartialSigMsgBody invalid")
	}

	if found := s.PartialSignatures[types.OperatorID(pMsg.Header.Sender)]; found == nil {
		s.PartialSignatures[types.OperatorID(pMsg.Header.Sender)] = pMsg.Body.Signature
	} else if bytes.Compare(found, pMsg.Body.Signature) != 0 {
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

		s.output = &dkgtypes.SignedDepositDataMsgBody{
			RequestID:             s.Identifier[:],
			OperatorID:            uint64(s.config.Operator.OperatorID),
			EncryptedShare:        encryptedShare,
			Committee:             s.key.Committee,
			Threshold:             s.InitMsg.Threshold,
			ValidatorPublicKey:    s.key.PublicKey,
			WithdrawalCredentials: s.InitMsg.WithdrawalCredentials,
			DepositDataSignature:  sig.Serialize(),
		}
		return nil, nil
	}
	return nil, nil
}

//func (s *SignDepositData) generateSignedOutput(o *dkgtypes.Output) (*dkgtypes.SignedOutput, error) {
//	sig, err := s.config.Signer.SignDKGOutput(o, s.config.Operator.ETHAddress)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not sign output")
//	}
//
//	return &dkgtypes.SignedOutput{
//		Data:      o,
//		Signer:    s.config.Operator.OperatorID,
//		Signature: sig,
//	}, nil
//}

func (s *SignDepositData) partialSign() (*dkgtypes.PartialSigMsgBody, error) {
	share := bls.SecretKey{}
	err := share.Deserialize(s.key.SecretShare)
	if err != nil {
		return nil, err
	}

	root, depData, err := types.GenerateETHDepositData(
		s.key.PublicKey,
		s.InitMsg.WithdrawalCredentials,
		s.InitMsg.Fork,
		types.DomainDeposit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate deposit data")
	}

	copy(s.DepositDataRoot[:], root)
	//root := make([]byte, len(depData.DepositDataRoot))
	//copy(root, depData.DepositDataRoot[:])
	rawSig := share.SignByte(root[:])
	sigBytes := rawSig.Serialize()
	var sig spec.BLSSignature
	copy(sig[:], sigBytes)

	copy(depData.DepositData.Signature[:], sigBytes)

	return &dkgtypes.PartialSigMsgBody{
		Signer:    uint64(s.config.Operator.OperatorID),
		Root:      root,
		Signature: sig[:],
	}, nil
}

func (s *SignDepositData) validateDepositDataSig(msg *dkgtypes.PartialSigMsgBody) error {
	if !bytes.Equal(s.DepositDataRoot[:], msg.Root) {
		return errors.New("deposit data roots not equal")
	}

	index := -1
	for i, d := range s.InitMsg.OperatorIDs {
		if d == types.OperatorID(msg.Signer) {
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

	root := make([]byte, 32)
	copy(root, s.DepositDataRoot[:])
	if !sig.VerifyByte(sharePk, root) {
		return errors.New("partial deposit data sig invalid")
	}

	return nil
}
