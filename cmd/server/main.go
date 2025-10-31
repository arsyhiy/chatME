package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/arsyhiy/sshme/internal/chat"
	"github.com/arsyhiy/sshme/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func handleConnection(conn net.Conn, c *chat.Chat, storage *database.Storage) {
	defer conn.Close()

	fmt.Fprintln(conn, "Введите имя:")

	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Проверяем есть ли пользователь в базе
	var id int64
	err := storage.DB.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, _ := storage.DB.Exec("INSERT INTO users(name) VALUES(?)", name)
		id, _ = res.LastInsertId()
	}

	colors := []string{"\033[91m", "\033[92m", "\033[93m", "\033[94m", "\033[95m", "\033[96m"}
	user := &chat.User{
		ID:    id,
		Name:  name,
		Color: colors[len(c.Users)%len(colors)],
		Out:   make(chan string, 10),
	}
	c.AddUser(user)

	fmt.Fprintf(conn, "Привет, %s! Введите сообщение:\n", user.Name)

	// Загружаем последние 20 сообщений из базы
	rows, err := storage.DB.Query("SELECT username, text FROM messages ORDER BY id DESC LIMIT 20")
	if err == nil {
		defer rows.Close()
		var history []string
		for rows.Next() {
			var uname, text string
			rows.Scan(&uname, &text)
			history = append([]string{fmt.Sprintf("[%s]: %s", uname, text)}, history...)
		}
		for _, msg := range history {
			fmt.Fprintln(conn, msg)
		}
	}

	go func() {
		for msg := range user.Out {
			fmt.Fprintln(conn, msg)
		}
	}()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(name, "disconnected")
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
			// Сохраняем общее сообщение через Storage
			_ = storage.AddMessage(user.Name, "", line)
			msg := chat.Message{
				Text:       fmt.Sprintf("%s[%s] %s: %s\033[0m", user.Color, now, user.Name, line),
				Recipients: nil,
			}
			c.Broadcast(msg)
		}
	}
}

func main() {
	storage, err := database.NewStorage("chat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	ch := chat.NewChat() // создаём объект чата

	// Загружаем историю из базы
	rows, err := storage.DB.Query("SELECT sender, text FROM messages ORDER BY id ASC")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var uname, text string
			rows.Scan(&uname, &text)
			msg := chat.Message{
				Text:       fmt.Sprintf("[%s]: %s", uname, text),
				Recipients: nil,
			}
			ch.History = append(ch.History, msg)
		}
	}

	// запускаем сервер
	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Сервер слушает :2222 ...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка подключения:", err)
			continue
		}
		go handleConnection(conn, ch, storage)
	}
}
