import React from "react";
import Chat from "./Chat";
import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";

const theme = createTheme();

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Chat />
    </ThemeProvider>
  );
};

export default App;
