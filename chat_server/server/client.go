package server

//这个是必需的 独立运行的执行文件必需使用                                                                 
import "fmt"
import "net"
import "container/list"
import "bytes"

type Client struct {
	//定义一个结构 struct 用于抽象数据类型                                                                
	Name       string
	Incoming   chan string
	Outgoing   chan string
	Conn       net.Conn
	Quit       chan bool
	ClientList *list.List
}

//go可以返回多个值  使用括号()
func (c *Client) Read(buffer []byte) (int, bool) {
	bytesRead, error := c.Conn.Read(buffer)
	//:= 同时定义和赋值，并且可以同时多个变量
	if error != nil {
		//nil is like C's NULL
		c.Close()
		Log(error)
		return 0, false
		//返回的值为bool类型
	}
	Log("Read ", bytesRead, " bytes")
	return bytesRead, true
}

func (c *Client) Close() {
	c.Quit <- true
	//管道传递true给Quit
	c.Conn.Close()
	c.RemoveMe()
	//调用RemoveMe func
}

//Equal 帮助对比client
func (c *Client) Equal(other *Client) bool {
	if bytes.Equal([]byte(c.Name), []byte(other.Name)) {
		if c.Conn == other.Conn {
			return true
		}
	}
	return false
}

func (c *Client) RemoveMe() {
	for entry := c.ClientList.Front(); entry != nil; entry = entry.Next() {
		//双向链表队列的使用
		client := entry.Value.(Client)
		if c.Equal(&client) {
			Log("RemoveMe: ", c.Name)
			c.ClientList.Remove(entry)
		}
	}
}
