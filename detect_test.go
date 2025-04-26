package main

import (
	"math"
	"testing"
)

func floatEquals(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func TestCosineDistance_Orthogonal(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{0, 1, 0}
	dist, err := cosineDistance(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !floatEquals(dist, 1.0, 1e-6) {
		t.Errorf("expected distance 1.0, got %v", dist)
	}
}

func TestCosineDistance_SameVector(t *testing.T) {
	v := []float64{2, 3, 4}
	dist, err := cosineDistance(v, v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !floatEquals(dist, 0.0, 1e-6) {
		t.Errorf("expected distance 0.0, got %v", dist)
	}
}

func TestCosineDistance_MismatchedLengths(t *testing.T) {
	a := []float64{1, 2}
	b := []float64{1, 2, 3}
	_, err := cosineDistance(a, b)
	if err == nil {
		t.Fatal("expected length-mismatch error, got nil")
	}
}
