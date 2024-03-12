# app.py
from flask import Flask
from xbot.models import db
from xbot.api import api  # 假设你的Blueprint仍然定义在api.py文件中

app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///chat.db"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

# 初始化数据库
db.init_app(app)

# 注册Blueprint
app.register_blueprint(api)

# 在第一次运行时创建所有数据库表
with app.app_context():
    db.create_all()

if __name__ == "__main__":
    app.run()
