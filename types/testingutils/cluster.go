package testingutils

import (
	"github.com/bloxapp/ssv-spec/ssv"
)

var BaseCluster = func(keySet *TestKeySet) *ssv.Cluster {

	return ssv.NewCluster(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingKeyManager(),
		func() *ssv.ClusterRunner {
			return &ssv.ClusterRunner{}
		},
	)
}
