import React from "react";
import Chat from "./Chat";
import { Box, Paper } from "@mui/material";
import PromptInput from "./Prompt";

const Layout: React.FC = () => {
  return (
    <Box display="flex" flexDirection="row" height="100vh" overflow="hidden">
      {/* 横向排版 PromptInput，占据可用空间的部分 */}
      <Box width="40%" height="100%">
        <Paper elevation={3} style={{ height: "100%", overflowY: "auto" }}>
          <PromptInput />
        </Paper>
      </Box>
      {/* 横向排版 Chat，占据剩余的空间 */}
      <Box width="60%" height="100%">
        <Paper elevation={3} style={{ height: "100%", overflowY: "auto" }}>
          <Chat />
        </Paper>
      </Box>
    </Box>
  );
};

export default Layout;
