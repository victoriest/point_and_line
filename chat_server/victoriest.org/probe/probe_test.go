package probe

import (
	"../protocol"
	// "bufio"
	// "bytes"
	"testing"
)

// 构造struct
func genVictoriestMsg() *VictoriestMsg {
	msgObj := &VictoriestMsg{MsgContext: "estest"}
	return msgObj
}

func testUtil(probe *JsonProbe, t *testing.T) {
	msgObj := genVictoriestMsg()
	probe = new(JsonProbe)
	// 序列化
	bt, err := probe.Serialize(msgObj, protocol.MSG_TYPE_TEST_MESSGAE)
	if err != nil {
		t.Error("error on probe.Serialize")
	}
	t.Log("success on probe.Serialize")

	// 反序列化
	var dest VictoriestMsg
	var mst int32
	mst, err = probe.Deserialize(bt, &dest)
	// reader := bufio.NewReader(bytes.NewBuffer(bt))
	// dest, mst, err = probe.DeserializeByReader(reader)
	if err != nil {
		t.Error("error on probe.Deserialize")
	}
	t.Log("success on probe.Deserialize  ", dest.MsgContext, "  ", mst)
}

func TestJsonProbe(t *testing.T) {
	probe := new(JsonProbe)
	testUtil(probe, t)
}

func TestGobProbe(t *testing.T) {
	// probe := new(GobProbe)
	// testUtil(probe, t)
}
