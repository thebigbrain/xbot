import logging
from websockets import WebSocketServer, serve
from chat_records import ChatRecordsManager

logger = logging.getLogger("websocket")
logging.basicConfig(level=logging.INFO)

connected = set()
chat_manager = ChatRecordsManager()


async def chat_handler(websocket, path):
    connected.add(websocket)
    logger.info(f"WebSocket connection established with {websocket.remote_address}")

    try:
        # 保持连接，持续监听消息
        async for message in websocket:
            # 在这里处理接收到的消息
            # 例如：await process_message(message)
            chat_manager.add_chat_record(websocket.remote_address[0], message)
            logger.info("got message %s", message)
    except Exception as e:
        logger.error(f"Error in chat_handler: {e}")
    finally:
        # 无论如何，连接最终关闭时进行清理
        connected.remove(websocket)
        logger.info(f"WebSocket connection closed with {websocket.remote_address}")


server = None


async def websocket_server_runner(event):
    global server
    server = await serve(chat_handler, "localhost", 6789)
    logger.info("WebSocket server started on ws://localhost:6789")
    await event.wait()
    stop_websocket_server()


async def stop_websocket_server():
    if server:
        server.close()
        await server.wait_closed()
    for websocket in connected:
        await websocket.close(reason="Server shutdown")
    chat_manager.close()
    logger.info("WebSocket server stopped gracefully")
