using protocol;
using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Threading;

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

        private static ServerConnector _instance;

        private Socket _socket;

        public string Ip { get; set; }

        public int Port { get; set; }

        public int BuffLength { get; set; }

        private SocketError _socketError;

        // 获取连接器的引用
        public static ServerConnector GetInstance()
        {
            if (_instance == null)
            {
                _instance = new ServerConnector();
            }
            return _instance;
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
            _socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            IPEndPoint ip = new IPEndPoint(IPAddress.Parse(Ip), Port);

            IAsyncResult ar = _socket.BeginConnect(ip, new AsyncCallback(OnConnected), _socket);

            bool success = ar.AsyncWaitHandle.WaitOne(5000, true);
            if (!success)
            {
                CloseConnect();
                Console.Out.WriteLine("connection error");
                return false;
            }
            Thread thread = new Thread(new ThreadStart(BeginReceive));
            thread.IsBackground = true;
            thread.Start();
            return true;
        }

        public void SendMessage<T>(int messageType, T dto) 
        {
            MobileSuiteModel msm = new MobileSuiteModel();
            msm.type = messageType;

            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize<T>(ms, dto);
                msm.message = ms.ToArray();
                ms.Close();
            }
            byte[] bytes;
            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize(ms, msm);
                bytes = ms.ToArray();
                ms.Close();
            }

            int length = bytes.Length;
            byte[] data = new byte[length + 4];
            byte[] lengthBytes = BitConverter.GetBytes(length);
            Array.Copy(lengthBytes, data, 4);
            Array.Copy(bytes, 0, data, 4, length);

            SendMessage(data);
        }

        public void SendMessage(byte[] data)
        {
            if (_socket == null) return;
            _socket.BeginSend(data, 0, data.Length, SocketFlags.None, SendCallBack, _socket);
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
            if (_socket != null && _socket.Connected)
            {
                _socket.Shutdown(SocketShutdown.Both);
                _socket.Close();
                OnDisconnected();
            }
            _socket = null;
            return true;
        }

        // 接收消息线程
        private void BeginReceive()
        {
            try
            {
                byte[] data = new byte[BuffLength];
                _socket.BeginReceive(data, 0, data.Length, SocketFlags.None, out _socketError, ReceivedResult, data);
            }
            catch (Exception)
            {
                if (_socket != null)
                {
                    _socket.Close();
                    OnDisconnected();
                }
            }
        }

        private void ReceivedResult(IAsyncResult result)
        {
            int count = 0;
            try
            {
                count = _socket.EndReceive(result);
            }
            catch (SocketException e)
            {
                _socketError = e.SocketErrorCode;
            }
            catch
            {
                _socketError = SocketError.HostDown;
            }


            if (_socketError == SocketError.Success && count > 0)
            {
                byte[] buffer = result.AsyncState as byte[];
                byte[] data = new byte[count];
                if (buffer != null) Array.Copy(buffer, 0, data, 0, data.Length);
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
                if (_socket != null) {
                    _socket.Close();
                    OnDisconnected();
                }
            }
        }

        private void OnConnected(IAsyncResult asyncConnected)
        {
            try
            {
                _socket.EndConnect(asyncConnected);
                if (_socket.Connected)
                {
                    if (OnConnectedEvent != null)
                    {
                        OnConnectedEvent(asyncConnected);
                    }
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
