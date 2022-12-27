package ir

// irFrame represents a stack frame (pointed to by the fp).
//
// The frame manages function arguments and function locals. Note: the
// frame does not manage return values nor the saved pc. Both are
// considered to be outside of the frame. The caller manages the
// return values. The 'call' and 'return' opcodes manage the saved pc.
//
// The frame allocation and deallocation process is asymetric. To
// create a frame, the caller allocates spaces for the function's
// arguments and the 'enter' opcode allocates the remaining space for
// locals. To destroy a frame, the 'leave' opcode deallocates the
// space for arguments and also for locals.
//
// Because the frame allocation / deallocation process is asymetric,
// the entire frame size is used when deallocating the frame, but only
// the size of the locals is used when allocating the frame.
type irFrame struct {
	frameSize  uint16
	localsSize uint16
}

func (f irFrame) enterSize() uint16 { return f.localsSize }
func (f irFrame) leaveSize() uint16 { return f.frameSize }
