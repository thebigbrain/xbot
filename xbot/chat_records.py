# chat_records.py
from database import Session, ChatRecord

class ChatRecordsManager:
    def __init__(self):
        self.session = Session()
    
    def add_chat_record(self, username, message):
        new_record = ChatRecord(username=username, message=message)
        self.session.add(new_record)
        self.session.commit()
    
    def get_chat_history(self):
        records = self.session.query(ChatRecord).all()
        return [record.as_dict() for record in records]
    
    def close(self):
        self.session.close()