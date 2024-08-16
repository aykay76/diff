package diff

import (
	"os"
	"testing"
)

func TestDiff(t *testing.T) {
	a, err := os.Open("v1.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.Open("v2.txt")
	if err != nil {
		t.Fatal(err)
	}

	diff := NewDiff(a, b)
	t.Log(diff)
	for i := 0; i < len(diff.targetHashList); i++ {
		if diff.targetHashList[i].marked {
			t.Log("Copy ", diff.targetHashList[i].sourceLineNumber, diff.targetHashList[i].numberOfLines)
			i += diff.targetHashList[i].numberOfLines - 1
		} else {
			t.Log("Insert ", diff.targetHashList[i].line)
		}
	}
}
