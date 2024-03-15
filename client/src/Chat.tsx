import React, {
  useState,
  useEffect,
  useCallback,
  useRef,
  useLayoutEffect,
} from "react";
import {
  Container,
  List,
  ListItem,
  TextField,
  Button,
  Box,
  Typography,
  CircularProgress,
} from "@mui/material";
import ReactMarkdown from "react-markdown";
// import "github-markdown-css"; // 导入样式文件
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
// 导入一个样式主题，例如：prism
import { materialLight } from "react-syntax-highlighter/dist/esm/styles/prism";

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
  const [isLoading, setIsLoading] = React.useState(false);

  const messagesEndRef = useRef<HTMLDivElement>();

  const scrollToBottom = () => {
    setTimeout(() => {
      messagesEndRef.current?.scrollIntoView({ behavior: "auto" });
    });
  };

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
    scrollToBottom();
  }, []);

  useLayoutEffect(() => {
    scrollToBottom();
  }, []);

  useEffect(() => {
    fetchChatHistory();
  }, []);

  // 发送消息

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) {
      return; // 如果消息为空，则不发送
    }

    setIsLoading(true);

    setMessages((prev) => [
      ...prev,
      {
        sessionID: "",
        sender: "human",
        text: newMessage,
        timestamp: new Date(),
      },
    ]);
    setNewMessage("");
    scrollToBottom();

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
                scrollToBottom();
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
            })
            .finally(() => {
              fetchChatHistory();
            });
        }
      } else {
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || "Something went wrong");
        }
      }
    } catch (error: any) {
      console.error("Error sending message:", error.message);
    }

    setIsLoading(false);
  };

  // 假设我们有一个状态来存储合并后的文本消息
  const [combinedMessages, setCombinedMessages] = useState("");

  const renderers = {
    // 此方法用于渲染 Markdown 中的代码块
    code({ node, inline, className, children, ...props }: any) {
      const match = /language-(\w+)/.exec(className || "");
      return !inline && match ? (
        <SyntaxHighlighter
          style={materialLight}
          language={match[1]}
          PreTag="div"
          children={String(children).replace(/\n$/, "")}
          {...props}
        />
      ) : (
        <code className={className} {...props}>
          {children}
        </code>
      );
    },
  };

  const StyledListItem = ({ message }: { message: IMessage }) => (
    <Box
      sx={{
        mb: 2,
        p: 2,
        backgroundColor: message.sender === "human" ? "#eee" : "aliceblue",
        borderRadius: "10px",
      }}
    >
      <Typography variant="caption" display="block">
        {message.sender === "human" ? "You" : "Bot"} -{" "}
        {new Date(message.timestamp).toLocaleTimeString()}
      </Typography>
      <ReactMarkdown components={renderers} children={message.text} />
    </Box>
  );

  return (
    <Container maxWidth="md">
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          height: "calc(100vh - 16px)",
          mt: 1,
          mb: 1,
        }}
      >
        <Box
          sx={{
            overflowY: "auto",
            overflowX: "hidden",
            flexGrow: 1,
            "& .markdown-body": {
              // 应用样式
              padding: (theme) => theme.spacing(2),
            },
          }}
        >
          <List sx={{ paddingTop: 0, paddingBottom: 0 }}>
            {/* 显示消息 */}
            {messages.map((message, index) => (
              <ListItem
                key={index}
                alignItems="flex-start"
                sx={{ display: "block" }}
              >
                <StyledListItem message={message} />
              </ListItem>
            ))}

            {isLoading ? (
              <CircularProgress />
            ) : (
              combinedMessages && (
                <ListItem alignItems="flex-start" sx={{ display: "block" }}>
                  <StyledListItem
                    message={{
                      sessionID: "",
                      text: combinedMessages,
                      sender: "bot",
                      timestamp: new Date(),
                    }}
                  />
                </ListItem>
              )
            )}

            <div
              style={{ float: "left", clear: "both" }}
              ref={(el) => (el ? (messagesEndRef.current = el) : null)}
            ></div>
          </List>
        </Box>
        <Box
          component="form"
          sx={{
            display: "flex",
            alignItems: "center",
            pt: 1,
            pb: 1,
          }}
          onSubmit={handleSend}
        >
          <TextField
            fullWidth
            label="Message"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            variant="outlined"
            size="small"
            sx={{ mr: 1, flex: 1 }}
          />
          <Button
            type="submit"
            variant="contained"
            color="primary"
            sx={{ height: "40px", flexShrink: 0 }}
          >
            Send
          </Button>
        </Box>
      </Box>
    </Container>
  );
};

export default ChatApp;
