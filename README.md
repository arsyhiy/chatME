да первые три коммита я генерил chatgpt


# SSH Chat на Go

Простой SSH-чат с:

- Общими и приватными сообщениями
- Цветами и метками времени
- Смайлами (:smile:, :heart:, :thumbs:, :wink:)

## Структура проекта

ssh-chat-go/
├── cmd/server/ # сервер
├── cmd/client/ # клиент
├── internal/chat/ # логика чата
├── go.mod
└── README.md


## Установка

```bash
git clone <repo>
cd ssh-chat-go
go build ./cmd/server


Использование

Настройте ForceCommand в sshd_config для пользователя:

Match User chatuser
    ForceCommand /home/chatuser/ssh-chat-go/cmd/server/main


Подключение клиента:

ssh chatuser@IP_сервера

Введите имя пользователя.

Общие сообщения — просто пишем текст.

Приватные сообщения:

/msg Alice Привет, это только для тебя!

:smile: :heart: :thumbs: :wink:


---

Если хочешь, я могу **добавить клиентскую часть на Go**, чтобы было полноценное подключение через SSH или TCP-туннель с интерактивным вводом и выводом, как в оригинальном Python-чате.  

Хочешь, чтобы я сделал клиентскую часть?


Пример использования

go build ./cmd/client
./client chatuser@192.168.1.10 22 mypassword


После запуска клиент подключается к серверу и начинает переписку.

Общие сообщения: просто пишем текст.

Приватные: /msg Alice Привет!

Смайлы: :smile: :heart: :thumbs:
