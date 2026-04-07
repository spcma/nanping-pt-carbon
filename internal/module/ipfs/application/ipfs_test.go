package application

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIpfsTime(t *testing.T) {
	type tim struct {
		Time *time.Time
		Name string `json:"name"`
	}

	var ti tim
	err := json.Unmarshal([]byte(`{"name":"test"}`), &ti)
	if err != nil {
		t.Fatal(err)
	}

	if ti.Time == nil {
		t.Log("Time is nil")
	}

	if ti.Time.IsZero() {
		t.Log("Time is zero")
	} else {
		t.Log("Time is not zero")
	}
}
