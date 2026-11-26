package gloas

import "fmt"

// The list limits of a Gloas block, from consensus-specs at the snapshot SIP #94 pins. The block's lists are
// progressive (EIP-7916), so their types carry no bound and decoding checks none: the beacon chain asserts the
// limits in its state transition (process_operations, apply_parent_execution_payload) and rejects a block over
// any of them.
//
// Not limited here: a payload's deposit requests, which the state transition does not limit either, and the
// attestations' aggregation bits and the slashings' attesting indices, whose bound is the beacon state's
// committees, which the spec does not have.
const (
	MaxProposerSlashings     = 16 // MAX_PROPOSER_SLASHINGS
	MaxAttesterSlashings     = 1  // MAX_ATTESTER_SLASHINGS_ELECTRA
	MaxAttestations          = 8  // MAX_ATTESTATIONS_ELECTRA
	MaxVoluntaryExits        = 16 // MAX_VOLUNTARY_EXITS
	MaxBLSToExecutionChanges = 16 // MAX_BLS_TO_EXECUTION_CHANGES
	MaxPayloadAttestations   = 4  // MAX_PAYLOAD_ATTESTATIONS
	// MaxDeposits is the number of deposits a block body may carry: none. Deposits reach the chain as execution
	// requests, and process_operations asserts that the body's list is empty.
	MaxDeposits = 0

	MaxWithdrawalRequestsPerPayload     = 16 // MAX_WITHDRAWAL_REQUESTS_PER_PAYLOAD
	MaxConsolidationRequestsPerPayload  = 2  // MAX_CONSOLIDATION_REQUESTS_PER_PAYLOAD
	MaxBuilderDepositRequestsPerPayload = 64 // MAX_BUILDER_DEPOSIT_REQUESTS_PER_PAYLOAD
	MaxBuilderExitRequestsPerPayload    = 16 // MAX_BUILDER_EXIT_REQUESTS_PER_PAYLOAD

	// MaxBlobCommitmentsPerBlock is MAX_BLOB_COMMITMENTS_PER_BLOCK, the most blob commitments a bid can carry
	// under any blob schedule. The chain's limit for an epoch (max_blobs_per_block) is at most this, and is the
	// beacon node's to check: the spec has no blob schedule.
	MaxBlobCommitmentsPerBlock = 4096
)

// listLimit is one of a block's lists: its name, its length and the most it may hold.
type listLimit struct {
	name       string
	len, limit int
}

// ValidateListLimits checks each of the block's lists against its Gloas limit (see the constants above). A block
// over one is invalid on chain, so a cluster that signed it would lose the slot at publication.
func (b *BeaconBlock) ValidateListLimits() error {
	if b == nil || b.Body == nil {
		return fmt.Errorf("nil block body")
	}
	body := b.Body
	lists := []listLimit{
		{"proposer slashings", len(body.ProposerSlashings), MaxProposerSlashings},
		{"attester slashings", len(body.AttesterSlashings), MaxAttesterSlashings},
		{"attestations", len(body.Attestations), MaxAttestations},
		{"deposits", len(body.Deposits), MaxDeposits},
		{"voluntary exits", len(body.VoluntaryExits), MaxVoluntaryExits},
		{"BLS to execution changes", len(body.BLSToExecutionChanges), MaxBLSToExecutionChanges},
		{"payload attestations", len(body.PayloadAttestations), MaxPayloadAttestations},
	}
	if requests := body.ParentExecutionRequests; requests != nil {
		lists = append(lists,
			listLimit{"withdrawal requests", len(requests.Withdrawals), MaxWithdrawalRequestsPerPayload},
			listLimit{"consolidation requests", len(requests.Consolidations), MaxConsolidationRequestsPerPayload},
			listLimit{"builder deposit requests", len(requests.BuilderDeposits), MaxBuilderDepositRequestsPerPayload},
			listLimit{"builder exit requests", len(requests.BuilderExits), MaxBuilderExitRequestsPerPayload},
		)
	}
	if bid := body.SignedExecutionPayloadBid; bid != nil && bid.Message != nil {
		lists = append(lists, listLimit{"blob KZG commitments", len(bid.Message.BlobKZGCommitments), MaxBlobCommitmentsPerBlock})
	}
	for _, list := range lists {
		if list.len > list.limit {
			return fmt.Errorf("%d %s, over the limit of %d", list.len, list.name, list.limit)
		}
	}
	return nil
}
