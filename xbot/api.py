# api.py (修改部分)
from flask import Blueprint, request, jsonify
from xbot.models import db, Message

api = Blueprint("api", __name__, url_prefix="/api")


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

    return jsonify({"status": "success"}), 200


@api.route("/history", methods=["GET"])
def get_history():
    # 从数据库中检索所有的消息
    all_messages = Message.query.order_by(Message.timestamp.asc()).all()

    # 将数据库模型转换成字典列表
    messages_dict = [message.to_dict() for message in all_messages]

    return jsonify(messages_dict), 200


# ...其它代码继续...
