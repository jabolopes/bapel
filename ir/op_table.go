package ir

import "fmt"

func UnaryOpCode(base OpCode, mode OpMode, typ IrIntType) OpCode {
	return base + uint64(mode)*uint64(maxIrIntType) + uint64(typ)
}

func BinaryOpCode(base OpCode, mode1, mode2 OpMode, typ IrIntType) OpCode {
	return base + uint64(mode1)*uint64(maxIrIntType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxIrIntType) + uint64(typ)
}

type OpTable struct {
	baseOpcodes []OpCode
	length      int
}

func (t OpTable) Len() int { return t.length }

func (t OpTable) Halt() OpCode    { return t.baseOpcodes[haltOpFamily] }
func (t OpTable) Call() OpCode    { return t.baseOpcodes[callOpFamily] }
func (t OpTable) Return() OpCode  { return t.baseOpcodes[returnOpFamily] }
func (t OpTable) IfThen() OpCode  { return t.baseOpcodes[ifThenOpFamily] }
func (t OpTable) IfElse() OpCode  { return t.baseOpcodes[ifElseOpFamily] }
func (t OpTable) Else() OpCode    { return t.baseOpcodes[elseOpFamily] }
func (t OpTable) Syscall() OpCode { return t.baseOpcodes[syscallOpFamily] }
func (t OpTable) IOWait() OpCode  { return t.baseOpcodes[ioWaitOpFamily] }
func (t OpTable) IODo() OpCode    { return t.baseOpcodes[ioDoOpFamily] }

func (t OpTable) Print(mode OpMode, typ IrIntType, sign Sign) OpCode {
	switch sign {
	case Unsigned:
		return UnaryOpCode(t.baseOpcodes[printUOpFamily], mode, typ)
	case Signed:
		return UnaryOpCode(t.baseOpcodes[printSOpFamily], mode, typ)
	default:
		panic(fmt.Errorf("Unhandled sign %d", sign))
	}
}

func (t OpTable) Push(mode OpMode, typ IrIntType) OpCode {
	return UnaryOpCode(t.baseOpcodes[pushOpFamily], mode, typ)
}

func (t OpTable) Pop(mode OpMode, typ IrIntType) OpCode {
	return UnaryOpCode(t.baseOpcodes[popOpFamily], mode, typ)
}

func (t OpTable) Neg(mode OpMode, typ IrIntType) OpCode {
	return UnaryOpCode(t.baseOpcodes[negOpFamily], mode, typ)
}

func (t OpTable) Add(mode1, mode2 OpMode, typ IrIntType) OpCode {
	return BinaryOpCode(t.baseOpcodes[addOpFamily], mode1, mode2, typ)
}

func NewOpTable() OpTable {
	baseOpcodes := make([]OpCode, maxOpFamily)

	family := haltOpFamily
	base := haltOpFamily
	for ; family < printUOpFamily; family++ {
		baseOpcodes[family] = family
		base++
	}

	if base != printUOpFamily {
		panic("Invalid op table")
	}

	baseOpcodes[printUOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxIrIntType; typ++ {
			familyBase := baseOpcodes[printUOpFamily]
			opcode := UnaryOpCode(familyBase, mode, typ)
			if opcode != base {
				panic(fmt.Errorf("Invalid op table: family:%d base:%d mode:%d type:%d; want %d; got %d", family, familyBase, mode, typ, base, opcode))
			}
			base++
		}
	}

	baseOpcodes[printSOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxIrIntType; typ++ {
			familyBase := baseOpcodes[printSOpFamily]
			opcode := UnaryOpCode(familyBase, mode, typ)
			if opcode != base {
				panic(fmt.Errorf("Invalid op table: family:%d base:%d mode:%d type:%d; want %d; got %d", family, familyBase, mode, typ, base, opcode))
			}
			base++
		}
	}

	baseOpcodes[pushOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxIrIntType; typ++ {
			if UnaryOpCode(baseOpcodes[pushOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	baseOpcodes[popOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxIrIntType; typ++ {
			if UnaryOpCode(baseOpcodes[popOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	baseOpcodes[negOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxIrIntType; typ++ {
			if UnaryOpCode(baseOpcodes[negOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	baseOpcodes[addOpFamily] = base
	for mode1 := ImmediateMode; mode1 < maxOpMode; mode1++ {
		for mode2 := ImmediateMode; mode2 < maxOpMode; mode2++ {
			for typ := I8; typ < maxIrIntType; typ++ {
				if BinaryOpCode(baseOpcodes[addOpFamily], mode1, mode2, typ) != base {
					panic("Invalid op table")
				}
				base++
			}
		}
	}

	return OpTable{baseOpcodes, int(base)}
}

func init() {
	_ = NewOpTable()
}
