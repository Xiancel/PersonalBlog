package main

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2,3) = %d; want %d", result, expected)
	}
}

func TestMultyply(t *testing.T) {
	result := Multiply(2, 3)
	expected := 6

	if result != expected {
		t.Errorf("Multiply(2,3) = %d; want %d", result, expected)
	}
}

//AAA
//Arrange - підготовка даних
//Act - виконання функції
//Assert - перевірка результату
func TestAddTableDriven(t *testing.T) {
	//Arrange
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{0, 0, 0},
		{-1, 1, 0},
		{100, 200, 300},
	}

	for _, test := range tests {
		//Act
		result := Add(test.a, test.b)
		//Assert
		if result != test.expected {
			t.Errorf("Add(%d,%d) = %d; want %d", test.a, test.b, result, test.expected)
		}
	}
}
