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

func TestClient_CreateSegmentFromFile(t *testing.T) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.CreateSegmentFromFile("upload_from_my_lib", "../test-files/macs_for_uploads.csv", mac); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SaveUploadedSegment(t *testing.T) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	segment, err := client.CreateSegmentFromFile("upload_from_my_lib", "../test-files/macs_for_uploads.csv", mac)
	if err != nil {
		t.Fatal(err)
	}
	err = client.SaveUploadedSegment(segment)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_RemoveSegment(t *testing.T) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	segment, err := client.CreateSegmentFromFile("upload_from_my_lib_to_delete", "../test-files/macs_for_uploads.csv", mac)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SaveUploadedSegment(segment); err != nil {
		t.Fatal(err)
	}
	if ok, err := client.RemoveSegment(segment.Id); !ok || err != nil {
		if err != nil {
			t.Fatal(err)
		} else {
			t.Fatal("can't remove without error")
		}
	}
}
