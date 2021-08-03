package gpio

import (
	"testing"
)

func TestDataMarshalling(t *testing.T) {
	id := &ItemData{
		Pin:         0,
		State:       0,
		Mode:        0,
		Key:         "",
		Name:        "",
		Description: "",
	}
	data, err := id.marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}
