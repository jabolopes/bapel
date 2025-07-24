package stlc

import (
	"testing"

	"github.com/jabolopes/bapel/ir"
)

func TestRenameTypeVars(t *testing.T) {
	tests := []struct {
		name               string
		input              ir.IrType
		expected           ir.IrType
		setupContext       func() Context
		setupSubstitutions func() []substitution
		expectedError      error
	}{
		{
			name:     "VarType - Single variable",
			input:    ir.NewVarType("a"),
			expected: ir.NewVarType("a"),
		},
		{
			name:     "ForallType - Single quantified variable - No rename",
			input:    ir.NewForallType("a", ir.NewTypeKind(), ir.NewVarType("a")),
			expected: ir.NewForallType("a", ir.NewTypeKind(), ir.NewVarType("a")),
		},
		{
			name:  "ForallType - Single quantified variable - Renamed",
			input: ir.NewForallType("a", ir.NewTypeKind(), ir.NewVarType("a")),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewForallType("b", ir.NewTypeKind(), ir.NewVarType("b")),
		},
		{
			name:     "ForallType - Multiple quantified variables - No rename",
			input:    ir.NewForallType("a", ir.NewTypeKind(), ir.NewForallType("b", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")))),
			expected: ir.NewForallType("a", ir.NewTypeKind(), ir.NewForallType("b", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")))),
		},
		{
			name:  "ForallType - Partially bound variables - Renamed",
			input: ir.NewForallType("a", ir.NewTypeKind(), ir.NewForallType("b", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")))),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewForallType("b", ir.NewTypeKind(), ir.NewForallType("c", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("b"), ir.NewVarType("c")))),
		},
		{
			name:  "ForallType - Multiple bound variables - Renamed",
			input: ir.NewForallType("a", ir.NewTypeKind(), ir.NewForallType("b", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")))),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}
				context, err = context.AddBind(NewTypeVarBind("b", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewForallType("c", ir.NewTypeKind(), ir.NewForallType("d", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("c"), ir.NewVarType("d")))),
		},
		{
			name:     "LambdaType - Single abstracted variable - No rename",
			input:    ir.NewLambdaType("a", ir.NewTypeKind(), ir.NewVarType("a")),
			expected: ir.NewLambdaType("a", ir.NewTypeKind(), ir.NewVarType("a")),
		},
		{
			name:  "LambdaType - Single abstracted variable - Renamed",
			input: ir.NewLambdaType("a", ir.NewTypeKind(), ir.NewVarType("a")),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewLambdaType("b", ir.NewTypeKind(), ir.NewVarType("b")),
		},
		{
			name:  "LambdaType - Partially bound variables - Renamed",
			input: ir.NewLambdaType("a", ir.NewTypeKind(), ir.NewLambdaType("b", ir.NewTypeKind(), ir.NewVarType("a"))),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewLambdaType("b", ir.NewTypeKind(), ir.NewLambdaType("c", ir.NewTypeKind(), ir.NewVarType("b"))),
		},
		{
			name:  "LambdaType - Multiple bound variables - Renamed",
			input: ir.NewLambdaType("a", ir.NewTypeKind(), ir.NewLambdaType("b", ir.NewTypeKind(), ir.NewVarType("a"))),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}
				context, err = context.AddBind(NewTypeVarBind("b", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewLambdaType("c", ir.NewTypeKind(), ir.NewLambdaType("d", ir.NewTypeKind(), ir.NewVarType("c"))),
		},
		{
			name:     "FunType - no substitutions",
			input:    ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")),
			expected: ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")),
		},
		{
			name:  "FunType - argument substitutions",
			input: ir.NewFunctionType(ir.NewVarType("a"), ir.NewVarType("b")),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("a"), ir.NewVarType("c")},
					{ir.NewVarType("b"), ir.NewVarType("d")},
				}
			},
			expected: ir.NewFunctionType(ir.NewVarType("c"), ir.NewVarType("d")),
		},
		{
			name:     "AppType - no substitutions",
			input:    ir.NewAppType(ir.NewVarType("f"), ir.NewVarType("arg")),
			expected: ir.NewAppType(ir.NewVarType("f"), ir.NewVarType("arg")),
		},
		{
			name:  "AppType - Function and argument substituted",
			input: ir.NewAppType(ir.NewVarType("fun"), ir.NewVarType("arg")),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("fun"), ir.NewVarType("a")},
					{ir.NewVarType("arg"), ir.NewVarType("b")},
				}
			},
			expected: ir.NewAppType(ir.NewVarType("a"), ir.NewVarType("b")),
		},
		{
			name:  "ArrayType - Element type freshened",
			input: ir.NewArrayType(ir.NewVarType("elem"), 10),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("elem"), ir.NewVarType("a")},
				}
			},
			expected: ir.NewArrayType(ir.NewVarType("a"), 10),
		},
		{
			name: "StructType - Fields freshened",
			input: ir.NewStructType([]ir.StructField{
				{ID: "f1", Type: ir.NewVarType("a")},
				{ID: "f2", Type: ir.NewFunctionType(ir.NewVarType("b"), ir.NewVarType("c"))},
			}),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("a"), ir.NewVarType("b")},
					{ir.NewVarType("b"), ir.NewVarType("c")},
				}
			},
			expected: ir.NewStructType([]ir.StructField{
				{ID: "f1", Type: ir.NewVarType("b")},
				{ID: "f2", Type: ir.NewFunctionType(ir.NewVarType("c"), ir.NewVarType("c"))},
			}),
		},
		{
			name: "TupleType - Elements freshened",
			input: ir.NewTupleType([]ir.IrType{
				ir.NewVarType("a"),
				ir.NewVarType("b"),
			}),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("a"), ir.NewVarType("b")},
					{ir.NewVarType("b"), ir.NewVarType("c")},
				}
			},
			expected: ir.NewTupleType([]ir.IrType{
				ir.NewVarType("b"),
				ir.NewVarType("c"),
			}),
		},
		{
			name: "VariantType - Tags freshened",
			input: ir.NewVariantType([]ir.VariantTag{
				{ID: "TagA", Type: ir.NewVarType("a")},
				{ID: "TagB", Type: ir.NewFunctionType(ir.NewVarType("b"), ir.NewVarType("c"))},
			}),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("a"), ir.NewVarType("b")},
					{ir.NewVarType("b"), ir.NewVarType("c")},
				}
			},
			expected: ir.NewVariantType([]ir.VariantTag{
				{ID: "TagA", Type: ir.NewVarType("b")},
				{ID: "TagB", Type: ir.NewFunctionType(ir.NewVarType("c"), ir.NewVarType("c"))},
			}),
		},
		{
			name:     "NameType - Remains unchanged",
			input:    ir.NewNameType("i8"),
			expected: ir.NewNameType("i8"),
		},
		{
			name:  "Complex nested type",
			input: ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewArrayType(ir.NewAppType(ir.NewNameType("List"), ir.NewVarType("a")), 5))),
			setupSubstitutions: func() []substitution {
				return []substitution{
					{ir.NewVarType("a"), ir.NewVarType("b")},
				}
			},
			expected: ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewArrayType(ir.NewAppType(ir.NewNameType("List"), ir.NewVarType("a")), 5))),
		},
		{
			name:  "Complex nested type",
			input: ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("a"), ir.NewArrayType(ir.NewAppType(ir.NewNameType("List"), ir.NewVarType("a")), 5))),
			setupContext: func() Context {
				context := NewContext()

				var err error
				context, err = context.AddBind(NewTypeVarBind("a", ir.NewTypeKind()))
				if err != nil {
					t.Fatal(err)
				}

				return context
			},
			expected: ir.NewForallType("b", ir.NewTypeKind(), ir.NewFunctionType(ir.NewVarType("b"), ir.NewArrayType(ir.NewAppType(ir.NewNameType("List"), ir.NewVarType("b")), 5))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var context Context
			if tt.setupContext == nil {
				context = NewContext()
			} else {
				context = tt.setupContext()
			}

			var substitutions []substitution
			if tt.setupSubstitutions != nil {
				substitutions = tt.setupSubstitutions()
			}

			result, err := renameTypeVarsWithSubstitutions(context, tt.input, substitutions)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if !ir.EqualsType(result, tt.expected) {
					t.Errorf("renameTypeVars(%s) = %s, want %s", tt.input.String(), result.String(), tt.expected.String())
				}
			}
		})
	}
}
