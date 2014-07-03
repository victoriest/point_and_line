using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace connectToGoServer
{
    class ServerConnector
    {

        public delegate void ConnectedDelegate(IAsyncResult result);
        public event ConnectedDelegate OnConnectedEvent;

        public delegate void RecivedDelegate(byte[] data);
        public event RecivedDelegate OnRecivedMessageEvent;

        public delegate void DisconnectedDelegate();
        public event DisconnectedDelegate OnDisconnectedEvent;

        private static ServerConnector instance = null;

        private Socket socket = null;

        public string Ip { get; set; }

        public int Port { get; set; }

        public int BuffLength { get; set; }

        private static ManualResetEvent connectDone = new ManualResetEvent(false);
        private static ManualResetEvent sendDone = new ManualResetEvent(false);
        private static ManualResetEvent receiveDone = new ManualResetEvent(false);

        private SocketError socketError;

        // 获取连接器的引用
        public static ServerConnector GetInstance()
        {
            if (instance == null)
            {
                instance = new ServerConnector();
            }
            return instance;
        }

        private ServerConnector()
        {
            BuffLength = 4096;
        }

        // 初始化连接
        public bool InitSocket(string ipParam, int portParam)
        {
            Ip = ipParam;
            Port = portParam;
            return InitSocket();
        }

        // 初始化连接
        public bool InitSocket()
        {
            socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            IPEndPoint ip = new IPEndPoint(IPAddress.Parse(Ip), Port);

            IAsyncResult ar = socket.BeginConnect(ip, new AsyncCallback(OnConnected), socket);

            bool success = ar.AsyncWaitHandle.WaitOne(5000, true);
            if (!success)
            {
                CloseConnect();
                Console.Out.WriteLine("connection error");
                return false;
            }
            else
            {
                Thread thread = new Thread(new ThreadStart(BeginReceive));
                thread.IsBackground = true;
                thread.Start();
            }
            return true;
        }

        public void SendMessage(byte[] data)
        {
            if (socket == null) return;
            socket.BeginSend(data, 0, data.Length, SocketFlags.None, SendCallBack, socket);
        }

        private void SendCallBack(IAsyncResult result)
        {
            try
            {
                Socket sock = result.AsyncState as Socket;
                if (sock != null)
                {
                    sock.EndSend(result);
                }
            }
            catch
            {
            }
        }

        // 关闭连接
        public bool CloseConnect()
        {
            if (socket != null && socket.Connected)
            {
                socket.Shutdown(SocketShutdown.Both);
                socket.Close();
                OnDisconnected();
            }
            socket = null;
            return true;
        }

        // 接收消息线程
        private void BeginReceive()
        {
            try
            {
                byte[] data = new byte[BuffLength];
                socket.BeginReceive(data, 0, data.Length, SocketFlags.None, out socketError, ReceivedResult, data);
            }
            catch (Exception e)
            {
                if (socket != null)
                {
                    socket.Close();
                    OnDisconnected();
                }
            }
        }

        private void ReceivedResult(IAsyncResult result)
        {
            int count = 0;
            try
            {
                count = socket.EndReceive(result);
            }
            catch (SocketException e)
            {
                socketError = e.SocketErrorCode;
            }
            catch
            {
                socketError = SocketError.HostDown;
            }


            if (socketError == SocketError.Success && count > 0)
            {
                byte[] buffer = result.AsyncState as byte[];
                byte[] data = new byte[count];
                Array.Copy(buffer, 0, data, 0, data.Length);
                if (OnRecivedMessageEvent != null)
                {
                    OnRecivedMessageEvent(data);
                }
                //if (this.BinaryInput != null)
                //    DataMM(data);
                BeginReceive();
            }
            else
            {
                if (socket != null) {
                    socket.Close();
                    OnDisconnected();
                }
            }
        }

        private void OnConnected(IAsyncResult asyncConnected)
        {
            try
            {
                socket.EndConnect(asyncConnected);
                if (socket.Connected)
                {

                    if (OnConnectedEvent != null)
                    {
                        OnConnectedEvent(asyncConnected);
                    }
                }
                else
                {
                }
            }
            catch (Exception)
            {
            }
        }

        private void OnDisconnected()
        {
            if (OnDisconnectedEvent != null)
            {
                OnDisconnectedEvent();
            }
        }
    }
}
