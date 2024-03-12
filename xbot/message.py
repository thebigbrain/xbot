from langchain_community.llms import Ollama

from xbot.models import Message, db

llm = Ollama(model="codellama")


def get_new_messages(last_checked: Message):
    msg = Message(user="bot", content="")
    for chunk in llm.stream(last_checked.content):
        msg.content += str(chunk)
        yield chunk

    db.session.add(msg)
    db.session.commit()
