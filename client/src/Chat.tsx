import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Container,
  List,
  ListItem,
  ListItemText,
  TextField,
  Button,
  Box,
} from "@mui/material";

// Message接口
interface Message {
  sessionID: string;
  sender: string;
  text: string;
  timestamp: Date;
}

// 设置axios默认值
axios.defaults.baseURL = "http://localhost:5000";
axios.defaults.withCredentials = false;
axios.defaults.headers["content-type"] = "application/json";

const ChatApp = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState("");

  // 获取聊天历史
  const fetchChatHistory = async () => {
    try {
      const response = await axios.get<any>("/api/history", {
        params: { session_id: 1 },
      });
      setMessages(response.data?.history || []);
    } catch (error) {
      console.error("Error fetching chat history:", error);
    }
  };

  // 处理SSE数据并更新消息列表
  const handleSSEData = (data: string) => {
    // 根据SSE协议格式提取消息内容
    if (data.startsWith("data:")) {
      try {
        // 提取“data: ”后面的内容
        const jsonData = data.replace("data: ", "");
        const parsedData: Message = JSON.parse(jsonData);
        setMessages((prevMessages) => [...prevMessages, parsedData]);
      } catch (error) {
        console.error("Error parsing SSE data:", error);
      }
    }
  };

  // 处理发送消息
  const handleSend = async () => {
    try {
      // 使用fetch发送请求
      const response = await fetch("http://localhost:5000/api/send", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: newMessage }),
      });

      if (response.body) {
        const reader = response.body.getReader();

        // 处理接收到的数据
        const processStream = () => {
          let buffer = "";
          const decoder = new TextDecoder();

          reader.read().then(function processResult({ done, value }): any {
            if (done) {
              return;
            }

            // 更新缓冲区
            buffer += decoder.decode(value, { stream: true });
            // 处理完整SSE消息
            let dataIndex = buffer.indexOf("\n\n");
            while (dataIndex !== -1) {
              const message = buffer.slice(0, dataIndex + 1);
              handleSSEData(message.trim());
              buffer = buffer.slice(dataIndex + 2);
              dataIndex = buffer.indexOf("\n\n");
            }
            return reader.read().then(processResult);
          });
        };
        // 开始处理流
        processStream();
      }
      setNewMessage("");
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

  // 在组件挂载时获取聊天历史
  useEffect(() => {
    fetchChatHistory();
  }, []);

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          my: 4,
          display: "flex",
          flexDirection: "column",
          height: "80vh",
          justifyContent: "space-between",
        }}
      >
        <List dense>
          {messages.map((message, index) => (
            <ListItem key={index}>
              <ListItemText primary={`${message.sender}: ${message.text}`} />
            </ListItem>
          ))}
        </List>
        <Box
          component="form"
          sx={{
            display: "flex",
            alignItems: "center",
          }}
          onSubmit={(e) => {
            e.preventDefault();
            handleSend();
          }}
        >
          <TextField
            fullWidth
            label="Message"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            variant="outlined"
            margin="normal"
          />
          <Button variant="contained" color="primary" onClick={handleSend}>
            Send
          </Button>
        </Box>
      </Box>
    </Container>
  );
};

export default ChatApp;
