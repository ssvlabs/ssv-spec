package drand

import "github.com/drand/kyber/share/dkg"

func (d *DRand) processDealBundle(bundle dkg.DealBundle) error {
	d.board.DealsC <- bundle
	return nil
}
