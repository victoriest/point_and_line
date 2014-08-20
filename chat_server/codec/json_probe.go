package codec

import (
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"encoding/json"
)

// Josn的序列化实现
type JsonProbe struct{}

/**
 * Json的序列化实现
 */
func (self JsonProbe) Serialize(src *VictoriestMsg) ([]byte, error) {
	var v []byte
	var err error

	// vMsg := &VictoriestMsg{MsgType: msgType, MsgContext: src}

	v, err = json.Marshal(src)
	if err != nil {
		log.Error("when Encoding:", err.Error())
		return nil, err
	}

	// 序列化后的bytes长度
	var length int32 = int32(len(v))

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

/**
 * Json的反序列化实现
 */
func (self JsonProbe) Deserialize(src []byte, dst *VictoriestMsg) (int32, error) {
	// msg 序列化后的对象
	msg := src[4:]

	// 开始饭反序列化
	err := json.Unmarshal(msg, dst)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	return dst.MsgType, nil
}
