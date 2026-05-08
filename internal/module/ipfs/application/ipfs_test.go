package application

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
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

func TestA(t *testing.T) {
	saveContent := strings.Builder{}
	saveContent.WriteString(carbon.Now().ToDateString())
	saveContent.WriteString("\t")
	saveContent.WriteString(fmt.Sprintf("%25.4f", 2.2222222))
	saveContent.WriteString("\t")
	saveContent.WriteString(fmt.Sprintf("%20.4f", baseline))
	saveContent.WriteString("\t")
	saveContent.WriteString(fmt.Sprintf("%25.4f", 2.2222222))
	t.Log(saveContent.String())
}
