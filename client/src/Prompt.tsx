import React, { useState } from "react";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { CssBaseline } from "@mui/material";
import OptimizeButton from "./components/OptimizePreviewBtn";

const PromptInput: React.FC = () => {
  const [prompt, setPrompt] = useState("");

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setPrompt(event.target.value);
  };

  return (
    <>
      <CssBaseline />
      <Box sx={{ display: "flex", flexDirection: "column", height: "100vh" }}>
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            p: 2,
            boxShadow: 1,
            backgroundColor: "background.paper",
          }}
        >
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            人设 & 提示语
          </Typography>
          <OptimizeButton />
        </Box>
        <Box
          sx={{
            flexGrow: 1,
            overflow: "auto",
            px: 2,
            py: 1,
            backgroundColor: "background.default",
          }}
        >
          <TextField
            fullWidth
            multiline
            placeholder="在此输入您的提示..."
            variant="standard"
            value={prompt}
            onChange={handleInputChange}
            InputProps={{
              disableUnderline: true,
              style: { backgroundColor: "transparent" },
            }}
            sx={{
              width: "100%",
              height: "100%",
              ".MuiInputBase-inputMultiline": {
                height: "100%",
                overflow: "auto",
              },
            }}
          />
        </Box>
      </Box>
    </>
  );
};

export default PromptInput;
