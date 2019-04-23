package command

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	c := Command{
		Name:   "command01",
		Period: "100ms",
	}
	c.Attribution = make(map[string]interface{})
	c.Attribution["bytes"] = []byte{1, 2, 3, 4}

	b, err := json.Marshal(c)
	if err != nil {
		t.Errorf("fail %s", err)
		return
	}
	result := `{"name":"command01","attribution":{"bytes":[1,2,3,4]},"period":"100ms"}`
	if result != string(b) {
		t.Errorf("fail %s", string(b))
	}
	t.Logf("%s", string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	cmd := `{"name":"command01","attribution":{"bytes":[1,2,3,4],"funCode":"1"},"period":"100ms"}`
	var c Command
	err := json.Unmarshal([]byte(cmd), &c)
	if err != nil {
		t.Errorf("fail %s", err)
	}
	if c.Name != "command01" || c.Period != "100ms" {
		t.Errorf("fail %s %s", c.Name, c.Period)
	}

	bytes := []byte{1, 2, 3, 4}
	res := c.Attribution["bytes"].([]byte)

	t.Logf("%v\n", c.Attribution)

	for i := 0; i < len(bytes); i++ {
		if res[i] != bytes[i] {
			t.Error("fail")
		}
	}
	if c.Attribution["funCode"] != "1" {
		t.Errorf("funCode: %v", c.Attribution["funCode"])
	}
}
