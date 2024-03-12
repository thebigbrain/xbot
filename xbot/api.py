# api.py (修改部分)
import json
import time
from flask import Blueprint, Response, request, jsonify, stream_with_context
from xbot.message import get_new_messages
from xbot.models import db, Message
from langchain_community.llms import Ollama

api = Blueprint("api", __name__, url_prefix="/api")
llm = Ollama(model="codellama")


@api.route("/send", methods=["POST"])
def send_message():
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
        # 检索最新的消息，这里我们仅获取时间戳大于last_checked的消息
        # 假设我们有一种机制来存储last_checked，例如通过客户端提供，
        # 或者使用数据库中的某个值等
        yield f"data: { json.dumps(last_checked.to_dict())}\n\n"

        # 用实际的方法获取新消息，代码将根据你的应用逻辑有所不同
        new_messages = get_new_messages(last_checked)

        for chunk in new_messages:
            # 将消息对象转换为字典，准备作为JSON发送
            print(chunk)
            yield f"data: {json.dumps(chunk)}\n\n"

    # 使用Flask的stream_with_context来确保请求上下文随生成器一起工作
    return Response(generate(new_message), minetype="text/event-stream")


@api.route("/history", methods=["GET"])
def get_history():
    # 从数据库中检索所有的消息
    all_messages = Message.query.order_by(Message.timestamp.asc()).all()

    # 将数据库模型转换成字典列表
    messages_dict = [message.to_dict() for message in all_messages]

    return jsonify(messages_dict), 200


# ...其它代码继续...
