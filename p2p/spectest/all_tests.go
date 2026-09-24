package spectest

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation/consensus"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation/general"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation/postconsensus"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation/preconsensus"
)

var AllTests = []tests.TestF{
	general.MalformedPubsubPayload,
	general.NoSigners,
	general.NoSignatures,
	general.EmptySignature,
	general.SignatureCountMismatch,
	general.ZeroSigner,
	general.NonUniqueSigners,
	general.UnsupportedSSVMessageType,
	consensus.ValidExistingInstanceConsensus,
	consensus.UnknownInstanceConsensus,
	consensus.InvalidConsensusSignature,
	consensus.WrongConsensusIdentifier,
	consensus.FutureConsensusMultiSigner,
	consensus.ValidFutureDecidedConsensus,
	consensus.DecidedConsensusBadFullData,
	preconsensus.ValidPreConsensusPartialSignature,
	preconsensus.InvalidPreConsensusPartialSignatureSignature,
	preconsensus.PreConsensusWithoutRunningDuty,
	preconsensus.PreConsensusSignerMismatch,
	preconsensus.PreConsensusPastSlot,
	preconsensus.PreConsensusFutureSlot,
	preconsensus.PreConsensusUnknownValidatorIndex,
	preconsensus.MultiSignerPartialSignature,
	postconsensus.ValidPostConsensusPartialSignature,
	postconsensus.PostConsensusWithoutDecidedValue,
	postconsensus.PostConsensusWithoutRunningInstance,
	postconsensus.PostConsensusNotDecided,
}
