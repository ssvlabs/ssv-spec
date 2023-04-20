package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	ssz "github.com/ferranbt/fastssz"
)

// GetRootNoFulldata returns the signed message root omitting fullData
func GetRootNoFulldata(message *qbft.SignedMessage) ([32]byte, error) {
	sw := &signedMessageTestWrapper{SignedMessage: message}
	return ssz.HashWithDefaultHasher(sw)
}

type signedMessageTestWrapper struct {
	*qbft.SignedMessage
}

// GetTree ssz hashes the SignedMessage object
func (sw *signedMessageTestWrapper) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(sw)
}

// HashTreeRoot ssz hashes the SignedMessage object
func (sw *signedMessageTestWrapper) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(sw)
}

// HashTreeRootWith ssz hashes the SignedMessageTestWrapper object with a hasher
// This is a shameless copy from messages_encoding.go omitting full_data merkelizing.
func (sw *signedMessageTestWrapper) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	s := sw.SignedMessage
	indx := hh.Index()

	// Field (0) 'Signature'
	if size := len(s.Signature); size != 96 {
		err = ssz.ErrBytesLengthFn("SignedMessage.Signature", size, 96)
		return
	}
	hh.PutBytes(s.Signature)

	// Field (1) 'Signers'
	{
		if size := len(s.Signers); size > 13 {
			err = ssz.ErrListTooBigFn("SignedMessage.Signers", size, 13)
			return
		}
		subIndx := hh.Index()
		for _, i := range s.Signers {
			hh.AppendUint64(i)
		}
		hh.FillUpTo32()
		numItems := uint64(len(s.Signers))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(13, numItems, 8))
	}

	// Field (2) 'Message'
	if err = s.Message.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}
