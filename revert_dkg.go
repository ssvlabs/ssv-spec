diff --git a/dkg/README.md b/dkg/README.md
new file mode 100644
index 0000000..3b12db4
--- /dev/null
+++ b/dkg/README.md
@@ -0,0 +1,10 @@
+
+# DKG
+
+## Introduction
+This is a spec implementation for a generalized DKG protocol for SSV.Network following [SIP - DKG](https://docs.google.com/document/d/1TRVUHjFyxINWW2H9FYLNL2pQoLy6gmvaI62KL_4cREQ/edit).
+
+## TODO
+- [//] Generalized message processing flow
+- [ ] spec tests
+- [ ] specific dkg implementation
\ No newline at end of file
diff --git a/dkg/frost/frost.go b/dkg/frost/frost.go
new file mode 100644
index 0000000..e279b62
--- /dev/null
+++ b/dkg/frost/frost.go
@@ -0,0 +1,385 @@
+package frost
+
+import (
+	"math/rand"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/coinbase/kryptology/pkg/core/curves"
+	"github.com/coinbase/kryptology/pkg/dkg/frost"
+	ecies "github.com/ecies/go/v2"
+	"github.com/ethereum/go-ethereum/common"
+	"github.com/ethereum/go-ethereum/crypto"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+)
+
+var thisCurve = curves.BLS12381G1()
+
+func init() {
+	types.InitBLS()
+}
+
+type FROST struct {
+	network dkg.Network
+	signer  types.DKGSigner
+	storage dkg.Storage
+
+	state *State
+}
+
+type State struct {
+	identifier dkg.RequestID
+	operatorID types.OperatorID
+	sessionSK  *ecies.PrivateKey
+
+	threshold    uint32
+	currentRound ProtocolRound
+	participant  *frost.DkgParticipant
+
+	operators      []uint32
+	operatorsOld   []uint32
+	operatorShares map[uint32]*bls.SecretKey
+
+	msgs            ProtocolMessageStore
+	oldKeyGenOutput *dkg.KeyGenOutput
+}
+
+type ProtocolRound int
+
+const (
+	Uninitialized ProtocolRound = iota
+	Preparation
+	Round1
+	Round2
+	KeygenOutput
+	Blame
+)
+
+var rounds = []ProtocolRound{
+	Uninitialized,
+	Preparation,
+	Round1,
+	Round2,
+	KeygenOutput,
+	Blame,
+}
+
+type ProtocolMessageStore map[ProtocolRound]map[uint32]*dkg.SignedMessage
+
+func newProtocolMessageStore() ProtocolMessageStore {
+	m := make(map[ProtocolRound]map[uint32]*dkg.SignedMessage)
+	for _, round := range rounds {
+		m[round] = make(map[uint32]*dkg.SignedMessage)
+	}
+	return m
+}
+
+func New(
+	network dkg.Network,
+	operatorID types.OperatorID,
+	requestID dkg.RequestID,
+	signer types.DKGSigner,
+	storage dkg.Storage,
+	init *dkg.Init,
+) dkg.Protocol {
+
+	fr := &FROST{
+		network: network,
+		signer:  signer,
+		storage: storage,
+		state: &State{
+			identifier:     requestID,
+			operatorID:     operatorID,
+			threshold:      uint32(init.Threshold),
+			currentRound:   Uninitialized,
+			operators:      types.OperatorList(init.OperatorIDs).ToUint32List(),
+			operatorShares: make(map[uint32]*bls.SecretKey),
+			msgs:           newProtocolMessageStore(),
+		},
+	}
+
+	return fr
+}
+
+// Temporary, TODO: Remove and use interface with Reshare
+func NewResharing(
+	network dkg.Network,
+	operatorID types.OperatorID,
+	requestID dkg.RequestID,
+	signer types.DKGSigner,
+	storage dkg.Storage,
+	operatorsOld []types.OperatorID,
+	init *dkg.Reshare,
+	output *dkg.KeyGenOutput,
+) dkg.Protocol {
+
+	return &FROST{
+		network: network,
+		signer:  signer,
+		storage: storage,
+
+		state: &State{
+			identifier:      requestID,
+			operatorID:      operatorID,
+			threshold:       uint32(init.Threshold),
+			currentRound:    Uninitialized,
+			operators:       types.OperatorList(init.OperatorIDs).ToUint32List(),
+			operatorsOld:    types.OperatorList(operatorsOld).ToUint32List(),
+			operatorShares:  make(map[uint32]*bls.SecretKey),
+			msgs:            newProtocolMessageStore(),
+			oldKeyGenOutput: output,
+		},
+	}
+}
+
+// TODO: If Reshare, confirm participating operators using qbft before kick-starting this process.
+func (fr *FROST) Start() error {
+	fr.state.currentRound = Preparation
+
+	ctx := make([]byte, 16)
+	if _, err := rand.Read(ctx); err != nil {
+		return err
+	}
+	participant, err := frost.NewDkgParticipant(uint32(fr.state.operatorID), fr.state.threshold, string(ctx), thisCurve, fr.state.operators...)
+	if err != nil {
+		return errors.Wrap(err, "failed to initialize a dkg participant")
+	}
+	fr.state.participant = participant
+
+	if !fr.needToRunCurrentRound() {
+		return nil
+	}
+
+	k, err := ecies.GenerateKey()
+	if err != nil {
+		return errors.Wrap(err, "failed to generate session sk")
+	}
+	fr.state.sessionSK = k
+
+	msg := &ProtocolMsg{
+		Round: fr.state.currentRound,
+		PreparationMessage: &PreparationMessage{
+			SessionPk: k.PublicKey.Bytes(true),
+		},
+	}
+	_, err = fr.broadcastDKGMessage(msg)
+	return err
+}
+
+func (fr *FROST) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
+
+	if err := fr.validateSignedMessage(msg); err != nil {
+		return false, nil, errors.Wrap(err, "failed to Validate signed message")
+	}
+
+	protocolMessage := &ProtocolMsg{}
+	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
+		return false, nil, errors.Wrap(err, "failed to decode protocol msg")
+	}
+	if err := protocolMessage.Validate(); err != nil {
+		fr.state.currentRound = Blame
+		outcome, err := fr.createAndBroadcastBlameOfInvalidMessage(uint32(msg.Signer), msg)
+		return true, outcome, err
+	}
+
+	existingMessage, ok := fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)]
+
+	if isBlameTypeInconsisstent := ok && !fr.haveSameRoot(existingMessage, msg); isBlameTypeInconsisstent {
+		fr.state.currentRound = Blame
+		outcome, err := fr.createAndBroadcastBlameOfInconsistentMessage(existingMessage, msg)
+		if err != nil {
+			return false, nil, err
+		}
+		return true, outcome, nil
+
+	}
+
+	if protocolMessage.Round == Blame {
+		fr.state.currentRound = Blame
+		valid, err := fr.checkBlame(uint32(msg.Signer), protocolMessage)
+		if err != nil {
+			return false, nil, err
+		}
+		return valid, &dkg.ProtocolOutcome{
+			BlameOutput: &dkg.BlameOutput{
+				Valid:        valid,
+				BlameMessage: msg,
+			},
+		}, nil
+	}
+
+	fr.state.msgs[protocolMessage.Round][uint32(msg.Signer)] = msg
+
+	switch fr.state.currentRound {
+	case Preparation:
+		if fr.canProceedThisRound() {
+			fr.state.currentRound = Round1
+			if err := fr.processRound1(); err != nil {
+				return false, nil, err
+			}
+		}
+	case Round1:
+		if fr.canProceedThisRound() {
+			fr.state.currentRound = Round2
+			outcome, err := fr.processRound2()
+			return outcome != nil, outcome, err
+		}
+	case Round2:
+		if fr.canProceedThisRound() {
+			fr.state.currentRound = KeygenOutput
+			out, err := fr.processKeygenOutput()
+			if err != nil {
+				return false, nil, err
+			}
+			return true, &dkg.ProtocolOutcome{ProtocolOutput: out}, nil
+		}
+	default:
+		return true, nil, dkg.ErrInvalidRound{}
+	}
+
+	return false, nil, nil
+}
+
+func (fr *FROST) canProceedThisRound() bool {
+	// Note: for Resharing, Preparation (New Committee) -> Round1 (Old Committee) -> Round2 (New Committee)
+	if fr.isResharing() && fr.state.currentRound == Round1 {
+		return fr.allMessagesReceivedFor(Round1, fr.state.operatorsOld)
+	}
+	return fr.allMessagesReceivedFor(fr.state.currentRound, fr.state.operators)
+}
+
+func (fr *FROST) allMessagesReceivedFor(round ProtocolRound, operators []uint32) bool {
+	for _, operatorID := range operators {
+		if _, ok := fr.state.msgs[round][operatorID]; !ok {
+			return false
+		}
+	}
+	return true
+}
+
+func (fr *FROST) isResharing() bool {
+	return len(fr.state.operatorsOld) > 0
+}
+
+func (fr *FROST) inOldCommittee() bool {
+	for _, id := range fr.state.operatorsOld {
+		if types.OperatorID(id) == fr.state.operatorID {
+			return true
+		}
+	}
+	return false
+}
+
+func (fr *FROST) inNewCommittee() bool {
+	for _, id := range fr.state.operators {
+		if types.OperatorID(id) == fr.state.operatorID {
+			return true
+		}
+	}
+	return false
+}
+
+func (fr *FROST) needToRunCurrentRound() bool {
+	if !fr.isResharing() {
+		return true // always run for new keygen
+	}
+	switch fr.state.currentRound {
+	case Preparation, Round2, KeygenOutput:
+		return fr.inNewCommittee()
+	case Round1:
+		return fr.inOldCommittee()
+	default:
+		return false
+	}
+}
+
+func (fr *FROST) validateSignedMessage(msg *dkg.SignedMessage) error {
+	if msg.Message.Identifier != fr.state.identifier {
+		return errors.New("got mismatching identifier")
+	}
+
+	found, operator, err := fr.storage.GetDKGOperator(msg.Signer)
+	if !found {
+		return errors.New("unable to find signer")
+	}
+	if err != nil {
+		return errors.Wrap(err, "unable to find signer")
+	}
+
+	root, err := msg.Message.GetRoot()
+	if err != nil {
+		return errors.Wrap(err, "failed to get root")
+	}
+
+	pk, err := crypto.Ecrecover(root, msg.Signature)
+	if err != nil {
+		return errors.Wrap(err, "unable to recover public key")
+	}
+
+	addr := common.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
+	if addr != operator.ETHAddress {
+		return errors.New("invalid signature")
+	}
+	return nil
+}
+
+func (fr *FROST) encryptByOperatorID(operatorID uint32, data []byte) ([]byte, error) {
+	msg, ok := fr.state.msgs[Preparation][operatorID]
+	if !ok {
+		return nil, errors.New("no session pk found for the operator")
+	}
+
+	protocolMessage := &ProtocolMsg{}
+	if err := protocolMessage.Decode(msg.Message.Data); err != nil {
+		return nil, errors.Wrap(err, "failed to decode protocol msg")
+	}
+
+	sessionPK, err := ecies.NewPublicKeyFromBytes(protocolMessage.PreparationMessage.SessionPk)
+	if err != nil {
+		return nil, err
+	}
+
+	return ecies.Encrypt(sessionPK, data)
+}
+
+func (fr *FROST) toSignedMessage(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
+	msgBytes, err := msg.Encode()
+	if err != nil {
+		return nil, err
+	}
+
+	bcastMessage := &dkg.SignedMessage{
+		Message: &dkg.Message{
+			MsgType:    dkg.ProtocolMsgType,
+			Identifier: fr.state.identifier,
+			Data:       msgBytes,
+		},
+		Signer: fr.state.operatorID,
+	}
+
+	exist, operator, err := fr.storage.GetDKGOperator(fr.state.operatorID)
+	if err != nil {
+		return nil, err
+	}
+	if !exist {
+		return nil, errors.Errorf("operator with id %d not found", fr.state.operatorID)
+	}
+
+	sig, err := fr.signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
+	if err != nil {
+		return nil, err
+	}
+	bcastMessage.Signature = sig
+	return bcastMessage, nil
+}
+
+func (fr *FROST) broadcastDKGMessage(msg *ProtocolMsg) (*dkg.SignedMessage, error) {
+	bcastMessage, err := fr.toSignedMessage(msg)
+	if err != nil {
+		return bcastMessage, err
+	}
+	fr.state.msgs[fr.state.currentRound][uint32(fr.state.operatorID)] = bcastMessage
+	err = fr.network.BroadcastDKGMessage(bcastMessage)
+	return bcastMessage, err
+}
diff --git a/dkg/frost/frost_blame.go b/dkg/frost/frost_blame.go
new file mode 100644
index 0000000..70380a0
--- /dev/null
+++ b/dkg/frost/frost_blame.go
@@ -0,0 +1,246 @@
+package frost
+
+import (
+	"bytes"
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/coinbase/kryptology/pkg/sharing"
+	ecies "github.com/ecies/go/v2"
+	"github.com/pkg/errors"
+)
+
+func (fr *FROST) checkBlame(blamerOID uint32, protocolMessage *ProtocolMsg) (bool, error) {
+	switch protocolMessage.BlameMessage.Type {
+	case InvalidShare:
+		return fr.processBlameTypeInvalidShare(blamerOID, protocolMessage.BlameMessage)
+	case InconsistentMessage:
+		return fr.processBlameTypeInconsistentMessage(protocolMessage.BlameMessage)
+	case InvalidMessage:
+		return fr.processBlameTypeInvalidMessage(protocolMessage.BlameMessage)
+	default:
+		return false, errors.New("unrecognized blame type")
+	}
+}
+
+func (fr *FROST) processBlameTypeInvalidShare(blamerOID uint32, blameMessage *BlameMessage) (bool /*valid*/, error) {
+	if err := blameMessage.Validate(); err != nil {
+		return false, errors.Wrap(err, "invalid blame message")
+	}
+	if len(blameMessage.BlameData) != 1 {
+		return false, errors.New("invalid blame data")
+	}
+	signedMessage, protocolMessage, err := fr.decodeMessage(blameMessage.BlameData[0])
+	if err != nil {
+		return false, errors.Wrap(err, "failed to decode signed message")
+	}
+
+	if err := fr.validateSignedMessage(signedMessage); err != nil {
+		return false, errors.Wrap(err, "failed to Validate signature for blame data")
+	}
+
+	round1Message := protocolMessage.Round1Message
+
+	blamerPrepSignedMessage := fr.state.msgs[Preparation][blamerOID]
+	blamerPrepProtocolMessage := &ProtocolMsg{}
+	err = blamerPrepProtocolMessage.Decode(blamerPrepSignedMessage.Message.Data)
+	if err != nil || blamerPrepProtocolMessage.PreparationMessage == nil {
+		return false, errors.New("unable to decode blamer's PreparationMessage")
+	}
+
+	blamerSessionSK := ecies.NewPrivateKeyFromBytes(blameMessage.BlamerSessionSk)
+	blamerSessionPK := blamerSessionSK.PublicKey.Bytes(true)
+	if !bytes.Equal(blamerSessionPK, blamerPrepProtocolMessage.PreparationMessage.SessionPk) {
+		return false, errors.New("blame's session pubkey is invalid")
+	}
+
+	verifiers := new(sharing.FeldmanVerifier)
+	for _, commitmentBytes := range round1Message.Commitment {
+		commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
+		if err != nil {
+			return false, err
+		}
+		verifiers.Commitments = append(verifiers.Commitments, commitment)
+	}
+
+	shareBytes, err := ecies.Decrypt(blamerSessionSK, round1Message.Shares[blamerOID])
+	if err != nil {
+		return true, nil
+	}
+
+	share := &sharing.ShamirShare{
+		Id:    blamerOID,
+		Value: shareBytes,
+	}
+
+	if err = verifiers.Verify(share); err != nil {
+		return true, nil
+	}
+	return false, err
+}
+
+func (fr *FROST) decodeMessage(data []byte) (*dkg.SignedMessage, *ProtocolMsg, error) {
+	signedMsg := &dkg.SignedMessage{}
+	if err := signedMsg.Decode(data); err != nil {
+		return nil, nil, errors.Wrap(err, "failed to decode signed message")
+	}
+	pMsg := &ProtocolMsg{}
+	if err := pMsg.Decode(signedMsg.Message.Data); err != nil {
+		return signedMsg, nil, errors.Wrap(err, "failed to decode protocol msg")
+	}
+	return signedMsg, pMsg, nil
+}
+
+func (fr *FROST) processBlameTypeInconsistentMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
+	if err := blameMessage.Validate(); err != nil {
+		return false, errors.Wrap(err, "invalid blame message")
+	}
+
+	if len(blameMessage.BlameData) != 2 {
+		return false, errors.New("invalid blame data")
+	}
+
+	signedMsg1, protocolMessage1, err := fr.decodeMessage(blameMessage.BlameData[0])
+
+	if err != nil {
+		return false, err
+	} else if err := fr.validateSignedMessage(signedMsg1); err != nil {
+		return false, errors.Wrap(err, "failed to validate signed message in blame data")
+	} else if err := protocolMessage1.Validate(); err != nil {
+		return false, errors.New("invalid protocol message")
+	}
+
+	signedMsg2, protocolMessage2, err := fr.decodeMessage(blameMessage.BlameData[1])
+
+	if err != nil {
+		return false, err
+	} else if err := fr.validateSignedMessage(signedMsg2); err != nil {
+		return false, errors.Wrap(err, "failed to validate signed message in blame data")
+	} else if err := protocolMessage2.Validate(); err != nil {
+		return false, errors.New("invalid protocol message")
+	}
+
+	if fr.haveSameRoot(signedMsg1, signedMsg2) {
+		return false, errors.New("the two messages are consistent")
+	}
+
+	if protocolMessage1.Round != protocolMessage2.Round {
+		return false, errors.New("the two messages don't belong the the same round")
+	}
+
+	return true, nil
+}
+
+func (fr *FROST) processBlameTypeInvalidMessage(blameMessage *BlameMessage) (bool /*valid*/, error) {
+	if err := blameMessage.Validate(); err != nil {
+		return false, errors.Wrap(err, "invalid blame message")
+	}
+	if len(blameMessage.BlameData) != 1 {
+		return false, errors.New("invalid blame data")
+	}
+	signedMsg, pMsg, err := fr.decodeMessage(blameMessage.BlameData[0])
+	if err != nil {
+		return false, err
+	} else if err := fr.validateSignedMessage(signedMsg); err != nil {
+		return false, errors.Wrap(err, "failed to validate signed message in blame data")
+	}
+
+	err = pMsg.Validate()
+	if err != nil {
+		return true, nil
+	}
+	return false, errors.New("message is valid")
+}
+
+func (fr *FROST) createAndBroadcastBlameOfInconsistentMessage(existingMessage, newMessage *dkg.SignedMessage) (*dkg.ProtocolOutcome, error) {
+	existingMessageBytes, err := existingMessage.Encode()
+	if err != nil {
+		return nil, err
+	}
+	newMessageBytes, err := newMessage.Encode()
+	if err != nil {
+		return nil, err
+	}
+	msg := &ProtocolMsg{
+		Round: Blame,
+		BlameMessage: &BlameMessage{
+			Type:             InconsistentMessage,
+			TargetOperatorID: uint32(newMessage.Signer),
+			BlameData:        [][]byte{existingMessageBytes, newMessageBytes},
+			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
+		},
+	}
+	signedMessage, err := fr.broadcastDKGMessage(msg)
+	return &dkg.ProtocolOutcome{
+		BlameOutput: &dkg.BlameOutput{
+			Valid:        true,
+			BlameMessage: signedMessage,
+		},
+	}, err
+}
+
+func (fr *FROST) createAndBroadcastBlameOfInvalidShare(culpritOID uint32) (*dkg.ProtocolOutcome, error) {
+	round1Bytes, err := fr.state.msgs[Round1][culpritOID].Encode()
+	if err != nil {
+		return nil, err
+	}
+	msg := &ProtocolMsg{
+		Round: Blame,
+		BlameMessage: &BlameMessage{
+			Type:             InvalidShare,
+			TargetOperatorID: culpritOID,
+			BlameData:        [][]byte{round1Bytes},
+			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
+		},
+	}
+	signedMessage, err := fr.broadcastDKGMessage(msg)
+	return &dkg.ProtocolOutcome{
+		BlameOutput: &dkg.BlameOutput{
+			Valid:        true,
+			BlameMessage: signedMessage,
+		},
+	}, err
+}
+
+func (fr *FROST) createAndBroadcastBlameOfInvalidMessage(culpritOID uint32, message *dkg.SignedMessage) (*dkg.ProtocolOutcome, error) {
+	bytes, err := message.Encode()
+	if err != nil {
+		return nil, err
+	}
+
+	msg := &ProtocolMsg{
+		Round: Blame,
+		BlameMessage: &BlameMessage{
+			Type:             InvalidMessage,
+			TargetOperatorID: culpritOID,
+			BlameData:        [][]byte{bytes},
+			BlamerSessionSk:  fr.state.sessionSK.Bytes(),
+		},
+	}
+	signedMsg, err := fr.broadcastDKGMessage(msg)
+
+	return &dkg.ProtocolOutcome{
+		BlameOutput: &dkg.BlameOutput{
+			Valid:        true,
+			BlameMessage: signedMsg,
+		},
+	}, err
+}
+
+func (fr *FROST) haveSameRoot(existingMessage, newMessage *dkg.SignedMessage) bool {
+	r1, err := existingMessage.GetRoot()
+	if err != nil {
+		return false
+	}
+	r2, err := newMessage.GetRoot()
+	if err != nil {
+		return false
+	}
+	return bytes.Equal(r1, r2)
+}
+
+type ErrBlame struct {
+	BlameOutput *dkg.BlameOutput
+}
+
+func (e ErrBlame) Error() string {
+	return "detected and processed blame"
+}
diff --git a/dkg/frost/frost_keygen_output.go b/dkg/frost/frost_keygen_output.go
new file mode 100644
index 0000000..1e44487
--- /dev/null
+++ b/dkg/frost/frost_keygen_output.go
@@ -0,0 +1,138 @@
+package frost
+
+import (
+	"bytes"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+)
+
+func (fr *FROST) processKeygenOutput() (*dkg.KeyGenOutput, error) {
+
+	if fr.state.currentRound != KeygenOutput {
+		return nil, dkg.ErrInvalidRound{}
+	}
+
+	if !fr.needToRunCurrentRound() {
+		return nil, nil
+	}
+
+	reconstructed, err := fr.verifyShares()
+	if err != nil {
+		return nil, errors.Wrap(err, "failed to verify shares")
+	}
+
+	reconstructedBytes := reconstructed.Serialize()
+
+	out := &dkg.KeyGenOutput{
+		Threshold: uint64(fr.state.threshold),
+	}
+
+	operatorPubKeys := make(map[types.OperatorID]*bls.PublicKey)
+	for _, operatorID := range fr.state.operators {
+		protocolMessage := &ProtocolMsg{}
+		if err := protocolMessage.Decode(fr.state.msgs[Round2][operatorID].Message.Data); err != nil {
+			return nil, errors.Wrap(err, "failed to decode protocol msg")
+		}
+
+		if operatorID == uint32(fr.state.operatorID) {
+			sk := &bls.SecretKey{}
+			if err := sk.Deserialize(fr.state.participant.SkShare.Bytes()); err != nil {
+				return nil, err
+			}
+
+			out.Share = sk
+			out.ValidatorPK = protocolMessage.Round2Message.Vk
+		}
+
+		pk := &bls.PublicKey{}
+		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
+			return nil, err
+		}
+
+		operatorPubKeys[types.OperatorID(operatorID)] = pk
+	}
+
+	out.OperatorPubKeys = operatorPubKeys
+
+	if !bytes.Equal(out.ValidatorPK, reconstructedBytes) {
+		return nil, errors.New("can't reconstruct to the validator pk")
+	}
+
+	return out, nil
+}
+
+func (fr *FROST) verifyShares() (*bls.G1, error) {
+
+	var (
+		quorumStart       = 0
+		quorumEnd         = int(fr.state.threshold)
+		prevReconstructed = (*bls.G1)(nil)
+	)
+
+	// Sliding window of quorum 0...threshold until n-threshold...n
+	for quorumEnd < len(fr.state.operators) {
+		quorum := fr.state.operators[quorumStart:quorumEnd]
+		currReconstructed, err := fr.verifyShare(quorum)
+		if err != nil {
+			return nil, err
+		}
+		if prevReconstructed != nil && !currReconstructed.IsEqual(prevReconstructed) {
+			return nil, errors.New("failed to create consistent public key from tshares")
+		}
+		prevReconstructed = currReconstructed
+		quorumStart++
+		quorumEnd++
+	}
+	return prevReconstructed, nil
+}
+
+func (fr *FROST) verifyShare(operators []uint32) (*bls.G1, error) {
+	xVec, err := fr.getXVec(operators)
+	if err != nil {
+		return nil, err
+	}
+
+	yVec, err := fr.getYVec(operators)
+	if err != nil {
+		return nil, err
+	}
+
+	reconstructed := &bls.G1{}
+	if err := bls.G1LagrangeInterpolation(reconstructed, xVec, yVec); err != nil {
+		return nil, err
+	}
+	return reconstructed, nil
+}
+
+func (fr *FROST) getXVec(operators []uint32) ([]bls.Fr, error) {
+	xVec := make([]bls.Fr, 0)
+	for _, operator := range operators {
+		x := bls.Fr{}
+		x.SetInt64(int64(operator))
+		xVec = append(xVec, x)
+	}
+	return xVec, nil
+}
+
+func (fr *FROST) getYVec(operators []uint32) ([]bls.G1, error) {
+	yVec := make([]bls.G1, 0)
+	for _, operator := range operators {
+
+		protocolMessage := &ProtocolMsg{}
+		if err := protocolMessage.Decode(fr.state.msgs[Round2][operator].Message.Data); err != nil {
+			return nil, errors.Wrap(err, "failed to decode protocol msg")
+		}
+
+		pk := &bls.PublicKey{}
+		if err := pk.Deserialize(protocolMessage.Round2Message.VkShare); err != nil {
+			return nil, errors.Wrap(err, "failed to deserialize public key")
+		}
+
+		y := bls.CastFromPublicKey(pk)
+		yVec = append(yVec, *y)
+	}
+	return yVec, nil
+}
diff --git a/dkg/frost/frost_round1.go b/dkg/frost/frost_round1.go
new file mode 100644
index 0000000..5b16ae1
--- /dev/null
+++ b/dkg/frost/frost_round1.go
@@ -0,0 +1,103 @@
+package frost
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+)
+
+func (fr *FROST) processRound1() error {
+
+	if fr.state.currentRound != Round1 {
+		return dkg.ErrInvalidRound{}
+	}
+
+	if !fr.needToRunCurrentRound() {
+		return fr.state.participant.SkipRound1()
+	}
+
+	var (
+		skI []byte // secret to be shared, nil if new keygen, lagrange interpolation of own part of secret if resharing
+		err error
+	)
+
+	if fr.isResharing() {
+		skI, err = fr.partialInterpolate()
+		if err != nil {
+			return err
+		}
+	}
+
+	bCastMessage, p2pMessages, err := fr.state.participant.Round1(skI)
+	if err != nil {
+		return err
+	}
+
+	// get bytes representation of commitment points
+	commitments := make([][]byte, 0)
+	for _, commitment := range bCastMessage.Verifiers.Commitments {
+		commitments = append(commitments, commitment.ToAffineCompressed())
+	}
+
+	// encrypted shares by operators
+	shares := make(map[uint32][]byte)
+	for _, operatorID := range fr.state.operators {
+		if uint32(fr.state.operatorID) == operatorID {
+			continue
+		}
+
+		share := &bls.SecretKey{}
+		shamirShare := p2pMessages[operatorID]
+
+		if err := share.Deserialize(shamirShare.Value); err != nil {
+			return err
+		}
+
+		fr.state.operatorShares[operatorID] = share
+
+		encryptedShare, err := fr.encryptByOperatorID(operatorID, shamirShare.Value)
+		if err != nil {
+			return err
+		}
+		shares[operatorID] = encryptedShare
+	}
+
+	msg := &ProtocolMsg{
+		Round: Round1,
+		Round1Message: &Round1Message{
+			Commitment: commitments,
+			ProofS:     bCastMessage.Wi.Bytes(),
+			ProofR:     bCastMessage.Ci.Bytes(),
+			Shares:     shares,
+		},
+	}
+	_, err = fr.broadcastDKGMessage(msg)
+	return err
+}
+
+func (fr *FROST) partialInterpolate() ([]byte, error) {
+	if !fr.isResharing() {
+		return nil, nil
+	}
+
+	skI := new(bls.Fr)
+
+	indices := make([]bls.Fr, fr.state.oldKeyGenOutput.Threshold)
+	values := make([]bls.Fr, fr.state.oldKeyGenOutput.Threshold)
+	for i, id := range fr.state.operatorsOld {
+		(&indices[i]).SetInt64(int64(id))
+		if types.OperatorID(id) == fr.state.operatorID {
+			err := (&values[i]).Deserialize(fr.state.oldKeyGenOutput.Share.Serialize())
+			if err != nil {
+				return nil, err
+			}
+		} else {
+			(&values[i]).SetInt64(0)
+		}
+	}
+
+	if err := bls.FrLagrangeInterpolation(skI, indices, values); err != nil {
+		return nil, err
+	}
+	return skI.Serialize(), nil
+}
diff --git a/dkg/frost/frost_round2.go b/dkg/frost/frost_round2.go
new file mode 100644
index 0000000..f915621
--- /dev/null
+++ b/dkg/frost/frost_round2.go
@@ -0,0 +1,94 @@
+package frost
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/coinbase/kryptology/pkg/dkg/frost"
+	"github.com/coinbase/kryptology/pkg/sharing"
+	ecies "github.com/ecies/go/v2"
+	"github.com/pkg/errors"
+)
+
+func (fr *FROST) processRound2() (*dkg.ProtocolOutcome, error) {
+
+	if fr.state.currentRound != Round2 {
+		return nil, dkg.ErrInvalidRound{}
+	}
+
+	if !fr.needToRunCurrentRound() {
+		return nil, nil
+	}
+
+	bcast := make(map[uint32]*frost.Round1Bcast)
+	p2psend := make(map[uint32]*sharing.ShamirShare)
+
+	for peerOID, dkgMessage := range fr.state.msgs[Round1] {
+
+		protocolMessage := &ProtocolMsg{}
+		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
+			return nil, errors.Wrap(err, "failed to decode protocol msg")
+		}
+		verifiers := new(sharing.FeldmanVerifier)
+		for _, commitmentBytes := range protocolMessage.Round1Message.Commitment {
+			commitment, err := thisCurve.Point.FromAffineCompressed(commitmentBytes)
+			if err != nil {
+				return nil, errors.Wrap(err, "failed to decode commitment point")
+			}
+			verifiers.Commitments = append(verifiers.Commitments, commitment)
+		}
+
+		Wi, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofS)
+		if err != nil {
+			return nil, errors.Wrap(err, "failed to decode scalar")
+		}
+		Ci, err := thisCurve.Scalar.SetBytes(protocolMessage.Round1Message.ProofR)
+		if err != nil {
+			return nil, errors.Wrap(err, "failed to decode scalar")
+		}
+
+		bcastMessage := &frost.Round1Bcast{
+			Verifiers: verifiers,
+			Wi:        Wi,
+			Ci:        Ci,
+		}
+		bcast[peerOID] = bcastMessage
+
+		if uint32(fr.state.operatorID) == peerOID {
+			continue
+		}
+
+		encryptedShare := protocolMessage.Round1Message.Shares[uint32(fr.state.operatorID)]
+		shareBytes, err := ecies.Decrypt(fr.state.sessionSK, encryptedShare)
+		if err != nil {
+			fr.state.currentRound = Blame
+			return fr.createAndBroadcastBlameOfInvalidShare(peerOID)
+		}
+
+		share := &sharing.ShamirShare{
+			Id:    uint32(fr.state.operatorID),
+			Value: shareBytes,
+		}
+
+		p2psend[peerOID] = share
+
+		err = verifiers.Verify(share)
+		if err != nil {
+			fr.state.currentRound = Blame
+			return fr.createAndBroadcastBlameOfInvalidShare(peerOID)
+		}
+	}
+
+	bCastMessage, err := fr.state.participant.Round2(bcast, p2psend)
+	if err != nil {
+		return nil, err
+	}
+
+	msg := &ProtocolMsg{
+		Round: Round2,
+		Round2Message: &Round2Message{
+			Vk:      bCastMessage.VerificationKey.ToAffineCompressed(),
+			VkShare: bCastMessage.VkShare.ToAffineCompressed(),
+		},
+	}
+	_, err = fr.broadcastDKGMessage(msg)
+	return nil, err
+}
diff --git a/dkg/frost/frost_test.go b/dkg/frost/frost_test.go
new file mode 100644
index 0000000..b79f9ae
--- /dev/null
+++ b/dkg/frost/frost_test.go
@@ -0,0 +1,238 @@
+package frost
+
+import (
+	"encoding/hex"
+	"testing"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+	"github.com/pkg/errors"
+	"github.com/stretchr/testify/require"
+)
+
+var (
+	expectedFrostOutput = testingutils.TestKeygenOutcome{
+		Share: map[uint32]string{
+			1: "13b4682b21fe50088beff43530787d1dac1e50c8e0686ec55849c8c9c9c5c044",
+			2: "2b2becc7a00babd145cd75772126d9a8a10f3ca975ad4fa862abe06f4c7e8b59",
+			3: "5bfa5c11d0a68028a94abf35e43955c2c7f4f4d18ef62fb9eed6b905dfe4d2ef",
+			4: "32320eb68a314fc6832df969700e1966cd11d53e2c44b2fafcca528e83f89705",
+		},
+		ValidatorPK: "ab2caf206286eb161d47124885b05b0e92d5d77ba29ce7aa77d9cd38ea24cfc6198d037d7b2011388b475c24ab40091e",
+		OperatorPubKeys: map[uint32]string{
+			1: "b9dbb91742532eb1e8641491bd3a2ee149584d4d6c68169daad84addfa848088c38c3c6302abbcb4f648441b0c67c6e4",
+			2: "ad68795bfe98239f64eaaea753ad6cb5fbdc51fdecf1b42abcee65906eabe4f376b1fd85dbc11e15bf0b04d28fbda199",
+			3: "8ec8fa19ece71538a6435a9784d7565496c57ffbaa1160a020ca14e4c64bdae5d6073bdb43bad401fec264b4cc554295",
+			4: "ad67ab94ab4f560a414f3fdc7b15bf0cf091ff72791c37515631a14c6446462a93b0559851238f28550cb82ab0808e22",
+		},
+	}
+)
+
+func TestFrostDKG(t *testing.T) {
+
+	operators := []types.OperatorID{
+		1, 2, 3, 4,
+	}
+
+	outputs, err := TestingFrost(
+		3,
+		operators,
+		nil,
+		false,
+		nil,
+	)
+	if err != nil {
+		t.Error(err)
+	}
+
+	for _, operatorID := range operators {
+		output := outputs[uint32(operatorID)].ProtocolOutput
+		require.Equal(t, expectedFrostOutput.ValidatorPK, hex.EncodeToString(output.ValidatorPK))
+		require.Equal(t, expectedFrostOutput.Share[uint32(operatorID)], output.Share.SerializeToHexStr())
+		for opID, publicKey := range output.OperatorPubKeys {
+			require.Equal(t, expectedFrostOutput.OperatorPubKeys[uint32(opID)], publicKey.SerializeToHexStr())
+		}
+	}
+}
+
+func TestResharing(t *testing.T) {
+
+	tests := map[string]struct {
+		input, expected testingutils.TestKeygenOutcome
+	}{
+		"test_1": {
+			input: testingutils.TestKeygenOutcome{
+				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
+				Share: map[uint32]string{
+					1: "5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912",
+					2: "533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7",
+					3: "442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186",
+					4: "2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f",
+				},
+				OperatorPubKeys: map[uint32]string{
+					1: "add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa",
+					2: "8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5",
+					3: "b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051",
+					4: "a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037",
+				},
+			},
+			expected: testingutils.TestKeygenOutcome{
+				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
+				Share: map[uint32]string{
+					5: "437a5713f74cbfc67bf2781dd6c8cb74db8bb0a1e598b9bb372e84c56da99dad",
+					6: "72475d8d45ee21f23d27a30f839782814114a24e5fe356bc9301af2289127ee5",
+					7: "0e7645e7292c9aabd34be54e52f4bd42ae63932c2370ad076fc0e50496ec1460",
+					8: "73d0061b1de0a1cbd80cc6f261c603c91eb16f44303bd098cd6c266897365e21",
+				},
+				OperatorPubKeys: map[uint32]string{
+					5: "a7bc4dbb7be0e2ce4b5f6121d7a0fc2902eb9abe4c8b17875cc74a5ce3bd61784eeca8561ce169b2f92e5a157d8d0e49",
+					6: "83678e9ffc680a5ba164226eae678b575c87a8d67696c0c997ae0ad283a942cad21e576994db93c5b55ce42852e0fe87",
+					7: "8f3dc952da48c29313099f36d11b66b307896ab7f8a8396e0ac7b64908533641818241bf963b04c4113d3c7f01d6e72b",
+					8: "8153e8b4f1820e625cc695c7fde34e9f4498e889e9efcc4b86f2ba645d3026df284a6d133ac28f40e6569f658bf4c39c",
+				},
+			},
+		},
+		"test_2": {
+			input: expectedFrostOutput,
+			expected: testingutils.TestKeygenOutcome{
+				Share: map[uint32]string{
+					5: "1459f89fc085be6d93fafbfaa197635491a6816be7d636de0cfb9931ac2c47c6",
+					6: "4326ff190f272099553026ec4e661a60f72f73186220d3df68cec38ec79528fe",
+					7: "53438ec61c03169b1e8e413327652d27b83c07f925ac8629458df96fd56ebe7a",
+					8: "44afa7a6e719a072f0154acf2c949ba8d4cc400e32794dbba3393ad4d5b9083a",
+				},
+				ValidatorPK: "ab2caf206286eb161d47124885b05b0e92d5d77ba29ce7aa77d9cd38ea24cfc6198d037d7b2011388b475c24ab40091e",
+				OperatorPubKeys: map[uint32]string{
+					5: "a841b23b62457fe942389e5289b78848c4c7c709623da1e769ed2d5f7197f831f926bbba945940cb80e27b476e33a003",
+					6: "b4a227e0baf87ce87ae7dd773615d04bd0732a14bc3db4db8f3a68c02ad639d1b1ed4420b63927fa37da7aa4fd66eef6",
+					7: "84df18d6eef01054a3cb21104570b261b2705ee6ff82ed542f0120bc558e359255e3cf22d6280995983e5f5a216313ea",
+					8: "8746cd13d7b9f4c472f7e851a23eb30dca95b5fd02bccdfbe9e27742eed30c1a0a9dda775079016701e4e56362c5049b",
+				},
+			},
+		},
+	}
+
+	for name, test := range tests {
+		t.Run(name, func(t *testing.T) {
+			operatorsOld := []types.OperatorID{
+				1, 2, 3, // 4,
+			}
+
+			operators := []types.OperatorID{
+				5, 6, 7, 8,
+			}
+
+			outcomes, err := TestingFrost(
+				3,
+				operators,
+				operatorsOld,
+				true,
+				&test.input,
+			)
+			if err != nil {
+				t.Fatalf("failed to run frost: %s", err.Error())
+			}
+
+			for _, operatorID := range operators {
+				outcome := outcomes[uint32(operatorID)].ProtocolOutput
+
+				require.Equal(t, test.expected.ValidatorPK, hex.EncodeToString(outcome.ValidatorPK))
+				require.Equal(t, test.expected.Share[uint32(operatorID)], outcome.Share.SerializeToHexStr())
+				for opID, publicKey := range outcome.OperatorPubKeys {
+					require.Equal(t, test.expected.OperatorPubKeys[uint32(opID)], publicKey.SerializeToHexStr())
+				}
+			}
+		})
+	}
+}
+
+func TestingFrost(
+	threshold uint64,
+	operators, operatorsOld []types.OperatorID,
+	isResharing bool,
+	oldKeygenOutcomes *testingutils.TestKeygenOutcome,
+) (map[uint32]*dkg.ProtocolOutcome, error) {
+
+	testingutils.ResetRandSeed()
+	requestID := testingutils.GetRandRequestID()
+	dkgsigner := testingutils.NewTestingKeyManager()
+	storage := testingutils.NewTestingStorage()
+	network := testingutils.NewTestingNetwork()
+
+	init := &dkg.Init{
+		OperatorIDs: operators,
+		Threshold:   uint16(threshold),
+	}
+
+	kgps := make(map[types.OperatorID]dkg.Protocol)
+	for _, operatorID := range operators {
+		p := New(network, operatorID, requestID, dkgsigner, storage, init)
+		kgps[operatorID] = p
+	}
+
+	if isResharing {
+		operatorsOldList := types.OperatorList(operatorsOld).ToUint32List()
+		keygenOutcomeOld := oldKeygenOutcomes.ToKeygenOutcomeMap(threshold, operatorsOldList)
+
+		reshare := &dkg.Reshare{
+			ValidatorPK: keygenOutcomeOld[operatorsOldList[0]].ValidatorPK,
+			OperatorIDs: operators,
+			Threshold:   uint16(threshold),
+		}
+
+		for _, operatorID := range operatorsOld {
+			p := NewResharing(network, operatorID, requestID, dkgsigner, storage, operatorsOld, reshare, keygenOutcomeOld[uint32(operatorID)])
+			kgps[operatorID] = p
+
+		}
+
+		for _, operatorID := range operators {
+			p := NewResharing(network, operatorID, requestID, dkgsigner, storage, operatorsOld, reshare, nil)
+			kgps[operatorID] = p
+		}
+	}
+
+	alloperators := operators
+	if isResharing {
+		alloperators = append(alloperators, operatorsOld...)
+	}
+
+	for _, operatorID := range alloperators {
+		if err := kgps[operatorID].Start(); err != nil {
+			return nil, errors.Wrapf(err, "failed to start dkg protocol for operator %d", operatorID)
+		}
+	}
+
+	outcomes := make(map[uint32]*dkg.ProtocolOutcome)
+	for i := 0; i < 3; i++ {
+
+		messages := network.BroadcastedMsgs
+		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)
+
+		for _, msg := range messages {
+
+			dkgMsg := &dkg.SignedMessage{}
+			if err := dkgMsg.Decode(msg.Data); err != nil {
+				return nil, err
+			}
+
+			for _, operatorID := range alloperators {
+
+				if operatorID == dkgMsg.Signer {
+					continue
+				}
+
+				finished, outcome, err := kgps[operatorID].ProcessMsg(dkgMsg)
+				if err != nil {
+					return nil, err
+				}
+				if finished {
+					outcomes[uint32(operatorID)] = outcome
+				}
+			}
+		}
+	}
+
+	return outcomes, nil
+}
diff --git a/dkg/frost/messages.go b/dkg/frost/messages.go
new file mode 100644
index 0000000..f575232
--- /dev/null
+++ b/dkg/frost/messages.go
@@ -0,0 +1,207 @@
+package frost
+
+import (
+	"encoding/json"
+	"github.com/bloxapp/ssv-spec/dkg"
+	ecies "github.com/ecies/go/v2"
+	"github.com/pkg/errors"
+)
+
+type ProtocolMsg struct {
+	Round              ProtocolRound       `json:"round,omitempty"`
+	PreparationMessage *PreparationMessage `json:"preparation,omitempty"`
+	Round1Message      *Round1Message      `json:"round1,omitempty"`
+	Round2Message      *Round2Message      `json:"round2,omitempty"`
+	BlameMessage       *BlameMessage       `json:"blame,omitempty"`
+}
+
+func (msg *ProtocolMsg) hasOnlyOneMsg() bool {
+	var count = 0
+	if msg.PreparationMessage != nil {
+		count++
+	}
+	if msg.Round1Message != nil {
+		count++
+	}
+	if msg.Round2Message != nil {
+		count++
+	}
+	if msg.BlameMessage != nil {
+		count++
+	}
+	return count == 1
+}
+
+func (msg *ProtocolMsg) msgMatchesRound() bool {
+	switch msg.Round {
+	case Preparation:
+		return msg.PreparationMessage != nil
+	case Round1:
+		return msg.Round1Message != nil
+	case Round2:
+		return msg.Round2Message != nil
+	case Blame:
+		return msg.BlameMessage != nil
+	default:
+		return false
+	}
+}
+
+func (msg *ProtocolMsg) Validate() error {
+	if !msg.hasOnlyOneMsg() {
+		return errors.New("need to contain one and only one message round")
+	}
+	if !msg.msgMatchesRound() {
+		return errors.New("")
+	}
+	switch msg.Round {
+	case Preparation:
+		return msg.PreparationMessage.Validate()
+	case Round1:
+		return msg.Round1Message.Validate()
+	case Round2:
+		return msg.Round2Message.Validate()
+	}
+	return nil
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *ProtocolMsg) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *ProtocolMsg) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+type PreparationMessage struct {
+	SessionPk []byte
+}
+
+func (msg *PreparationMessage) Validate() error {
+	_, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
+	return err
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *PreparationMessage) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *PreparationMessage) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+type Round1Message struct {
+	// Commitment bytes representation of commitment points to pre-selected polynomials
+	Commitment [][]byte
+	// ProofS the S value of the Schnorr's proof
+	ProofS []byte
+	// ProofR the R value of the Schnorr's proof
+	ProofR []byte
+	// Shares the encrypted shares by operator
+	Shares map[uint32][]byte
+}
+
+func (msg *Round1Message) Validate() error {
+	var err error
+	for _, bytes := range msg.Commitment {
+		_, err = thisCurve.Point.FromAffineCompressed(bytes)
+		if err != nil {
+			return errors.Wrap(err, "invalid commitment")
+		}
+	}
+
+	_, err = thisCurve.Scalar.SetBytes(msg.ProofS)
+	if err != nil {
+		return errors.Wrap(err, "invalid ProofS")
+	}
+	_, err = thisCurve.Scalar.SetBytes(msg.ProofR)
+	if err != nil {
+		return errors.Wrap(err, "invalid ProofR")
+	}
+
+	return nil
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *Round1Message) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *Round1Message) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+type Round2Message struct {
+	Vk      []byte
+	VkShare []byte
+}
+
+func (msg *Round2Message) Validate() error {
+	var err error
+	_, err = thisCurve.Point.FromAffineCompressed(msg.Vk)
+	if err != nil {
+		return errors.Wrap(err, "invalid vk")
+	}
+	_, err = thisCurve.Point.FromAffineCompressed(msg.VkShare)
+	if err != nil {
+		return errors.Wrap(err, "invalid vk share")
+	}
+	return nil
+}
+
+type BlameMessage struct {
+	Type             BlameType
+	TargetOperatorID uint32
+	BlameData        [][]byte // SignedMessages received from the bad participant
+	BlamerSessionSk  []byte
+}
+
+func (msg *BlameMessage) Validate() error {
+	if len(msg.BlameData) < 1 {
+		return errors.New("no blame data")
+	}
+	for _, datum := range msg.BlameData {
+		signedMsg := &dkg.SignedMessage{}
+		err := signedMsg.Decode(datum)
+		if err != nil {
+			return errors.Wrap(err, "contained data is not SignedMessage")
+		}
+	}
+	return nil
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *BlameMessage) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *BlameMessage) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+type BlameType int
+
+const (
+	// InconsistentMessage refers to an operator sending multiple messages for same round
+	InconsistentMessage BlameType = iota
+	// InvalidShare refers to an operator sending invalid share
+	InvalidShare
+	//// InvalidMessage refers to messages containing invalid values
+	InvalidMessage
+)
+
+func (t BlameType) ToString() string {
+	m := map[BlameType]string{
+		InconsistentMessage: "Inconsistent Message",
+		InvalidShare:        "Invalid Share",
+		//FailedEcies:         "Failed Ecies",
+		InvalidMessage: "Invalid Message",
+	}
+	return m[t]
+}
diff --git a/dkg/messages.go b/dkg/messages.go
new file mode 100644
index 0000000..5e42253
--- /dev/null
+++ b/dkg/messages.go
@@ -0,0 +1,335 @@
+package dkg
+
+import (
+	"crypto/ecdsa"
+	"crypto/sha256"
+	"encoding/binary"
+	"encoding/json"
+
+	"github.com/attestantio/go-eth2-client/spec/phase0"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/ethereum/go-ethereum/accounts/abi"
+	"github.com/ethereum/go-ethereum/common"
+	"github.com/ethereum/go-ethereum/crypto"
+	"github.com/pkg/errors"
+)
+
+type RequestID [24]byte
+
+const (
+	blsPubkeySize      = 48
+	ethAddressSize     = 20
+	ethAddressStartPos = 0
+	indexSize          = 4
+	indexStartPos      = ethAddressStartPos + ethAddressSize
+)
+
+func (msg RequestID) GetETHAddress() common.Address {
+	ret := common.Address{}
+	copy(ret[:], msg[ethAddressStartPos:ethAddressStartPos+ethAddressSize])
+	return ret
+}
+
+func (msg RequestID) GetRoleType() uint32 {
+	indexByts := msg[indexStartPos : indexStartPos+indexSize]
+	return binary.LittleEndian.Uint32(indexByts)
+}
+
+func NewRequestID(ethAddress common.Address, index uint32) RequestID {
+	indexByts := make([]byte, 4)
+	binary.LittleEndian.PutUint32(indexByts, index)
+
+	ret := RequestID{}
+	copy(ret[ethAddressStartPos:ethAddressStartPos+ethAddressSize], ethAddress[:])
+	copy(ret[indexStartPos:indexStartPos+indexSize], indexByts[:])
+	return ret
+}
+
+type MsgType int
+
+const (
+	// InitMsgType sent when DKG instance is started by requester
+	InitMsgType MsgType = iota
+	// ProtocolMsgType is the DKG itself
+	ProtocolMsgType
+	// DepositDataMsgType post DKG deposit data signatures
+	DepositDataMsgType
+	// OutputMsgType final output msg used by requester to make deposits and register validator with SSV
+	OutputMsgType
+	// ReshareMsgType sent when Resharing is requested
+	ReshareMsgType
+)
+
+type Message struct {
+	MsgType    MsgType
+	Identifier RequestID
+	Data       []byte
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *Message) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *Message) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+func (msg *Message) Validate() error {
+	// TODO msg type
+	// TODO len(data)
+	return nil
+}
+
+func (msg *Message) GetRoot() ([]byte, error) {
+	marshaledRoot, err := msg.Encode()
+	if err != nil {
+		return nil, errors.Wrap(err, "could not encode PartialSignatureMessage")
+	}
+	ret := sha256.Sum256(marshaledRoot)
+	return ret[:], nil
+}
+
+type SignedMessage struct {
+	Message   *Message
+	Signer    types.OperatorID
+	Signature types.Signature
+}
+
+// Encode returns a msg encoded bytes or error
+func (signedMsg *SignedMessage) Encode() ([]byte, error) {
+	return json.Marshal(signedMsg)
+}
+
+// Decode returns error if decoding failed
+func (signedMsg *SignedMessage) Decode(data []byte) error {
+	return json.Unmarshal(data, signedMsg)
+}
+
+func (signedMsg *SignedMessage) Validate() error {
+	// TODO len(sig) == ecdsa sig lenth
+
+	return signedMsg.Message.Validate()
+}
+
+func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
+	return signedMsg.Message.GetRoot()
+}
+
+// Init is the first message in a DKG which initiates a DKG
+type Init struct {
+	// OperatorIDs are the operators selected for the DKG
+	OperatorIDs []types.OperatorID
+	// Threshold DKG threshold for signature reconstruction
+	Threshold uint16
+	// WithdrawalCredentials used when signing the deposit data
+	WithdrawalCredentials []byte
+	// Fork is eth2 fork version
+	Fork phase0.Version
+}
+
+func (msg *Init) Validate() error {
+	if len(msg.WithdrawalCredentials) != phase0.HashLength {
+		return errors.New("invalid WithdrawalCredentials")
+	}
+	contains := func(container []int, elem int) bool {
+		for _, n := range container {
+			if elem == n {
+				return true
+			}
+		}
+		return false
+	}
+	validSizes := []int{4, 7, 10, 13}
+	validN := contains(validSizes, len(msg.OperatorIDs))
+
+	if !validN {
+		return errors.New("invalid number of operators which has to be 3f+1")
+	}
+
+	f := len(msg.OperatorIDs) / 3
+
+	if int(msg.Threshold) != (2*f + 1) {
+		return errors.New("invalid threshold which has to be 2f+1")
+	}
+
+	return nil
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *Init) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *Init) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+// Reshare triggers the resharing protocol
+type Reshare struct {
+	// ValidatorPK is the the public key to be reshared
+	ValidatorPK types.ValidatorPK
+	// OperatorIDs are the operators in the new set
+	OperatorIDs []types.OperatorID
+	// Threshold is the threshold of the new set
+	Threshold uint16
+}
+
+func (msg *Reshare) Validate() error {
+
+	if len(msg.ValidatorPK) != blsPubkeySize {
+		return errors.New("invalid validator pubkey size")
+	}
+
+	if len(msg.OperatorIDs) < 4 || (len(msg.OperatorIDs)-1)%3 != 0 {
+		return errors.New("invalid number of operators which has to be 3f+1")
+	}
+
+	if int(msg.Threshold) != (len(msg.OperatorIDs)-1)*2/3+1 {
+		return errors.New("invalid threshold which has to be 2f+1")
+	}
+
+	return nil
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *Reshare) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *Reshare) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+// Output is the last message in every DKG which marks a specific node's end of process
+type Output struct {
+	// RequestID for the DKG instance (not used for signing)
+	RequestID RequestID
+	// EncryptedShare standard SSV encrypted shares
+	EncryptedShare []byte
+	// SharePubKey is the share's BLS pubkey
+	SharePubKey []byte
+	// ValidatorPubKey the resulting public key corresponding to the shared private key
+	ValidatorPubKey types.ValidatorPK
+	// DepositDataSignature reconstructed signature of DepositMessage according to eth2 spec
+	DepositDataSignature types.Signature
+}
+
+func (o *Output) GetRoot() ([]byte, error) {
+	bytesSolidity, _ := abi.NewType("bytes", "", nil)
+
+	arguments := abi.Arguments{
+		{
+			Type: bytesSolidity,
+		},
+		{
+			Type: bytesSolidity,
+		},
+		{
+			Type: bytesSolidity,
+		},
+		{
+			Type: bytesSolidity,
+		},
+	}
+
+	bytes, err := arguments.Pack(
+		[]byte(o.EncryptedShare),
+		[]byte(o.SharePubKey),
+		[]byte(o.ValidatorPubKey),
+		[]byte(o.DepositDataSignature),
+	)
+	if err != nil {
+		return nil, err
+	}
+	return crypto.Keccak256(bytes), nil
+}
+
+type SignedOutput struct {
+	// Blame Data
+	BlameData *BlameData
+	// Data signed
+	Data *Output
+	// Signer Operator ID which signed
+	Signer types.OperatorID
+	// Signature over Data.GetRoot()
+	Signature types.Signature
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *SignedOutput) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *SignedOutput) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+type BlameData struct {
+	RequestID    RequestID
+	Valid        bool
+	BlameMessage []byte
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *BlameData) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *BlameData) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
+
+func (msg *BlameData) GetRoot() ([]byte, error) {
+	bytesSolidity, _ := abi.NewType("bytes", "", nil)
+	boolSolidity, _ := abi.NewType("bool", "", nil)
+
+	arguments := abi.Arguments{
+		{
+			Type: boolSolidity,
+		},
+		{
+			Type: bytesSolidity,
+		},
+	}
+
+	bytes, err := arguments.Pack(
+		msg.Valid,
+		[]byte(msg.BlameMessage),
+	)
+	if err != nil {
+		return nil, err
+	}
+	return crypto.Keccak256(bytes), nil
+}
+
+func SignOutput(output *Output, privKey *ecdsa.PrivateKey) (types.Signature, error) {
+	root, err := output.GetRoot()
+	if err != nil {
+		return nil, errors.Wrap(err, "could not get root from output message")
+	}
+
+	return crypto.Sign(root, privKey)
+}
+
+// PartialDepositData contains a partial deposit data signature
+type PartialDepositData struct {
+	Signer    types.OperatorID
+	Root      []byte
+	Signature types.Signature
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *PartialDepositData) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *PartialDepositData) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
diff --git a/dkg/messages_test.go b/dkg/messages_test.go
new file mode 100644
index 0000000..1415a62
--- /dev/null
+++ b/dkg/messages_test.go
@@ -0,0 +1,84 @@
+package dkg
+
+import (
+	spec "github.com/attestantio/go-eth2-client/spec/phase0"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/ethereum/go-ethereum/common"
+	"github.com/stretchr/testify/require"
+	"testing"
+)
+
+func TestInit_Validate(t *testing.T) {
+	t.Run("valid", func(t *testing.T) {
+		init := Init{
+			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
+			Threshold:             3,
+			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
+			Fork:                  spec.Version{},
+		}
+		require.NoError(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7}
+		init.Threshold = 5
+		require.NoError(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
+		init.Threshold = 7
+		require.NoError(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
+		init.Threshold = 9
+		require.NoError(t, init.Validate())
+	})
+	t.Run("invalid number of operators", func(t *testing.T) {
+		init := Init{
+			OperatorIDs:           []types.OperatorID{1, 2, 3},
+			Threshold:             3,
+			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
+			Fork:                  spec.Version{},
+		}
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6}
+		init.Threshold = 3
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8}
+		init.Threshold = 5
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
+		init.Threshold = 7
+		require.Error(t, init.Validate())
+	})
+	t.Run("invalid threshold", func(t *testing.T) {
+		init := Init{
+			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
+			Threshold:             2,
+			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
+			Fork:                  spec.Version{},
+		}
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7}
+		init.Threshold = 6
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
+		init.Threshold = 8
+		require.Error(t, init.Validate())
+		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
+		init.Threshold = 8
+		require.Error(t, init.Validate())
+	})
+	t.Run("short WithdrawalCredentials", func(t *testing.T) {
+		init := Init{
+			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
+			Threshold:             3,
+			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd6680"),
+			Fork:                  spec.Version{},
+		}
+		require.Error(t, init.Validate())
+	})
+	t.Run("long WithdrawalCredentials", func(t *testing.T) {
+		init := Init{
+			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
+			Threshold:             3,
+			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808faa"),
+			Fork:                  spec.Version{},
+		}
+		require.Error(t, init.Validate())
+	})
+}
diff --git a/dkg/node.go b/dkg/node.go
new file mode 100644
index 0000000..a0d8535
--- /dev/null
+++ b/dkg/node.go
@@ -0,0 +1,259 @@
+package dkg
+
+import (
+	"encoding/hex"
+
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/pkg/errors"
+)
+
+// Runners is a map of dkg runners mapped by dkg ID.
+type Runners map[string]Runner
+
+func (runners Runners) AddRunner(id RequestID, runner Runner) {
+	runners[hex.EncodeToString(id[:])] = runner
+}
+
+// RunnerForID returns a Runner from the provided msg ID, or nil if not found
+func (runners Runners) RunnerForID(id RequestID) Runner {
+	return runners[hex.EncodeToString(id[:])]
+}
+
+func (runners Runners) Exists(id RequestID) bool {
+	_, ok := runners[hex.EncodeToString(id[:])]
+	return ok
+}
+
+func (runners Runners) DeleteRunner(id RequestID) {
+	delete(runners, hex.EncodeToString(id[:]))
+}
+
+// Node is responsible for receiving and managing DKG session and messages
+type Node struct {
+	operator *Operator
+	// runners holds all active running DKG runners
+	operatorsOld []types.OperatorID
+	runners      Runners
+	config       *Config
+}
+
+func NewNode(operator *Operator, config *Config) *Node {
+	return &Node{
+		operator: operator,
+		config:   config,
+		runners:  make(Runners, 0),
+	}
+}
+
+func NewResharingNode(operator *Operator, operatorsOld []types.OperatorID, config *Config) *Node {
+	return &Node{
+		operator:     operator,
+		operatorsOld: operatorsOld,
+		config:       config,
+		runners:      make(Runners, 0),
+	}
+}
+
+func (n *Node) newRunner(id RequestID, initMsg *Init) (Runner, error) {
+	r := &runner{
+		Operator:              n.operator,
+		InitMsg:               initMsg,
+		Identifier:            id,
+		KeygenOutcome:         nil,
+		DepositDataRoot:       nil,
+		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
+		OutputMsgs:            map[types.OperatorID]*SignedOutput{},
+		protocol:              n.config.KeygenProtocol(n.config.Network, n.operator.OperatorID, id, n.config.Signer, n.config.Storage, initMsg),
+		config:                n.config,
+	}
+
+	if err := r.protocol.Start(); err != nil {
+		return nil, errors.Wrap(err, "could not start dkg protocol")
+	}
+
+	return r, nil
+}
+
+func (n *Node) newResharingRunner(id RequestID, reshareMsg *Reshare) (Runner, error) {
+	kgOutput, err := n.config.Storage.GetKeyGenOutput(reshareMsg.ValidatorPK)
+	if err != nil {
+		return nil, errors.Wrap(err, "could not find the keygen output from storage")
+	}
+	r := &runner{
+		Operator:              n.operator,
+		ReshareMsg:            reshareMsg,
+		Identifier:            id,
+		KeygenOutcome:         nil,
+		DepositDataRoot:       nil,
+		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},
+		OutputMsgs:            map[types.OperatorID]*SignedOutput{},
+		protocol:              n.config.ReshareProtocol(n.config.Network, n.operator.OperatorID, id, n.config.Signer, n.config.Storage, n.operatorsOld, reshareMsg, kgOutput),
+		config:                n.config,
+	}
+
+	if err := r.protocol.Start(); err != nil {
+		return nil, errors.Wrap(err, "could not start resharing protocol")
+	}
+
+	return r, nil
+}
+
+// ProcessMessage processes network Messages of all types
+func (n *Node) ProcessMessage(msg *types.SSVMessage) error {
+	if msg.MsgType != types.DKGMsgType {
+		return errors.New("not a DKGMsgType")
+	}
+	signedMsg := &SignedMessage{}
+	if err := signedMsg.Decode(msg.GetData()); err != nil {
+		return errors.Wrap(err, "could not get dkg Message from network Messages")
+	}
+
+	if err := n.validateSignedMessage(signedMsg); err != nil {
+		return errors.Wrap(err, "signed message doesn't pass validation")
+	}
+
+	switch signedMsg.Message.MsgType {
+	case InitMsgType:
+		return n.startNewDKGMsg(signedMsg)
+	case ReshareMsgType:
+		return n.startResharing(signedMsg)
+	case ProtocolMsgType:
+		return n.processDKGMsg(signedMsg)
+	case DepositDataMsgType:
+		return n.processDKGMsg(signedMsg)
+	case OutputMsgType:
+		return n.processDKGMsg(signedMsg)
+	default:
+		return errors.New("unknown msg type")
+	}
+}
+
+func (n *Node) validateSignedMessage(message *SignedMessage) error {
+	if err := message.Validate(); err != nil {
+		return errors.Wrap(err, "message invalid")
+	}
+
+	return nil
+}
+
+func (n *Node) startNewDKGMsg(message *SignedMessage) error {
+	initMsg, err := n.validateInitMsg(message)
+	if err != nil {
+		return errors.Wrap(err, "could not start new dkg")
+	}
+
+	runner, err := n.newRunner(message.Message.Identifier, initMsg)
+	if err != nil {
+		return errors.Wrap(err, "could not start new dkg")
+	}
+
+	// add runner to runners
+	n.runners.AddRunner(message.Message.Identifier, runner)
+
+	return nil
+}
+
+func (n *Node) startResharing(message *SignedMessage) error {
+	reshareMsg, err := n.validateReshareMsg(message)
+	if err != nil {
+		return errors.Wrap(err, "could not start resharing")
+	}
+
+	r, err := n.newResharingRunner(message.Message.Identifier, reshareMsg)
+	if err != nil {
+		return errors.Wrap(err, "could not start resharing")
+	}
+
+	// add runner to runners
+	n.runners.AddRunner(message.Message.Identifier, r)
+
+	return nil
+}
+
+func (n *Node) validateInitMsg(message *SignedMessage) (*Init, error) {
+	// validate identifier.GetEthAddress is the signer for message
+	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, message.Message.Identifier.GetETHAddress()); err != nil {
+		return nil, errors.Wrap(err, "signed message invalid")
+	}
+
+	initMsg := &Init{}
+	if err := initMsg.Decode(message.Message.Data); err != nil {
+		return nil, errors.Wrap(err, "could not get dkg init Message from signed Messages")
+	}
+
+	if err := initMsg.Validate(); err != nil {
+		return nil, errors.Wrap(err, "init message invalid")
+	}
+
+	// check instance not running already
+	if n.runners.RunnerForID(message.Message.Identifier) != nil {
+		return nil, errors.New("dkg started already")
+	}
+
+	return initMsg, nil
+}
+
+func (n *Node) validateReshareMsg(message *SignedMessage) (*Reshare, error) {
+	// validate identifier.GetEthAddress is the signer for message
+	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, message.Message.Identifier.GetETHAddress()); err != nil {
+		return nil, errors.Wrap(err, "signed message invalid")
+	}
+
+	reshareMsg := &Reshare{}
+	if err := reshareMsg.Decode(message.Message.Data); err != nil {
+		return nil, errors.Wrap(err, "could not get reshare Message from signed Messages")
+	}
+
+	if err := reshareMsg.Validate(); err != nil {
+		return nil, errors.Wrap(err, "reshare message invalid")
+	}
+
+	// check instance not running already
+	if n.runners.RunnerForID(message.Message.Identifier) != nil {
+		return nil, errors.New("dkg started already")
+	}
+
+	return reshareMsg, nil
+}
+
+func (n *Node) processDKGMsg(message *SignedMessage) error {
+	if !n.runners.Exists(message.Message.Identifier) {
+		return errors.New("could not find dkg runner")
+	}
+
+	if err := n.validateDKGMsg(message); err != nil {
+		return errors.Wrap(err, "dkg msg not valid")
+	}
+
+	r := n.runners.RunnerForID(message.Message.Identifier)
+	finished, err := r.ProcessMsg(message)
+	if err != nil {
+		return errors.Wrap(err, "could not process dkg message")
+	}
+	if finished {
+		n.runners.DeleteRunner(message.Message.Identifier)
+	}
+
+	return nil
+}
+
+func (n *Node) validateDKGMsg(message *SignedMessage) error {
+
+	// find signing operator and verify sig
+	found, signingOperator, err := n.config.Storage.GetDKGOperator(message.Signer)
+	if err != nil {
+		return errors.Wrap(err, "can't fetch operator")
+	}
+	if !found {
+		return errors.New("can't find operator")
+	}
+	if err := message.Signature.ECRecover(message, n.config.SignatureDomainType, types.DKGSignatureType, signingOperator.ETHAddress); err != nil {
+		return errors.Wrap(err, "signed message invalid")
+	}
+
+	return nil
+}
+
+func (n *Node) GetConfig() *Config {
+	return n.config
+}
diff --git a/dkg/protocol.go b/dkg/protocol.go
new file mode 100644
index 0000000..d238e46
--- /dev/null
+++ b/dkg/protocol.go
@@ -0,0 +1,43 @@
+package dkg
+
+import (
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+)
+
+type ProtocolOutcome struct {
+	ProtocolOutput *KeyGenOutput
+	BlameOutput    *BlameOutput
+}
+
+func (o *ProtocolOutcome) IsFailedWithBlame() (bool, error) {
+	if o.ProtocolOutput == nil && o.BlameOutput == nil {
+		return false, errors.New("invalid outcome - missing KeyGenOutput and BlameOutput")
+	}
+	if o.ProtocolOutput != nil && o.BlameOutput != nil {
+		return false, errors.New("invalid outcome - has both KeyGenOutput and BlameOutput")
+	}
+	return o.BlameOutput != nil, nil
+}
+
+// KeyGenOutput is the bare minimum output from the protocol
+type KeyGenOutput struct {
+	Share           *bls.SecretKey
+	OperatorPubKeys map[types.OperatorID]*bls.PublicKey
+	ValidatorPK     types.ValidatorPK
+	Threshold       uint64
+}
+
+// BlameOutput is the output of blame round
+type BlameOutput struct {
+	Valid        bool
+	BlameMessage *SignedMessage
+}
+
+// Protocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
+type Protocol interface {
+	Start() error
+	// ProcessMsg returns true and a bls share if finished
+	ProcessMsg(msg *SignedMessage) (bool, *ProtocolOutcome, error)
+}
diff --git a/dkg/runner.go b/dkg/runner.go
new file mode 100644
index 0000000..68860eb
--- /dev/null
+++ b/dkg/runner.go
@@ -0,0 +1,336 @@
+package dkg
+
+import (
+	"bytes"
+
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/ethereum/go-ethereum/common"
+	"github.com/ethereum/go-ethereum/crypto"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+)
+
+type Runner interface {
+	ProcessMsg(msg *SignedMessage) (bool, error)
+}
+
+// Runner manages the execution of a DKG, start to finish.
+type runner struct {
+	Operator *Operator
+	// InitMsg holds the init method which started this runner
+	InitMsg *Init
+	// ReshareMsg holds the reshare method which started this runner
+	ReshareMsg *Reshare
+	// Identifier unique for DKG session
+	Identifier RequestID
+	// KeygenOutcome holds the protocol outcome once it finishes
+	KeygenOutcome *ProtocolOutcome
+	// DepositDataRoot is the signing root for the deposit data
+	DepositDataRoot []byte
+	// DepositDataSignatures holds partial sigs on deposit data
+	DepositDataSignatures map[types.OperatorID]*PartialDepositData
+	// OutputMsgs holds all output messages received
+	OutputMsgs map[types.OperatorID]*SignedOutput
+
+	protocol Protocol
+	config   *Config
+}
+
+// ProcessMsg processes a DKG signed message and returns true and stream keygen output or blame if finished
+func (r *runner) ProcessMsg(msg *SignedMessage) (bool, error) {
+	// TODO - validate message
+
+	switch msg.Message.MsgType {
+	case ProtocolMsgType:
+		if r.DepositDataSignatures[r.Operator.OperatorID] != nil {
+			return false, errors.New("keygen has already completed")
+		}
+
+		finished, o, err := r.protocol.ProcessMsg(msg)
+		if err != nil {
+			return false, errors.Wrap(err, "failed to process dkg msg")
+		}
+
+		if finished {
+			r.KeygenOutcome = o
+			isBlame, err := r.KeygenOutcome.IsFailedWithBlame()
+			if err != nil {
+				return true, errors.Wrap(err, "invalid KeygenOutcome")
+			}
+			if isBlame {
+				err := r.config.Network.StreamDKGBlame(r.KeygenOutcome.BlameOutput)
+				return true, errors.Wrap(err, "failed to stream blame output")
+			}
+			if r.KeygenOutcome.ProtocolOutput == nil {
+				return true, errors.Wrap(err, "protocol finished without blame or keygen result")
+			}
+
+			if r.isResharing() {
+				if err := r.prepareAndBroadcastOutput(); err != nil {
+					return false, err
+				}
+			} else {
+				if err := r.prepareAndBroadcastDepositData(); err != nil {
+					return false, err
+				}
+			}
+
+		}
+		return false, nil
+	case DepositDataMsgType:
+		depSig := &PartialDepositData{}
+		if err := depSig.Decode(msg.Message.Data); err != nil {
+			return false, errors.Wrap(err, "could not decode PartialDepositData")
+		}
+
+		if err := r.validateDepositDataSig(depSig); err != nil {
+			return false, errors.Wrap(err, "PartialDepositData invalid")
+		}
+
+		if found := r.DepositDataSignatures[msg.Signer]; found == nil {
+			r.DepositDataSignatures[msg.Signer] = depSig
+		} else if !bytes.Equal(found.Signature, msg.Signature) {
+			return false, errors.New("inconsistent partial signature received")
+		}
+
+		if len(r.DepositDataSignatures) == int(r.InitMsg.Threshold) {
+			if err := r.prepareAndBroadcastOutput(); err != nil {
+				return false, err
+			}
+		}
+		return false, nil
+	case OutputMsgType:
+		output := &SignedOutput{}
+		if err := output.Decode(msg.Message.Data); err != nil {
+			return false, errors.Wrap(err, "could not decode SignedOutput")
+		}
+
+		if err := r.validateSignedOutput(output); err != nil {
+			return false, errors.Wrap(err, "signed output invali")
+		}
+
+		r.OutputMsgs[msg.Signer] = output
+		// GLNOTE: Actually we need every operator to sign instead only the quorum!
+		finished := false
+		if !r.isResharing() {
+			finished = len(r.OutputMsgs) == len(r.InitMsg.OperatorIDs)
+		} else {
+			finished = len(r.OutputMsgs) == len(r.ReshareMsg.OperatorIDs)
+		}
+		if finished {
+			err := r.config.Network.StreamDKGOutput(r.OutputMsgs)
+			return true, errors.Wrap(err, "failed to stream dkg output")
+		}
+
+		return false, nil
+	default:
+		return false, errors.New("msg type invalid")
+	}
+}
+
+func (r *runner) prepareAndBroadcastDepositData() error {
+	// generate deposit data
+	root, _, err := types.GenerateETHDepositData(
+		r.KeygenOutcome.ProtocolOutput.ValidatorPK,
+		r.InitMsg.WithdrawalCredentials,
+		r.InitMsg.Fork,
+		types.DomainDeposit,
+	)
+	if err != nil {
+		return errors.Wrap(err, "could not generate deposit data")
+	}
+
+	r.DepositDataRoot = root
+
+	// sign
+	sig := r.KeygenOutcome.ProtocolOutput.Share.SignByte(root)
+
+	// broadcast
+	pdd := &PartialDepositData{
+		Signer:    r.Operator.OperatorID,
+		Root:      r.DepositDataRoot,
+		Signature: sig.Serialize(),
+	}
+	if err := r.signAndBroadcastMsg(pdd, DepositDataMsgType); err != nil {
+		return errors.Wrap(err, "could not broadcast partial deposit data")
+	}
+	r.DepositDataSignatures[r.Operator.OperatorID] = pdd
+	return nil
+}
+
+func (r *runner) prepareAndBroadcastOutput() error {
+	var (
+		depositSig types.Signature
+		err        error
+	)
+	if r.isResharing() {
+		depositSig = nil
+	} else {
+		// reconstruct deposit data sig
+		depositSig, err = r.reconstructDepositDataSignature()
+		if err != nil {
+			return errors.Wrap(err, "could not reconstruct deposit data sig")
+		}
+	}
+
+	// encrypt Operator's share
+	encryptedShare, err := r.config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.KeygenOutcome.ProtocolOutput.Share.Serialize())
+	if err != nil {
+		return errors.Wrap(err, "could not encrypt share")
+	}
+
+	ret, err := r.generateSignedOutput(&Output{
+		RequestID:            r.Identifier,
+		EncryptedShare:       encryptedShare,
+		SharePubKey:          r.KeygenOutcome.ProtocolOutput.Share.GetPublicKey().Serialize(),
+		ValidatorPubKey:      r.KeygenOutcome.ProtocolOutput.ValidatorPK,
+		DepositDataSignature: depositSig,
+	})
+	if err != nil {
+		return errors.Wrap(err, "could not generate dkg SignedOutput")
+	}
+
+	r.OutputMsgs[r.Operator.OperatorID] = ret
+	if err := r.signAndBroadcastMsg(ret, OutputMsgType); err != nil {
+		return errors.Wrap(err, "could not broadcast SignedOutput")
+	}
+	return nil
+}
+
+func (r *runner) signAndBroadcastMsg(msg types.Encoder, msgType MsgType) error {
+	data, err := msg.Encode()
+	if err != nil {
+		return err
+	}
+	signedMessage := &SignedMessage{
+		Message: &Message{
+			MsgType:    msgType,
+			Identifier: r.Identifier,
+			Data:       data,
+		},
+		Signer:    r.Operator.OperatorID,
+		Signature: nil,
+	}
+	// GLNOTE: Should we use SignDKGOutput?
+	sig, err := r.config.Signer.SignDKGOutput(signedMessage, r.Operator.ETHAddress)
+	if err != nil {
+		return errors.Wrap(err, "failed to sign message")
+	}
+	signedMessage.Signature = sig
+	if err = r.config.Network.BroadcastDKGMessage(signedMessage); err != nil {
+		return errors.Wrap(err, "failed to broadcast message")
+	}
+	return nil
+}
+
+func (r *runner) reconstructDepositDataSignature() (types.Signature, error) {
+	sigBytes := map[types.OperatorID][]byte{}
+	for id, d := range r.DepositDataSignatures {
+		if err := r.validateDepositDataRoot(d); err != nil {
+			return nil, errors.Wrap(err, "PartialDepositData invalid")
+		}
+		sigBytes[id] = d.Signature
+	}
+
+	sig, err := types.ReconstructSignatures(sigBytes)
+	if err != nil {
+		return nil, err
+	}
+	return sig.Serialize(), nil
+}
+
+func (r *runner) validateSignedOutput(msg *SignedOutput) error {
+	// TODO: Separate fields match and signature validation
+	output := r.ownOutput()
+	if output != nil {
+		if output.BlameData == nil {
+			if output.Data.RequestID != msg.Data.RequestID {
+				return errors.New("got mismatching RequestID")
+			}
+			if !bytes.Equal(output.Data.ValidatorPubKey, msg.Data.ValidatorPubKey) {
+				return errors.New("got mismatching ValidatorPubKey")
+			}
+		} else {
+			if output.BlameData.RequestID != msg.BlameData.RequestID {
+				return errors.New("got mismatching RequestID")
+			}
+		}
+	}
+
+	found, operator, err := r.config.Storage.GetDKGOperator(msg.Signer)
+	if !found {
+		return errors.New("unable to find signer")
+	}
+	if err != nil {
+		return errors.Wrap(err, "unable to find signer")
+	}
+
+	var (
+		root []byte
+	)
+
+	if msg.BlameData == nil {
+		root, err = msg.Data.GetRoot()
+	} else {
+		root, err = msg.BlameData.GetRoot()
+	}
+	if err != nil {
+		return errors.Wrap(err, "fail to get root")
+	}
+
+	pk, err := crypto.Ecrecover(root, msg.Signature)
+	if err != nil {
+		return errors.New("unable to recover public key")
+	}
+	addr := common.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
+	if addr != operator.ETHAddress {
+		return errors.New("invalid signature")
+	}
+	return nil
+}
+
+func (r *runner) validateDepositDataRoot(msg *PartialDepositData) error {
+	if !bytes.Equal(r.DepositDataRoot, msg.Root) {
+		return errors.New("deposit data roots not equal")
+	}
+	return nil
+}
+
+func (r *runner) validateDepositDataSig(msg *PartialDepositData) error {
+
+	// find operator and verify msg
+	sharePK, found := r.KeygenOutcome.ProtocolOutput.OperatorPubKeys[msg.Signer]
+	if !found {
+		return errors.New("signer not part of committee")
+	}
+	sig := &bls.Sign{}
+	if err := sig.Deserialize(msg.Signature); err != nil {
+		return errors.Wrap(err, "could not deserialize partial sig")
+	}
+	if !sig.VerifyByte(sharePK, r.DepositDataRoot) {
+		return errors.New("partial deposit data sig invalid")
+	}
+
+	return nil
+}
+
+func (r *runner) generateSignedOutput(o *Output) (*SignedOutput, error) {
+	sig, err := r.config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
+	if err != nil {
+		return nil, errors.Wrap(err, "could not sign output")
+	}
+
+	return &SignedOutput{
+		Data:      o,
+		Signer:    r.Operator.OperatorID,
+		Signature: sig,
+	}, nil
+}
+
+func (r *runner) ownOutput() *SignedOutput {
+	return r.OutputMsgs[r.Operator.OperatorID]
+}
+
+func (r *runner) isResharing() bool {
+	return r.ReshareMsg != nil
+}
diff --git a/dkg/spectest/all_tests.go b/dkg/spectest/all_tests.go
new file mode 100644
index 0000000..2e07e53
--- /dev/null
+++ b/dkg/spectest/all_tests.go
@@ -0,0 +1,25 @@
+package spectest
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
+	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost"
+	"testing"
+)
+
+type SpecTest interface {
+	TestName() string
+	Run(t *testing.T)
+}
+
+var AllTests = []SpecTest{
+	tests.HappyFlow(),
+
+	frost.Keygen(),
+	frost.Resharing(),
+	frost.BlameTypeInvalidCommitment(),
+	frost.BlameTypeInvalidScalar(),
+	frost.BlameTypeInvalidShare_FailedShareDecryption(),
+	frost.BlameTypeInvalidShare_FailedValidationAgainstCommitment(),
+	frost.BlameTypeInconsistentMessage(),
+	tests.ResharingHappyFlow(),
+}
diff --git a/dkg/spectest/generate/main.go b/dkg/spectest/generate/main.go
new file mode 100644
index 0000000..3c44e26
--- /dev/null
+++ b/dkg/spectest/generate/main.go
@@ -0,0 +1,46 @@
+package main
+
+import (
+	"encoding/json"
+	"fmt"
+	"os"
+	"reflect"
+
+	"github.com/bloxapp/ssv-spec/dkg/spectest"
+)
+
+//go:generate go run main.go
+
+func main() {
+	all := map[string]spectest.SpecTest{}
+	for _, t := range spectest.AllTests {
+		n := reflect.TypeOf(t).String() + "_" + t.TestName()
+		if all[n] != nil {
+			panic(fmt.Sprintf("duplicate test: %s\n", n))
+		}
+		all[n] = t
+	}
+
+	byts, err := json.Marshal(all)
+	if err != nil {
+		panic(err.Error())
+	}
+
+	if len(all) != len(spectest.AllTests) {
+		panic("did not generate all tests\n")
+	}
+
+	fmt.Printf("found %d tests\n", len(all))
+	writeJson(byts)
+}
+
+func writeJson(data []byte) {
+	basedir, _ := os.Getwd()
+	fileName := "tests.json"
+	fullPath := basedir + "/" + fileName
+
+	fmt.Printf("writing spec tests json to: %s\n", fullPath)
+	if err := os.WriteFile(fullPath, data, 0644); err != nil {
+		panic(err.Error())
+	}
+}
diff --git a/dkg/spectest/generate/tests.json b/dkg/spectest/generate/tests.json
new file mode 100644
index 0000000..7bb1dfd
--- /dev/null
+++ b/dkg/spectest/generate/tests.json
@@ -0,0 +1 @@
+{"*frost.FrostSpecTest_Blame Type Inconsisstent Message - Happy Flow":{"Name":"Blame Type Inconsisstent Message - Happy Flow","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":true,"BlameMessage":null}},"ExpectedError":"could not find dkg runner","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"1biIEj3eAjC8zW1urPNyy3JVy/1b5sUaJ5MTJ2rKMpoxZ/flAPkCAHct3z7NDW1cEDwh4AxDA0Y+K0dnpwo7zQE="}],"2":[{"Message":{"MsgType":0,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"sUbv1BKIohRxhwehoBypvrvqFIXe+KKHGorlXAK8eT5KNsur+kOLQwoQiNR7Y2r+ByC6BK0fiAy3T/bUtg1lZQA="}],"3":[{"Message":{"MsgType":0,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"k0JR218crG23mMSMXB3CgQPSbkHGvN+2sVEAa/hmSU8pGvMP7VIhutMpqzfSp3RXBjgHN5wzsVqjpekpwMCT4wE="}],"4":[{"Message":{"MsgType":0,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"LZ1nnnq+Y1lfqzAY8uDH1b5s3hicBC9xDyLSPkYcQ+MHGDTt7ceKC9rFW/mT4/bHFhXxAFraEaoXMLlkEUfo0QE="}]},"2":{"2":[{"Message":{"MsgType":1,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsicjRwT2QxMTlnTFh6dDA2eHd2bVh1ZElZckVGSGw3WnlUN3lYRE16M1d0L0NtSytLa1BSZW02bnE0b3Y1U2YzcSIsImdRaDhCZDhsSm1va1Q5enpGVUsvamF2V3A4ejhWT0lwNVIva0N5WHhDb1lxcElDeU93bWcwWFZZWk1Md0lqL1EiLCJxekhyUkFJcG1hN2xtYmJtMzdTYXpOQ1lYNldFMi9SWVFGK2xRZHIrcytTTy8zQWtuTG9NSDBvY3VBRmp4K0ZhIl0sIlByb29mUyI6IlJZNk1qbmFwQ1B0dDZjUjlZWFdkYmRkM01lOEJTSkNybFRKcFg5eTViTDg9IiwiUHJvb2ZSIjoiTkFSRFZ2S1RINnBTamlZdE13VXFpWlNZdjFsS05WazZkZURCOUZjZGRrQT0iLCJTaGFyZXMiOnsiMSI6IkJPU0R4clk4Ym13TytXZHdZcy9UZ0RDK3ZpWENZUXFsZE5vRW1TdXRPSHJsakJvSUdtUzlLS214YkFZRXBkVHRrK2FoeXlPbkcwbEhuM1dUck45UEpFZVlFNlFwY0dnclJrVVdPcS9SU3dIUVg1MFIwMGlVbUNuWEg1QjNXVlVkeVRUQU96a3ZlbmZXcnE2Vyt1VlE0VnUwMGs1OTBXOXhDYkJ2dEdNNFVYSisiLCIzIjoiQk9sUG9DZUpEYVVzcjNiUlZHUGxVMEpaMU9nUG04U3RiQTkzRFl5RWFMNWU1WTdQTnpFeUNuckRQV1ZvVnFuTmJQazZHaWtXSG9HZC9zT0pDQjRsN2ZCaXlkMEgwSDZZcHd6NDRNRmhFdThxZ0JXeEZHZUc3MzBIWkt2NCs2bWowNDhUZmtqMWwrdFRIZHFJOE8zR2p6d1dENTFVT2wxYUlWNjhzd3NsUXFlTCIsIjQiOiJCSFphLytyaUpremhNN1BGa0lGemhsVWtxWDJLM1AxaVFaTzF3SlJUbXl2UHVxcVluQWMwS3Nia1NuRFNxN0dUd3dBNUwranRsZTNZNE5seFZGSDVscTlSWW50d05uRHlSbGlEd3p4aXM4eGxSRFF0cm5BRmZJeVN3K3JESmE3Y2x4V1VUYXZNakhlRWF3RFdZdjlNSUtiUElkMEF3cmxYTVJiN3B5Y0RNb0FXIn19fQ=="},"Signer":2,"Signature":"Qyt4RRZ3F3ta16kTINyBERwUskevXrQZA837fuHSDupJod7BRbo93BAqsBgrWn5Q/VyBPQHura9PgbwCiwhfLQE="},{"Message":{"MsgType":1,"Identifier":[136,142,148,134,123,195,209,183,237,26,115,54,192,194,209,187,78,17,191,37,54,175,179,207],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsicjRwT2QxMTlnTFh6dDA2eHd2bVh1ZElZckVGSGw3WnlUN3lYRE16M1d0L0NtSytLa1BSZW02bnE0b3Y1U2YzcSIsImdRaDhCZDhsSm1va1Q5enpGVUsvamF2V3A4ejhWT0lwNVIva0N5WHhDb1lxcElDeU93bWcwWFZZWk1Md0lqL1EiLCJxekhyUkFJcG1hN2xtYmJtMzdTYXpOQ1lYNldFMi9SWVFGK2xRZHIrcytTTy8zQWtuTG9NSDBvY3VBRmp4K0ZhIl0sIlByb29mUyI6IlJZNk1qbmFwQ1B0dDZjUjlZWFdkYmRkM01lOEJTSkNybFRKcFg5eTViTDg9IiwiUHJvb2ZSIjoiTkFSRFZ2S1RINnBTamlZdE13VXFpWlNZdjFsS05WazZkZURCOUZjZGRrQT0iLCJTaGFyZXMiOnsiMSI6IkJPU0R4clk4Ym13TytXZHdZcy9UZ0RDK3ZpWENZUXFsZE5vRW1TdXRPSHJsakJvSUdtUzlLS214YkFZRXBkVHRrK2FoeXlPbkcwbEhuM1dUck45UEpFZVlFNlFwY0dnclJrVVdPcS9SU3dIUVg1MFIwMGlVbUNuWEg1QjNXVlVkeVRUQU96a3ZlbmZXcnE2Vyt1VlE0VnUwMGs1OTBXOXhDYkJ2dEdNNFVYSisiLCIzIjoiQk9sUG9DZUpEYVVzcjNiUlZHUGxVMEpaMU9nUG04U3RiQTkzRFl5RWFMNWU1WTdQTnpFeUNuckRQV1ZvVnFuTmJQazZHaWtXSG9HZC9zT0pDQjRsN2ZCaXlkMEgwSDZZcHd6NDRNRmhFdThxZ0JXeEZHZUc3MzBIWkt2NCs2bWowNDhUZmtqMWwrdFRIZHFJOE8zR2p6d1dENTFVT2wxYUlWNjhzd3NsUXFlTCIsIjQiOiJCSFphLytyaUpremhNN1BGa0lGemhsVWtxWDJLM1AxaVFaTzF3SlJUbXl2UHVxcVluQWMwS3Nia1NuRFNxN0dUd3dBNUwranRsZTNZNE5seFZGSDVscTlSWW50d05uRHlSbGlEd3p4aXM4eGxSRFF0cm5BRmZJeVN3K3JESmE3Y2x4V1VUYXZNakhlRWF3RFdZdjlNSUtiUElkMEF3cmxYTVJiN3B5Y0RNb1dBIn19fQ=="},"Signer":2,"Signature":"RsQ/8HDTvsn2SDCMzpOlBE58DLewgqqMPArTLEVh1TlWWXRCA2hCubub/2b6eGU3bSYWOjSLuxOY6wyANc7b1QA="}]}}},"*frost.FrostSpecTest_Blame Type Invalid Commitment - Happy Flow":{"Name":"Blame Type Invalid Commitment - Happy Flow","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":true,"BlameMessage":null}},"ExpectedError":"could not find dkg runner","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"FPVU+w2FnMcN+DpxgIeqAqOCgLzIqQkIuIRHM0WoTN04aAocBR9E1TtbzInCe0GMeip7BroFYXjWszt1MD6OdgE="}],"2":[{"Message":{"MsgType":0,"Identifier":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"WioWZIwC9/tScrTq1dHf7OC5Ja6A3nSY5LWDRkQ2M9FEJhLAJ3lyqrfueLsbkl17IguSeGyGDQ/SUHeqWz7GBQE="}],"3":[{"Message":{"MsgType":0,"Identifier":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"dYNk851ZvXb7KLaUg7Wr8/MufxXdmE84yz6VaHtl9iNziPNXaMTkoepnK/Qyqd0KJDNHY3CFRKCjL/kHTjsT0QA="}],"4":[{"Message":{"MsgType":0,"Identifier":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"YVmu7r/1YVTyKCFJKNy5Xy1aLp+oKiIX0rvCPQI81WhZWAiFa6KVOLGQCOcJ5s4iwb0N3530jj6B2eshK48KcAA="}]},"2":{"2":[{"Message":{"MsgType":1,"Identifier":[141,53,178,36,20,73,203,23,231,17,231,123,133,206,95,207,168,58,76,136,0,130,208,151],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsiazRNR2RqVUtiNXZhZVlwOVJLQWpmbGkzSDV1WmZIOEwzR01PeUE1Uzd3aUxpNFkzM3N2RTloN1RFWVl5V1Z5WiIsImdUekVFaThnL25uVlBNTkRWeHdlUmcxOWgycGdiSXZ2dkxXOHVVQ0tUcVdvTjF6aURqSWhnSkVnOG9GaG0wNHgiLCJjblZpWW1semFDMTJZV3gxWlE9PSJdLCJQcm9vZlMiOiJPU0ZQZ09sbi9IamJWTlJlb0twdlNjdVBBcUNpUDlhbndIWm5Ca3VaRURJPSIsIlByb29mUiI6IkMzTGxIc2NBNnlHbENFcFBJbVRjQjJkcks3VVVaTjZUZXVyN3VGMk9PUGM9IiwiU2hhcmVzIjp7IjEiOiJCQVpSM1lIUlJPZDhDOVNkbTg2TFVkM01va05BNjdWbTZQVjlXR081ellVdlpCTjNWekZTTmkzTFdGQ3p3TGdsOHMvNUViVFcwbXRJdlR6NHhybkdWNTFEWE85Ry9EcG01YVRib0lseTk4YU0vQU9FUU1qM09IbUVoc1VzV1BNNUFEVFRmOUUrRkJ1RW5zNzhKOXE2Y2tSMEJEQ2VaTWc2MEJQSmhWMWhvUEdPIiwiMyI6IkJDOVZOMTJwL2FCM241Mmo0QS85VUxyN2JybDRvR3ZGN3BxNEluNXZoRjdmTlY1TVV4NElsc1lSSU9tNXc2R2lNT3ZoZDNXYmd0REZURU5tWmV0M0FMVDRkMnNEbFhwamplU1NqTkgwdW5sTlVDRnF4dDd1MzFlNGtLRytOTEp3RU1LemtEMkcrOFF5RVV3NDcvU1FtaXBRRXJOekRNTUN1Y1FyR3pyUjdFSEsiLCI0IjoiQk14Q0MvM0hXc0pNcUNxa25sQzQ2SXJNcU8wRU1UQ1VXa1NrTmNxLzBjWnV3K2NHNFRlVHBBakEzSkZFMUFFRXF5VlJ1K0oxRG1hSW9ZbkZ3ZnhzK1NQTHdvUlljR2hxQnN5WFRKYnZyODlKOHJRNm9Ic0ZVRzgwUlB0YjFuSW9sc0NjaDMybndSMmN2bUhzSXJweUN1Z1dXR2tPUnhJazBEcTNSQXJvTVJZVCJ9fX0="},"Signer":2,"Signature":"mQuSut1eo0vpKOWmxFmh9zAm9gWCyCLpcYcGfrRwp5920+SO6w3I5x2cfpu5NnNwpQiVeKvIok3sMw4XXaZnsQE="}]}}},"*frost.FrostSpecTest_Blame Type Invalid Scalar - Happy Flow":{"Name":"Blame Type Invalid Scalar - Happy Flow","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":true,"BlameMessage":null}},"ExpectedError":"could not find dkg runner","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"/TDrswVDScvL0OiZDiuYGJnFbTftJYgF7S+ug+wgkfw5/jRkdfipHvQ88Af4MOxU+TunnnvkLnEcvzFvKZ01wgA="}],"2":[{"Message":{"MsgType":0,"Identifier":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"HMiRihga9SNL/MJ+w2nyzRXQfFLeSQH69kn4yLxh+GNa90TqWEEufv9E/lBBkbGCH/dBpPIroE83P5opK0v4GwA="}],"3":[{"Message":{"MsgType":0,"Identifier":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"wVfKfXoZuUeVA/HNc4vXg6RFTbnONTt1LPVolhYLqzUCE4YD2LTOWFveTmv6jQfmYPNk3fLr6kTyN7qLlL/AxQA="}],"4":[{"Message":{"MsgType":0,"Identifier":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"BlyhrDUnQ8iKVoqP2LFuhL2nYJFbYGFxDnSxoZIBl7Br+WlInQRJ/fIFlFu6W/lvtmMKGIW3aIFlbBnZSBOpwwA="}]},"2":{"2":[{"Message":{"MsgType":1,"Identifier":[224,44,243,3,210,155,181,206,127,250,119,220,35,58,48,58,230,182,177,213,85,237,80,65],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsiazRNR2RqVUtiNXZhZVlwOVJLQWpmbGkzSDV1WmZIOEwzR01PeUE1Uzd3aUxpNFkzM3N2RTloN1RFWVl5V1Z5WiIsImdUekVFaThnL25uVlBNTkRWeHdlUmcxOWgycGdiSXZ2dkxXOHVVQ0tUcVdvTjF6aURqSWhnSkVnOG9GaG0wNHgiLCJ0cDFYdlZ6K3RtMlExVndveVgvZjQ0eFVDTndIZVlGRC94aEVVd0hHOStvTDcvZm94SWpNcG1STTNwN1BSWXZ1Il0sIlByb29mUyI6Ik9TRlBnT2xuL0hqYlZOUmVvS3B2U2N1UEFxQ2lQOWFud0habkJrdVpFREk9IiwiUHJvb2ZSIjoiY25WaVltbHphQzEyWVd4MVpRPT0iLCJTaGFyZXMiOnsiMSI6IkJBWlIzWUhSUk9kOEM5U2RtODZMVWQzTW9rTkE2N1ZtNlBWOVdHTzV6WVV2WkJOM1Z6RlNOaTNMV0ZDendMZ2w4cy81RWJUVzBtdEl2VHo0eHJuR1Y1MURYTzlHL0RwbTVhVGJvSWx5OThhTS9BT0VRTWozT0htRWhzVXNXUE01QURUVGY5RStGQnVFbnM3OEo5cTZja1IwQkRDZVpNZzYwQlBKaFYxaG9QR08iLCIzIjoiQkM5Vk4xMnAvYUIzbjUyajRBLzlVTHI3YnJsNG9HdkY3cHE0SW41dmhGN2ZOVjVNVXg0SWxzWVJJT201dzZHaU1PdmhkM1diZ3RERlRFTm1aZXQzQUxUNGQyc0RsWHBqamVTU2pOSDB1bmxOVUNGcXh0N3UzMWU0a0tHK05MSndFTUt6a0QyRys4UXlFVXc0Ny9TUW1pcFFFck56RE1NQ3VjUXJHenJSN0VISyIsIjQiOiJCTXhDQy8zSFdzSk1xQ3FrbmxDNDZJck1xTzBFTVRDVVdrU2tOY3EvMGNadXcrY0c0VGVUcEFqQTNKRkUxQUVFcXlWUnUrSjFEbWFJb1luRndmeHMrU1BMd29SWWNHaHFCc3lYVEpidnI4OUo4clE2b0hzRlVHODBSUHRiMW5Jb2xzQ2NoMzJud1IyY3ZtSHNJcnB5Q3VnV1dHa09SeElrMERxM1JBcm9NUllUIn19fQ=="},"Signer":2,"Signature":"0P9LgMIyiH8ZUEMN2AxqarJvG3NexC1Ys2i3jne5TlxJp5b8ntbiqbCNMjO1oB2io2XGfJJvLIQohE26CIT6zgE="}]}}},"*frost.FrostSpecTest_Blame Type Invalid Share (Unable to Decrypt) - Happy Flow":{"Name":"Blame Type Invalid Share (Unable to Decrypt) - Happy Flow","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":true,"BlameMessage":null}},"ExpectedError":"could not find dkg runner","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"vLpXCbiznmZTADQQZRyZt1OvZJQQn0A5IkQ94FewMywEfxu8qclRQegj1li6VdzMJR8YChqeEuPkddoNdJGijwE="}],"2":[{"Message":{"MsgType":0,"Identifier":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"y+9leL5Rt4tjdjwbyb7LHPoujHZZOZg/udg406yDWGxbIcZqSDYyAfAIJpiDk/yjztPsgHDRRnB7QSL6iak1tQE="}],"3":[{"Message":{"MsgType":0,"Identifier":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"tsHt/CWkpqIiI8srYnn1UzieoEjhxvyVbaQl/elEY9wVncnHnlLgysFLjg9qBNoga/fvcHdD6zLsLiMpB0a1HwA="}],"4":[{"Message":{"MsgType":0,"Identifier":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"TZlrRzJ6XAq7fyx/ruR/xXXOw/Hy40kQMx+kE9RApeBL6EZNeZZ6JYcAkMWZ1oN0nVWYq2uvCJDOdrEckff2+AA="}]},"2":{"2":[{"Message":{"MsgType":1,"Identifier":[65,155,116,205,59,219,82,125,159,25,154,132,74,70,90,153,205,152,144,236,54,78,106,188],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsiazRNR2RqVUtiNXZhZVlwOVJLQWpmbGkzSDV1WmZIOEwzR01PeUE1Uzd3aUxpNFkzM3N2RTloN1RFWVl5V1Z5WiIsImdUekVFaThnL25uVlBNTkRWeHdlUmcxOWgycGdiSXZ2dkxXOHVVQ0tUcVdvTjF6aURqSWhnSkVnOG9GaG0wNHgiLCJ0cDFYdlZ6K3RtMlExVndveVgvZjQ0eFVDTndIZVlGRC94aEVVd0hHOStvTDcvZm94SWpNcG1STTNwN1BSWXZ1Il0sIlByb29mUyI6Ik9TRlBnT2xuL0hqYlZOUmVvS3B2U2N1UEFxQ2lQOWFud0habkJrdVpFREk9IiwiUHJvb2ZSIjoiQzNMbEhzY0E2eUdsQ0VwUEltVGNCMmRySzdVVVpONlRldXI3dUYyT09QYz0iLCJTaGFyZXMiOnsiMSI6IkJBWlIzWUhSUk9kOEM5U2RtODZMVWQzTW9rTkE2N1ZtNlBWOVdHTzV6WVV2WkJOM1Z6RlNOaTNMV0ZDendMZ2w4cy81RWJUVzBtdEl2VHo0eHJuR1Y1MURYTzlHL0RwbTVhVGJvSWx5OThhTS9BT0VRTWozT0htRWhzVXNXUE01QURUVGY5RStGQnVFbnM3OEo5cTZja1IwQkRDZVpNZzYwQlBKaFYxaG9QR08iLCIyIjoiY25WaVltbHphQzEyWVd4MVpRPT0iLCIzIjoiQkM5Vk4xMnAvYUIzbjUyajRBLzlVTHI3YnJsNG9HdkY3cHE0SW41dmhGN2ZOVjVNVXg0SWxzWVJJT201dzZHaU1PdmhkM1diZ3RERlRFTm1aZXQzQUxUNGQyc0RsWHBqamVTU2pOSDB1bmxOVUNGcXh0N3UzMWU0a0tHK05MSndFTUt6a0QyRys4UXlFVXc0Ny9TUW1pcFFFck56RE1NQ3VjUXJHenJSN0VISyIsIjQiOiJCTXhDQy8zSFdzSk1xQ3FrbmxDNDZJck1xTzBFTVRDVVdrU2tOY3EvMGNadXcrY0c0VGVUcEFqQTNKRkUxQUVFcXlWUnUrSjFEbWFJb1luRndmeHMrU1BMd29SWWNHaHFCc3lYVEpidnI4OUo4clE2b0hzRlVHODBSUHRiMW5Jb2xzQ2NoMzJud1IyY3ZtSHNJcnB5Q3VnV1dHa09SeElrMERxM1JBcm9NUllUIn19fQ=="},"Signer":2,"Signature":"5/EgrPQem/ElG6d3NBeqWM4oBalvz5YWu9Lz3c/rlU0wWlmtm1pXYkUTlXZnz8kbF95A2tUmYbD8cBAw8YxNiwE="}]}}},"*frost.FrostSpecTest_Blame Type Invalid Share - Happy Flow":{"Name":"Blame Type Invalid Share - Happy Flow","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":true,"BlameMessage":null}},"ExpectedError":"could not find dkg runner","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"hMBV2Zihzp66yTdwmcm32Gx6eRGVH5RDrrvmt+APeYF+MY05MWIN85yHYEESPZNqGqi3y1uIQYjtvYAKFf6EuQA="}],"2":[{"Message":{"MsgType":0,"Identifier":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"X81kkO7g+GgAbEfxQZ4H9S0cLyJ4al3nXfyiGl2R0BwNwlFal7wvHUAvADdhK4quxZVUMRs9/Mps9dM7wDqY0wE="}],"3":[{"Message":{"MsgType":0,"Identifier":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"/BlcWnhVIHs5SYyRNKisEME5QJ4S0X3vjlhRuYo2W1BkxifBpCc5tsmfUxXxSl1lMfpCdjkPb3TL4XPKbZnLrAA="}],"4":[{"Message":{"MsgType":0,"Identifier":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"yWXAqU7pbSUUBIim958r9oz+H+hYmxDjmAWyXvc/jcwzMAo/EYyX/pqyfiH5HkftzCT3+bUor52UgnzAsRK21gE="}]},"2":{"2":[{"Message":{"MsgType":1,"Identifier":[30,128,116,184,176,33,80,73,126,200,106,242,102,177,182,160,46,22,155,54,152,143,87,19],"Data":"eyJyb3VuZCI6Miwicm91bmQxIjp7IkNvbW1pdG1lbnQiOlsicXpIclJBSXBtYTdsbWJibTM3U2F6TkNZWDZXRTIvUllRRitsUWRyK3MrU08vM0FrbkxvTUgwb2N1QUZqeCtGYSIsImx1Rzh1UGdWZVRtdlJvRTRNQlBNcHQvVmdwL29DQUE5VEJUekcwN2JSSjQ1TDZVbzl1REdROWdLRmt0OSswN24iLCJrQmIwT2JiYy9DWWFCSC81NnJadE9VdzZiejZETUZyYm91UlVOd2g4bEJmSDhPaldnM05Rd0JJWEUzSXI4bG1uIiwicDdDTitBb3cwVFB6S1cwd0ZtYkwzcU5ZdVNuakxFbzNHdGphcGc3MW1EWW45K0lHcG1HdG9xSEFkMUxVenhEQSJdLCJQcm9vZlMiOiJZSTVwOVYwMTNqZ1lYLzc4VEozUFNkSS81UUViWkZLblJORDBwVG82WEFFPSIsIlByb29mUiI6Ik1XNCtUS0k3QUFmL3EwT2xqaUppTkxTa29BUFh5NFB6VFhDMmRGcWhsQWM9IiwiU2hhcmVzIjp7IjEiOiJCTGgycDRiKy9zbGl0UGlNWG9vUEVrYStTNlRxQ2NkUVNCN0J6djFYVFpOcDBONXdwbkkvamdBNHFBd3pnMllDYlZkQnpjRzI2RkY1cC80RlJIazFzeUR6MGxqdUprdjMwYWhweHQvYmJ5MUl0TW5CS2d5N3ArellPRTlSa0FsZWNwbm93WW9oUjN3ai9GeHEvbG41Z05SV0RtTWNXTWVQcmZsbTVkcE1DemlZIiwiMyI6IkJNcXdpbk96cGpCdExlZDNiL3BDRHVHMng5WFFQelhsS0lYSHRSKzhwSzRSK3FQYlU2aEI0WGdmLzlEL2IyUEtzL2puSDZYT0tmTFg3cTFiQzlEWmtINTljbWVlZUFIRmp5M1llT2JYeUYzTDdFMU1YNE5XSHhta2pqV1NMaUgwOE0yTWtRQ3RmcnN3V3pJZk9WVDdZZ0ZKU1JSRHkyc2Y5NENBMVdkYkZBYzMiLCI0IjoiQktMdS8zRHlhWk9Jelh4NlNrS3JSaHhvamgzMFk1dUxYT3FCR2J0NWhRaVBacGJEdFVMQlNkVHI0WHJVeVhMZG1pM0pXK2ppWmloa3NIeEhqWldncTdtZGxoTXNwRll5RW5xeG1vYk1CdUR4eWRtRWE1Vnpkb0N0c1NyUmM3OWt4OVNrd01zTmttWTUyVmh0a2d3TGxJU0NYU2RtQVFjOEJJbjNrVDd4UWN5MSJ9fX0="},"Signer":2,"Signature":"fBK6Go3LC24f2aXdgc+EjFGhaK7BTWuB8MLpeCzvcsJZPr6ELPuAPX9vFlSjFcAAIpAgMHSe9WDiBM86mB+hsQA="}]}}},"*frost.FrostSpecTest_Simple Keygen":{"Name":"Simple Keygen","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"RequestID":[228,44,192,127,171,41,46,147,248,202,104,167,187,247,15,214,145,76,204,242,19,13,181,70],"Threshold":3,"Operators":[1,2,3,4],"IsResharing":false,"OperatorsOld":null,"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"","Share":null,"OperatorPubKeys":null},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812","Share":{"1":"5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912","2":"533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7","3":"442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186","4":"2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f"},"OperatorPubKeys":{"1":"add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa","2":"8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5","3":"b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051","4":"a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037"}},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedError":"","InputMessages":{"0":{"1":[{"Message":{"MsgType":0,"Identifier":[228,44,192,127,171,41,46,147,248,202,104,167,187,247,15,214,145,76,204,242,19,13,181,70],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"qCWvsmz+r8gfOqgy/gMNA3DnQ0sALTb2YlWJWXy3QTgRWj33m1e9GzljoKhJMlUVXC7f6pNdzu87/9ywtm79/wA="}],"2":[{"Message":{"MsgType":0,"Identifier":[228,44,192,127,171,41,46,147,248,202,104,167,187,247,15,214,145,76,204,242,19,13,181,70],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":2,"Signature":"00WLhomPY43OhHek+0kjInMQM5HxME6A3FIYcfY3HywYDH5Izl+AAYb9aL7XBjdVFSKmz96PDu1zl56Lv1CnpQA="}],"3":[{"Message":{"MsgType":0,"Identifier":[228,44,192,127,171,41,46,147,248,202,104,167,187,247,15,214,145,76,204,242,19,13,181,70],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":3,"Signature":"qRuua0NWv6JLqLp8GAMUkuBLmJESHBhB3U2+HAIeWvMkLsNHS/j/7nX7k/rjRspL72WOCmx2QrUsQwPfhh8KCgA="}],"4":[{"Message":{"MsgType":0,"Identifier":[228,44,192,127,171,41,46,147,248,202,104,167,187,247,15,214,145,76,204,242,19,13,181,70],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":4,"Signature":"V6EiOCoBIFgCXbwTRiozoV3eKsSEPnZU64SLGJclhi99T86G7f2wTQzp4cc2PZ1EIeq72iTAQoOo7IUjDLZ+yQA="}]}}},"*frost.FrostSpecTest_Simple Resharing":{"Name":"Simple Resharing","Keyset":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":13,"Threshold":9,"PartialThreshold":5,"Shares":{"1":{},"10":{},"11":{},"12":{},"13":{},"2":{},"3":{},"4":{},"5":{},"6":{},"7":{},"8":{},"9":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"10":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":5752531036591206989219466662491825132495410102509371470286419989861316574159,"Y":57775677257089512089392402793629146585925428208451298678042380497944182960548,"D":102851925695855034964344707731093447669058532883530396293021789007092655539951},"ETHAddress":"0xc414624fb41d49d0e08bb153a32f029c9760f956","EncryptionKey":{"N":23400218098004632769275808204594636181936915294826301785298025325650963881998575588170590209667677259909655187164324689545148269033391019762288436463821239791160854831105787409209305470272318841394092663824840235468440832651872869596737841376383900307004775039719785159318306767025258053856642674330894919239785327969170731775704993718500636938666811616248925643191903485349106262413264101234764362717216013204662325376967186949411180376531896928034693107881663631013263909140282203568822853198285728332416543177232229570066985560668049019460116402957766300317448534162006295129687589551788160674673640074073029057703,"E":65537,"D":9263397139549143118636671926697945299665400161877619871488959962291356134009354182515186725050249764729787664613748565627955608157259206504918003503629678580677742614678861536910847629899828188979177259720030747043551413130295706681992867489645603392967817929268811605251295463077395897725814998288611590166608417290858706200087758584787230897104901287325908052887684104210075165445302103015646456143271578528661520621164008689279786727346825805476228898181065504747350834698741313133777952607172659149199595910915590945438592060736003500016972587669548251960917164411730248164320257034537779230685486830567049568769,"Primes":[138837323323832246337453571266465683636549137302016056216317361298407530580936341437445947542787090505910792434302353718192105228623735986625168612648890605267125937594378703974727107988360290519518569976216474129823530555406314989645745137331234643619733691183285206995464176344645125446315694087797090174497,168544146039351416695961579118200267381329559293886303669228078351666349863904851790209958642606942971921048940613274586857011366296654511657149637844070821589988894547313866582670919430698080325190067048777277033243416156801378595779051635727549857870909870574636703538746359349892710010911339818573474822599],"Precomputed":{"Dp":118665374765162629943275734233960591252289545233569272150953156127535557402553506654539187811605945873149412365647593626875364670054451245720872481766446545524483259751729909869910721485084786814331334928632158908434396814945522897676219751706808498441481647186043317055265359078769194529431826362634145690593,"Dq":100156427170949859212722945797920402416434680969534798608390489515932160389547499167946760445798974534102342966449242088974873852334150326632096260827278303047471161300412568992509864003975414472199324980005661749235772803555696620422283076024374463505845322634988272568419014676924051930587944660792284922503,"Qinv":2073945438224409168226421702132585515312547265999901240512108534257558179775065053222869834189529340502708915361615811996182457296880449918936189132979974924218387080202762496339679035633551882433595684973041024456650561106651370339399948326411998078041852774474298154684352199043354189480434982935885627696,"CRTValues":[]}}},"11":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":78798287434531098844992901520196002100151572293177951790075975608559845583090,"Y":85221328845385140210433216425331105141017745637854596506954890988973725146389,"D":13064744025495022682023357934593847624009425267937008005374207955571209767969},"ETHAddress":"0x0220377e265ad7efe65d823dd075553ac11618fa","EncryptionKey":{"N":22205407347425074684931152456814111855371940767794175605086250616929494298647288296156852753164870620798633174611614957477155096021085869060751483563424372210766415325807980563212121184863248783347867477375826942141952203780266931686945117750866415049972848562165415490002661195259279625325531822112591261478137938667519519006417741840077795748384621259955361515071178756427358324856530494430273063464404683742099265178962391581582780463300408403646526683483323139401827090262328212497926440915844459899195276955305096726661303031346560304067342969643758394548869855889430322473387879004230925717323407611228747665969,"E":65537,"D":16696831845732338303534221308868923175013715444651558350901709631989083258968822496396135722160179125110027071908780894601778018170241093450184969689213568490690116740779887302966722917876558841488007711942923476537120445703766942743808359516684408307484505001677670780350659032930742033605701857910590731561811565937782040919793511959221458593955117945839935136115003173402341957331738212283981077576409145292684724436947357528671578552333842607965248268020821423808396738085779656077920312425982919489766235017666114255904650649373065088721859160097838398169673936798098957748921934544083238517298150095113941853405,"Primes":[140049310445174661877574298762851832742184041856316168952439609383174477987177981193401047980481386085303342018507319469214289119875178858371331200796299753503352266398194363501970641318227649212641733282574108275937731437563299883499503471619894881603743502571564722858471600392510390768322799596887825308407,158554206920696426349713326101633339696566465516569463581168113231215140977387838254521158503766919584053400271917987951639798724169501080853870461049397345351164670981004294411271379957168696995106685084837887332818924794782416922940767443490856854919478089433918859746187269820641511374898680970702000966167],"Precomputed":{"Dp":67675078390653163358116796123117547204056127402042765987988616957684247728396775659761932778336893906273884040833548991722955309784198838301443273235241419254736453065663934842057284280741434959109374114253619100518363930240812731900236132913170438157208187480947910465919055080800335157880506285515671142895,"Dq":139555378976966486985319186152045677511279735063518210128875313847588601738554405326362184828894990452515927510952854372375759488609388594961238135330781641961587849319746886778762810037830830179375528669816877380812775346174358572759131928686491399044435572438564101280178623866424843098549758351983981929785,"Qinv":100067327924910658841119711145924370358267843185601956179295641374575837122464093767665384361721964268875813868261737696722156169525456443717979601115930990620233367845751397517523084122843470912475614457848771832129186063663661476042782249637183073192728605304154040175770148927942917311928975019097218180193,"CRTValues":[]}}},"12":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":81563779267824983728263871861015781144926785758743742260521563865030876528972,"Y":22056235952619773370444860649320235303033479586069714831557148226942739114207,"D":37395777834316488348287446187096004746482299601418146853952905937361211367210},"ETHAddress":"0xd49ba2cfeccda87648453d16953060e6fb335959","EncryptionKey":{"N":20191299494337239393975326479332415118660164885626018828498308594230937733005708250154149466972880888205918400572384472599106755461273230823777415029245052277520127325888886510612769378323861365146006143607611130934417991115629573643340108104569280103914525098755294708741586572535001328323863919022005322276290442824752028820222578343182428398496038689772491336506779756420100359851411393023627570604731897211207040918983442172017978100701137357225421956683099285353388450655035049445449762031431746290213047993602957179491756235138156896858125657377238382470629695560653198003013156069124350069391059812828426664713,"E":65537,"D":4698373701705035335125558521290558471696408357199700735990344478111933723367518360847319519833627318082003381429251616752923966931419149031274021990570014605981078500996467938520907777582722520385073984009278266425833870401656331508322575775434968362675992305964848014225692284223549601857560229566284406742894700131762450250949364511165929846164863997003947259015941551856540197957565084139508306777621582861255560203069534941923067650648548169010711741346675992066361388205857741433926328835346172002289747047826067034007810310665424837772437423002651465461537695191817485446536515886681766033135870521106702458473,"Primes":[135202763709839835005735303452723629594345821134523641058994263975739939319466088648364989705117946655206141331072157662819655634375020169065969783805958543538718244482530549894225337493068217815160890936632406672252438507457447519002872063641193032327115239497632031080683344156696810630776521545983256761191,149340878398536388727105205744029730993057342628283005996784739393732554725719539178478587372729643934485775061447801788620480851863566279801746245491101102991514383297271791331209354841177818670851294612814684591975720066691860355357402354960650655589232755749175754220056786229446365634442371293614352918543],"Precomputed":{"Dp":38914349632401678558877962510165644216826300014656139910215737695872293137969220900775268944377066813512267023635718906781927220214779520714884538992199765121547251868013088219704480539721772321701937619019745900158654615959989510520029536543690192849632646636920106234263544575244090192232441611942599940283,"Dq":58125747687411175422295796086465818975096612535853684437279477737084852921916688364811202934275108068432536586468872972887547882407431659447074831469048124496201365475179041660077333770065211997132533118447243153823285586175381447799817191362684848601097839901887104057939919482745897654811114438523504497401,"Qinv":14788758218489420390628856633469451841354567653468233516046170167943442101400010760708377376934863308081175606981452006681287161723272577660531578668187102729532618930002048474041069319791509259461158032689899814075686743484887588449044893983043497254903582547805349381184308554819044283413086040975717364886,"CRTValues":[]}}},"13":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":114951801003422958901429389937212739040830142197550679050594224390921003661855,"Y":4442775891681942667754687163844440338340572288962343960450360327279148096700,"D":7905407607637422753896855115784349266545339149305311264332470027809669654876},"ETHAddress":"0x27bbc2c2cd860b40ce721b6b22376dd757f224a7","EncryptionKey":{"N":31399617507570831250368777867004837050726270822013594842124427640541285361888214826878050847809736628987568601767161910623478170233899571597106878497576612526394484730230927456632769198290388063403093577171440513847774019268505119113317036081119245480309907512189224504266694130566969794291182328828946840948510877581092635066217240407007232972002972100623943269154084011937706206777738140403807807632180016617623374983357059051437699978802640636338277376598953633427139016542509819326532129060643577249182046979695388854904091829233795976134689047602345588602484926206715597201917303761375775882047159056522381605327,"E":65537,"D":24963694868768338365644060637234234549812191813940161216626503773177947304512314346699634731133230347998356263278395474176811672518838263553346305098459371333372876884568293394577014607133716521286521930282752285480330462181152489834448797594803502853442596106246965127642580999726282774032191953572843526569155513652025068007584653944790645454448257042696298276479557725205575351846103583925879180662309960410088912838828577932920637610186400058288930314236812213116460980426302072550592860697027795790712152627219741979332819376277606540723346579403590180965740557827072750212520981078374439065440826024258149461953,"Primes":[176234763567508649457524566867756134717598475136727138523855160067654246032359506933167600695822201690767034249047284434348796556505123205400925066868983852944559540042175887208556286323583619768792155072855010195888824522166413960922063108912870272628786384317507011857360268913226854985088090502378348284747,178169260547411039570335647068992207541487900644837423704833975530465905280575728132271869321859736556750454800056962702652011301562332020602605917818831249640158486284639188360333352496366140254535785758941940341428614008604586455190632142159974899691805809309889016200576005811348951354662979671298064848141],"Precomputed":{"Dp":97460618215316782615759996291989038109311099598400928964097877021102474006298817611102024078286831192737989567545367193403778531171936163287085573012688737389103428746335362926281420956492380354308773308291257374148933688708322645005086184236830451361598866698497133386427639748875915974557084762465454275967,"Dq":15599359400354678197881897903197846817416994581687857808845954981061283923584288692234554315407039815106491136956632924726753449934611915928677735576002162296645091998432330787365805218794710053565563554706636765172610695963701681185953693817445656261830442861591821031766866370836630954621910941207383555833,"Qinv":107817911542002994222613774115488714946791537777653746556822047991431305980999861434968948449061475388607965455598312783226798765142442693988454424430565260112106063205525564268465540680376093258956301525797338243988810538738133795746822792069129642298745642446325326075001118187673000394185432990042856263471,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}},"5":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":61647188515213997628479099203157789667628270558658620417427470702362745997514,"Y":56371254570619516703840080939084966392891137828146363958637931423618288539973,"D":14508474625686945816214364974519687324260362558081141572415600147316062668428},"ETHAddress":"0x72d3c5cbd1c7e8edd7a332d6cc66e241fc048f4a","EncryptionKey":{"N":22659332576714543423444525302043067670597969979179003135432995952919864334752252976014905118799228325883171874917662167077873448913648084019665304020817165963835248026286743111919454402442071300548882588002121910256855781865702553806581673454690522081481671827083644744042731244991376505342624649648936770137258321168196987412737549860706531250877741789081920827050931494348962955641438030373983791664717691380305741144594915842407876242109486917965318809661353327880572258075974223963929543430112808349074198243308541863903228994380012152036478699104776969277840573203317960267681630715762403054898293278688070619199,"E":65537,"D":13819229805433626594010318320908179680579371379492585506204455730995530121850447826393331136507413475704174690153410121853846283467805976321497818581992482207127745065576007986943385766989748505602000216058971430353178643291427214750374031559432465280295725182060609987864014604438108664449690474416258812054320377018209113853393280981739126207720956122444284978181135163991488809571360951124099209842689004629891628920198947771921063428647061349160043757521356249355829444452568658475668023142871523911719463568138711389725748449850899039262416275378685490350890943179140885525292770094870182811863907831920739200641,"Primes":[137392690156919424286193879248190345081223814219061072313109196882358929590910780030626816075543366713485718370956963854496739853356839641542109789684560559176001221747923527488087183042363512040456129320653896132342462542760211239907068688149497845939404122760388552788090862788038429049458139883607801262273,164923858400580020029942993887045294812798057010746751784106786883146793873633091024382522501523662993292847098913079208297167214561587714264991708143193616510796003150138009880734771188260415861212367407352618451707470743499048070273447253831185176001091514539320204980876457494450507795599403624493327352063],"Precomputed":{"Dp":67378748823464459718300674108318014112799386438204721976033837188138242474508635887885406238734818593640706599974927419374173655902601981768518678616076054319188843965672248858921251552276779177872957205331587068274207640330091234731726927340660402039953743770982621825979833224095627045021661288419590957313,"Dq":18659847262467016319361998560697634329262822416874241489222146645994376864564892655229815285240367442746333540419007313876489538672416694405830500570392002478409941916143145753782570583959457766008357177251318580640110099074499007294774118240966752827381381209226228236464881400145727074848857559479653055763,"Qinv":75479896897792471976982626992234536479447652480205110218310832027117354589957030085183886138604246145644875871424109731379607264956390917691333765529812127333421644092841162975603765894284509697700679898269008850356909174628494900418551512595298834457030298321738680736981463510814582308028061906394600449080,"CRTValues":[]}}},"6":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":21337329131034664109062324470670568981747126982649708451748900250353921901023,"Y":15429239909490967527664370763399702592217836027669736132365224652894126125045,"D":1189967332525881379558525835614636113250210475501213154092172210369513409677},"ETHAddress":"0x34e59dd81b0f832122df98c8cbe8fe4f4f3c31b6","EncryptionKey":{"N":20874620858425994465014770052329677361009100734723572210490173530656718110729765197891188682895545590018027951853719794009446350305856174481365808596882674334061472403133604501714454193877212565228470779170866619307534829678745303884046866441551817318652530816044435539079324981153653426305459169009847225640140216167742980787228283846607901983394164276494774734902190175318964465449946220142267580335676322937313068601093870883382837801767157667088671762209640570551017091401126034669783697715810785651489469795460077995390726917724540310139917466155982328321222736552254296481710083632415486484128469207643793319669,"E":65537,"D":238250399043328865523765933734266119383475095740928513869213570974124924040250001800854618533131942281970259669905281076629474495762400148192038464233459578587332062156399227417831327906681035122005830947706001666875750379170567577174223051074671702310940278779944730201738484915436055401933006674387990368472902834963661323349278703715943728278516484173491172965400848874637288858384064527963647008949867371055681961634665086420490429256082877160183373811661802278250233537912648457627452416236096352419794378142123232504479030355696279313706947357768109865997744739826850362178424609910972841153926320923879247873,"Primes":[149071068622591147242835719070306415185222990186865206316615281410459910881158076902567723857884491291437922654472427768625768102753626599646625352873822659009496932664195119094115985554286057761848043041499866332958330189907165382064836035471121149722072177280234758287220400612733404753488577464179563770869,140031335733394789261279136318606102715445533262333321856329663030023271264087024296981228101391486863955054075430205456014180311104569141710846182738385376902903106097651180324958689430421724348817811566392925180285281470569691696280412702093887092652767942943099361475121482939821885085312526376634985875201],"Precomputed":{"Dp":113989769640852232371446194142077539261673023333758990331428936518360437217119426673410118153318330051407907836889562311318976192095995291433708942031505084663654565871698036197606367396546859006015728961986729655154702826449298655061983201574679422267908587396763426412648780937585209054676389994047555713769,"Dq":36453829423797984338414687042918331910649194241247979068172793703636526405489856440343572831192153400154678694186714913484870077784382167733169762480729861228625507623663377748815481336839114379895031556742446198343640800912301601083969682933668121637378486573267287274776807306350629132253795756775096419073,"Qinv":109298074695996507519249132048275190624170680182999123135016906725341363691702648192513714395021552253258405194926510902235940749241077470735444487110274674652674510184983390830698031139309843393142690312973900623674164210602906406282062923033550107844520740522339551436330634163552595143259545478253921272678,"CRTValues":[]}}},"7":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":84686251956080094017111693890923034335088704044501989111146319063496175417054,"Y":47279239187028622446719099085264629473270174191198734660952278928687026649171,"D":95637268851330345032489228872949866197317280190964721959370079616758909443051},"ETHAddress":"0x389d05783094371063d9bf55e46fb67f86c2ee07","EncryptionKey":{"N":25594384358443283245082975689390796369601033879263410600781606058952071281325973099814034261464898495529689148397343936576131615206708441929803320427894851845682080606201164346638694616393776296752233916949346725884497267265958551764325092593875982657780602350150406679601023501991032352041976119756734682560263542863789451415373734048094032458182326649520364541334044894299719337116488608635455419330978075681874222560244503022720771625793748205191452664840746369062563836540017745256614225447625735726219711388681340390211358958573465925076338144576848803960411606209403330552993930569446509604196352169533225992817,"E":65537,"D":4677808502760786223196116435105710650702369406683509043382548443995421056925439153296191500737393444610748380601849118395851267023909451415767950502545501554962539657004097632544032899204489868051608371717034271978340605571076667119842012284502594263925661152480302290435342775941049872785736148472560050470772002611396031758847830502704633776945040232137806237595418142377617898038395859583629265198744439128149333533695379534920122484738185680325165721024546524197432051980623303455016294059010003483118310905581813856679829375346431442753601354834913329781584503039902263887210853262135930807451790423562100892673,"Primes":[144659295817377199623741951906548627381110588402441592646721490145306146471020937336252157094329485145347106407617274930099582998841375726385957831440057145976179161871880122909217754078894522807858471743606862314326326454831749350718811947280764618077336314428083726446025233778088186839107592476918256544801,176928722166285822060706100549348596233917539360082959265807391714778511755640534086620931853436509459425523917045821465466399239551598148243599203179808812639974222952149190344440814018143585745784453719284251433271409380342227722817805710418832984188660232862796088630583850386232049820693369060707354971217],"Precomputed":{"Dp":57215215937290910347544366317799822548266895982112219402714005005668257971456333318476466180052711971139087007214968391337889908188768179862526709320806891997627905685048211940275468872225230414918281979126805896663757999842741885041767619596025140991846355589824958011618171182108181501388649230117309731873,"Dq":54549670569174227086357743987368017081992428709123644276132629766220831111196310965626478920770528253309613443822998741645392114902720481306897866845464651543453608632850550682816898668700265400908199518620894829801213023775806847540421169487815116323848645272527850937174073895726151619343732780576663786081,"Qinv":129246060617987305608981594123425389754311639431906846129118977666546076423236123971108413452560338656668696675379964458454593111196368617877700176735006010468253887680173746748499700905492929069593124677298593078446791162230123250950193243479811572676520435012989901372631102378882971545984497656558018853451,"CRTValues":[]}}},"8":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":92700592012992630948569047095095806452589468573027273985711815836306277150172,"Y":107431301766463047387920998330044903238058481763878759335023269923441261580301,"D":112711874566949502491766762738367591310775262050094043527054080580845969374145},"ETHAddress":"0x4f1bee79773989347e553cfbb853736d027fb84f","EncryptionKey":{"N":25755960177229090907575975728159180971962469219482809575969964447683280109322627535012945301220234022660409781137238291855951823346911371791609970472690320373888814285835907660224115944347592470960275732261477174133242313672670133535454163961914932659437606455356484185819440715527805982402461705537732733936459858772075526160265029411141509984280152693069738640103143653537392876127220442198987116926124318963293781304645057611084802344477454628937799551686790801243130130180206654246628816680640460718932858217235458806418166258904814344847999840558305495920072353080744975026243514869229359556891668368725633914099,"E":65537,"D":10386957405498647675164152135362423563620062887695968034741995519359584559705159615948123110780944889435198903145661352423101556236307239520457932459423000251489713620926240741256441161620258312212644576402197868568008824791624145587104285419128304166942886287365486321180521203463690924438058850380125366707846325246170328834928645651400052027140386509242949085734089157768477422977441600432060649569562468199165098105340444034332921842689490852524565925917163333444733650060368561134716659243759395962038102369422201004991147981343028214810589467032861878924174614080305114294888968572636381118299137034644861952513,"Primes":[173585878384379006182919466989622809376669853146539323791562974090436196184977439003692362074708537163694446337392446586823237093797409861874379257675713954072249678147856428347400818483463910747149953858941294104641519296293975366147513158242161895239861442924480502680437709308977057360990230349916208210547,148375895648587908727639643137690818388177858604426927671588203749919038259587302570400362462448113302865205775814976492363090078362936859310772655462717657463445369777095803232337128092815614172914530730415890495794624062373304354635348509749640802380792382612527950380458157684926138391513600567446568428417],"Precomputed":{"Dp":61020056856727093770573854166454526182457544238979094885484959596799809080511928220197211307767143433132316931212096746375844731478473662173458494260236173061270550760186099398315761420602737015622329935796413225853049751255330645060140197744524859888855576576501546008391045471416351793376152811409884565277,"Dq":66360161090006564119746793109371127420783392794213309687381507241922531248405685729298030488682982880209378618107523627380785408042724602640616099367803183818469677182911257587973718075720867098315421979022844596521601933140423332449100510692764721586026298529926709394717626072338832146629463451675028872321,"Qinv":121299055587614159411969163158848545433483060113050930089103071412089162242157191215370166237623020820589661811480780517675241460818174418496923689122252166975940489639334085853039820539802307345772740521136571241750382727021236369147408904215588016662966117124355759389616221218975047516728878402188904810713,"CRTValues":[]}}},"9":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":32200814053536123300844119660835150442739815291100337450143125843168116249484,"Y":61550536938257152131767073321067027677395754299782333694755508040485838779292,"D":20938562768359604350968999947654998691985284294835737895012916351084466681050},"ETHAddress":"0x1e3fff6bba6b65aaeabd2b275491d7979bea3c30","EncryptionKey":{"N":22635054255258547360323239576778930055377462501171466680371110113701351833051019310526597754914514729058613188899552405228017534573897712330923010877881820784730696417517698746932793566618961144456941496011358146994259075958116129979867896471261683857737152450458015465019676936346022185760464933513390401694574773428053812794331305464045369318862182934680839429733780402301543262576993568420641897512560928719132305860881063880780395216850648192527993686269992602536827683055839652824203420683022396015829588996899549422337448179257119942050611861195094249951681078825719972998137234325031157413531702353372496279183,"E":65537,"D":5483225069143914092687973992110445909321033838421017211919552987856060877084974633778175167569813019188161542136034514631429671466426599950649766096971966778741543499466118151674703307500230818154605843884772295675433069714986216634273505414922112591748707330263384859280290386215869634269697137257710698037723263014257151932372057210449671859155058363447586996989555782869456395379681787595815951781625201271304713161065673446101967703290277688939642154739340325052486164344911473096068454962413822645348766129676870414864250488695016933368065570172035600096600717568295755553030662974401508892720531274150115810817,"Primes":[136526459930408190115172728176833053070284927543426551041306642302211121335823783049321876124194759818051460326096368397737886932315057755785554142830604392620542383618996930128912768789050743448946440184523316699937037148251858931192968412409948336486537394485341794143867347389809197353247930038873138309327,165792435157231374695629246246456104955096548000066631419576472148767528996384870174356678486352386102044754174444493880801096255653319782923005939562331842886505395640391273193083651839063949448451422568950114280793055004449118105708807893710524023597990671788513796662346823231937021265736943257181636888129],"Precomputed":{"Dp":26833656260792955076880844586199956613795873349205447365046780589510985457478159657267117603121178894614047949409453001072092429851080442380238993451043153964099767206254473915353546466450442137508263973279900551015288699614449774824856586664823603800648308258932933310452954845783790394848506749328240452911,"Dq":146637143962783109300580576753220747420879372031430524592760581169466252930031112190923979865746312165723548159385920744016899529410023058991593134973077273493084597376982165046265374961492619890740905278074879446513716424811833483336299335626142866287097811773825506561681700545329042184848893234394824026433,"Qinv":41683808126710022583303557828235433556458852254035344859062441502026475229781184291409997620992664879011157181249716945638551789966730507634212091882847314360521808642333581277076961706380747591807984460897454793908423539807155200381938696379061907053058751541142167512955345957081802575348750305363593347219,"CRTValues":[]}}}}},"RequestID":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Threshold":3,"Operators":[5,6,7,8],"IsResharing":true,"OperatorsOld":[1,2,3],"OldKeygenOutcomes":{"KeygenOutcome":{"ValidatorPK":"8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812","Share":{"1":"5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912","2":"533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7","3":"442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186","4":"2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f"},"OperatorPubKeys":{"1":"add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa","2":"8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5","3":"b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051","4":"a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037"}},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedOutcome":{"KeygenOutcome":{"ValidatorPK":"8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812","Share":{"5":"52046f0837c928ea5d5bbc893b90f3cd75a07a9d25092e2fbb0129825100c3be","6":"0f1d824d53df922ca8c15d639c802f84463a78cf69ef57e0b1cbb8b95cd1f458","7":"213989136198ba32e82eb8e449a843b7fb6a52007ba72794212d25d135a84679","8":"146adc07375723b4e869f703396758634172622d5a32414b092570cadb83ba20"},"OperatorPubKeys":{"5":"81e52afe4656f4544715cc2a37724c939afa8462d57549ba242681b52c80d8ac7e6b259d03ba37ce688aeca5e1a346b3","6":"ac6d0b0ba2f3f581f520c59049c6dfb98ce12d87a3ee9ccc00b9e0ef13153b036c777a946d9ec78409a047d92ce942e7","7":"8d9b4d117564b4852ee7d060626e27bc93ec5dddde0fbbbe053aed7e54b0772b334ad74149fb1c6d3f1ff3d5b4d87fc8","8":"b11cb28641e5d6440e214d45abfc6a2158cbf163312144609e08236fee95aa096a61a0d70b4401d8daf4af69a1cca9ad"}},"BlameOutcome":{"Valid":false,"BlameMessage":null}},"ExpectedError":"","InputMessages":{"0":{"1":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":1,"Signature":"I8VnvSFG5T25YOyNNMExlh5zNvn/JKh60Y09lxMNCShgch2oHER4Dv5jjQcZo2aCvawTO/1j6kaNUduNj35w4wA="}],"2":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":2,"Signature":"30Pw8CW+k3PmfFpRxbxEGUeWlpEpK9C9SeVquAW7Exc8ETtX21uUCNJdJwz2ZOsI1solwFZygXRTLOkr7GxawQA="}],"3":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":3,"Signature":"8DFLZYHOZQBe27VbJHnnSaWtGqG0BK2QTRGrpAdFvs9pp7kVfkZX8rZBtzYdR/A/NDIeXN3R8+HBYffvVEeZ2gA="}],"5":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":5,"Signature":"GJyUrZ+RtCfdRnj9Ad3SGVh892h0rSB/K6pmxf73JKZ5nf8fyvhHqUeSVJ45JNk4H7i6oEMKwVePBIHTFwA7dAE="}],"6":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":6,"Signature":"gLwQjg/b/CmAtJk3Uz0VcmxoOCPWYZxfCZPJcM7B4rtufmU/n1en235ZSMcS9v5fVvBSvHXHJRIVGzGVC9d5nwE="}],"7":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":7,"Signature":"LsrJn/RmdMGwHvnOiIGHkYPAsKD6h7qrRmRABHNHBPFx4n8QWUVtpNjbmSnktd0t4BW0Gq+jBEIsIW8lo54WAQA="}],"8":[{"Message":{"MsgType":4,"Identifier":[45,203,140,62,4,239,76,55,248,192,162,1,156,193,103,142,196,52,131,128,9,40,107,131],"Data":"eyJWYWxpZGF0b3JQSyI6Iml1YmxKVlJ5NVVqUU9kVXpNQUhHWVFuSmlXUnoyWnhXOWsrT0o5b2IyUFpGN0U1cURGZHJlTWNpaVd2T055Z1MiLCJPcGVyYXRvcklEcyI6WzUsNiw3LDhdLCJUaHJlc2hvbGQiOjN9"},"Signer":8,"Signature":"Kucp6St+3pl2g0yJJD41SJecW0I9zzV0QMNEcSLXOMNAP3aIn29dAEh/brcXwezltF8nVfmCtB7GO9ggwMHlDAE="}]}}},"*tests.MsgProcessingSpecTest_happy flow":{"Name":"happy flow","InputMessages":[{"Message":{"MsgType":0,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjMsIldpdGhkcmF3YWxDcmVkZW50aWFscyI6IkFRQUFBQUFBQUFBQUFBQUFVMWxUdGFZRUFIU1VqUEdGNnFmU3E3MW1nSTg9IiwiRm9yayI6WzAsMCwxNiwzMl19"},"Signer":1,"Signature":"NTe9goo5TjnylaqwjA+QhEE6IZWUqBnzDwFCPPME9VwxIFIT9Xbxu1TRk3inE57zWDi02YOgB4tVfy6sALrHpQE="},{"Message":{"MsgType":1,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":null},"Signer":1,"Signature":"4+wsZ9iFxzIMgiK12GXE96RSSKERSrgPKMzU8LDv6p0g9z76/+xK3ZSBVnkLK5ctFsoAvo3/ZxIMS/4zgYzZTQA="},{"Message":{"MsgType":2,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJTaWduZXIiOjIsIlJvb3QiOiJINkFHZ2pQR3dQL3QyUHNjYmVvUDBUMW4xYVZZdUFPczdQRldBZU50Mk5rPSIsIlNpZ25hdHVyZSI6InN2emxxRTk4akdRZXdPV2FFb1FLai94T2U1VUhqTGlSbjVZTnJDOWF1aFY3Y2N4dFZpMHFweEp0RGhuaFFyR21Gc2pBT21KMDMxR2xDb1ZzRHVjeVlGekw2blFKR3g1dmY0R2M5OVBGZmdkS25TSjZyU3ZndktWRWdBaDRMa002In0="},"Signer":2,"Signature":"3ssTTegT6B2wINZNDv8DH6TwAOvY7LmjqiKD6KkmWJpDt/HmWLOna8SRfV0zZZWPYwv7R9O6gHwWJNdXMGxV6AE="},{"Message":{"MsgType":2,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJTaWduZXIiOjMsIlJvb3QiOiJINkFHZ2pQR3dQL3QyUHNjYmVvUDBUMW4xYVZZdUFPczdQRldBZU50Mk5rPSIsIlNpZ25hdHVyZSI6ImhseWVYcEJCcU96Nk9nVFpRTjVBak5hcG0rU01nY3BwdFBJdTFYcWV3N1FHdkdEbVhuRVl4cmR3NUIxRmJGUzBDcXJReGpZaGRoTnJzTkhOaEZ2RXZPRSt3UlAvS3lrZVZXcEU4VG9yN1laS0llc3dOTUVSRjJmc08yd3FrSjRlIn0="},"Signer":3,"Signature":"O/tNIKTjNEsk2uqc0oDPFlhEniIZxLgjA+wgRMlFi3I0dCCPvBNEEZMLKx4ZS9KKr+KqhobsMaYXNCRVy8UW7gA="},{"Message":{"MsgType":2,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJTaWduZXIiOjQsIlJvb3QiOiJINkFHZ2pQR3dQL3QyUHNjYmVvUDBUMW4xYVZZdUFPczdQRldBZU50Mk5rPSIsIlNpZ25hdHVyZSI6ImxLY2xlc2ErbTlJUkFGbi9LWkJ3RHRjNUF4VjNsMWlTRHlyeDBhaHR2ZXlucndJbUpPYW1KSVFLZUc2SXRsT0ZCcjZmMXZIaVcyTTNCOTBNeGNOcnpPczdUTGxSTDA4MnA3ajFPbzRkZnZvV3RyaGdzMU1CTkg3b1dueFBVRVl2In0="},"Signer":4,"Signature":"7RPzaNmzXiuH2B27mYUAkE0SyfHoF+ldySILLauauaUN6vBO9xKAZ3M2KJBaLaCS5AVsfrGrxztfXF5quqJEnAA="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6ImpJTUlkT05LNDNLSW5pVklGY3ZSNE9ueDZMWWRRUWxZUkFXK0ovZ2U5S3dBRE9iOVY5TTk4VUZCMm5VdHQ4RW15QS9hTVA0UWVvNUpmNFVkVmx2U2lVQktPWWtTQXhBa2daVG0yaG0yc280Zm5QMndVd012Z1BwMUNReHpiaDgvdzc3azN5cnczanBYdEg5U1FIRUFuNSthcnE4UUdQWlJTMGNZZlNPdWluMDZNM3lXZTkrRUtoS245eTQ1bHNLMlZiNHNmejlmTXEzQ2N6Qm5naGwxNkFneU5UT09wcCtOQ1pYR25DblRPZU1GUmtRb3RTOGpxMWtwaGRYVUdMR2d0T08wY3ZrMFI3bDVQSUg2V0orVmVRdjk4N0RucExCcTEzQjRVei9QY0hzcVNMWlUyTHNJMXFTVFA4QUZMZ3A2Yk5uVWlaekMrS3lzbVZydWFaQ0pvZz09IiwiU2hhcmVQdWJLZXkiOiJwcnpyNHdsOWRCY2JRTWNTb0RIT3NEY2RzOVBFQXM4czVwRzVFZzg3cTNYVTFXMzZEemRaRlVTWm0vR01VMVB0IiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjoib2p2Zm5Kc3JYdDkybDVaS05QN2M4M0l4VFRqdnVrVGtNRGhLbG92UHAzUFR5MkIvNnUzNXJkbkljVzlUWnZ0Y0R5RGNheUV3SFJGbU5KR2dUeEI0WnlZR2c1U1FacDh0cFpYR2JGUjF1UHZ0S3BYMDlZYTBUTUY5RExoWkpSNmYifSwiU2lnbmVyIjoyLCJTaWduYXR1cmUiOiJLK3JZNFhObU03OGg3cFFwaXVIYzFKell4TjZWNE5DMjdDUFFJYjJzRFljcmg0cllrOGRrOUk5eS9odXhjMGNOTllzT3djVTd1NHgzQmpkN2RKN3pzQUU9In0="},"Signer":2,"Signature":"LTMw/6NdQ0grZLg1PwR8v0syRleff7YeRjgaNTFxSk8mycuauk3/FUgBHjTwVFeEPUHwL2TIYFziD8afa9u9qwE="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6ImNsc08ybnlVZ2liWG9wQ1loNzc1NFZITnhsWDF2Rnp3Rko2dVpSZCsyaVU2V05CVjBrMlhOb3FPaWZUVFBHWDFRRTN1YVMybmtnWjByV0wvd1FWYnNnTjNWa1pTREwrSlBUaThHS0xyOWtXWmx3c0h2UTg2ejZHVEFORFY4YTdSSEwyMVRLQVVYZkt4MlhvVWx0Z0tOR21vWFFkRUJmOW1EbTFBUFE4NHZjZk0xbzNlTXZRc1pnaXc0Vjc2anZ3T2o3eUZYV2E5R0NDcCt6dFd3U3E4YlpXVkpCNUtKOEVCTTVHVEg1bW1Yem5YZnpjZUltampNQVZmK1FVbVZ0U2dmR0NpMU96a3grejErd2Vjdmg2MVUwbDVkdjBDem1MYnpwUlJsMlJ4ZU84UlFmWkJIeXFsZFd2THcrSFdWSysrOWFaWGZoaUY4LzhJT0R6Ync2OS81dz09IiwiU2hhcmVQdWJLZXkiOiJnSkRndDJacVJlekYxTzkwR0t5WjhKNXNza1FDbitwcUNuL012cDdnaThVNTNnMzZacjVycThoSlBkbWQwYW1OIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjoib2p2Zm5Kc3JYdDkybDVaS05QN2M4M0l4VFRqdnVrVGtNRGhLbG92UHAzUFR5MkIvNnUzNXJkbkljVzlUWnZ0Y0R5RGNheUV3SFJGbU5KR2dUeEI0WnlZR2c1U1FacDh0cFpYR2JGUjF1UHZ0S3BYMDlZYTBUTUY5RExoWkpSNmYifSwiU2lnbmVyIjozLCJTaWduYXR1cmUiOiJhVVhKcDQxMEtDMjhZWHNMYXlLWVFEaVQxK01qbTBaU05naEovekgzK2pFd0xHeXhYeFJQbi9naWxNY2c4Y3BTUURNUE1uM0FFemdQZDNKZ1ZjQ0dqd0E9In0="},"Signer":3,"Signature":"AZ2lp6D0QF6Pt0m11G8/JgIs747Q3lJlAYDZ+2ZhGckKVtVapWMhRlj+tRfGs8gpSCISIqVM3je8+DzQhlnSmQA="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6IlhkT0dyRzNhZ0VPNk5PNlVxVUlnVmJGNFBhYlVuQzVUb25oVFJuSWFYYnlmdWdUK0l2REFVeHM3aldSaTNuZ0tiOFlRazZ5alc1Tm51Z0YrSmZrS2xMd1daOGs5MktvaVgyUHgxM0dkM2VqRlN6eHNqNVppMy90TWtpTEVQSDhTVmRUcFJpMVYyMUhCYVFVcHNaVityU3BFMHRUZk1yd3J4VmRvdFFjTE1GM1FiOGc4RnVWUzhUSlV6L1Q1RlgrdmF0RFZqd3k3Y3JxeWFYQlprOXlOcHRSb3dyTWsyaFhxaEx5ekJZYmxneXk1ek52eUlPbjdNbEhiSTZ6Z1ZpVnFsc0VtdWJVN1RRTnp5ektOZVRGcWRiRklTRlJNVXREbndaaXRCSzZpeDNtaVBkSk1kNlUvUllXSEN1S1FhY2dOUlEzYXRzaC9vU3VwdW0yRS9jNHhVdz09IiwiU2hhcmVQdWJLZXkiOiJwOENpZHJjS1h1TTVYSDF0SmxYdFlGS0tvbExVMGg3S1g4eFNJK1VNeEN2UmFMS0FxM3ExTVhOVTNkL1BQZm5rIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjoib2p2Zm5Kc3JYdDkybDVaS05QN2M4M0l4VFRqdnVrVGtNRGhLbG92UHAzUFR5MkIvNnUzNXJkbkljVzlUWnZ0Y0R5RGNheUV3SFJGbU5KR2dUeEI0WnlZR2c1U1FacDh0cFpYR2JGUjF1UHZ0S3BYMDlZYTBUTUY5RExoWkpSNmYifSwiU2lnbmVyIjo0LCJTaWduYXR1cmUiOiJkdm9STG52K3Z5QnBoanc5SlMzV0NDVnFIR2hIQWRReWZsRXhveHFGR1VoMGNheWM4ZGJRdm1tR0NDWUdZQVhuQSt1clNTU1dDTndHZFovOXpwa3hKQUE9In0="},"Signer":4,"Signature":"AZbU9FH0tL4AW6QsQokrl7Y4t1VEVO+Mz1Thq8CR+CMeif/MRxUb4iuj6eqzHCVqzPT9Yc0wPyT3F98VZ9WRBgA="}],"OutputMessages":[{"Message":{"MsgType":2,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJTaWduZXIiOjEsIlJvb3QiOiJINkFHZ2pQR3dQL3QyUHNjYmVvUDBUMW4xYVZZdUFPczdQRldBZU50Mk5rPSIsIlNpZ25hdHVyZSI6Im91dWpRL2VMdHcvS0x3L2tXSUR4ZHcwSGw3bHZKdk9USkRvNVFhNi9SNlZTQmFHZ3phWlZnU0NnWjhzUnBsUDBHUG9ubEd3eFhXZWoyR3hJMnlmU2wzNVl4WnJOaitZQ0dkaWVhTVNvRlJNSHFEY3I0QlhEMGxvQjBPV2dKajBrIn0="},"Signer":1,"Signature":"F31F7XhJzTfOOFbqEkkdwnAzH7PUgV78spDtkO0w81huzCW+BLR9nXPeY9zFwTAY+96flDysJ6b7J2zxvf268gE="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6InRlOEx2dVdQUFhIVG1ORkVuZXhWRUVMdVlkSVkvMloyTXRPQzNOck5aVkIrdXZmN1VQR3dId2dlRVpjTWxtUEtUMlU0RnM5a1F0VmpMdElmUlJSRkc5OHhqODR1U0hSazRCbnB2cHF0REFnU0xDR1E1RWp5ZUtLaWVpS1M4N2VtSmJkNlZpQUU2NEVzNDRvZXRkaDAvWnZHOXN3TmdEOVFDRWlISXN2b3pOS0R2SGJGczRpYXZCeTJ1S2ozSWo4cGFwK0k2K0dOYnRQc1lnZURUM2h0cjFCVlVpMGpxZ1BiZjhzcmdRbDFVTSthRjNhb2FxVHZDbm9rRHhQTks2ZnEvVVBPSk1YNi9WbEZzZkJiQVZIc0ZtMUR6MjFLYndHaVFjcFVxSHliRm45eHBLdjh4RUo1TFJ4aTYvUWwvS0ROZ3FaODV4UW9xcE9ISkZPampoQ2N1Zz09IiwiU2hhcmVQdWJLZXkiOiJsOWxLZ1Ixa1NUWUZLcDB0U3Mxa2NZbDB6MmVOdnYwbWN5VEk2ZmpuQTBwS2EzMkhlZUo2QVpVNHc4UWx3K1huIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjoib2p2Zm5Kc3JYdDkybDVaS05QN2M4M0l4VFRqdnVrVGtNRGhLbG92UHAzUFR5MkIvNnUzNXJkbkljVzlUWnZ0Y0R5RGNheUV3SFJGbU5KR2dUeEI0WnlZR2c1U1FacDh0cFpYR2JGUjF1UHZ0S3BYMDlZYTBUTUY5RExoWkpSNmYifSwiU2lnbmVyIjoxLCJTaWduYXR1cmUiOiJkTG5tMFh2dlVxU1RBeVp6WHlLSm1kLzNSbmowMjc1TjNDVjBiZUVWTnJWYXoyaWJGRUhsa0RlaFA5cG5Ib1hRV3BjWm9hT0kzcGdOa25jcWl6K3kvUUE9In0="},"Signer":1,"Signature":"tl2hxcEL1KtUtzAzXG3yx8ezxEkBTIR970dsIUFukXB0cYF3m0JjZ4AEYvox2aywa5iYiH/nTy48P6pCa5FawwE="}],"Output":{"1":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"te8LvuWPPXHTmNFEnexVEELuYdIY/2Z2MtOC3NrNZVB+uvf7UPGwHwgeEZcMlmPKT2U4Fs9kQtVjLtIfRRRFG98xj84uSHRk4BnpvpqtDAgSLCGQ5EjyeKKieiKS87emJbd6ViAE64Es44oetdh0/ZvG9swNgD9QCEiHIsvozNKDvHbFs4iavBy2uKj3Ij8pap+I6+GNbtPsYgeDT3htr1BVUi0jqgPbf8srgQl1UM+aF3aoaqTvCnokDxPNK6fq/UPOJMX6/VlFsfBbAVHsFm1Dz21KbwGiQcpUqHybFn9xpKv8xEJ5LRxi6/Ql/KDNgqZ85xQoqpOHJFOjjhCcug==","SharePubKey":"l9lKgR1kSTYFKp0tSs1kcYl0z2eNvv0mcyTI6fjnA0pKa32HeeJ6AZU4w8Qlw+Xn","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":"ojvfnJsrXt92l5ZKNP7c83IxTTjvukTkMDhKlovPp3PTy2B/6u35rdnIcW9TZvtcDyDcayEwHRFmNJGgTxB4ZyYGg5SQZp8tpZXGbFR1uPvtKpX09Ya0TMF9DLhZJR6f"},"Signer":1,"Signature":"dLnm0XvvUqSTAyZzXyKJmd/3Rnj0275N3CV0beEVNrVaz2ibFEHlkDehP9pnHoXQWpcZoaOI3pgNkncqiz+y/QA="},"2":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"jIMIdONK43KIniVIFcvR4Onx6LYdQQlYRAW+J/ge9KwADOb9V9M98UFB2nUtt8EmyA/aMP4Qeo5Jf4UdVlvSiUBKOYkSAxAkgZTm2hm2so4fnP2wUwMvgPp1CQxzbh8/w77k3yrw3jpXtH9SQHEAn5+arq8QGPZRS0cYfSOuin06M3yWe9+EKhKn9y45lsK2Vb4sfz9fMq3CczBnghl16AgyNTOOpp+NCZXGnCnTOeMFRkQotS8jq1kphdXUGLGgtOO0cvk0R7l5PIH6WJ+VeQv987DnpLBq13B4Uz/PcHsqSLZU2LsI1qSTP8AFLgp6bNnUiZzC+KysmVruaZCJog==","SharePubKey":"przr4wl9dBcbQMcSoDHOsDcds9PEAs8s5pG5Eg87q3XU1W36DzdZFUSZm/GMU1Pt","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":"ojvfnJsrXt92l5ZKNP7c83IxTTjvukTkMDhKlovPp3PTy2B/6u35rdnIcW9TZvtcDyDcayEwHRFmNJGgTxB4ZyYGg5SQZp8tpZXGbFR1uPvtKpX09Ya0TMF9DLhZJR6f"},"Signer":2,"Signature":"K+rY4XNmM78h7pQpiuHc1JzYxN6V4NC27CPQIb2sDYcrh4rYk8dk9I9y/huxc0cNNYsOwcU7u4x3Bjd7dJ7zsAE="},"3":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"clsO2nyUgibXopCYh7754VHNxlX1vFzwFJ6uZRd+2iU6WNBV0k2XNoqOifTTPGX1QE3uaS2nkgZ0rWL/wQVbsgN3VkZSDL+JPTi8GKLr9kWZlwsHvQ86z6GTANDV8a7RHL21TKAUXfKx2XoUltgKNGmoXQdEBf9mDm1APQ84vcfM1o3eMvQsZgiw4V76jvwOj7yFXWa9GCCp+ztWwSq8bZWVJB5KJ8EBM5GTH5mmXznXfzceImjjMAVf+QUmVtSgfGCi1Ozkx+z1+wecvh61U0l5dv0CzmLbzpRRl2RxeO8RQfZBHyqldWvLw+HWVK++9aZXfhiF8/8IODzbw69/5w==","SharePubKey":"gJDgt2ZqRezF1O90GKyZ8J5sskQCn+pqCn/Mvp7gi8U53g36Zr5rq8hJPdmd0amN","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":"ojvfnJsrXt92l5ZKNP7c83IxTTjvukTkMDhKlovPp3PTy2B/6u35rdnIcW9TZvtcDyDcayEwHRFmNJGgTxB4ZyYGg5SQZp8tpZXGbFR1uPvtKpX09Ya0TMF9DLhZJR6f"},"Signer":3,"Signature":"aUXJp410KC28YXsLayKYQDiT1+Mjm0ZSNghJ/zH3+jEwLGyxXxRPn/gilMcg8cpSQDMPMn3AEzgPd3JgVcCGjwA="},"4":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"XdOGrG3agEO6NO6UqUIgVbF4PabUnC5TonhTRnIaXbyfugT+IvDAUxs7jWRi3ngKb8YQk6yjW5NnugF+JfkKlLwWZ8k92KoiX2Px13Gd3ejFSzxsj5Zi3/tMkiLEPH8SVdTpRi1V21HBaQUpsZV+rSpE0tTfMrwrxVdotQcLMF3Qb8g8FuVS8TJUz/T5FX+vatDVjwy7crqyaXBZk9yNptRowrMk2hXqhLyzBYblgyy5zNvyIOn7MlHbI6zgViVqlsEmubU7TQNzyzKNeTFqdbFISFRMUtDnwZitBK6ix3miPdJMd6U/RYWHCuKQacgNRQ3atsh/oSupum2E/c4xUw==","SharePubKey":"p8CidrcKXuM5XH1tJlXtYFKKolLU0h7KX8xSI+UMxCvRaLKAq3q1MXNU3d/PPfnk","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":"ojvfnJsrXt92l5ZKNP7c83IxTTjvukTkMDhKlovPp3PTy2B/6u35rdnIcW9TZvtcDyDcayEwHRFmNJGgTxB4ZyYGg5SQZp8tpZXGbFR1uPvtKpX09Ya0TMF9DLhZJR6f"},"Signer":4,"Signature":"dvoRLnv+vyBphjw9JS3WCCVqHGhHAdQyflExoxqFGUh0cayc8dbQvmmGCCYGYAXnA+urSSSWCNwGdZ/9zpkxJAA="}},"KeySet":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"ExpectedError":""},"*tests.MsgProcessingSpecTest_resharing happy flow":{"Name":"resharing happy flow","InputMessages":[{"Message":{"MsgType":4,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJWYWxpZGF0b3JQSyI6IkFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUEiLCJPcGVyYXRvcklEcyI6WzEsMiwzLDRdLCJUaHJlc2hvbGQiOjN9"},"Signer":1,"Signature":"Di3zf473P2LYkFCUTp9muLRDXnDeMHRh3zZ6ycd9pIFyOltYuKitRKR29o3AhtkvjnpiYFaKHz5yh/3kY0i8sgE="},{"Message":{"MsgType":1,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":null},"Signer":1,"Signature":"4+wsZ9iFxzIMgiK12GXE96RSSKERSrgPKMzU8LDv6p0g9z76/+xK3ZSBVnkLK5ctFsoAvo3/ZxIMS/4zgYzZTQA="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6ImpJTUlkT05LNDNLSW5pVklGY3ZSNE9ueDZMWWRRUWxZUkFXK0ovZ2U5S3dBRE9iOVY5TTk4VUZCMm5VdHQ4RW15QS9hTVA0UWVvNUpmNFVkVmx2U2lVQktPWWtTQXhBa2daVG0yaG0yc280Zm5QMndVd012Z1BwMUNReHpiaDgvdzc3azN5cnczanBYdEg5U1FIRUFuNSthcnE4UUdQWlJTMGNZZlNPdWluMDZNM3lXZTkrRUtoS245eTQ1bHNLMlZiNHNmejlmTXEzQ2N6Qm5naGwxNkFneU5UT09wcCtOQ1pYR25DblRPZU1GUmtRb3RTOGpxMWtwaGRYVUdMR2d0T08wY3ZrMFI3bDVQSUg2V0orVmVRdjk4N0RucExCcTEzQjRVei9QY0hzcVNMWlUyTHNJMXFTVFA4QUZMZ3A2Yk5uVWlaekMrS3lzbVZydWFaQ0pvZz09IiwiU2hhcmVQdWJLZXkiOiJwcnpyNHdsOWRCY2JRTWNTb0RIT3NEY2RzOVBFQXM4czVwRzVFZzg3cTNYVTFXMzZEemRaRlVTWm0vR01VMVB0IiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjpudWxsfSwiU2lnbmVyIjoyLCJTaWduYXR1cmUiOiJGMUdGQTIwVlkybUxYY3FmT3JLYjQzcTY1TnlsaUpRaGVIWGx4MkxQVTl3Q3VoN004YTF1OUllSXFDejNQOGFmMVJrclJ1cllNZW5Tb252cmZSSFpqQUU9In0="},"Signer":2,"Signature":"w77UNx+/3oh00oz+8wPHcbkR8zYwR8TZCPKrGAtFFFgdXwZGqRqsq5FCu3eNLYY55faHgrMhBMsV6I/jBsdjcgA="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6ImNsc08ybnlVZ2liWG9wQ1loNzc1NFZITnhsWDF2Rnp3Rko2dVpSZCsyaVU2V05CVjBrMlhOb3FPaWZUVFBHWDFRRTN1YVMybmtnWjByV0wvd1FWYnNnTjNWa1pTREwrSlBUaThHS0xyOWtXWmx3c0h2UTg2ejZHVEFORFY4YTdSSEwyMVRLQVVYZkt4MlhvVWx0Z0tOR21vWFFkRUJmOW1EbTFBUFE4NHZjZk0xbzNlTXZRc1pnaXc0Vjc2anZ3T2o3eUZYV2E5R0NDcCt6dFd3U3E4YlpXVkpCNUtKOEVCTTVHVEg1bW1Yem5YZnpjZUltampNQVZmK1FVbVZ0U2dmR0NpMU96a3grejErd2Vjdmg2MVUwbDVkdjBDem1MYnpwUlJsMlJ4ZU84UlFmWkJIeXFsZFd2THcrSFdWSysrOWFaWGZoaUY4LzhJT0R6Ync2OS81dz09IiwiU2hhcmVQdWJLZXkiOiJnSkRndDJacVJlekYxTzkwR0t5WjhKNXNza1FDbitwcUNuL012cDdnaThVNTNnMzZacjVycThoSlBkbWQwYW1OIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjpudWxsfSwiU2lnbmVyIjozLCJTaWduYXR1cmUiOiJESWZmMmJwQklsdE1nc2IyK2QzSndMZjJSVmc1SWpndnBoTkl6Y0dzcERBN1djR21EQWxudU5IQU9CamNFUHJrMVoyTEdqYjJ0Z0xnaldRbVgvaU0xd0E9In0="},"Signer":3,"Signature":"uKNnQQOXxgdE0HmWVGtqfqxZVV7RTjzJm4vzBjpaAflNejvKv10sc/jNQbldBkoX1E44UJtFs24WHMiBALTkJwA="},{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6IlhkT0dyRzNhZ0VPNk5PNlVxVUlnVmJGNFBhYlVuQzVUb25oVFJuSWFYYnlmdWdUK0l2REFVeHM3aldSaTNuZ0tiOFlRazZ5alc1Tm51Z0YrSmZrS2xMd1daOGs5MktvaVgyUHgxM0dkM2VqRlN6eHNqNVppMy90TWtpTEVQSDhTVmRUcFJpMVYyMUhCYVFVcHNaVityU3BFMHRUZk1yd3J4VmRvdFFjTE1GM1FiOGc4RnVWUzhUSlV6L1Q1RlgrdmF0RFZqd3k3Y3JxeWFYQlprOXlOcHRSb3dyTWsyaFhxaEx5ekJZYmxneXk1ek52eUlPbjdNbEhiSTZ6Z1ZpVnFsc0VtdWJVN1RRTnp5ektOZVRGcWRiRklTRlJNVXREbndaaXRCSzZpeDNtaVBkSk1kNlUvUllXSEN1S1FhY2dOUlEzYXRzaC9vU3VwdW0yRS9jNHhVdz09IiwiU2hhcmVQdWJLZXkiOiJwOENpZHJjS1h1TTVYSDF0SmxYdFlGS0tvbExVMGg3S1g4eFNJK1VNeEN2UmFMS0FxM3ExTVhOVTNkL1BQZm5rIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjpudWxsfSwiU2lnbmVyIjo0LCJTaWduYXR1cmUiOiJmaXhMTlF5Y0hra0dtQk5HNEtBb0UyOGJIZ3NROEU1N0U5bjArMElQSzAwZ3V2NHB3N2d6L3hhVzh3THdWbitSaUtDVVVrKzNUSDRQa3ROTVVDVm1EQUE9In0="},"Signer":4,"Signature":"TasGYXiVa1XleK5m/j8e6D62oyzn/ICuvLt4/tzb+JcqvCejNhdQwg7kkpKWhHqSWLUKFZjBObrzHrnY8LHULgA="}],"OutputMessages":[{"Message":{"MsgType":3,"Identifier":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"Data":"eyJCbGFtZURhdGEiOm51bGwsIkRhdGEiOnsiUmVxdWVzdElEIjpbODMsODksODMsMTgxLDE2Niw0LDAsMTE2LDE0OCwxNDAsMjQxLDEzMywyMzQsMTY3LDIxMCwxNzEsMTg5LDEwMiwxMjgsMTQzLDEsMCwwLDBdLCJFbmNyeXB0ZWRTaGFyZSI6InRlOEx2dVdQUFhIVG1ORkVuZXhWRUVMdVlkSVkvMloyTXRPQzNOck5aVkIrdXZmN1VQR3dId2dlRVpjTWxtUEtUMlU0RnM5a1F0VmpMdElmUlJSRkc5OHhqODR1U0hSazRCbnB2cHF0REFnU0xDR1E1RWp5ZUtLaWVpS1M4N2VtSmJkNlZpQUU2NEVzNDRvZXRkaDAvWnZHOXN3TmdEOVFDRWlISXN2b3pOS0R2SGJGczRpYXZCeTJ1S2ozSWo4cGFwK0k2K0dOYnRQc1lnZURUM2h0cjFCVlVpMGpxZ1BiZjhzcmdRbDFVTSthRjNhb2FxVHZDbm9rRHhQTks2ZnEvVVBPSk1YNi9WbEZzZkJiQVZIc0ZtMUR6MjFLYndHaVFjcFVxSHliRm45eHBLdjh4RUo1TFJ4aTYvUWwvS0ROZ3FaODV4UW9xcE9ISkZPampoQ2N1Zz09IiwiU2hhcmVQdWJLZXkiOiJsOWxLZ1Ixa1NUWUZLcDB0U3Mxa2NZbDB6MmVOdnYwbWN5VEk2ZmpuQTBwS2EzMkhlZUo2QVpVNHc4UWx3K1huIiwiVmFsaWRhdG9yUHViS2V5Ijoiam9BR1pWR29HekdDV0hDZTJ2ZmRIMlBOYUdvT1RiaXltN3Q2eitaV0NHZDY5YVVuMlVTTzVIZzFTRjRDdFF2QSIsIkRlcG9zaXREYXRhU2lnbmF0dXJlIjpudWxsfSwiU2lnbmVyIjoxLCJTaWduYXR1cmUiOiJBdTRBTGloR3hmRE5TRHNwVnNYTnd3aHZJaDdNc085N0pWTUluQVRvWi81VkdIUXpZTnJjaGh4QzMwN0thbTd6M2czZFpLZmtLT0ZNNkttZk1wb2tNQUE9In0="},"Signer":1,"Signature":"vHZQ6jg76wtO4Fcm7xl78osdsIc2E2oTn/lpaPG/jKpalG3vvUtYtAhE9XDFcyPBcgipI048eZaCGW5xGq2CYgE="}],"Output":{"1":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"te8LvuWPPXHTmNFEnexVEELuYdIY/2Z2MtOC3NrNZVB+uvf7UPGwHwgeEZcMlmPKT2U4Fs9kQtVjLtIfRRRFG98xj84uSHRk4BnpvpqtDAgSLCGQ5EjyeKKieiKS87emJbd6ViAE64Es44oetdh0/ZvG9swNgD9QCEiHIsvozNKDvHbFs4iavBy2uKj3Ij8pap+I6+GNbtPsYgeDT3htr1BVUi0jqgPbf8srgQl1UM+aF3aoaqTvCnokDxPNK6fq/UPOJMX6/VlFsfBbAVHsFm1Dz21KbwGiQcpUqHybFn9xpKv8xEJ5LRxi6/Ql/KDNgqZ85xQoqpOHJFOjjhCcug==","SharePubKey":"l9lKgR1kSTYFKp0tSs1kcYl0z2eNvv0mcyTI6fjnA0pKa32HeeJ6AZU4w8Qlw+Xn","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":null},"Signer":1,"Signature":"Au4ALihGxfDNSDspVsXNwwhvIh7MsO97JVMInAToZ/5VGHQzYNrchhxC307Kam7z3g3dZKfkKOFM6KmfMpokMAA="},"2":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"jIMIdONK43KIniVIFcvR4Onx6LYdQQlYRAW+J/ge9KwADOb9V9M98UFB2nUtt8EmyA/aMP4Qeo5Jf4UdVlvSiUBKOYkSAxAkgZTm2hm2so4fnP2wUwMvgPp1CQxzbh8/w77k3yrw3jpXtH9SQHEAn5+arq8QGPZRS0cYfSOuin06M3yWe9+EKhKn9y45lsK2Vb4sfz9fMq3CczBnghl16AgyNTOOpp+NCZXGnCnTOeMFRkQotS8jq1kphdXUGLGgtOO0cvk0R7l5PIH6WJ+VeQv987DnpLBq13B4Uz/PcHsqSLZU2LsI1qSTP8AFLgp6bNnUiZzC+KysmVruaZCJog==","SharePubKey":"przr4wl9dBcbQMcSoDHOsDcds9PEAs8s5pG5Eg87q3XU1W36DzdZFUSZm/GMU1Pt","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":null},"Signer":2,"Signature":"F1GFA20VY2mLXcqfOrKb43q65NyliJQheHXlx2LPU9wCuh7M8a1u9IeIqCz3P8af1RkrRurYMenSonvrfRHZjAE="},"3":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"clsO2nyUgibXopCYh7754VHNxlX1vFzwFJ6uZRd+2iU6WNBV0k2XNoqOifTTPGX1QE3uaS2nkgZ0rWL/wQVbsgN3VkZSDL+JPTi8GKLr9kWZlwsHvQ86z6GTANDV8a7RHL21TKAUXfKx2XoUltgKNGmoXQdEBf9mDm1APQ84vcfM1o3eMvQsZgiw4V76jvwOj7yFXWa9GCCp+ztWwSq8bZWVJB5KJ8EBM5GTH5mmXznXfzceImjjMAVf+QUmVtSgfGCi1Ozkx+z1+wecvh61U0l5dv0CzmLbzpRRl2RxeO8RQfZBHyqldWvLw+HWVK++9aZXfhiF8/8IODzbw69/5w==","SharePubKey":"gJDgt2ZqRezF1O90GKyZ8J5sskQCn+pqCn/Mvp7gi8U53g36Zr5rq8hJPdmd0amN","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":null},"Signer":3,"Signature":"DIff2bpBIltMgsb2+d3JwLf2RVg5IjgvphNIzcGspDA7WcGmDAlnuNHAOBjcEPrk1Z2LGjb2tgLgjWQmX/iM1wA="},"4":{"BlameData":null,"Data":{"RequestID":[83,89,83,181,166,4,0,116,148,140,241,133,234,167,210,171,189,102,128,143,1,0,0,0],"EncryptedShare":"XdOGrG3agEO6NO6UqUIgVbF4PabUnC5TonhTRnIaXbyfugT+IvDAUxs7jWRi3ngKb8YQk6yjW5NnugF+JfkKlLwWZ8k92KoiX2Px13Gd3ejFSzxsj5Zi3/tMkiLEPH8SVdTpRi1V21HBaQUpsZV+rSpE0tTfMrwrxVdotQcLMF3Qb8g8FuVS8TJUz/T5FX+vatDVjwy7crqyaXBZk9yNptRowrMk2hXqhLyzBYblgyy5zNvyIOn7MlHbI6zgViVqlsEmubU7TQNzyzKNeTFqdbFISFRMUtDnwZitBK6ix3miPdJMd6U/RYWHCuKQacgNRQ3atsh/oSupum2E/c4xUw==","SharePubKey":"p8CidrcKXuM5XH1tJlXtYFKKolLU0h7KX8xSI+UMxCvRaLKAq3q1MXNU3d/PPfnk","ValidatorPubKey":"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA","DepositDataSignature":null},"Signer":4,"Signature":"fixLNQycHkkGmBNG4KAoE28bHgsQ8E57E9n0+0IPK00guv4pw7gz/xaW8wLwVn+RiKCUUk+3TH4PktNMUCVmDAA="}},"KeySet":{"ValidatorSK":{},"ValidatorPK":{},"ShareCount":4,"Threshold":3,"PartialThreshold":2,"Shares":{"1":{},"2":{},"3":{},"4":{}},"DKGOperators":{"1":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":53041238799374731371326291537220174474216562670347883252464550957073162636898,"Y":14698682645889972921328461223239899604299337271079813034849521020827685508393,"D":68257473393804446436778947025748892447965294342786393602221588243716150424833},"ETHAddress":"0x535953b5a6040074948cf185eaa7d2abbd66808f","EncryptionKey":{"N":25896833610471254564895404070857711037012918451830051295794631697757991199624079874475728854047197796820899014015994942258523510944399304501832381929700140187495257262581601359754008850155163173378995624727982137478477943464916944301595928086826737676096931009534817553384351681646106671991634927263619419233946414603004991157731382886392257229705216868482381372181577959203601011552180533295771004842184156625586135437848852256668549035153835642352113915505509531693398604035745274840314769582271100422242864110759758920861139373340329984704934075866339054716097083755922954387795160682576632729200104234487702166581,"E":65537,"D":5801172074329439679375458552638388326203315009709278469773730685792530460681464159744543926427314507150580869200128503704737532449985903983877827778353109817242441244975517487259847169219340991323634508241626955160009379251544099658082149323934622210702031587515154134324666479061392679730674159133881573672411613416405441206371990603774840525862573148646187606861616632638277695623897108841629163058102199071657566355036593711058012297605942363185562372171286232269732750791785572960514237780500296462297754368004641255582942431792802693133800270976776665299014194472194177555599224548444483694755515011090572351329,"Primes":[166054689456372045285169997337396030892288237437713204868128290879676839961190393842081801790549657090052298587347253469482001789608588550643167132467556544863522271976018457825583442823475535889507264182849696546801240656642661313935804217485137429358686529114867093411409714231631841100192239283852064716297,155953642111836853608180861291648642212189141716611686767862225151469615491041578582232867697173371508885657148162981387717236439183238736958439243702104897172395554957264376096999455536082675134901124891065148950782240644876613558725681294999619150478213411181878024694026861938213633812565375523359372877773],"Precomputed":{"Dp":131919944558737973019399360840006780115156126801570685436609845806954463472227563901124387906449301866023359719703903930429839986205825150514007305063144964040454813165561453937302622192109095260497058298061697220031533252790029469003275196962993122351648902732281844125685441376167841171880143099527865929993,"Dq":139367639006423074069156789344461693828543898300453128140338845849613515578434046917399825479047756979430036196293106656307664167319906970222086930831456696411121816182951656543219358719223553635743994712834163580885049481186026615406365509624222878466477331166959889409876424357772828959221757911967108523925,"Qinv":103990109935850922939952580571852026358340239656737020792470113288479826931926340522638802380112107102570872029577872471667392281237506643565004766017209812965452808042878171302604901859192917523403065267751088850912474781181517958842109140089457179251013089713706412453432562936940768311650043201486734136052,"CRTValues":[]}}},"2":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":94125038009610038565383212209811193818233918925264499382590849783322448284000,"Y":111171420883821742863295546917007988422701979908606848028330923489925287048864,"D":44387670994153663870544682707827633584660022672745599673663750571073877789078},"ETHAddress":"0x01139beb7ceb4c6a9fd3da779ed69612e489f4e6","EncryptionKey":{"N":24317377749361936601155508429702946010480136106507533969263487918266775357171582810281394462195143636207193494159426414459085276532605994927450224673463708888433041689704597393376070765040625259791319529057927480174121202937607833553013222775826146399261535688309549328081006491432584740367259193716664501850426189312584167963463824692984964486659573623202040188582111952897447064087302444370200227702125510337219524758020720255181361368528183061123486758670329181392608695517207958154156543659575810411436616305043715632146339592292549818798729512633813064568572803686676971975818134491964315688366300188293478078801,"E":65537,"D":23272877415741937951045604768723441409727865127673917720068732001915386483246349650220022635424322125672331512590865367162103036676657662295184903061919080044864690798505451236817887766069214299474055014748482954842016419589737683081497403893149939033744023092942178509176448252981286602763556939565384910585514263342007099898680043930088199244271359606584386885442261340149830823698444310808505361868398458825942838810231983523657975031836896805015686540587917762225492428399180835520114805133225159554966936107211223952909687936944558559952632776387423364688268155681137583919039758668446136190570560868445879781505,"Primes":[163406027290948663907893668790402242058038860174933232130561347047083047157150796775937169472508485701495772275359680321539362819004302860358001284935082894180064647024198418091800303254877945556002072356977407988470315320483355004586376409232305474611839146165249476196375460040095049665322513276423157243433,148815671933962482812299777643267421862908849683767086689274507361803979100921954704096674811211524235930763944250697403270315413439337202928792625763077568374303920946102323814759480352145693659460753385611316026670153631148253926746779711830725965461739395239059047112274716417145334587161276028363136382697],"Precomputed":{"Dp":108624852906868936506268147375111220798945953924975833361307896996402338105931483167378408002186622641734666142001004514826493134759623699808607107122416671007960334044222779232912278737232594963040515804874769312383810019564182738450189582138557155605831579746330449669214234586672886060079659023157167225649,"Dq":119080694379676526597839768972766988683257791707220555719043192625047290416261793316633929237675736667541711136676916447538060651411961511756591587656855117577631661843775242465990458346082739015967176237060418314397000117867414322084541886992491738723843590111337634445609513379433535810609451721614043947409,"Qinv":28805848433652365147443931333797444492534107073720398042138556439263654047144376521572761673349950984375528514358830428608073903683918252403992311150679244746311467147592686305042532629660711280412127653716251256407200552802260232894828128370629359059940454801550645075377139992084142549857964971262848820682,"CRTValues":[]}}},"3":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":50750364279652102928167696539394611227386326218496618610013733881994024004418,"Y":91863337803764030036769191672233324798321760346182590049947391152047968433470,"D":10939879519921015761001744320992943834490175148861685812970765826246835969090},"ETHAddress":"0xac8dea7a377f42f31a72cbbf0029048bda105c37","EncryptionKey":{"N":27499073670632760320623366385029487404514692606001946763007429417067863574018363264592065905324876271219368092745514019890087438301715784632682968801466090404703134387757854888819938269686973815348027339399349259939073863032853170142467515612501634838106592295114737667190447791698479977666329448135203226419096014975012523517884550721642484738306325221375896185609017114382569221062967813440108214857883380696969823222477063565208512446792048613691262812686076429726354051727197615265213590346859430075219203652004868518910757532557294449359293353973627359880902529965608950519668989005891472868519963379777303399757,"E":65537,"D":20978549312764180810079900653864523578490335020252366332149510178450981508311276197259708547362983336621370318034048925834943644853607642770957632957976412133053734683991172480832666336108452169704980741992298471842987563209386452654423430704460750980374678349311862598936796286701388581158482588742477618921509302071143816584053366569458624004180448300393356685049635356042466886040029342250216849750583869130296420071308752621009800599816097788691665822817368273814769209971863528000631111714741569426267578107150911311153157645177946509795242633120702126734943564996540110213724958936276926756362811040048588215213,"Primes":[170530157979976024040438300825902050644751968755755961778499604956288892349447439101137153775488572008257747152180390897802882734679103362556930878954394460337675239174477870247842220336399480855303951888518543251560809041437746617387785059673345712874138761261838592986821967104335032352246357207854721547791,161256366594475060115979035561270541693334353028576388638446938612928481394573177446108255671989883769372186032627178087268435442253801686970768059299578751975663437351687282378065788306577351554269347681074395059556591037831808959390137083090979983310382057417528198941793706026468316252182769779105726323427],"Precomputed":{"Dp":29743965025391853926884816466131900162048304848361176115629781409819465774242545071716721925135570237062946239476540707581743939150660398513637744744612708487113625265170156320903984324356965769824365992304430595062203154747316501874662725134290779923772528189939682262422172299153970349856235550040241726243,"Dq":64387393396771646140576153967489014374035634070094556325295321568321902148911163272804077611496273532156359981411843633178821408561537490953981088175402853070771637832353522518107718516357418488367186324036113491897353773867933790825352200899106828253750975456640949538563829581772478137344076164298619809137,"Qinv":163027213165549568291660986613407265669650857012741329205142474998647418165597928346776184779163000114385492067675429801989482645193971854833418407928603398428856406868945275453927676353736613540802354925665741042925013327898648629306554823769171522789096830523943031609879575434579443691872573440145421116833,"CRTValues":[]}}},"4":{"SK":{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":64949903245160046251401446495645361564584726999976810552286449245156967535302,"Y":68338799682262770740603469450564556414933717526454783745184309748523517095535,"D":34938880395882804457709048647740716294681363478614601290879240504729463375611},"ETHAddress":"0xaaaa953af60e4423ad0dadacacab635b095ba255","EncryptionKey":{"N":21875024364809584008994940350655860736099486615616500802014231577190201200340631711768862068402316814678008739885772877191575638796423608125902893902464164004311974207666791949208837591703929424060954432313915740736690368931757768282899887343026305325374067186655452051890113930934795368540295081897633717967821288283758997143301932502146703898431477758533315993859026443757978794835316195778173961417623589102881529618250759029787482474012615626234881744522762150828958038355015676814550627353738677667387240372423352616696968876586610676322475470502781893890329625115013132576432115235769991447724735785565843567207,"E":65537,"D":9818176628330163321857670787715979424635952191866569588038033810565783729923854949138365774174193953091438837355082002267271883290306245830957071946243852849334524334628052629598211052687353464588751005180490890852033922854687501015327222579537036653277967961540353115131112215671254493883039807040586169858473190914534530444186449875281595016328254918446327239411468276209275030264016525612246271897656870745902917074380957282989462339252283935726855680122383653363196669473671573449976718210442608346278500650419642801210561454064329261987888148198811206546837374622905025897025095729754451941067171757928407772993,"Primes":[152791666589048934445316674380789349740560164196070174857816658863039927021194023205129202452892456035093768303048915806016276848217690640533112734181137117961138026237521129352546226783403901368204056749525340418373334114563765877662782715006858970841373770260254963565790916839520927440791847396159871983239,143168962373092124523882002291500386355209874119093573911650158244343672899525842841622562080259321472539096474008411387793629785999686763851320575591305213372194646579639238633477316383430910135999294548393435272845594101984974771704050206847087633453168636746264996064411937854317480907775341712562386330913],"Precomputed":{"Dp":88037580507738618834002903061894310464364144229549749652575990234312124817650009983247462395686786468669772474475993082789670664546690174524488503717718233188099762680337410723878886976744405808415423210943068875270669130936065536602255228101515303674442777554171657753198904461510128050096613231820026638861,"Dq":105697347688476575872614047009657974783869791863805537027042453217179978784055699529259289312773959902457110392578588836642003502842803979159593736811415100535579379283599568519190174647846577597695803400055967960714729466262432204008861638587202466652396528988682508652829809132906556544269681758710629113857,"Qinv":135770051041387167489011774540676657014450447331791408161056908496147566070699800251299551498615096769203903360215262850792784854812857830793226628507676094349856100083071927750942688973600588380496018065843083493141545925612480849243271138945926545007935292844335184696535927368108249788422009727757591852151,"CRTValues":[]}}}}},"ExpectedError":""}}
\ No newline at end of file
diff --git a/dkg/spectest/run_test.go b/dkg/spectest/run_test.go
new file mode 100644
index 0000000..c2d7e95
--- /dev/null
+++ b/dkg/spectest/run_test.go
@@ -0,0 +1,11 @@
+package spectest
+
+import "testing"
+
+func TestAll(t *testing.T) {
+	for _, test := range AllTests {
+		t.Run(test.TestName(), func(t *testing.T) {
+			test.Run(t)
+		})
+	}
+}
diff --git a/dkg/spectest/tests/frost/blame.go b/dkg/spectest/tests/frost/blame.go
new file mode 100644
index 0000000..6bea0c9
--- /dev/null
+++ b/dkg/spectest/tests/frost/blame.go
@@ -0,0 +1,261 @@
+package frost
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/dkg/frost"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+)
+
+var (
+	maliciousOperatorID uint32 = 2
+)
+
+func GetBlameSpecTest(testName string, data []byte) *FrostSpecTest {
+
+	requestID := testingutils.GetRandRequestID()
+	ks := testingutils.Testing4SharesSet()
+
+	threshold := 3
+	operators := []types.OperatorID{1, 2, 3, 4}
+
+	initMessages := make(map[uint32][]*dkg.SignedMessage)
+	initMsgBytes := testingutils.InitMessageDataBytes(
+		operators,
+		uint16(threshold),
+		testingutils.TestingWithdrawalCredentials,
+		testingutils.TestingForkVersion,
+	)
+	for _, operatorID := range operators {
+		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
+				MsgType:    dkg.InitMsgType,
+				Identifier: requestID,
+				Data:       initMsgBytes,
+			}),
+		}
+	}
+
+	blameProtocolMessageBytes := []byte(data)
+	blameSignedMessage := &dkg.SignedMessage{
+		Message: &dkg.Message{
+			MsgType:    dkg.ProtocolMsgType,
+			Identifier: requestID,
+			Data:       blameProtocolMessageBytes,
+		},
+		Signer: types.OperatorID(maliciousOperatorID),
+	}
+	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(blameSignedMessage, ks.DKGOperators[2].ETHAddress)
+	blameSignedMessage.Signature = sig
+
+	return &FrostSpecTest{
+		Name:   testName,
+		Keyset: ks,
+
+		RequestID: requestID,
+		Threshold: uint64(threshold),
+		Operators: operators,
+
+		ExpectedOutcome: testingutils.TestOutcome{
+			BlameOutcome: testingutils.TestBlameOutcome{
+				Valid: true,
+			},
+		},
+		ExpectedError: "could not find dkg runner",
+
+		InputMessages: map[int]MessagesForNodes{
+			0: initMessages,
+			int(frost.Round1): {
+				maliciousOperatorID: []*dkg.SignedMessage{blameSignedMessage},
+			},
+		},
+	}
+}
+
+var validRound1MsgBytes = []byte(`{"round":2,"round1":{"Commitment":["k4MGdjUKb5vaeYp9RKAjfli3H5uZfH8L3GMOyA5S7wiLi4Y33svE9h7TEYYyWVyZ","gTzEEi8g/nnVPMNDVxweRg19h2pgbIvvvLW8uUCKTqWoN1ziDjIhgJEg8oFhm04x","tp1XvVz+tm2Q1VwoyX/f44xUCNwHeYFD/xhEUwHG9+oL7/foxIjMpmRM3p7PRYvu"],"ProofS":"OSFPgOln/HjbVNReoKpvScuPAqCiP9anwHZnBkuZEDI=","ProofR":"C3LlHscA6yGlCEpPImTcB2drK7UUZN6Teur7uF2OOPc=","Shares":{"1":"BAZR3YHRROd8C9Sdm86LUd3MokNA67Vm6PV9WGO5zYUvZBN3VzFSNi3LWFCzwLgl8s/5EbTW0mtIvTz4xrnGV51DXO9G/Dpm5aTboIly98aM/AOEQMj3OHmEhsUsWPM5ADTTf9E+FBuEns78J9q6ckR0BDCeZMg60BPJhV1hoPGO","3":"BC9VN12p/aB3n52j4A/9ULr7brl4oGvF7pq4In5vhF7fNV5MUx4IlsYRIOm5w6GiMOvhd3WbgtDFTENmZet3ALT4d2sDlXpjjeSSjNH0unlNUCFqxt7u31e4kKG+NLJwEMKzkD2G+8QyEUw47/SQmipQErNzDMMCucQrGzrR7EHK","4":"BMxCC/3HWsJMqCqknlC46IrMqO0EMTCUWkSkNcq/0cZuw+cG4TeTpAjA3JFE1AEEqyVRu+J1DmaIoYnFwfxs+SPLwoRYcGhqBsyXTJbvr89J8rQ6oHsFUG80RPtb1nIolsCch32nwR2cvmHsIrpyCugWWGkORxIk0Dq3RAroMRYT"}}}`)
+
+func BlameTypeInvalidCommitment() *FrostSpecTest {
+	return GetBlameSpecTest(
+		"Blame Type Invalid Commitment - Happy Flow",
+		makeInvalidForInvalidCommitment(validRound1MsgBytes),
+	)
+}
+
+func BlameTypeInvalidScalar() *FrostSpecTest {
+	return GetBlameSpecTest(
+		"Blame Type Invalid Scalar - Happy Flow",
+		makeInvalidForInvalidScalar(validRound1MsgBytes),
+	)
+}
+
+func BlameTypeInvalidShare_FailedShareDecryption() *FrostSpecTest {
+	return GetBlameSpecTest(
+		"Blame Type Invalid Share (Unable to Decrypt) - Happy Flow",
+		makeInvalidForFailedEcies(validRound1MsgBytes),
+	)
+}
+
+func BlameTypeInvalidShare_FailedValidationAgainstCommitment() *FrostSpecTest {
+
+	requestID := testingutils.GetRandRequestID()
+	ks := testingutils.Testing4SharesSet()
+
+	threshold := 3
+	operators := []types.OperatorID{1, 2, 3, 4}
+
+	initMessages := make(map[uint32][]*dkg.SignedMessage)
+	initMsgBytes := testingutils.InitMessageDataBytes(
+		operators,
+		uint16(threshold),
+		testingutils.TestingWithdrawalCredentials,
+		testingutils.TestingForkVersion,
+	)
+	for _, operatorID := range operators {
+		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
+				MsgType:    dkg.InitMsgType,
+				Identifier: requestID,
+				Data:       initMsgBytes,
+			}),
+		}
+	}
+
+	pmData := `{"round":2,"round1":{"Commitment":["qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa","luG8uPgVeTmvRoE4MBPMpt/Vgp/oCAA9TBTzG07bRJ45L6Uo9uDGQ9gKFkt9+07n","kBb0Obbc/CYaBH/56rZtOUw6bz6DMFrbouRUNwh8lBfH8OjWg3NQwBIXE3Ir8lmn","p7CN+Aow0TPzKW0wFmbL3qNYuSnjLEo3Gtjapg71mDYn9+IGpmGtoqHAd1LUzxDA"],"ProofS":"YI5p9V013jgYX/78TJ3PSdI/5QEbZFKnRND0pTo6XAE=","ProofR":"MW4+TKI7AAf/q0OljiJiNLSkoAPXy4PzTXC2dFqhlAc=","Shares":{"1":"BLh2p4b+/slitPiMXooPEka+S6TqCcdQSB7Bzv1XTZNp0N5wpnI/jgA4qAwzg2YCbVdBzcG26FF5p/4FRHk1syDz0ljuJkv30ahpxt/bby1ItMnBKgy7p+zYOE9RkAlecpnowYohR3wj/Fxq/ln5gNRWDmMcWMePrflm5dpMCziY","3":"BMqwinOzpjBtLed3b/pCDuG2x9XQPzXlKIXHtR+8pK4R+qPbU6hB4Xgf/9D/b2PKs/jnH6XOKfLX7q1bC9DZkH59cmeeeAHFjy3YeObXyF3L7E1MX4NWHxmkjjWSLiH08M2MkQCtfrswWzIfOVT7YgFJSRRDy2sf94CA1WdbFAc3","4":"BKLu/3DyaZOIzXx6SkKrRhxojh30Y5uLXOqBGbt5hQiPZpbDtULBSdTr4XrUyXLdmi3JW+jiZihksHxHjZWgq7mdlhMspFYyEnqxmobMBuDxydmEa5VzdoCtsSrRc79kx9SkwMsNkmY52VhtkgwLlISCXSdmAQc8BIn3kT7xQcy1"}}}`
+
+	blameProtocolMessageBytes := []byte(pmData)
+	blameSignedMessage := &dkg.SignedMessage{
+		Message: &dkg.Message{
+			MsgType:    dkg.ProtocolMsgType,
+			Identifier: requestID,
+			Data:       blameProtocolMessageBytes,
+		},
+		Signer: 2,
+	}
+	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(blameSignedMessage, ks.DKGOperators[2].ETHAddress)
+	blameSignedMessage.Signature = sig
+
+	return &FrostSpecTest{
+		Name:   "Blame Type Invalid Share - Happy Flow",
+		Keyset: ks,
+
+		RequestID: requestID,
+		Threshold: uint64(threshold),
+		Operators: operators,
+
+		ExpectedOutcome: testingutils.TestOutcome{
+			BlameOutcome: testingutils.TestBlameOutcome{
+				Valid: true,
+			},
+		},
+		ExpectedError: "could not find dkg runner",
+
+		InputMessages: map[int]MessagesForNodes{
+			0: initMessages,
+			2: {
+				2: []*dkg.SignedMessage{blameSignedMessage},
+			},
+		},
+	}
+}
+
+func BlameTypeInconsistentMessage() *FrostSpecTest {
+
+	requestID := testingutils.GetRandRequestID()
+	ks := testingutils.Testing4SharesSet()
+
+	threshold := 3
+	operators := []types.OperatorID{1, 2, 3, 4}
+
+	initMessages := make(map[uint32][]*dkg.SignedMessage)
+	initMsgBytes := testingutils.InitMessageDataBytes(
+		operators,
+		uint16(threshold),
+		testingutils.TestingWithdrawalCredentials,
+		testingutils.TestingForkVersion,
+	)
+	for _, operatorID := range operators {
+		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
+				MsgType:    dkg.InitMsgType,
+				Identifier: requestID,
+				Data:       initMsgBytes,
+			}),
+		}
+	}
+
+	pmData1 := `{"round":2,"round1":{"Commitment":["r4pOd119gLXzt06xwvmXudIYrEFHl7ZyT7yXDMz3Wt/CmK+KkPRem6nq4ov5Sf3q","gQh8Bd8lJmokT9zzFUK/javWp8z8VOIp5R/kCyXxCoYqpICyOwmg0XVYZMLwIj/Q","qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa"],"ProofS":"RY6MjnapCPtt6cR9YXWdbdd3Me8BSJCrlTJpX9y5bL8=","ProofR":"NARDVvKTH6pSjiYtMwUqiZSYv1lKNVk6deDB9FcddkA=","Shares":{"1":"BOSDxrY8bmwO+WdwYs/TgDC+viXCYQqldNoEmSutOHrljBoIGmS9KKmxbAYEpdTtk+ahyyOnG0lHn3WTrN9PJEeYE6QpcGgrRkUWOq/RSwHQX50R00iUmCnXH5B3WVUdyTTAOzkvenfWrq6W+uVQ4Vu00k590W9xCbBvtGM4UXJ+","3":"BOlPoCeJDaUsr3bRVGPlU0JZ1OgPm8StbA93DYyEaL5e5Y7PNzEyCnrDPWVoVqnNbPk6GikWHoGd/sOJCB4l7fBiyd0H0H6Ypwz44MFhEu8qgBWxFGeG730HZKv4+6mj048Tfkj1l+tTHdqI8O3GjzwWD51UOl1aIV68swslQqeL","4":"BHZa/+riJkzhM7PFkIFzhlUkqX2K3P1iQZO1wJRTmyvPuqqYnAc0KsbkSnDSq7GTwwA5L+jtle3Y4NlxVFH5lq9RYntwNnDyRliDwzxis8xlRDQtrnAFfIySw+rDJa7clxWUTavMjHeEawDWYv9MIKbPId0AwrlXMRb7pycDMoAW"}}}`
+
+	pmData2 := `{"round":2,"round1":{"Commitment":["r4pOd119gLXzt06xwvmXudIYrEFHl7ZyT7yXDMz3Wt/CmK+KkPRem6nq4ov5Sf3q","gQh8Bd8lJmokT9zzFUK/javWp8z8VOIp5R/kCyXxCoYqpICyOwmg0XVYZMLwIj/Q","qzHrRAIpma7lmbbm37SazNCYX6WE2/RYQF+lQdr+s+SO/3AknLoMH0ocuAFjx+Fa"],"ProofS":"RY6MjnapCPtt6cR9YXWdbdd3Me8BSJCrlTJpX9y5bL8=","ProofR":"NARDVvKTH6pSjiYtMwUqiZSYv1lKNVk6deDB9FcddkA=","Shares":{"1":"BOSDxrY8bmwO+WdwYs/TgDC+viXCYQqldNoEmSutOHrljBoIGmS9KKmxbAYEpdTtk+ahyyOnG0lHn3WTrN9PJEeYE6QpcGgrRkUWOq/RSwHQX50R00iUmCnXH5B3WVUdyTTAOzkvenfWrq6W+uVQ4Vu00k590W9xCbBvtGM4UXJ+","3":"BOlPoCeJDaUsr3bRVGPlU0JZ1OgPm8StbA93DYyEaL5e5Y7PNzEyCnrDPWVoVqnNbPk6GikWHoGd/sOJCB4l7fBiyd0H0H6Ypwz44MFhEu8qgBWxFGeG730HZKv4+6mj048Tfkj1l+tTHdqI8O3GjzwWD51UOl1aIV68swslQqeL","4":"BHZa/+riJkzhM7PFkIFzhlUkqX2K3P1iQZO1wJRTmyvPuqqYnAc0KsbkSnDSq7GTwwA5L+jtle3Y4NlxVFH5lq9RYntwNnDyRliDwzxis8xlRDQtrnAFfIySw+rDJa7clxWUTavMjHeEawDWYv9MIKbPId0AwrlXMRb7pycDMoWA"}}}` // root changed
+
+	data1SignedMessage := &dkg.SignedMessage{
+		Message: &dkg.Message{
+			MsgType:    dkg.ProtocolMsgType,
+			Identifier: requestID,
+			Data:       []byte(pmData1),
+		},
+		Signer: 2,
+	}
+	sig, _ := testingutils.NewTestingKeyManager().SignDKGOutput(data1SignedMessage, ks.DKGOperators[2].ETHAddress)
+	data1SignedMessage.Signature = sig
+
+	data2SignedMessage := &dkg.SignedMessage{
+		Message: &dkg.Message{
+			MsgType:    dkg.ProtocolMsgType,
+			Identifier: requestID,
+			Data:       []byte(pmData2),
+		},
+		Signer: 2,
+	}
+	sig, _ = testingutils.NewTestingKeyManager().SignDKGOutput(data2SignedMessage, ks.DKGOperators[2].ETHAddress)
+	data2SignedMessage.Signature = sig
+
+	return &FrostSpecTest{
+		Name:   "Blame Type Inconsisstent Message - Happy Flow",
+		Keyset: ks,
+
+		RequestID: requestID,
+		Threshold: uint64(threshold),
+		Operators: operators,
+
+		ExpectedOutcome: testingutils.TestOutcome{
+			BlameOutcome: testingutils.TestBlameOutcome{
+				Valid: true,
+			},
+		},
+		ExpectedError: "could not find dkg runner",
+
+		InputMessages: map[int]MessagesForNodes{
+			0: initMessages,
+			2: {
+				2: []*dkg.SignedMessage{data1SignedMessage, data2SignedMessage},
+			},
+		},
+	}
+}
+
+func makeInvalidForFailedEcies(data []byte) []byte {
+	protocolMessage := &frost.ProtocolMsg{}
+	_ = protocolMessage.Decode(data)
+
+	protocolMessage.Round1Message.Shares[maliciousOperatorID] = []byte("rubbish-value")
+	d, _ := protocolMessage.Encode()
+	return d
+}
+
+func makeInvalidForInvalidScalar(data []byte) []byte {
+	protocolMessage := &frost.ProtocolMsg{}
+	_ = protocolMessage.Decode(data)
+
+	protocolMessage.Round1Message.ProofR = []byte("rubbish-value")
+	d, _ := protocolMessage.Encode()
+	return d
+}
+
+func makeInvalidForInvalidCommitment(data []byte) []byte {
+	protocolMessage := &frost.ProtocolMsg{}
+	_ = protocolMessage.Decode(data)
+
+	protocolMessage.Round1Message.Commitment[len(protocolMessage.Round1Message.Commitment)-1] = []byte("rubbish-value")
+	d, _ := protocolMessage.Encode()
+	return d
+}
diff --git a/dkg/spectest/tests/frost/keygen.go b/dkg/spectest/tests/frost/keygen.go
new file mode 100644
index 0000000..f62a91b
--- /dev/null
+++ b/dkg/spectest/tests/frost/keygen.go
@@ -0,0 +1,65 @@
+package frost
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+)
+
+func Keygen() *FrostSpecTest {
+
+	requestID := testingutils.GetRandRequestID()
+	ks := testingutils.Testing4SharesSet()
+
+	threshold := 3
+	operators := []types.OperatorID{1, 2, 3, 4}
+	initMsgBytes := testingutils.InitMessageDataBytes(
+		operators,
+		uint16(threshold),
+		testingutils.TestingWithdrawalCredentials,
+		testingutils.TestingForkVersion,
+	)
+
+	initMessages := make(map[uint32][]*dkg.SignedMessage)
+	for _, operatorID := range operators {
+		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
+				MsgType:    dkg.InitMsgType,
+				Identifier: requestID,
+				Data:       initMsgBytes,
+			}),
+		}
+	}
+
+	return &FrostSpecTest{
+		Name:   "Simple Keygen",
+		Keyset: ks,
+
+		RequestID: requestID,
+		Threshold: uint64(threshold),
+		Operators: operators,
+
+		ExpectedOutcome: testingutils.TestOutcome{
+			KeygenOutcome: testingutils.TestKeygenOutcome{
+				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
+				Share: map[uint32]string{
+					1: "5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912",
+					2: "533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7",
+					3: "442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186",
+					4: "2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f",
+				},
+				OperatorPubKeys: map[uint32]string{
+					1: "add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa",
+					2: "8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5",
+					3: "b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051",
+					4: "a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037",
+				},
+			},
+		},
+		ExpectedError: "",
+
+		InputMessages: map[int]MessagesForNodes{
+			0: initMessages,
+		},
+	}
+}
diff --git a/dkg/spectest/tests/frost/resharing.go b/dkg/spectest/tests/frost/resharing.go
new file mode 100644
index 0000000..218de83
--- /dev/null
+++ b/dkg/spectest/tests/frost/resharing.go
@@ -0,0 +1,90 @@
+package frost
+
+import (
+	"encoding/hex"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+)
+
+func Resharing() *FrostSpecTest {
+
+	requestID := testingutils.GetRandRequestID()
+	ks := testingutils.Testing13SharesSet()
+
+	threshold := 3
+	operators := []types.OperatorID{5, 6, 7, 8}
+	operatorsOld := []types.OperatorID{1, 2, 3} //4}
+
+	vk, _ := hex.DecodeString("8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812")
+	reshareMsgBytes := testingutils.ReshareMessageDataBytes(
+		operators,
+		uint16(threshold),
+		vk,
+	)
+
+	initMessages := make(map[uint32][]*dkg.SignedMessage)
+	for _, operatorID := range append(operators, operatorsOld...) {
+		initMessages[uint32(operatorID)] = []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[operatorID].SK, operatorID, &dkg.Message{
+				MsgType:    dkg.ReshareMsgType,
+				Identifier: requestID,
+				Data:       reshareMsgBytes,
+			}),
+		}
+	}
+
+	spectest := &FrostSpecTest{
+		Name:   "Simple Resharing",
+		Keyset: ks,
+
+		RequestID: requestID,
+		Threshold: uint64(threshold),
+		Operators: operators,
+
+		IsResharing:  true,
+		OperatorsOld: operatorsOld,
+		OldKeygenOutcomes: testingutils.TestOutcome{
+			KeygenOutcome: testingutils.TestKeygenOutcome{
+				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
+				Share: map[uint32]string{
+					1: "5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912",
+					2: "533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7",
+					3: "442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186",
+					4: "2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f",
+				},
+				OperatorPubKeys: map[uint32]string{
+					1: "add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa",
+					2: "8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5",
+					3: "b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051",
+					4: "a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037",
+				},
+			},
+		},
+
+		ExpectedOutcome: testingutils.TestOutcome{
+			KeygenOutcome: testingutils.TestKeygenOutcome{
+				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
+				Share: map[uint32]string{
+					5: "52046f0837c928ea5d5bbc893b90f3cd75a07a9d25092e2fbb0129825100c3be",
+					6: "0f1d824d53df922ca8c15d639c802f84463a78cf69ef57e0b1cbb8b95cd1f458",
+					7: "213989136198ba32e82eb8e449a843b7fb6a52007ba72794212d25d135a84679",
+					8: "146adc07375723b4e869f703396758634172622d5a32414b092570cadb83ba20",
+				},
+				OperatorPubKeys: map[uint32]string{
+					5: "81e52afe4656f4544715cc2a37724c939afa8462d57549ba242681b52c80d8ac7e6b259d03ba37ce688aeca5e1a346b3",
+					6: "ac6d0b0ba2f3f581f520c59049c6dfb98ce12d87a3ee9ccc00b9e0ef13153b036c777a946d9ec78409a047d92ce942e7",
+					7: "8d9b4d117564b4852ee7d060626e27bc93ec5dddde0fbbbe053aed7e54b0772b334ad74149fb1c6d3f1ff3d5b4d87fc8",
+					8: "b11cb28641e5d6440e214d45abfc6a2158cbf163312144609e08236fee95aa096a61a0d70b4401d8daf4af69a1cca9ad",
+				},
+			},
+		},
+		ExpectedError: "",
+
+		InputMessages: map[int]MessagesForNodes{
+			0: initMessages,
+		},
+	}
+	return spectest
+}
diff --git a/dkg/spectest/tests/frost/test.go b/dkg/spectest/tests/frost/test.go
new file mode 100644
index 0000000..35c3307
--- /dev/null
+++ b/dkg/spectest/tests/frost/test.go
@@ -0,0 +1,270 @@
+package frost
+
+import (
+	"encoding/hex"
+	"sort"
+	"testing"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/dkg/frost"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+	"github.com/stretchr/testify/require"
+)
+
+type MessagesForNodes map[uint32][]*dkg.SignedMessage
+
+type FrostSpecTest struct {
+	Name   string
+	Keyset *testingutils.TestKeySet
+
+	// Keygen Options
+	RequestID dkg.RequestID
+	Threshold uint64
+	Operators []types.OperatorID
+
+	// Resharing Options
+	IsResharing       bool
+	OperatorsOld      []types.OperatorID
+	OldKeygenOutcomes testingutils.TestOutcome
+
+	// Expected
+	ExpectedOutcome testingutils.TestOutcome
+	ExpectedError   string
+
+	InputMessages map[int]MessagesForNodes
+}
+
+func (test *FrostSpecTest) TestName() string {
+	return test.Name
+}
+
+func (test *FrostSpecTest) Run(t *testing.T) {
+
+	outcomes, blame, err := test.TestingFrost()
+
+	if len(test.ExpectedError) > 0 {
+		require.EqualError(t, err, test.ExpectedError)
+		return
+	} else {
+		require.NoError(t, err)
+	}
+
+	if blame != nil {
+		require.Equal(t, test.ExpectedOutcome.BlameOutcome.Valid, blame.Valid)
+		return
+	}
+
+	for _, operatorID := range test.Operators {
+
+		outcome := outcomes[uint32(operatorID)]
+		if outcome.ProtocolOutput != nil {
+			vk := hex.EncodeToString(outcome.ProtocolOutput.ValidatorPK)
+			sk := outcome.ProtocolOutput.Share.SerializeToHexStr()
+			pk := outcome.ProtocolOutput.OperatorPubKeys[operatorID].SerializeToHexStr()
+
+			t.Logf("printing outcome keys for operatorID %d\n", operatorID)
+			t.Logf("vk %s\n", vk)
+			t.Logf("sk %s\n", sk)
+			t.Logf("pk %s\n", pk)
+
+			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.ValidatorPK, vk)
+			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.Share[uint32(operatorID)], sk)
+			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.OperatorPubKeys[uint32(operatorID)], pk)
+		}
+	}
+}
+
+func (test *FrostSpecTest) TestingFrost() (map[uint32]*dkg.ProtocolOutcome, *dkg.BlameOutput, error) {
+
+	testingutils.ResetRandSeed()
+	dkgsigner := testingutils.NewTestingKeyManager()
+	storage := testingutils.NewTestingStorage()
+	network := testingutils.NewTestingNetwork()
+
+	nodes := test.TestingFrostNodes(
+		test.RequestID,
+		network,
+		storage,
+		dkgsigner,
+	)
+
+	alloperators := test.Operators
+	if test.IsResharing {
+		alloperators = append(alloperators, test.OperatorsOld...)
+	}
+
+	initMessages, exists := test.InputMessages[0]
+	if !exists {
+		return nil, nil, errors.New("init messages not found in spec")
+	}
+
+	for operatorID, messages := range initMessages {
+		for _, message := range messages {
+
+			messageBytes, _ := message.Encode()
+			startMessage := &types.SSVMessage{
+				MsgType: types.DKGMsgType,
+				Data:    messageBytes,
+			}
+			if err := nodes[types.OperatorID(operatorID)].ProcessMessage(startMessage); err != nil {
+				return nil, nil, errors.Wrapf(err, "failed to start dkg protocol for operator %d", operatorID)
+			}
+		}
+	}
+
+	for round := 1; round <= 5; round++ {
+
+		messages := network.BroadcastedMsgs
+		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)
+		for _, msg := range messages {
+
+			dkgMsg := &dkg.SignedMessage{}
+			if err := dkgMsg.Decode(msg.Data); err != nil {
+				return nil, nil, err
+			}
+
+			msgsToBroadcast := []*types.SSVMessage{}
+			if testMessages, ok := test.InputMessages[round][uint32(dkgMsg.Signer)]; ok {
+				for _, testMessage := range testMessages {
+					testMessageBytes, _ := testMessage.Encode()
+					msgsToBroadcast = append(msgsToBroadcast, &types.SSVMessage{
+						MsgType: msg.MsgType,
+						Data:    testMessageBytes,
+					})
+				}
+			} else {
+				msgsToBroadcast = append(msgsToBroadcast, msg)
+			}
+
+			operatorList := alloperators
+			if test.IsResharing && round > 2 {
+				operatorList = test.Operators
+			}
+
+			sort.SliceStable(operatorList, func(i, j int) bool {
+				return operatorList[i] < operatorList[j]
+			})
+
+			for _, operatorID := range operatorList {
+
+				if operatorID == dkgMsg.Signer {
+					continue
+				}
+
+				for _, msgToBroadcast := range msgsToBroadcast {
+					if err := nodes[operatorID].ProcessMessage(msgToBroadcast); err != nil {
+						return nil, nil, err
+					}
+				}
+			}
+		}
+
+	}
+
+	ret := make(map[uint32]*dkg.ProtocolOutcome)
+
+	outputs := network.DKGOutputs
+	blame := network.BlameOutput
+	if blame != nil {
+		return nil, blame, nil
+	}
+
+	for operatorID, output := range outputs {
+		if output.BlameData != nil {
+			signedMsg := &dkg.SignedMessage{}
+			_ = signedMsg.Decode(output.BlameData.BlameMessage)
+			ret[uint32(operatorID)] = &dkg.ProtocolOutcome{
+				BlameOutput: &dkg.BlameOutput{
+					Valid:        output.BlameData.Valid,
+					BlameMessage: signedMsg,
+				},
+			}
+			continue
+		}
+
+		pk := &bls.PublicKey{}
+		_ = pk.Deserialize(output.Data.SharePubKey)
+
+		share, _ := dkgsigner.Decrypt(test.Keyset.DKGOperators[operatorID].EncryptionKey, output.Data.EncryptedShare)
+		sk := &bls.SecretKey{}
+		_ = sk.Deserialize(share)
+
+		ret[uint32(operatorID)] = &dkg.ProtocolOutcome{
+			ProtocolOutput: &dkg.KeyGenOutput{
+				ValidatorPK: output.Data.ValidatorPubKey,
+				Share:       sk,
+				OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
+					operatorID: pk,
+				},
+				Threshold: test.Threshold,
+			},
+		}
+	}
+
+	return ret, nil, nil
+}
+
+func (test *FrostSpecTest) TestingFrostNodes(
+	requestID dkg.RequestID,
+	network dkg.Network,
+	storage dkg.Storage,
+	dkgsigner types.DKGSigner,
+) map[types.OperatorID]*dkg.Node {
+
+	nodes := make(map[types.OperatorID]*dkg.Node)
+	for _, operatorID := range test.Operators {
+		_, operator, _ := storage.GetDKGOperator(operatorID)
+		node := dkg.NewNode(
+			operator,
+			&dkg.Config{
+				KeygenProtocol: frost.New,
+				Network:        network,
+				Storage:        storage,
+				Signer:         dkgsigner,
+			})
+		nodes[operatorID] = node
+	}
+
+	if test.IsResharing {
+		operatorsOldList := types.OperatorList(test.OperatorsOld).ToUint32List()
+		keygenOutcomeOld := test.OldKeygenOutcomes.KeygenOutcome.ToKeygenOutcomeMap(test.Threshold, operatorsOldList)
+
+		for _, operatorID := range test.OperatorsOld {
+			storage := testingutils.NewTestingStorage()
+			_ = storage.SaveKeyGenOutput(keygenOutcomeOld[uint32(operatorID)])
+
+			_, operator, _ := storage.GetDKGOperator(operatorID)
+			node := dkg.NewResharingNode(
+				operator,
+				test.OperatorsOld,
+				&dkg.Config{
+					ReshareProtocol: frost.NewResharing,
+					Network:         network,
+					Storage:         storage,
+					Signer:          dkgsigner,
+				})
+			nodes[operatorID] = node
+		}
+
+		for _, operatorID := range test.Operators {
+			storage := testingutils.NewTestingStorage()
+			_ = storage.SaveKeyGenOutput(keygenOutcomeOld[operatorsOldList[0]])
+
+			_, operator, _ := storage.GetDKGOperator(operatorID)
+			node := dkg.NewResharingNode(
+				operator,
+				test.OperatorsOld,
+				&dkg.Config{
+					ReshareProtocol: frost.NewResharing,
+					Network:         network,
+					Storage:         storage,
+					Signer:          dkgsigner,
+				})
+			nodes[operatorID] = node
+		}
+	}
+	return nodes
+}
diff --git a/dkg/spectest/tests/happy_flow.go b/dkg/spectest/tests/happy_flow.go
new file mode 100644
index 0000000..732da61
--- /dev/null
+++ b/dkg/spectest/tests/happy_flow.go
@@ -0,0 +1,87 @@
+package tests
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+)
+
+// HappyFlow tests a simple full happy flow until decided
+func HappyFlow() *MsgProcessingSpecTest {
+	ks := testingutils.Testing4SharesSet()
+	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
+	init := &dkg.Init{
+		OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
+		Threshold:             3,
+		WithdrawalCredentials: testingutils.TestingWithdrawalCredentials,
+		Fork:                  testingutils.TestingForkVersion,
+	}
+	initBytes, _ := init.Encode()
+	root := testingutils.DespositDataSigningRoot(ks, init)
+
+	return &MsgProcessingSpecTest{
+		Name: "happy flow",
+		InputMessages: []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.InitMsgType,
+				Identifier: identifier,
+				Data:       initBytes,
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.ProtocolMsgType,
+				Identifier: identifier,
+				Data:       nil, // GLNOTE: Dummy message simulating the Protocol to complete
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
+				MsgType:    dkg.DepositDataMsgType,
+				Identifier: identifier,
+				Data:       testingutils.PartialDepositDataBytes(2, root, ks.Shares[2]),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
+				MsgType:    dkg.DepositDataMsgType,
+				Identifier: identifier,
+				Data:       testingutils.PartialDepositDataBytes(3, root, ks.Shares[3]),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
+				MsgType:    dkg.DepositDataMsgType,
+				Identifier: identifier,
+				Data:       testingutils.PartialDepositDataBytes(4, root, ks.Shares[4]),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 2, root),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 3, root),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 4, root),
+			}),
+		},
+		OutputMessages: []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.DepositDataMsgType,
+				Identifier: identifier,
+				Data:       testingutils.PartialDepositDataBytes(1, root, ks.Shares[1]),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 1, root),
+			}),
+		},
+		Output: map[types.OperatorID]*dkg.SignedOutput{
+			1: ks.SignedOutputObject(identifier, 1, root),
+			2: ks.SignedOutputObject(identifier, 2, root),
+			3: ks.SignedOutputObject(identifier, 3, root),
+			4: ks.SignedOutputObject(identifier, 4, root),
+		},
+		KeySet:        ks,
+		ExpectedError: "",
+	}
+}
diff --git a/dkg/spectest/tests/msg_process_spectest.go b/dkg/spectest/tests/msg_process_spectest.go
new file mode 100644
index 0000000..6a00513
--- /dev/null
+++ b/dkg/spectest/tests/msg_process_spectest.go
@@ -0,0 +1,95 @@
+package tests
+
+import (
+	"fmt"
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+	"github.com/stretchr/testify/require"
+	"testing"
+)
+
+type MsgProcessingSpecTest struct {
+	Name           string
+	InputMessages  []*dkg.SignedMessage
+	OutputMessages []*dkg.SignedMessage
+	Output         map[types.OperatorID]*dkg.SignedOutput
+	KeySet         *testingutils.TestKeySet
+	ExpectedError  string
+}
+
+func (test *MsgProcessingSpecTest) TestName() string {
+	return test.Name
+}
+
+func (test *MsgProcessingSpecTest) Run(t *testing.T) {
+	node := testingutils.TestingDKGNode(test.KeySet)
+
+	var lastErr error
+	for _, msg := range test.InputMessages {
+		byts, _ := msg.Encode()
+		err := node.ProcessMessage(&types.SSVMessage{
+			MsgType: types.DKGMsgType,
+			Data:    byts,
+		})
+
+		if err != nil {
+			lastErr = err
+		}
+	}
+
+	if len(test.ExpectedError) > 0 {
+		require.EqualError(t, lastErr, test.ExpectedError)
+	} else {
+		require.NoError(t, lastErr)
+	}
+
+	// test output message
+	broadcastedMsgs := node.GetConfig().Network.(*testingutils.TestingNetwork).BroadcastedMsgs
+	if len(test.OutputMessages) > 0 {
+		require.Len(t, broadcastedMsgs, len(test.OutputMessages))
+
+		for i, msg := range test.OutputMessages {
+			bMsg := broadcastedMsgs[i]
+			require.Equal(t, types.DKGMsgType, bMsg.MsgType)
+			sMsg := &dkg.SignedMessage{}
+			require.NoError(t, sMsg.Decode(bMsg.Data))
+			if sMsg.Message.MsgType == dkg.OutputMsgType {
+				require.Equal(t, dkg.OutputMsgType, msg.Message.MsgType, "OutputMsgType expected")
+				o1 := &dkg.SignedOutput{}
+				require.NoError(t, o1.Decode(msg.Message.Data))
+
+				o2 := &dkg.SignedOutput{}
+				require.NoError(t, o2.Decode(sMsg.Message.Data))
+
+				es1 := o1.Data.EncryptedShare
+				o1.Data.EncryptedShare = nil
+				es2 := o2.Data.EncryptedShare
+				o2.Data.EncryptedShare = nil
+
+				s1, _ := types.Decrypt(test.KeySet.DKGOperators[msg.Signer].EncryptionKey, es1)
+				s2, _ := types.Decrypt(test.KeySet.DKGOperators[msg.Signer].EncryptionKey, es2)
+				require.Equal(t, s1, s2, "shares don't match")
+				r1, _ := o1.Data.GetRoot()
+				r2, _ := o2.Data.GetRoot()
+				require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
+			} else {
+				r1, _ := msg.GetRoot()
+				r2, _ := sMsg.GetRoot()
+				require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
+			}
+
+		}
+	}
+	streamed := node.GetConfig().Network.(*testingutils.TestingNetwork).DKGOutputs
+	if len(test.Output) > 0 {
+		require.Len(t, streamed, len(test.Output))
+		for id, output := range test.Output {
+			s := streamed[id]
+			require.NotNilf(t, s, "output for operator %d not found", id)
+			r1, _ := output.Data.GetRoot()
+			r2, _ := s.Data.GetRoot()
+			require.EqualValues(t, r1, r2, fmt.Sprintf("output for operator %d roots not equal", id))
+		}
+	}
+}
diff --git a/dkg/spectest/tests/resharing.go b/dkg/spectest/tests/resharing.go
new file mode 100644
index 0000000..c2830e7
--- /dev/null
+++ b/dkg/spectest/tests/resharing.go
@@ -0,0 +1,67 @@
+package tests
+
+import (
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/bloxapp/ssv-spec/types/testingutils"
+)
+
+// ResharingHappyFlow tests a simple (dummy) resharing flow, the difference between this and keygen happy flow is
+// resharing doesn't sign deposit data
+func ResharingHappyFlow() *MsgProcessingSpecTest {
+	ks := testingutils.Testing4SharesSet()
+	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
+	reshare := &dkg.Reshare{
+		ValidatorPK: make([]byte, 48),
+		OperatorIDs: []types.OperatorID{1, 2, 3, 4},
+		Threshold:   3,
+	}
+	reshareBytes, _ := reshare.Encode()
+	var root []byte
+
+	return &MsgProcessingSpecTest{
+		Name: "resharing happy flow",
+		InputMessages: []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.ReshareMsgType,
+				Identifier: identifier,
+				Data:       reshareBytes,
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.ProtocolMsgType,
+				Identifier: identifier,
+				Data:       nil, // GLNOTE: Dummy message simulating the Protocol to complete
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 2, root),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 3, root),
+			}),
+			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 4, root),
+			}),
+		},
+		OutputMessages: []*dkg.SignedMessage{
+			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
+				MsgType:    dkg.OutputMsgType,
+				Identifier: identifier,
+				Data:       ks.SignedOutputBytes(identifier, 1, root),
+			}),
+		},
+		Output: map[types.OperatorID]*dkg.SignedOutput{
+			1: ks.SignedOutputObject(identifier, 1, root),
+			2: ks.SignedOutputObject(identifier, 2, root),
+			3: ks.SignedOutputObject(identifier, 3, root),
+			4: ks.SignedOutputObject(identifier, 4, root),
+		},
+		KeySet:        ks,
+		ExpectedError: "",
+	}
+}
diff --git a/dkg/stubdkg/messages.go b/dkg/stubdkg/messages.go
new file mode 100644
index 0000000..ef81e7c
--- /dev/null
+++ b/dkg/stubdkg/messages.go
@@ -0,0 +1,29 @@
+package stubdkg
+
+import (
+	"encoding/json"
+)
+
+type Stage int
+
+const (
+	StubStage1 Stage = iota
+	StubStage2
+	StubStage3
+)
+
+type ProtocolMsg struct {
+	Stage Stage
+	// Data is any data a real DKG implementation will need
+	Data interface{}
+}
+
+// Encode returns a msg encoded bytes or error
+func (msg *ProtocolMsg) Encode() ([]byte, error) {
+	return json.Marshal(msg)
+}
+
+// Decode returns error if decoding failed
+func (msg *ProtocolMsg) Decode(data []byte) error {
+	return json.Unmarshal(data, msg)
+}
diff --git a/dkg/stubdkg/stub_dkg.go b/dkg/stubdkg/stub_dkg.go
new file mode 100644
index 0000000..b624d1c
--- /dev/null
+++ b/dkg/stubdkg/stub_dkg.go
@@ -0,0 +1,97 @@
+package stubdkg
+
+import (
+	"fmt"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/pkg/errors"
+)
+
+// DKG is a stub dkg protocol simulating a real DKG protocol with 3 stages in it
+type DKG struct {
+	identifier dkg.RequestID
+	network    dkg.Network
+	operatorID types.OperatorID
+	operators  []types.OperatorID
+
+	validatorPK    []byte
+	operatorShares map[types.OperatorID]*bls.SecretKey
+
+	msgs map[Stage][]*ProtocolMsg
+}
+
+func New(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID) dkg.Protocol {
+	return &DKG{
+		identifier: identifier,
+		network:    network,
+		operatorID: operatorID,
+		msgs:       map[Stage][]*ProtocolMsg{},
+	}
+}
+
+func (s *DKG) SetOperators(validatorPK []byte, operatorShares map[types.OperatorID]*bls.SecretKey) {
+	s.validatorPK = validatorPK
+	s.operatorShares = operatorShares
+}
+
+func (s *DKG) Start() error {
+	//s.operators = initOrReshare.Init.OperatorIDs
+	// TODO send Stage 1 msg
+	return nil
+}
+
+func (s *DKG) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
+	// TODO validate msg
+
+	dataMsg := &ProtocolMsg{}
+	if err := dataMsg.Decode(msg.Message.Data); err != nil {
+		return false, nil, errors.Wrap(err, "could not decode protocol msg")
+	}
+
+	if s.msgs[dataMsg.Stage] == nil {
+		s.msgs[dataMsg.Stage] = []*ProtocolMsg{}
+	}
+	s.msgs[dataMsg.Stage] = append(s.msgs[dataMsg.Stage], dataMsg)
+
+	switch dataMsg.Stage {
+	case StubStage1:
+		if len(s.msgs[StubStage1]) == len(s.operators) {
+			fmt.Printf("stage 1 done\n")
+			// TODO send Stage 2 msg
+		}
+	case StubStage2:
+		if len(s.msgs[StubStage2]) == len(s.operators) {
+			fmt.Printf("stage 2 done\n")
+			// TODO send Stage 3 msg
+		}
+	case StubStage3:
+		if len(s.msgs[StubStage3]) == len(s.operators) {
+			ret := &dkg.KeyGenOutput{
+				Share:       s.operatorShares[s.operatorID],
+				ValidatorPK: s.validatorPK,
+				OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
+					1: s.operatorShares[1].GetPublicKey(),
+					2: s.operatorShares[2].GetPublicKey(),
+					3: s.operatorShares[3].GetPublicKey(),
+					4: s.operatorShares[4].GetPublicKey(),
+				},
+			}
+			return true, &dkg.ProtocolOutcome{ProtocolOutput: ret}, nil
+		}
+	}
+	return false, nil, nil
+}
+
+//func (s *DKG) signDKGMsg(data []byte) *dkg.SignedMessage {
+//	return &dkg.SignedMessage{
+//		Message: &dkg.Message{
+//			MsgType:    dkg.ProtocolMsgType,
+//			Identifier: s.identifier,
+//			Data:       data,
+//		},
+//		Signer: s.operatorID,
+//		// TODO - how do we sign?
+//	}
+//}
diff --git a/dkg/stubdkg/stub_dkg_test.go b/dkg/stubdkg/stub_dkg_test.go
new file mode 100644
index 0000000..3e711ee
--- /dev/null
+++ b/dkg/stubdkg/stub_dkg_test.go
@@ -0,0 +1,105 @@
+package stubdkg
+
+import (
+	"encoding/hex"
+	"fmt"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+	"github.com/stretchr/testify/require"
+	"testing"
+)
+
+func TestSimpleDKG(t *testing.T) {
+	types.InitBLS()
+
+	operators := []types.OperatorID{
+		1, 2, 3, 4,
+	}
+	k := 3
+	polyDegree := k - 1
+	payloadToSign := "hello"
+
+	// create polynomials for each operator
+	poly := make(map[types.OperatorID][]bls.Fr)
+	for _, id := range operators {
+		coeff := make([]bls.Fr, 0)
+		for i := 1; i <= polyDegree; i++ {
+			c := bls.Fr{}
+			c.SetByCSPRNG()
+			coeff = append(coeff, c)
+		}
+		poly[id] = coeff
+	}
+
+	// create points for each operator
+	points := make(map[types.OperatorID][]*bls.Fr)
+	for _, id := range operators {
+		for _, evalID := range operators {
+			if points[evalID] == nil {
+				points[evalID] = make([]*bls.Fr, 0)
+			}
+
+			res := &bls.Fr{}
+			x := &bls.Fr{}
+			x.SetInt64(int64(evalID))
+			require.NoError(t, bls.FrEvaluatePolynomial(res, poly[id], x))
+
+			points[evalID] = append(points[evalID], res)
+		}
+	}
+
+	// calculate shares
+	shares := make(map[types.OperatorID]*bls.SecretKey)
+	pks := make(map[types.OperatorID]*bls.PublicKey)
+	sigs := make(map[types.OperatorID]*bls.Sign)
+	for id, ps := range points {
+		var sum *bls.Fr
+		for _, p := range ps {
+			if sum == nil {
+				sum = p
+			} else {
+				bls.FrAdd(sum, sum, p)
+			}
+		}
+		shares[id] = bls.CastToSecretKey(sum)
+		pks[id] = shares[id].GetPublicKey()
+		sigs[id] = shares[id].Sign(payloadToSign)
+	}
+
+	// get validator pk
+	validatorPK := bls.PublicKey{}
+	idVec := make([]bls.ID, 0)
+	pkVec := make([]bls.PublicKey, 0)
+	for operatorID, pk := range pks {
+		blsID := bls.ID{}
+		err := blsID.SetDecString(fmt.Sprintf("%d", operatorID))
+		require.NoError(t, err)
+		idVec = append(idVec, blsID)
+
+		pkVec = append(pkVec, *pk)
+	}
+	require.NoError(t, validatorPK.Recover(pkVec, idVec))
+	fmt.Printf("validator pk: %s\n", hex.EncodeToString(validatorPK.Serialize()))
+
+	// reconstruct sig
+	reconstructedSig := bls.Sign{}
+	idVec = make([]bls.ID, 0)
+	sigVec := make([]bls.Sign, 0)
+	for operatorID, sig := range sigs {
+		blsID := bls.ID{}
+		err := blsID.SetDecString(fmt.Sprintf("%d", operatorID))
+		require.NoError(t, err)
+		idVec = append(idVec, blsID)
+
+		sigVec = append(sigVec, *sig)
+
+		if len(sigVec) >= k {
+			break
+		}
+	}
+	require.NoError(t, reconstructedSig.Recover(sigVec, idVec))
+	fmt.Printf("reconstructed sig: %s\n", hex.EncodeToString(reconstructedSig.Serialize()))
+
+	// verify
+	require.True(t, reconstructedSig.Verify(&validatorPK, payloadToSign))
+}
diff --git a/dkg/types.go b/dkg/types.go
new file mode 100644
index 0000000..98ed021
--- /dev/null
+++ b/dkg/types.go
@@ -0,0 +1,57 @@
+package dkg
+
+import (
+	"crypto/rsa"
+
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/ethereum/go-ethereum/common"
+)
+
+// Network is a collection of funcs for DKG
+type Network interface {
+	// StreamDKGBlame will stream to any subscriber the blame result of the DKG
+	StreamDKGBlame(blame *BlameOutput) error
+	// StreamDKGOutput will stream to any subscriber the result of the DKG
+	StreamDKGOutput(output map[types.OperatorID]*SignedOutput) error
+	// BroadcastDKGMessage will broadcast a msg to the dkg network
+	BroadcastDKGMessage(msg *SignedMessage) error
+}
+
+type Storage interface {
+	// GetDKGOperator returns true and operator object if found by operator ID
+	GetDKGOperator(operatorID types.OperatorID) (bool, *Operator, error)
+	SaveKeyGenOutput(output *KeyGenOutput) error
+	GetKeyGenOutput(pk types.ValidatorPK) (*KeyGenOutput, error)
+}
+
+// Operator holds all info regarding a DKG Operator on the network
+type Operator struct {
+	// OperatorID the node's Operator ID
+	OperatorID types.OperatorID
+	// ETHAddress the operator's eth address used to sign and verify messages against
+	ETHAddress common.Address
+	// EncryptionPubKey encryption pubkey for shares
+	EncryptionPubKey *rsa.PublicKey
+}
+
+type Config struct {
+	// Protocol the DKG protocol implementation
+	KeygenProtocol      func(network Network, operatorID types.OperatorID, identifier RequestID, signer types.DKGSigner, storage Storage, init *Init) Protocol
+	ReshareProtocol     func(network Network, operatorID types.OperatorID, identifier RequestID, signer types.DKGSigner, storage Storage, oldOperators []types.OperatorID, reshare *Reshare, output *KeyGenOutput) Protocol
+	Network             Network
+	Storage             Storage
+	SignatureDomainType types.DomainType
+	Signer              types.DKGSigner
+}
+
+type ErrInvalidRound struct{}
+
+func (e ErrInvalidRound) Error() string {
+	return "invalid dkg round"
+}
+
+type ErrMismatchRound struct{}
+
+func (e ErrMismatchRound) Error() string {
+	return "mismatch dkg round"
+}
diff --git a/types/testingutils/dkg.go b/types/testingutils/dkg.go
new file mode 100644
index 0000000..c42e6cc
--- /dev/null
+++ b/types/testingutils/dkg.go
@@ -0,0 +1,192 @@
+package testingutils
+
+import (
+	"crypto/ecdsa"
+	"crypto/rsa"
+	"encoding/hex"
+	"fmt"
+	"strconv"
+
+	spec "github.com/attestantio/go-eth2-client/spec/phase0"
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/ethereum/go-ethereum/crypto"
+	"github.com/herumi/bls-eth-go-binary/bls"
+)
+
+var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
+var TestingForkVersion = types.PraterNetwork.ForkVersion()
+
+var TestingDKGNode = func(keySet *TestKeySet) *dkg.Node {
+	network := NewTestingNetwork()
+	km := NewTestingKeyManager()
+	config := &dkg.Config{
+		KeygenProtocol: func(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID, signer types.DKGSigner, storage dkg.Storage, init *dkg.Init) dkg.Protocol {
+			return &TestingKeygenProtocol{
+				KeyGenOutput: keySet.KeyGenOutput(1),
+			}
+		},
+		ReshareProtocol: func(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID, signer types.DKGSigner, storage dkg.Storage, oldOperators []types.OperatorID, reshare *dkg.Reshare, output *dkg.KeyGenOutput) dkg.Protocol {
+			return &TestingKeygenProtocol{
+				KeyGenOutput: keySet.KeyGenOutput(1),
+			}
+		},
+		Network:             network,
+		Storage:             NewTestingStorage(),
+		SignatureDomainType: types.PrimusTestnet,
+		Signer:              km,
+	}
+
+	return dkg.NewNode(&dkg.Operator{
+		OperatorID:       1,
+		ETHAddress:       keySet.DKGOperators[1].ETHAddress,
+		EncryptionPubKey: &keySet.DKGOperators[1].EncryptionKey.PublicKey,
+	}, config)
+}
+
+var SignDKGMsg = func(sk *ecdsa.PrivateKey, id types.OperatorID, msg *dkg.Message) *dkg.SignedMessage {
+	domain := types.PrimusTestnet
+	sigType := types.DKGSignatureType
+
+	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
+	sig, _ := crypto.Sign(r, sk)
+
+	return &dkg.SignedMessage{
+		Message:   msg,
+		Signer:    id,
+		Signature: sig,
+	}
+}
+
+var InitMessageDataBytes = func(operators []types.OperatorID, threshold uint16, withdrawalCred []byte, fork spec.Version) []byte {
+	byts, _ := InitMessageData(operators, threshold, withdrawalCred, fork).Encode()
+	return byts
+}
+
+var InitMessageData = func(operators []types.OperatorID, threshold uint16, withdrawalCred []byte, fork spec.Version) *dkg.Init {
+	return &dkg.Init{
+		OperatorIDs:           operators,
+		Threshold:             threshold,
+		WithdrawalCredentials: withdrawalCred,
+		Fork:                  fork,
+	}
+}
+
+var ReshareMessageDataBytes = func(operators []types.OperatorID, threshold uint16, validatorPK types.ValidatorPK) []byte {
+	byts, _ := ReshareMessageData(operators, threshold, validatorPK).Encode()
+	return byts
+}
+
+var ReshareMessageData = func(operators []types.OperatorID, threshold uint16, validatorPK types.ValidatorPK) *dkg.Reshare {
+	return &dkg.Reshare{
+		ValidatorPK: validatorPK,
+		OperatorIDs: operators,
+		Threshold:   threshold,
+	}
+}
+
+var ProtocolMsgDataBytes = func(stage stubdkg.Stage) []byte {
+	d := &stubdkg.ProtocolMsg{
+		Stage: stage,
+	}
+
+	ret, _ := d.Encode()
+	return ret
+}
+
+var PartialDepositDataBytes = func(signer types.OperatorID, root []byte, sk *bls.SecretKey) []byte {
+	d := &dkg.PartialDepositData{
+		Signer:    signer,
+		Root:      root,
+		Signature: sk.SignByte(root).Serialize(),
+	}
+	ret, _ := d.Encode()
+	return ret
+}
+
+var DespositDataSigningRoot = func(keySet *TestKeySet, initMsg *dkg.Init) []byte {
+	root, _, _ := types.GenerateETHDepositData(
+		keySet.ValidatorPK.Serialize(),
+		initMsg.WithdrawalCredentials,
+		initMsg.Fork,
+		types.DomainDeposit,
+	)
+	return root
+}
+var (
+	encryptedDataCache = map[string][]byte{}
+	decryptedDataCache = map[string][]byte{}
+)
+
+func TestingEncryption(pk *rsa.PublicKey, data []byte) []byte {
+	id := hex.EncodeToString(pk.N.Bytes()) + fmt.Sprintf("%x", pk.E) + hex.EncodeToString(data)
+	if found := encryptedDataCache[id]; found != nil {
+		return found
+	}
+	cipherText, _ := types.Encrypt(pk, data)
+	encryptedDataCache[id] = cipherText
+	return cipherText
+}
+
+func TestingDecryption(sk *rsa.PrivateKey, data []byte) []byte {
+	id := hex.EncodeToString(sk.N.Bytes()) + fmt.Sprintf("%x", sk.E) + hex.EncodeToString(data)
+	if found := decryptedDataCache[id]; found != nil {
+		return found
+	}
+	plaintext, _ := types.Decrypt(sk, data)
+	decryptedDataCache[id] = plaintext
+	return plaintext
+}
+
+func (ks *TestKeySet) KeyGenOutput(opId types.OperatorID) *dkg.KeyGenOutput {
+	opPks := make(map[types.OperatorID]*bls.PublicKey)
+	for id, share := range ks.Shares {
+		opPks[id] = share.GetPublicKey()
+	}
+
+	return &dkg.KeyGenOutput{
+		Share:           ks.Shares[opId],
+		OperatorPubKeys: opPks,
+		ValidatorPK:     ks.ValidatorPK.Serialize(),
+		Threshold:       ks.Threshold,
+	}
+}
+
+var (
+	signedOutputCache = map[string]*dkg.SignedOutput{}
+)
+
+func (ks *TestKeySet) SignedOutputObject(requestID dkg.RequestID, opId types.OperatorID, root []byte) *dkg.SignedOutput {
+	id := hex.EncodeToString(requestID[:]) + strconv.FormatUint(uint64(opId), 10) + hex.EncodeToString(root)
+	if found := signedOutputCache[id]; found != nil {
+		return found
+	}
+	share := ks.Shares[opId]
+	o := &dkg.Output{
+		RequestID:       requestID,
+		EncryptedShare:  TestingEncryption(&ks.DKGOperators[opId].EncryptionKey.PublicKey, share.Serialize()),
+		SharePubKey:     share.GetPublicKey().Serialize(),
+		ValidatorPubKey: ks.ValidatorPK.Serialize(),
+	}
+	if root != nil {
+		o.DepositDataSignature = ks.ValidatorSK.SignByte(root).Serialize()
+	}
+	root1, _ := o.GetRoot()
+
+	sig, _ := crypto.Sign(root1, ks.DKGOperators[opId].SK)
+
+	ret := &dkg.SignedOutput{
+		Data:      o,
+		Signer:    opId,
+		Signature: sig,
+	}
+	signedOutputCache[id] = ret
+	return ret
+}
+
+func (ks *TestKeySet) SignedOutputBytes(requestID dkg.RequestID, opId types.OperatorID, root []byte) []byte {
+	d := ks.SignedOutputObject(requestID, opId, root)
+	ret, _ := d.Encode()
+	return ret
+}
diff --git a/types/testingutils/frost.go b/types/testingutils/frost.go
new file mode 100644
index 0000000..b9b29c6
--- /dev/null
+++ b/types/testingutils/frost.go
@@ -0,0 +1,74 @@
+package testingutils
+
+import (
+	crand "crypto/rand"
+	"encoding/hex"
+	"math/big"
+	mrand "math/rand"
+
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+	"github.com/herumi/bls-eth-go-binary/bls"
+)
+
+type TestOutcome struct {
+	KeygenOutcome TestKeygenOutcome
+	BlameOutcome  TestBlameOutcome
+}
+
+type TestKeygenOutcome struct {
+	ValidatorPK     string
+	Share           map[uint32]string
+	OperatorPubKeys map[uint32]string
+}
+
+func (o TestKeygenOutcome) ToKeygenOutcomeMap(threshold uint64, operators []uint32) map[uint32]*dkg.KeyGenOutput {
+	m := make(map[uint32]*dkg.KeyGenOutput)
+
+	opPublicKeys := make(map[types.OperatorID]*bls.PublicKey)
+	for _, operatorID := range operators {
+
+		pk := &bls.PublicKey{}
+		_ = pk.DeserializeHexStr(o.OperatorPubKeys[operatorID])
+		opPublicKeys[types.OperatorID(operatorID)] = pk
+
+		share := o.Share[operatorID]
+		sk := &bls.SecretKey{}
+		_ = sk.DeserializeHexStr(share)
+
+		vk, _ := hex.DecodeString(o.ValidatorPK)
+
+		m[operatorID] = &dkg.KeyGenOutput{
+			Share:           sk,
+			ValidatorPK:     vk,
+			OperatorPubKeys: opPublicKeys,
+			Threshold:       threshold,
+		}
+	}
+
+	return m
+}
+
+func ResetRandSeed() {
+	src := mrand.NewSource(1)
+	src.Seed(12345)
+	crand.Reader = mrand.New(src)
+}
+
+func GetRandRequestID() dkg.RequestID {
+	requestID := dkg.RequestID{}
+	for i := range requestID {
+		rndInt, _ := crand.Int(crand.Reader, big.NewInt(255))
+		if len(rndInt.Bytes()) == 0 {
+			requestID[i] = 0
+		} else {
+			requestID[i] = rndInt.Bytes()[0]
+		}
+	}
+	return requestID
+}
+
+type TestBlameOutcome struct {
+	Valid        bool
+	BlameMessage []byte
+}
diff --git a/types/testingutils/keygen_protocol.go b/types/testingutils/keygen_protocol.go
new file mode 100644
index 0000000..ec15785
--- /dev/null
+++ b/types/testingutils/keygen_protocol.go
@@ -0,0 +1,15 @@
+package testingutils
+
+import "github.com/bloxapp/ssv-spec/dkg"
+
+type TestingKeygenProtocol struct {
+	KeyGenOutput *dkg.KeyGenOutput
+}
+
+func (m TestingKeygenProtocol) Start() error {
+	return nil
+}
+
+func (m TestingKeygenProtocol) ProcessMsg(msg *dkg.SignedMessage) (bool, *dkg.ProtocolOutcome, error) {
+	return true, &dkg.ProtocolOutcome{ProtocolOutput: m.KeyGenOutput}, nil
+}
diff --git a/types/testingutils/keymanager.go b/types/testingutils/keymanager.go
index 2c7bf6c..a1ce40c 100644
--- a/types/testingutils/keymanager.go
+++ b/types/testingutils/keymanager.go
@@ -130,10 +130,7 @@ func (km *testingKeyManager) Encrypt(pk *rsa.PublicKey, data []byte) ([]byte, er
 
 // SignDKGOutput signs output according to the SIP https://docs.google.com/document/d/1TRVUHjFyxINWW2H9FYLNL2pQoLy6gmvaI62KL_4cREQ/edit
 func (km *testingKeyManager) SignDKGOutput(output types.Root, address common.Address) (types.Signature, error) {
-	domain := types.PrimusTestnet
-	sigType := types.DKGSignatureType
-
-	root, err := types.ComputeSigningRoot(output, types.ComputeSignatureDomain(domain, sigType))
+	root, err := output.GetRoot()
 	if err != nil {
 		return nil, err
 	}
@@ -161,28 +158,3 @@ func (km *testingKeyManager) RemoveShare(pubKey string) error {
 	delete(km.keys, pubKey)
 	return nil
 }
-
-var (
-	encryptedDataCache = map[string][]byte{}
-	decryptedDataCache = map[string][]byte{}
-)
-
-func TestingEncryption(pk *rsa.PublicKey, data []byte) []byte {
-	id := hex.EncodeToString(pk.N.Bytes()) + fmt.Sprintf("%x", pk.E) + hex.EncodeToString(data)
-	if found := encryptedDataCache[id]; found != nil {
-		return found
-	}
-	cipherText, _ := types.Encrypt(pk, data)
-	encryptedDataCache[id] = cipherText
-	return cipherText
-}
-
-func TestingDecryption(sk *rsa.PrivateKey, data []byte) []byte {
-	id := hex.EncodeToString(sk.N.Bytes()) + fmt.Sprintf("%x", sk.E) + hex.EncodeToString(data)
-	if found := decryptedDataCache[id]; found != nil {
-		return found
-	}
-	plaintext, _ := types.Decrypt(sk, data)
-	decryptedDataCache[id] = plaintext
-	return plaintext
-}
diff --git a/types/testingutils/network.go b/types/testingutils/network.go
index 06138de..b323f58 100644
--- a/types/testingutils/network.go
+++ b/types/testingutils/network.go
@@ -1,12 +1,15 @@
 package testingutils
 
 import (
+	"github.com/bloxapp/ssv-spec/dkg"
 	"github.com/bloxapp/ssv-spec/qbft"
 	"github.com/bloxapp/ssv-spec/types"
 )
 
 type TestingNetwork struct {
 	BroadcastedMsgs           []*types.SSVMessage
+	DKGOutputs                map[types.OperatorID]*dkg.SignedOutput
+	BlameOutput               *dkg.BlameOutput
 	SyncHighestDecidedCnt     int
 	SyncHighestChangeRoundCnt int
 	DecidedByRange            [2]qbft.Height
@@ -15,6 +18,7 @@ type TestingNetwork struct {
 func NewTestingNetwork() *TestingNetwork {
 	return &TestingNetwork{
 		BroadcastedMsgs: make([]*types.SSVMessage, 0),
+		DKGOutputs:      make(map[types.OperatorID]*dkg.SignedOutput, 0),
 	}
 }
 
@@ -23,6 +27,21 @@ func (net *TestingNetwork) Broadcast(message *types.SSVMessage) error {
 	return nil
 }
 
+// StreamDKGOutput will stream to any subscriber the result of the DKG
+func (net *TestingNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
+	for id, signedOutput := range output {
+		net.DKGOutputs[id] = signedOutput
+	}
+
+	return nil
+}
+
+func (net *TestingNetwork) StreamDKGBlame(blame *dkg.BlameOutput) error {
+	//TODO implement me
+	net.BlameOutput = blame
+	return nil
+}
+
 func (net *TestingNetwork) SyncHighestDecided(identifier types.MessageID) error {
 	net.SyncHighestDecidedCnt++
 	return nil
@@ -36,3 +55,17 @@ func (net *TestingNetwork) SyncHighestDecided(identifier types.MessageID) error
 func (net *TestingNetwork) SyncDecidedByRange(identifier types.MessageID, from, to qbft.Height) {
 	net.DecidedByRange = [2]qbft.Height{from, to}
 }
+
+// BroadcastDKGMessage will broadcast a msg to the dkg network
+func (net *TestingNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
+	data, err := msg.Encode()
+	if err != nil {
+		return err
+	}
+	net.BroadcastedMsgs = append(net.BroadcastedMsgs, &types.SSVMessage{
+		MsgType: types.DKGMsgType,
+		MsgID:   types.MessageID{}, // TODO: what should we use for the MsgID?
+		Data:    data,
+	})
+	return nil
+}
diff --git a/types/testingutils/storage.go b/types/testingutils/storage.go
new file mode 100644
index 0000000..cc03c05
--- /dev/null
+++ b/types/testingutils/storage.go
@@ -0,0 +1,46 @@
+package testingutils
+
+import (
+	"encoding/hex"
+	"github.com/bloxapp/ssv-spec/dkg"
+	"github.com/bloxapp/ssv-spec/types"
+)
+
+type testingStorage struct {
+	operators   map[types.OperatorID]*dkg.Operator
+	keygenoupts map[string]*dkg.KeyGenOutput
+}
+
+func NewTestingStorage() *testingStorage {
+	ret := &testingStorage{
+		operators:   make(map[types.OperatorID]*dkg.Operator),
+		keygenoupts: make(map[string]*dkg.KeyGenOutput),
+	}
+
+	for i, s := range Testing13SharesSet().DKGOperators {
+		ret.operators[i] = &dkg.Operator{
+			OperatorID:       i,
+			ETHAddress:       s.ETHAddress,
+			EncryptionPubKey: &s.EncryptionKey.PublicKey,
+		}
+	}
+
+	return ret
+}
+
+// GetDKGOperator returns true and operator object if found by operator ID
+func (s *testingStorage) GetDKGOperator(operatorID types.OperatorID) (bool, *dkg.Operator, error) {
+	if ret, found := s.operators[operatorID]; found {
+		return true, ret, nil
+	}
+	return false, nil, nil
+}
+
+func (s *testingStorage) SaveKeyGenOutput(output *dkg.KeyGenOutput) error {
+	s.keygenoupts[hex.EncodeToString(output.ValidatorPK)] = output
+	return nil
+}
+
+func (s *testingStorage) GetKeyGenOutput(pk types.ValidatorPK) (*dkg.KeyGenOutput, error) {
+	return s.keygenoupts[hex.EncodeToString(pk)], nil
+}
