```sh
openssl genpkey -algorithm RSA -out keyfile.pem
openssl req -new -key keyfile.pem -out csr.pem
openssl x509 -req -days 365 -in csr.pem -signkey keyfile.pem -out certfile.crt
```
