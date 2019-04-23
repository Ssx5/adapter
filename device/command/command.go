package command

import (
	"encoding/json"
	"fmt"
	"time"
)

type Command struct {
	Name        string
	Attribution map[string]interface{}
	Period      string
	NextTime    time.Time
}

func (c Command) MarshalJSON() ([]byte, error) {
	d := struct {
		Name        *string                `json:"name"`
		Attribution map[string]interface{} `json:"attribution"`
		Period      *string                `json:"period"`
	}{
		Name:   &c.Name,
		Period: &c.Period,
	}
	d.Attribution = make(map[string]interface{})
	for k, v := range c.Attribution {
		if k == "bytes" {
			bytes := v.([]byte)
			d.Attribution[k] = make([]int, len(bytes))
			for i := 0; i < len(bytes); i++ {
				(d.Attribution[k].([]int))[i] = int(bytes[i])
			}
		} else {
			d.Attribution[k] = v
		}
	}
	return json.Marshal(d)
}

func (c *Command) UnmarshalJSON(b []byte) error {
	d := struct {
		Name        string                 `json:"name"`
		Attribution map[string]interface{} `json:"attribution"`
		Period      string                 `json:"period"`
	}{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		return err
	}
	c.Name = d.Name
	c.Period = d.Period
	c.Attribution = make(map[string]interface{})
	for k, v := range d.Attribution {
		if k == "bytes" {
			arr := v.([]interface{})
			c.Attribution[k] = make([]byte, len(arr))
			for i := 0; i < len(arr); i++ {
				(c.Attribution[k].([]byte))[i] = byte(arr[i].(float64))
			}
		} else {
			c.Attribution[k] = v
		}
	}
	return nil
}

func (c *Command) Check(fields []string) error {
	s := ""
	for _, f := range fields {
		if _, ok := c.Attribution[f]; !ok {
			s += (f) + ","
		}
	}
	if s == "" {
		return nil
	} else {
		return fmt.Errorf("Attribution should have field '%s'", s)
	}
}
