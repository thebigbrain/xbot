import React from "react";
import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";
import Layout from "./Layout";

const theme = createTheme();

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Layout />
    </ThemeProvider>
  );
};

export default App;
