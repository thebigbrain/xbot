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

interface Message {
  user: string;
  content: string;
}

axios.defaults.baseURL = "https://localhost:5000";

const ChatApp = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState("");

  const fetchChatHistory = async () => {
    try {
      // 使用 axios 请求聊天历史
      const response = await axios.get<Message[]>("/api/history");
      setMessages(response.data);
    } catch (error) {
      console.error("Error fetching chat history:", error);
    }
  };

  useEffect(() => {
    fetchChatHistory();
    // Optional: Setup a polling mechanism if needed
  }, []);

  const handleSend = async () => {
    try {
      // 使用 axios 发送新消息
      const response = await axios.post("/api/send", {
        user: "current_user",
        content: newMessage,
      });

      const newMessageFromResponse: Message = response.data;
      if (newMessageFromResponse) {
        setMessages((msgs) => [...msgs, newMessageFromResponse]);
        setNewMessage(""); // 清空输入框
      }
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

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
              <ListItemText primary={`${message.user}: ${message.content}`} />
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
