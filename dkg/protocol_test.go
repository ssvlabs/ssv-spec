package dkg

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyGenOutput(t *testing.T) {
	types.InitBLS()
	bytes := []byte("{\"Share\":\"X0cRp5bBEWtRGOw1J5+2TVUdmziBPSk5lU3S31Fg09k=\",\"OperatorPubKeys\":{\"1\":\"l9lKgR1kSTYFKp0tSs1kcYl0z2eNvv0mcyTI6fjnA0pKa32HeeJ6AZU4w8Qlw+Xn\",\"2\":\"przr4wl9dBcbQMcSoDHOsDcds9PEAs8s5pG5Eg87q3XU1W36DzdZFUSZm/GMU1Pt\",\"3\":\"gJDgt2ZqRezF1O90GKyZ8J5sskQCn+pqCn/Mvp7gi8U53g36Zr5rq8hJPdmd0amN\",\"4\":\"p8CidrcKXuM5XH1tJlXtYFKKolLU0h7KX8xSI+UMxCvRaLKAq3q1MXNU3d/PPfnk\"},\"ValidatorPK\":\"joAGZVGoGzGCWHCe2vfdH2PNaGoOTbiym7t6z+ZWCGd69aUn2USO5Hg1SF4CtQvA\",\"Threshold\":3}")
	kgo := &KeyGenOutput{}
	err := json.Unmarshal(bytes, kgo)
	require.NoError(t, err)
	serialized, err := json.Marshal(kgo)
	require.Equal(t, bytes, serialized)
}
