Простой SSH-чат с:

- Общими и приватными сообщениями
- Цветами и метками времени
- Смайлами (:smile:, :heart:, :thumbs:, :wink:)




## compiling
`` 
go build -o bin/server ./cmd/server
go build -o bin/client ./cmd/client
``


## run 

- server 
``
./bin/server
``

- client

``
./bin/client 127.0.0.1 2222
``

