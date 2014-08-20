package codec

import (
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"encoding/gob"
)

// Gob的序列化实现
type GobProbe struct{}

// gob的序列化方法实现
func (self GobProbe) Serialize(src *VictoriestMsg) ([]byte, error) {
	// 序列化
	// vMsg := &VictoriestMsg{MsgType: msgType, MsgContext: src}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(src)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}
	v := buf.Bytes()
	// 序列化后的byte长度
	var length int32 = int32(len(v))

	// 将长度信息写入byte数组
	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	// 写入序列化后的对象
	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

// gob的反序列化方法实现
func (self GobProbe) Deserialize(src []byte, dst *VictoriestMsg) (int32, error) {
	// msg 序列化后的对象
	msg := src[4:]
	buf := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&dst)
	if err != nil {
		log.Error("when GobProbe.Deserialize:", err.Error())
		println("when GobProbe.Deserialize:", err.Error())
		return -1, err
	}

	return dst.MsgType, nil
}
