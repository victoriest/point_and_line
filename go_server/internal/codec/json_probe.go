package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"go_server/pkg/log"
)

// JSONProbe Json的序列化实现
type JSONProbe struct{}

// Serialize Json的序列化实现
func (jsonProbe *JSONProbe) Serialize(src interface{}) ([]byte, error) {
	var v []byte
	var err error

	// vMsg := &VictoriestMsg{MsgType: msgType, MsgContext: src}

	v, err = json.Marshal(src)
	if err != nil {
		log.Error("when Encoding:", err.Error())
		return nil, err
	}

	// 序列化后的bytes长度
	var length = int32(len(v))

	pkg := new(bytes.Buffer)
	// 写入长度信息
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when Write length:", err.Error())
		return nil, err
	}

	// 写入序列化后的对象
	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		log.Error("when Serialize:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

// Deserialize Json的反序列化实现
func (jsonProbe JSONProbe) Deserialize(src []byte, dst interface{}) (int32, error) {
	var dstObj = dst.(*VictoriestMsg)
	// msg 序列化后的对象
	msg := src[4:]

	// 开始饭反序列化
	err := json.Unmarshal(msg, dstObj)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	return dstObj.MsgType, nil
}
