using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
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

        private void ConnectedCallBack(IAsyncResult ar){
            //MessageBox.Show("connected");
            lbInfo.Items.Add("connected");
        }

        private void RecivedMessage(byte[] data)
        {
            //MessageBox.Show(BitConverter.ToString(data));
            var tmp = new byte[data.Length - 4];
            var lengthByte = new byte[4];
            Array.Copy(data, 0, lengthByte, 0, 4);
            Array.Copy(data, 4, tmp, 0, tmp.Length);
            int length = BitConverter.ToInt32(lengthByte, 0);
            lbInfo.Items.Add(length.ToString());
            var csData = Encoding.UTF8.GetString(tmp, 0, tmp.Length);
            lbInfo.Items.Add(csData);
        }

        private void DisconnectedCallBack() {
            //MessageBox.Show("disconnected");
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
            byte[] bytes = Encoding.UTF8.GetBytes("{\"MsgType\":10,\"MsgContext\":{\"ChatMessage\":\""+txtMsg.Text+"\"}}");
            int length = bytes.Length;
            byte[] data = new byte[length + 4];
            byte[] lengthBytes = BitConverter.GetBytes(length);
            Array.Copy(lengthBytes, data, 4);
            Array.Copy(bytes, 0, data, 4, length);

            connector.SendMessage(data);
        }
    }
}
