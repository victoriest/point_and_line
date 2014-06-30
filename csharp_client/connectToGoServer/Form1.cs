using protocol;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Formatters.Binary;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace connectToGoServer
{
    public partial class Form1 : Form
    {
        private ServerConnector connector;

        public Form1()
        {
            InitializeComponent();
            connector = ServerConnector.GetInstance();
            connector.OnConnectedEvent += ConnectedCallBack;
            connector.OnRecivedMessageEvent += RecivedMessage;
            connector.OnDisconnectedEvent += DisconnectedCallBack;
            Control.CheckForIllegalCrossThreadCalls = false;
        }

        private void ConnectedCallBack(IAsyncResult ar)
        {
            lbInfo.Items.Add("connected");
        }

        private void RecivedMessage(byte[] data)
        {
            var tmp = new byte[data.Length - 4];
            var lengthByte = new byte[4];
            Array.Copy(data, 0, lengthByte, 0, 4);
            Array.Copy(data, 4, tmp, 0, tmp.Length);

            Stream stream = new MemoryStream(tmp);

            MobileSuiteModel msm = ProtoBuf.Serializer.Deserialize<MobileSuiteModel>(stream);

            int length = BitConverter.ToInt32(lengthByte, 0);
            lbInfo.Items.Add(length.ToString());
            lbInfo.Items.Add(msm.type.ToString());

            stream = new MemoryStream(msm.message);
            ChatMsg chat = ProtoBuf.Serializer.Deserialize<ChatMsg>(stream);
            lbInfo.Items.Add(chat.chatContext);

        }

        private void DisconnectedCallBack()
        {
            lbInfo.Items.Add("disconnected");
        }

        private void btnConnect_Click(object sender, EventArgs e)
        {
            connector.InitSocket(txtIp.Text, Int32.Parse(txtPort.Text));
        }

        private void btnDisconnect_Click(object sender, EventArgs e)
        {
            connector.CloseConnect();
        }

        private void btnSend_Click(object sender, EventArgs e)
        {
            MobileSuiteModel msm = new MobileSuiteModel();
            msm.type = 10;

            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtMsg.Text;

            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize<ChatMsg>(ms, chat);
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

            //byte[] bytes = Encoding.UTF8.GetBytes("{\"MsgType\":10,\"MsgContext\":{\"ChatMessage\":\""+txtMsg.Text+"\"}}");
            int length = bytes.Length;
            byte[] data = new byte[length + 4];
            byte[] lengthBytes = BitConverter.GetBytes(length);
            Array.Copy(lengthBytes, data, 4);
            Array.Copy(bytes, 0, data, 4, length);

            connector.SendMessage(data);
        }
    }
}
