package vm

import "fmt"

func unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return base + uint64(mode1)*uint64(maxOpType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxOpType) + uint64(typ)
}

func merge(ops *[]Op, m map[OpCode]func(*Machine) error) {
	for opcode, f := range m {
		if opcode >= uint64(len(*ops)) {
			delta := opcode - uint64(len(*ops)) + 1
			*ops = append(*ops, make([]Op, delta)...)
		}
		(*ops)[opcode] = Op{f}
	}
}

type OpTable struct {
	baseOpcodes []OpCode
	ops         []Op
}

func (t OpTable) Halt() OpCode   { return haltOpcode }
func (t OpTable) Call() OpCode   { return callOpcode }
func (t OpTable) Return() OpCode { return returnOpcode }
func (t OpTable) IfThen() OpCode { return ifThenOpcode }
func (t OpTable) IfElse() OpCode { return ifElseOpcode }
func (t OpTable) Else() OpCode   { return elseOpcode }

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
	table := OpTable{
		make([]OpCode, maxOpFamily), /* baseOpcodes */
		nil,                         /* ops */
	}

	opFactories := []func(OpCode) opFamilyMap{
		haltOpFamily:   opHalt,
		callOpFamily:   opCall,
		returnOpFamily: opReturn,
		ifThenOpFamily: opIfThen,
		ifElseOpFamily: opIfElse,
		elseOpFamily:   opElse,
		printOpFamily:  opPrint,
		pushOpFamily:   opPush,
		popOpFamily:    opPop,
		addOpFamily:    opAdd,
	}

	for opFamily, factory := range opFactories {
		base := OpCode(len(table.ops))
		table.baseOpcodes[opFamily] = base
		merge(&table.ops, factory(base))
	}

	return table
}

func newOpTable() OpTable {
	table := OpTable{
		make([]OpCode, maxOpFamily),
		nil, /* ops */
	}

	family := haltOpFamily
	base := haltOpFamily
	for ; family < printOpFamily; family++ {
		table.baseOpcodes[family] = family
		base++
	}

	if base != printOpFamily {
		panic("Invalid op table")
	}

	table.baseOpcodes[printOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			familyBase := table.baseOpcodes[printOpFamily]
			opcode := unaryOpCode(familyBase, mode, typ)
			if opcode != base {
				panic(fmt.Errorf("Invalid op table: family:%d base:%d mode:%d type:%d; want %d; got %d", family, familyBase, mode, typ, base, opcode))
			}
			base++
		}
	}

	table.baseOpcodes[pushOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			if unaryOpCode(table.baseOpcodes[pushOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	table.baseOpcodes[popOpFamily] = base
	for mode := ImmediateMode; mode < maxOpMode; mode++ {
		for typ := I8; typ < maxOpType; typ++ {
			if unaryOpCode(table.baseOpcodes[popOpFamily], mode, typ) != base {
				panic("Invalid op table")
			}
			base++
		}
	}

	table.baseOpcodes[addOpFamily] = base
	for mode1 := ImmediateMode; mode1 < maxOpMode; mode1++ {
		for mode2 := ImmediateMode; mode2 < maxOpMode; mode2++ {
			for typ := I8; typ < maxOpType; typ++ {
				if binaryOpCode(table.baseOpcodes[addOpFamily], mode1, mode2, typ) != base {
					panic("Invalid op table")
				}
				base++
			}
		}
	}

	return table
}

func init() {
	_ = newOpTable()
}
