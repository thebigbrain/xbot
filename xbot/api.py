# api.py (修改部分)
import json
import time
from flask import Blueprint, Response, request, jsonify, stream_with_context
from flask_cors import CORS
from xbot.message import get_new_messages
from xbot.models import db, Message
from langchain_community.llms import Ollama

api = Blueprint("api", __name__, url_prefix="/api")
llm = Ollama(model="codellama")


@api.route("/send", methods=["POST", "OPTIONS"])
def send_message():
    if request.method == "OPTIONS":
        # 为 Preflight 请求创建一个空响应
        headers = {
            "Content-Type": "text/event-stream",
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Headers": "*",
        }
        return Response(status=200, headers=headers)

    data = request.get_json()
    user = data["user"]
    content = data["content"]

    # 创建一个新消息对象
    new_message = Message(user=user, content=content)

    # 将新消息保存到数据库中
    db.session.add(new_message)
    db.session.commit()

    @stream_with_context
    def generate(last_checked: Message):
        # yield f"data: { json.dumps(last_checked.to_dict())}\n\n"
        # new_messages = get_new_messages(last_checked)
        # for chunk in new_messages:
        #     print(chunk)
        #     yield f"data: {json.dumps(chunk)}\n\n"
        yield "Hello "
        yield "jjee"
        yield "!"

    headers = {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        "Connection": "keep-alive",
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Headers": "*",
    }

    # 使用Flask的stream_with_context来确保请求上下文随生成器一起工作
    return Response(
        generate(new_message), content_type="text/event-stream", headers=headers
    )


@api.route("/history", methods=["GET"])
def get_history():
    # 从数据库中检索所有的消息
    all_messages = Message.query.order_by(Message.timestamp.asc()).all()

    # 将数据库模型转换成字典列表
    messages_dict = [message.to_dict() for message in all_messages]

    return jsonify(messages_dict), 200


# ...其它代码继续...
