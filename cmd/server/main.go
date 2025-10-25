package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/arsyhiy/sshme/internal/chat"
)

func handleConnection(conn net.Conn, c *chat.Chat) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Отправка приглашения имени с переводом строки
	fmt.Fprintln(conn, "Введите имя:")

	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)

	// Создание пользователя
	colors := []string{"\033[91m", "\033[92m", "\033[93m", "\033[94m", "\033[95m", "\033[96m"}
	user := &chat.User{
		Name:  name,
		Color: colors[len(c.Users)%len(colors)],
		Out:   make(chan string, 10),
	}
	c.AddUser(user)

	// Отправляем приветствие
	fmt.Fprintf(conn, "Привет, %s! Введите сообщение:\n", user.Name)

	// Отправка сообщений пользователю через TCP
	go func() {
		for msg := range user.Out {
			fmt.Fprintln(conn, msg) // обязательно с переводом строки
		}
	}()

	// Чтение сообщений от клиента
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(name, "отключился")
			c.RemoveUser(user.Name)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = chat.ApplyEmojis(line)
		now := time.Now().Format("15:04:05")

		if strings.HasPrefix(line, "/msg ") {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) < 3 {
				fmt.Fprintln(conn, "Используй: /msg <имя> <текст>")
				continue
			}
			target := parts[1]
			text := parts[2]
			msg := chat.Message{
				Text:       fmt.Sprintf("%s[%s] [Приватно] %s → %s: %s\033[0m", user.Color, now, user.Name, target, text),
				Recipients: map[string]bool{user.Name: true, target: true},
			}
			c.Broadcast(msg)
		} else {
			msg := chat.Message{
				Text:       fmt.Sprintf("%s[%s] %s: %s\033[0m", user.Color, now, user.Name, line),
				Recipients: nil,
			}
			c.Broadcast(msg)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Сервер слушает :2222 ...")
	chat := chat.NewChat()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка подключения:", err)
			continue
		}
		go handleConnection(conn, chat)
	}
}

