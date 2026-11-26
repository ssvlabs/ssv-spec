package proposerpreferences

import (
	"bytes"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BuilderRequestAuthData tests which bytes a builder entry signs (SIP #94 §5): its configured data, byte for
// byte and up to MaxBuilderAuthDataSize; or, when none is configured, the default derived from its URL by the
// upstream get_default_auth_data rule. The URL cases are the upstream example table, which is the conformance
// set, followed by the hosts the rule rejects. An entry with nothing to sign is skipped without stopping the
// round: each such case configures a second, valid entry, which is still signed.
//
// Each case starts the duty and compares the broadcast auth container with one built from the expected bytes.
func BuilderRequestAuthData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	const hostname = "builder.example.com"
	// The valid entry behind an entry that signs nothing.
	valid := ssv.BuilderEntry{Data: []byte("builder-auth-token")}

	cases := []struct {
		name    string
		entries []ssv.BuilderEntry
		// signed is the data of each expected auth partial, in order.
		signed [][]byte
	}{
		// The upstream example table.
		{"plain host", []ssv.BuilderEntry{{URL: "https://builder.example.com/"}}, [][]byte{[]byte(hostname)}},
		{"uppercase scheme and host with port path and query", []ssv.BuilderEntry{{URL: "HTTPS://Builder.Example.com:443/bids?x=1"}}, [][]byte{[]byte(hostname)}},
		{"explicit port", []ssv.BuilderEntry{{URL: "https://builder.example.com:8080"}}, [][]byte{[]byte(hostname)}},
		{"userinfo", []ssv.BuilderEntry{{URL: "https://user:pw@builder.example.com/"}}, [][]byte{[]byte(hostname)}},
		{"ipv4 literal", []ssv.BuilderEntry{{URL: "https://10.0.0.5:18550/eth/v1/builder"}}, [][]byte{[]byte("10.0.0.5")}},
		{"ipv6 literal compressed", []ssv.BuilderEntry{{URL: "https://[0:0:0:0:0:0:0:1]:8443/"}}, [][]byte{[]byte("[::1]")}},
		{"ipv4-mapped ipv6 literal in hexadecimal groups", []ssv.BuilderEntry{{URL: "https://[::ffff:192.0.2.1]/"}}, [][]byte{[]byte("[::ffff:c000:201]")}},

		// Hosts with no default: the entry is skipped.
		{"zoned ipv6 literal", []ssv.BuilderEntry{{URL: "https://[fe80::1%25eth0]/"}, valid}, [][]byte{valid.Data}},
		{"non-ascii host", []ssv.BuilderEntry{{URL: "https://exämple.com/"}, valid}, [][]byte{valid.Data}},
		{"url without a scheme", []ssv.BuilderEntry{{URL: hostname}, valid}, [][]byte{valid.Data}},
		{"empty url", []ssv.BuilderEntry{{URL: ""}, valid}, [][]byte{valid.Data}},

		// Lowercasing differs across languages on two code points. U+212A (Kelvin sign) lowercases to ASCII "k"
		// in the upstream rule and here, so the host is accepted. U+0130 (dotted capital I) lowercases to "i" plus
		// a combining dot upstream, which is not ASCII, so the host is rejected.
		{"kelvin sign lowercases to ascii", []ssv.BuilderEntry{{URL: "https://Kevin.example/"}}, [][]byte{[]byte("kevin.example")}},
		{"dotted capital i is not ascii", []ssv.BuilderEntry{{URL: "https://İstanbul.example/"}, valid}, [][]byte{valid.Data}},

		// Configured data.
		{"configured data is signed as is", []ssv.BuilderEntry{{Data: []byte("Builder-Auth-Token"), URL: "https://builder.example.com/"}}, [][]byte{[]byte("Builder-Auth-Token")}},
		{"configured empty data is invalid", []ssv.BuilderEntry{{Data: []byte{}, URL: "https://builder.example.com/"}, valid}, [][]byte{valid.Data}},
		{"configured data at the size bound", []ssv.BuilderEntry{{Data: bytes.Repeat([]byte{0xa5}, gloas.MaxBuilderAuthDataSize)}}, [][]byte{bytes.Repeat([]byte{0xa5}, gloas.MaxBuilderAuthDataSize)}},
		{"configured data over the size bound", []ssv.BuilderEntry{{Data: bytes.Repeat([]byte{0xa5}, gloas.MaxBuilderAuthDataSize+1)}, valid}, [][]byte{valid.Data}},
	}

	subTests := make([]*newduty.StartNewRunnerDutySpecTest, 0, len(cases))
	for _, c := range cases {
		runner := testingutils.ProposerPreferencesRunner(ks)
		runner.(*ssv.ProposerPreferencesRunner).BuilderEntries = c.entries
		subTests = append(subTests, &newduty.StartNewRunnerDutySpecTest{
			Name:      c.name,
			Runner:    runner,
			Duty:      testingutils.TestingProposerPreferencesDuty(),
			Threshold: ks.Threshold,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusProposerPreferencesMsg(ks.Shares[1], 1), // the preference partial
				testingutils.PreConsensusBuilderRequestAuthMsg(ks.Shares[1], 1, c.signed),
			},
		})
	}

	return newduty.NewMultiStartNewRunnerDutySpecTest(
		"proposer preferences builder request auth data",
		testdoc.ProposerPreferencesBuilderRequestAuthDataDoc,
		subTests,
		ks,
	)
}
