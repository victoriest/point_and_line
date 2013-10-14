package main

//这个是必需的 独立运行的执行文件必需使用                                                                 
import "fmt"
import "net"
import "container/list"
import "bytes"

//加载其它的库                                                                                            
//fmt 实现格式化的I/O　类似printf                                                                         
//net 实现网络协议TCP UDP 在这里使用TCP                                                                   
//container/list 实现双向链表队列                                                                         
//bytes 字符数组操作相关功能                                                                              

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

//...可变数组 在屏幕输出日志
//interface 代表任意数据类型 Everything in Go is derived from interface, so in theory, this can be any datatype.
func Log(v ...interface{}) {
	fmt.Println(v...)
}

func IOHandler(Incoming <-chan string, clientList *list.List) {
	for {
		//在go里面没有while do ，for可以无限循环
		Log("IOHandler: Waiting for input")
		input := <-Incoming
		Log("IOHandler: Handling ", input)
		for e := clientList.Front(); e != nil; e = e.Next() {
			client := e.Value.(Client)
			client.Incoming <- input
		}
	}
}

func ClientReader(client *Client) {
	//make创建切片slice 长度2048 http://golang.org/doc/effective_go.html#allocation_make
	buffer := make([]byte, 2048)
	for {
		n, status := client.Read(buffer)
		if !status {
			break
		}

		if bytes.Equal(buffer, []byte("/quit")) {
			client.Close()
			break
		}
		Log("ClientReader received ", client.Name, "> ", string(buffer[:n]))
		send := client.Name + "> " + string(buffer[:n])
		client.Outgoing <- send
	}
	client.Outgoing <- client.Name + " has left chat"
	Log("ClientReader stopped for ", client.Name)
}

func ClientSender(client *Client) {
	for {
		select {
		//select http://golang.org/doc/go_spec.html#Select_statements
		case buffer := <-client.Incoming:
			Log("ClientSender sending ", string(buffer), " to ", client.Name)
			client.Conn.Write([]byte(buffer))
		case <-client.Quit:
			Log("Client ", client.Name, " quitting")
			client.Conn.Close()
			break
		}
	}
}

//该函数主要是接受新的连接和注册用户在client list
func ClientHandler(conn net.Conn, ch chan string, clientList *list.List) {
	buffer := make([]byte, 1024)
	bytesRead, error := conn.Read(buffer)
	if error != nil {
		Log("Client connection error: ", error)
		return
	}
	name := string(buffer[0:bytesRead])
	//      newClient := &Client{name,make(chan string),ch,conn,make(chan bool),clientList}
	//      初始化struct赋值的两种方法
	newClient := &Client{
		Name:       name,
		Incoming:   make(chan string),
		Outgoing:   ch,
		Conn:       conn,
		Quit:       make(chan bool),
		ClientList: clientList,
	}
	//创建go的线程 使用Goroutine
	go ClientSender(newClient)
	go ClientReader(newClient)
	clientList.PushBack(*newClient)
	ch <- string(name + " has joined the chat ")
}

func main() {
	Log("Hello Server!")
	clientList := list.New()
	in := make(chan string)
	//创建一个管道 chan map 需要make creates slices, maps, and channels only
	go IOHandler(in, clientList)
	netListen, error := net.Listen("tcp", ":9988")
	if error != nil {
		Log(error)
	} else {
		//defer函数退出时执行
		defer netListen.Close()
		for {
			Log("Waiting for clients")
			connection, error := netListen.Accept()
			if error != nil {
				Log("Client error: ", error)
			} else {
				go ClientHandler(connection, in, clientList)
			}
		}
	}
}
