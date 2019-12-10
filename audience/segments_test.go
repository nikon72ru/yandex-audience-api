package audience

import (
	"context"
	"testing"
)

func TestClient_SegmentsList(t *testing.T) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.SegmentsList()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateAppMetrikaSegment(t *testing.T) {
	t.Fatal("implement me")
}

func TestClient_CreateCircleGeoSegment(t *testing.T) {
	t.Fatal("implement me")
}

func TestClient_CreateFileSegment(t *testing.T) {
	t.Fatal("implement me")
}

func TestClient_CreateLookalikeSegment(t *testing.T) {
	t.Fatal("implement me")
}
