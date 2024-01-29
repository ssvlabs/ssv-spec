package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

var AttesterCommitExtraLoadManager = func(runner ssv.AttesterRunner) qbft.CommitExtraLoadManagerF {
	return ssv.NewCommitExtraLoadManagerF(runner.GetBaseRunner(), types.BNRoleAttester, runner.GetBeaconNode(), runner.GetSigner(), types.DomainAttester)
}

var AggregatorCommitExtraLoadManager = func(runner ssv.AggregatorRunner) qbft.CommitExtraLoadManagerF {
	return ssv.NewCommitExtraLoadManagerF(runner.GetBaseRunner(), types.BNRoleAggregator, runner.GetBeaconNode(), runner.GetSigner(), types.DomainAggregateAndProof)
}

var ProposerCommitExtraLoadManager = func(runner ssv.ProposerRunner) qbft.CommitExtraLoadManagerF {
	return ssv.NewCommitExtraLoadManagerF(runner.GetBaseRunner(), types.BNRoleProposer, runner.GetBeaconNode(), runner.GetSigner(), types.DomainProposer)
}

var SyncCommitteeCommitExtraLoadManager = func(runner ssv.ProposerRunner) qbft.CommitExtraLoadManagerF {
	return ssv.NewCommitExtraLoadManagerF(runner.GetBaseRunner(), types.BNRoleSyncCommittee, runner.GetBeaconNode(), runner.GetSigner(), types.DomainSyncCommittee)
}

var SyncCommitteeAggregatorCommitExtraLoadManager = func(runner ssv.ProposerRunner) qbft.CommitExtraLoadManagerF {
	return ssv.NewCommitExtraLoadManagerF(runner.GetBaseRunner(), types.BNRoleSyncCommitteeContribution, runner.GetBeaconNode(), runner.GetSigner(), types.DomainSyncCommitteeSelectionProof)
}
