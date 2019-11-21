package audience

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_SegmentsList(t *testing.T) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	segments, err := client.SegmentsList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(segments)
}
