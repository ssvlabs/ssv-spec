package dkg

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
}

func StartNewDKG(initMsg *Init) (*Runner, error) {
	panic("implement")
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, *SignedOutput, error) {
	panic("implement")
}
