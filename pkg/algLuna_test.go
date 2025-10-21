package pkg

import "testing"

func test_genRandomDigitsStr(t *testing.T, luna *Luna) {
	dataInput := 5
	if st := luna.genRandomDigitsStr(dataInput); len(st) != dataInput {
		t.Errorf("genRandomDigitsStr: expected/actual: %v/%v", dataInput, len(st))
	}
}

func test_preparCheckSum(t *testing.T, luna *Luna) {
	dataInput := "12345"
	expected := "5"
	if actual := luna.preparCheckSum(dataInput); expected != actual {
		t.Errorf("preparCheckSum: expected/actual: %v/%v", expected, actual)
	}
}

func test_CheckDigits(t *testing.T, luna *Luna) {
	dataTests := []struct {
		nameTest  string
		dataInput string
		expected  bool
	}{
		{
			nameTest:  "CheckDigits test positive",
			dataInput: "96357934",
			expected:  true,
		},
		{
			nameTest:  "CheckDigits test negative",
			dataInput: "96357933",
			expected:  false,
		},
	}

	for _, tt := range dataTests {
		t.Run(tt.nameTest, func(t *testing.T) {
			if actual := luna.CheckDigits(tt.dataInput); tt.expected != actual {
				t.Errorf("CheckDigits: expected/actual: %v/%v", tt.expected, actual)
			}
		})
	}
}

func test_GenDigits(t *testing.T, luna *Luna) {
	// GenDigits test positive
	result := luna.GenDigits(5)
	actual := luna.CheckDigits(result)
	expected := true

	if expected != actual {
		t.Errorf("GenDigits test positive: expected/actual: %v/%v", expected, actual)
	}

	// GenDigits test negative
	result = "12345"
	actual = luna.CheckDigits(result)
	expected = false

	if expected != actual {
		t.Errorf("GenDigits test negative: expected/actual: %v/%v", expected, actual)
	}
}

func TestLuna(t *testing.T) {
	// Инициализация
	luna := NewLuna()

	test_genRandomDigitsStr(t, luna)
	test_preparCheckSum(t, luna)
	test_CheckDigits(t, luna)
	test_GenDigits(t, luna)
}
