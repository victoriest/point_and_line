package probe

/**
 * 序列化与反序列化的包
 */
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// 序列化接口
type Serializable interface {
	Serialize(src interface{}) ([]byte, error)
	Deserialize(src []byte, dst interface{}) (interface{}, error)
}

// 
type Codecable interface {
}

// Josn的序列化实现
type JsonProbe struct{}

func (self JsonProbe) Encoding(obj interface{}) ([]byte, error) {
	var v []byte
	var err error

	v, err = json.Marshal(obj)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	var length int32 = int32(len(v))

	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

func (self JsonProbe) Decoding(reader *bufio.Reader) (interface{}, error) {
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println("when Encoding1:", err.Error())
		return nil, err
	}

	if int32(reader.Buffered()) < length+4 {
		fmt.Println("int32(reader.Buffered()) < length+4")
		_, err := reader.Peek(int(4 + length))
		if err != nil {
			return nil, err
		}
	}
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		fmt.Println("when Encoding2:", err.Error())
		return nil, err
	}

	msg := pack[4:]
	fmt.Println(msg)
	var dst interface{}
	err = json.Unmarshal(msg, &dst)
	if err != nil {
		fmt.Println("when Encoding3:", err.Error())
		return nil, err
	}
	return dst, nil
}

// Gob的序列化实现
type GobProbe struct{}

// gob的序列化方法实现
func (self GobProbe) Serialize(src interface{}) (v []byte, err error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(src)
	if err != nil {
		return
	}
	v = buf.Bytes()
	return
}

// gob的反序列化方法实现
func (self GobProbe) Deserialize(src []byte, dst interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	err = dec.Decode(dst)
	return
}
