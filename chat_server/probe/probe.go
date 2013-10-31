package probe

/**
 * 用于序列化的包
 */
import (
	"bufio"
	"bytes"
	"encoding/binary"
	// "encoding/gob"
	"encoding/json"
	"fmt"
)

func Encoding(obj interface{}) ([]byte, error) {
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

func Decoding(reader *bufio.Reader) (interface{}, error) {
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

// type LGISerialize interface {
// 	Serialize(src interface{}) (dst []byte, err error)
// 	Deserialize(src []byte, dst interface{}) (err error)
// }

// type LGGobSerialize struct {
// }

// // serialize encodes a value using gob.
// func (self LGGobSerialize) Serialize(src interface{}) (v []byte, err error) {
// 	buf := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buf)
// 	err = enc.Encode(src)
// 	if err != nil {
// 		return
// 	}
// 	v = buf.Bytes()
// 	return
// }

// // deserialize decodes a value using gob.
// func (self LGGobSerialize) Deserialize(src []byte, dst interface{}) (err error) {
// 	dec := gob.NewDecoder(bytes.NewBuffer(src))
// 	err = dec.Decode(dst)
// 	return
// }
