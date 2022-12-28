package vm

import "fmt"

func unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return base + uint64(mode1)*uint64(maxOpType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxOpType) + uint64(typ)
}

type OpTable struct {
	baseOpcodes []OpCode
	length      int
}

func (t OpTable) Len() int { return t.length }

func (t OpTable) Halt() OpCode   { return t.baseOpcodes[haltOpFamily] }
func (t OpTable) Call() OpCode   { return t.baseOpcodes[callOpFamily] }
func (t OpTable) Return() OpCode { return t.baseOpcodes[returnOpFamily] }
func (t OpTable) IfThen() OpCode { return t.baseOpcodes[ifThenOpFamily] }
func (t OpTable) IfElse() OpCode { return t.baseOpcodes[ifElseOpFamily] }
func (t OpTable) Else() OpCode   { return t.baseOpcodes[elseOpFamily] }

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.baseOpcodes[printOpFamily], mode, typ)
}

func (t OpTable) Push(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.baseOpcodes[pushOpFamily], mode, typ)
}

func (t OpTable) Pop(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.baseOpcodes[popOpFamily], mode, typ)
}

func (t OpTable) Add(mode1, mode2 OpMode, typ OpType) OpCode {
	return binaryOpCode(t.baseOpcodes[addOpFamily], mode1, mode2, typ)
}

func NewOpTable() OpTable {
	baseOpcodes := make([]OpCode, maxOpFamily)

	family := haltOpFamily
	base := haltOpFamily
	for ; family < printOpFamily; family++ {
		baseOpcodes[family] = family
		base++
	}

	if base != printOpFamily {
		panic("Invalid op table")
	}

	baseOpcodes[printOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			familyBase := baseOpcodes[printOpFamily]
			opcode := unaryOpCode(familyBase, mode, typ)
			if opcode != base {
				panic(fmt.Errorf("Invalid op table: family:%d base:%d mode:%d type:%d; want %d; got %d", family, familyBase, mode, typ, base, opcode))
			}
			base++
		}
	}

	baseOpcodes[pushOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			if unaryOpCode(baseOpcodes[pushOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	baseOpcodes[popOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			if unaryOpCode(baseOpcodes[popOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	baseOpcodes[addOpFamily] = base
	for mode1 := ImmediateMode; mode1 < maxOpMode; mode1++ {
		for mode2 := ImmediateMode; mode2 < maxOpMode; mode2++ {
			for typ := I8; typ < maxOpType; typ++ {
				if binaryOpCode(baseOpcodes[addOpFamily], mode1, mode2, typ) != base {
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
