package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: client <host> <порт>")
		return
	}
	host := os.Args[1]
	port := os.Args[2]

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Получаем приглашение имени
	greeting, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при подключении:", err)
		return
	}
	fmt.Print(greeting)

	// Вводим имя и отправляем на сервер
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		name := scanner.Text()
		conn.Write([]byte(name + "\n"))
	}

	// Чтение сообщений от сервера
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nСоединение закрыто сервером")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	// Отправка сообщений на сервер
	for scanner.Scan() {
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
	}
}

