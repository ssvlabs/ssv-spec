package msgcontainer

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumPostConsensusMsg tests adding a quorum of post consensus message to container
func QuorumPostConsensusMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "quorum post consensus message",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
			testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks),
			testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks),
		},
		PostMsgCount: 3,
		PostReconstructedSignature: []string{
			"8e264f992161837da62d51ae3131145526654b01b914f53182e4720307fd643798b4f2656e1a2750ac11577db71fe36d0d22b8fbf213a495541db24236c1be76544eb110c1fe082b6dcbd6bb86bff914bb62d58ebeb8df297b02caec9e0e1b83",
			"a550f7d58d7a95635fd86378d139ef1fd00057266c7a5472c466c936a80e591590a96fdca8b5a470a2a3fe9a169074c2072d15f12504a52cb588fcb192a8576fa4b75f51aa2b09aa6ba5a6add81a6c7a7fa8176fe7254540daa220193f6646dd",
			"82f747fd5abee289fb0f2788f6548cf96d0ab854ff47b9631a8ad7014cbf104c88842724190b26b1de115d5c99ecb622170a88da343e5f84cc364b58fccdad3ec1015bd3f6144dcd5f95d711057e0bc2f0cabdec0db7c5420beae9e3ed90727b",
		},
	}
}
