package parser

import (
	"testing"

	"encoding/json"
)

func TestMarshalNode(t *testing.T) {
	node := new(SSRNode)
	node.NodeName = "unreachable"
	node.Type = "ssr"
	node.IP = "0.0.0.0"
	node.Port = 1
	node.Passwd = "123"
	node.Crypto = "aes"
	node.Proto = "http"
	node.Minx = "non"

	data, err := json.Marshal(node)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestUnmarshalNode(t *testing.T) {
	node := new(SSRNode)
	node.NodeName = "test"
	node.Type = "ssr"

	data := `{"node_name":"a","type":"ss","server":"0.0.0.0","server_port":1,"password":"123","method":"aes","protocol":"http","obfs":"non"}`
	if err := json.Unmarshal([]byte(data), node); err != nil {
		t.Error(err)
	}
	if node.NodeName != "test" || node.Type != "ssr" {
		t.Error("unmarshal error")
	}
	t.Log(*node)
}
