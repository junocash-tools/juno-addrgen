package addrgen

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type vectorsV1 struct {
	Version   int      `json:"version"`
	UFVK      string   `json:"ufvk"`
	Addresses []string `json:"addresses"`
}

func loadVectors(t *testing.T) vectorsV1 {
	t.Helper()

	path := filepath.Join("..", "..", "vectors", "v1.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}

	var v vectorsV1
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}
	if v.Version != 1 {
		t.Fatalf("unexpected vectors version: %d", v.Version)
	}
	if v.UFVK == "" {
		t.Fatalf("missing ufvk")
	}
	if len(v.Addresses) != 100 {
		t.Fatalf("unexpected address count: %d", len(v.Addresses))
	}
	return v
}

func TestDerive_GoldenVectors(t *testing.T) {
	v := loadVectors(t)

	tests := []struct {
		name  string
		index uint32
		want  string
	}{
		{"index0", 0, v.Addresses[0]},
		{"index1", 1, v.Addresses[1]},
		{"index10", 10, v.Addresses[10]},
		{"index99", 99, v.Addresses[99]},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Derive(v.UFVK, tc.index)
			if err != nil {
				t.Fatalf("Derive error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected address\nwant: %s\ngot:  %s", tc.want, got)
			}
		})
	}
}

func TestBatch_GoldenVectors(t *testing.T) {
	v := loadVectors(t)

	got, err := Batch(v.UFVK, 0, 10)
	if err != nil {
		t.Fatalf("Batch error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("unexpected len: %d", len(got))
	}
	for i := 0; i < 10; i++ {
		if got[i] != v.Addresses[i] {
			t.Fatalf("mismatch at %d\nwant: %s\ngot:  %s", i, v.Addresses[i], got[i])
		}
	}
}

func TestErrors(t *testing.T) {
	_, err := Derive("", 0)
	var ae *Error
	if !errors.As(err, &ae) || ae.Code != ErrUFVKEmpty {
		t.Fatalf("expected %q, got %v", ErrUFVKEmpty, err)
	}

	v := loadVectors(t)
	_, err = Derive(v.Addresses[0], 0)
	if !errors.As(err, &ae) || ae.Code != ErrUFVKHrpMismatch {
		t.Fatalf("expected %q, got %v", ErrUFVKHrpMismatch, err)
	}

	_, err = Batch("j1not-a-ufvk", 0, 0)
	if !errors.As(err, &ae) || ae.Code != ErrCountZero {
		t.Fatalf("expected %q, got %v", ErrCountZero, err)
	}
}
