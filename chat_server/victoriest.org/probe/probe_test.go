package probe

import (
	"testing"
)

// 构造struct
func genVictoriestMsg() *VictoriestMsg {
	msgObj := &VictoriestMsg{MsgType: 1, MsgContext: "msg"}
	return msgObj
}

func testUtil(probe Serializable, t *testing.T) {
	msgObj := genVictoriestMsg()

	// 序列化
	bt, err := probe.Serialize(msgObj)
	if err != nil {
		t.Error("error on probe.Serialize")
	}
	t.Log("success on probe.Serialize")

	// 反序列化
	var dest VictoriestMsg
	err = probe.Deserialize(bt, &dest)
	if err != nil {
		t.Error("error on probe.Deserialize")
	}
	t.Log("success on probe.Deserialize  ", dest)
}

func TestJsonProbe(t *testing.T) {
	probe := new(JsonProbe)
	testUtil(probe, t)
}

func TestGobProbe(t *testing.T) {
	probe := new(GobProbe)
	testUtil(probe, t)
}
