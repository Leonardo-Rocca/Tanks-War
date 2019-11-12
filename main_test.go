package main

import (
	"testing"


)

func TestFib(t *testing.T) {
	got := Fibbonacci(10)
	if got != 55{
		t.Errorf("ERROR")
	}
}

func BenchmarkFib( b *testing.B){
	for n := 0 ; n < b.N ; n++{
		Fibbonacci(10)
	}
}



