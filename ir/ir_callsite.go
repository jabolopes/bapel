package ir

// irCallsite represents a callsite for a function that has been
// declared but not yet defined. This callsite has a placeholder
// address that needs to be fixed once the function has been defined.
type irCallsite struct {
	offsets []uint64
}
