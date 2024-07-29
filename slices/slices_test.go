package slices

import (
	"testing"
)

func TestFindContentIndexes(t *testing.T) {
	const (
		OpToken string = "("
		ClToken string = ")"
	)

	var (
		ContentTokens []string = []string{
			"(", "(", "a", "+", "b", ")", "*", "c", ")", "+", "d",
		}
	)

	indices, err := FindContentIndexes(OpToken, ClToken, ContentTokens)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if indices[0] != 1 {
		t.Errorf("expected 1, got %d instead", indices[0])
	}

	if indices[1] != 9 {
		t.Errorf("expected 9, got %d instead", indices[1])
	}
}
