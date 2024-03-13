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

axios.defaults.baseURL = "http://localhost:5000";
axios.defaults.withCredentials = false;
axios.defaults.headers["content-type"] = "application/json";

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

  function readChunk(r: ReadableStreamDefaultReader) {
    r?.read().then(({ done, value }) => {
      // If there is no more data to read
      if (done) {
        console.log("done", done);
        return;
      }
      // Get the data and send it to the browser via the controller
      // Check chunks by logging to the console
      console.log(done, value);

      readChunk(r);
    });
  }

  const handleSend = async () => {
    try {
      // 使用 axios 发送新消息
      const response = await fetch("http://localhost:5000/api/send", {
        method: "POST", // *GET, POST, PUT, DELETE, etc.
        mode: "no-cors", // no-cors, *cors, same-origin
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, *same-origin, omit
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({}), // body data type must match "Content-Type" header
      });

      console.log(response);
      if (!response.body) return;

      const reader: ReadableStreamDefaultReader = response.body
        .pipeThrough(new TextDecoderStream())
        .getReader();

      readChunk(reader);

      // const newMessageFromResponse: Message = response.data;
      // if (newMessageFromResponse) {
      //   setMessages((msgs) => [...msgs, newMessageFromResponse]);
      //   setNewMessage(""); // 清空输入框
      // }
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
