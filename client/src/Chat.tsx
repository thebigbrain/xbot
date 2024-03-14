import React, { useState, useEffect, useCallback } from "react";
import {
  Container,
  List,
  ListItem,
  ListItemText,
  TextField,
  Button,
  Box,
} from "@mui/material";

// 定义消息接口
interface IMessage {
  sessionID: string;
  sender: string;
  text: string;
  timestamp: Date;
}

const ChatApp: React.FC = () => {
  const [messages, setMessages] = useState<IMessage[]>([]);
  const [newMessage, setNewMessage] = useState("");

  // 获取聊天历史
  const fetchChatHistory = useCallback(async () => {
    try {
      const response = await fetch(
        `http://localhost:5000/api/history?session_id=1`
      );
      const data = await response.json();
      if (response.ok) {
        setMessages(data.history || []);
      } else {
        throw new Error(data.message || "Error fetching data");
      }
    } catch (error: any) {
      console.error("Error fetching chat history:", error.message); // TypeScript 需要 error 为 any 类型，或者自定义 Error 类型
    }
  }, []);

  useEffect(() => {
    fetchChatHistory();
    // 依赖数组中包含 fetchChatHistory 则每当 fetchChatHistory 改变时都会重新执行
  }, [fetchChatHistory]);

  // 发送消息

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) {
      return; // 如果消息为空，则不发送
    }
    try {
      const response = await fetch(`http://localhost:5000/api/send`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: newMessage }),
      });

      if (response.headers.get("Content-Type")?.includes("text/event-stream")) {
        // 处理事件流
        const reader = response.body?.getReader();
        const decoder = new TextDecoder();

        let combinedMessage = ""; // 用于累积合并后文本的字符串

        const processStream = async () => {
          let completeMessage = "";
          while (true && reader) {
            const { done, value } = await reader.read();
            if (done) {
              break;
            }
            const chunk = decoder.decode(value, { stream: true });

            completeMessage += chunk;

            const events = completeMessage.split("\n\n");
            completeMessage = events.pop() || ""; // 保存未完成部分的数据

            for (const event of events) {
              const dataMatch = event.match(/^data: (.*)$/m);
              if (dataMatch) {
                const data = dataMatch[1];
                // 解析数据并且合并text字段
                const newData = JSON.parse(data); // 假设数据是JSON格式且包含text字段
                combinedMessage += newData.text; // 连接新的文本消息

                // 更新消息显示，这里需要你实现一个更新界面的方法(setCombinedMessages)
                setCombinedMessages(combinedMessage);
              }
            }
          }
        };

        if (reader) {
          processStream()
            .then(() => {
              console.log("Finished processing the stream.");
            })
            .catch((error) => {
              console.error("Error while processing the stream:", error);
            });
        }
      } else {
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || "Something went wrong");
        }
      }

      setNewMessage("");
    } catch (error: any) {
      console.error("Error sending message:", error.message);
    }
  };

  // ...[其余代码保持不变]...

  // 假设我们有一个状态来存储合并后的文本消息
  const [combinedMessages, setCombinedMessages] = useState("");

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
          <ListItem>
            <ListItemText primary={`AI: ${combinedMessages}`} />
          </ListItem>
        </List>
        <Box
          component="form"
          sx={{
            display: "flex",
            alignItems: "center",
          }}
          onSubmit={handleSend}
        >
          <TextField
            fullWidth
            label="Message"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            variant="outlined"
            margin="normal"
          />
          <Button type="submit" variant="contained" color="primary">
            Send
          </Button>
        </Box>
      </Box>
    </Container>
  );
};

export default ChatApp;
