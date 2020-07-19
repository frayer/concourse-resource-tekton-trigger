package commands

type BridgeCommand struct {
	Check CheckCommand `command:"check" description:"Commands related to the 'check' operation on a Resource"`
}
