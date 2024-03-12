# database.py
from sqlalchemy import create_engine, Column, Integer, Text, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from datetime import datetime

engine = create_engine('sqlite:///chat.db')
Session = sessionmaker(bind=engine)
Base = declarative_base()

class ChatRecord(Base):
    __tablename__ = 'chats'
    id = Column(Integer, primary_key=True)
    username = Column(Text)
    message = Column(Text)
    timestamp = Column(DateTime, default=datetime.utcnow)

    def as_dict(self):
        return {c.name: getattr(self, c.name) for c in self.__table__.columns}

Base.metadata.create_all(engine)

def add_chat_record(username, message):
    session = Session()
    chat_record = ChatRecord(username=username, message=message)
    session.add(chat_record)
    session.commit()
    session.close()

def get_all_chats():
    session = Session()
    records = session.query(ChatRecord).all()
    session.close()
    return [record.as_dict() for record in records]

# 注意，实际生产环境中，需要增加错误处理和数据库连接管理。