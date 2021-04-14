package vm

import "testing"

func TestOpCode(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for idx, test := range tests {
		instruction := Make(test.op, test.operands...)
		if len(instruction) != len(test.expected) {
			t.Errorf("test[%02d] instruction has wrong length. expected=%d. got=%d",
				idx, len(test.expected), len(instruction))
		}

		for i, b := range test.expected {
			if instruction[i] != test.expected[i] {
				t.Errorf("wrong byte at pos %d. expected=%d- got=%d", i, b, instruction[i])
			}
		}
	}
}
