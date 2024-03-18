import React, { useState } from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import TuneIcon from "@mui/icons-material/Tune";

const OptimizePreviewDialog: React.FC = () => {
  const [open, setOpen] = useState<boolean>(false);
  const [optimized, setOptimized] = useState<boolean>(false);

  const handleOptimize = () => {
    // 这里应该是您的优化逻辑
    console.log("执行优化算法");
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleAccept = () => {
    // 应用优化逻辑
    console.log("应用了优化");
    setOptimized(true);
    handleClose();
  };

  const handleReject = () => {
    // 放弃优化逻辑
    console.log("放弃了优化");
    handleClose();
  };

  return (
    <>
      <Button startIcon={<TuneIcon />} color="inherit" onClick={handleOptimize}>
        优化
      </Button>
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>优化预览</DialogTitle>
        <DialogContent>
          {/* 这里放置优化后的预览内容，可能是一些文本、图片或组件 */}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleReject} color="primary">
            放弃
          </Button>
          <Button onClick={handleAccept} color="primary" autoFocus>
            使用
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default OptimizePreviewDialog;
