package cmd

import (
	"testing"
)

// UniqueSorted() のテスト
func TestUniqueSorted(t *testing.T) {
	input := []int{3, 1, 2, 2, 3, 4, 1}
	expected := []int{1, 2, 3, 4}

	result := uniqueSorted(input)

	// 長さが一致するか検証
	if len(result) != len(expected) {
		t.Fatalf("結果の長さが一致しません: got %d, want %d", len(result), len(expected))
	}

	// 値が一致するか検証
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("要素が一致しません: got %d, want %d", v, expected[i])
		}
	}
}

func TestUniqueSorted_EmptySlice(t *testing.T) {
	result := uniqueSorted([]int{})
	if len(result) != 0 {
		t.Errorf("空スライスに対して空スライスが返されませんでした: %v", result)
	}
}
