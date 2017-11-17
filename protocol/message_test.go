package protocol

import (
	"testing"
	"bytes"
	"fmt"
)

func TestMessage(t *testing.T) {
	req := NewMessage()
	req.SetVersion(0)
	req.SetMessageType(Request)
	req.SetHeartbeat(false)
	req.SetOneway(false)
	req.SetCompressType(None)
	req.SetMessageStatusType(Normal)
	req.SetSerializeType(JSON)
	req.SetModule(10)
	req.SetMessageID(2333)

	payload := `{
		"A": 1,
		"B": 2,
	}
	`
	req.Payload = []byte(payload)
	var buf bytes.Buffer
	err := req.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Read(&buf)
	if err != nil {
		t.Fatal(err)
	}
	res.SetMessageType(Response)

	if res.Version() != 0 {
		t.Errorf("expect 0 but got %d", res.Version())
	}
	if res.Module() != 10 {
		t.Errorf("expect 10 but got %d", res.Module())
	}

	if res.MessageID() != 2333 {
		t.Errorf("expect 2333 but got %d", res.MessageID())
	}

	if string(res.Payload) != payload {
		t.Errorf("got wrong payload: %v", string(res.Payload))
	}
	fmt.Println(res.Module(), "  ", res.MessageID())
	fmt.Println("测试结束", res.ModuleMessageID())
}