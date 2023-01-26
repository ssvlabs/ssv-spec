package tests

import (
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/stretchr/testify/require"
)

// MsgSpecTest tests encoding and decoding of a msg
type MsgSpecTest struct {
	Name            string
	Messages        []*alea.SignedMessage
	EncodedMessages [][]byte
	ExpectedRoots   [][]byte
	ExpectedError   string
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error

	for i, msg := range test.Messages {
		if err := msg.Validate(); err != nil {
			lastErr = err
			continue
		}

		switch msg.Message.MsgType {
		case alea.ProposalMsgType:
			rc := alea.ProposalData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCMsgType:
			rc := alea.VCBCData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.ABAMsgType:
			rc := alea.ABAData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.FillGapMsgType:
			rc := alea.FillGapData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.FillerMsgType:
			rc := alea.FillerData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.ABAInitMsgType:
			rc := alea.ABAInitData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.ABAAuxMsgType:
			rc := alea.ABAAuxData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.ABAConfMsgType:
			rc := alea.ABAConfData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.ABAFinishMsgType:
			rc := alea.ABAFinishData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCBroadcastMsgType:
			rc := alea.VCBCBroadcastData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCSendMsgType:
			rc := alea.VCBCSendData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCReadyMsgType:
			rc := alea.VCBCReadyData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCFinalMsgType:
			rc := alea.VCBCFinalData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCRequestMsgType:
			rc := alea.VCBCRequestData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		case alea.VCBCAnswerMsgType:
			rc := alea.VCBCAnswerData{}
			if err := rc.Decode(msg.Message.Data); err != nil {
				lastErr = err
			}
			if err := rc.Validate(); err != nil {
				lastErr = err
			}
		}

		if len(test.EncodedMessages) > 0 {
			byts, err := msg.Encode()
			require.NoError(t, err)
			require.EqualValues(t, test.EncodedMessages[i], byts)
		}

		if len(test.ExpectedRoots) > 0 {
			r, err := msg.GetRoot()
			require.NoError(t, err)
			require.EqualValues(t, test.ExpectedRoots[i], r)
		}
	}

	// check error
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *MsgSpecTest) TestName() string {
	return "alea message " + test.Name
}
