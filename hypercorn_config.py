# hypercorn_config.py
bind = "0.0.0.0:5000"  # 监听所有接口的5000端口
workers = 3
keep_alive_timeout = 5
keyfile = "keyfile.pem"
certfile = "certfile.crt"
