package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpContstant, []int{65534}, []byte{byte(OpContstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {
		t.Run(string(tt.op), func(t *testing.T) {
			instruction := Make(tt.op, tt.operands...)

			if len(instruction) != len(tt.expected) {
				t.Errorf("instruction has wrong len. want=%d, got=%d", len(tt.expected), len(instruction))
			}

			for i, b := range tt.expected {
				if instruction[i] != tt.expected[i] {
					t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
				}
			}
		})
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpContstant, 2),
		Make(OpContstant, 65535),
	}
	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{}
	for _, instruction := range instructions {
		concatted = append(concatted, instruction...)
	}
	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted. want=%s, got=%s", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpContstant, []int{65535}, 2},
	}
	for _, tt := range tests {
		t.Run(string(tt.op), func(t *testing.T) {
			instructions := Make(tt.op, tt.operands...)

			def, err := Lookup(byte(tt.op))
			if err != nil {
				t.Fatalf("definition not found: %q\n", err)
			}

			operandsRead, n := ReadOperands(def, instructions[1:])
			if n != tt.bytesRead {
				t.Fatalf("n wrong. want=%q, got=%q", tt.bytesRead, n)
			}

			for i, want := range tt.operands {
				if operandsRead[i] != want {
					t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
				}
			}
		})
	}
}
