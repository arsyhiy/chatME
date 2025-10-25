package chat

import (
    "sync"
)

type User struct {
    Name  string
    Color string
    Out   chan string
}

type Chat struct {
    Users   map[string]*User
    History []Message
    mu      sync.Mutex
}

type Message struct {
    Text       string
    Recipients map[string]bool // nil = всем
}

func NewChat() *Chat {
    return &Chat{
        Users: make(map[string]*User),
    }
}

func (c *Chat) AddUser(u *User) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Users[u.Name] = u
    // Отправить историю
    for _, msg := range c.History {
        if msg.Recipients == nil || msg.Recipients[u.Name] {
            u.Out <- msg.Text
        }
    }
}

func (c *Chat) RemoveUser(name string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.Users, name)
}

func (c *Chat) Broadcast(msg Message) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.History = append(c.History, msg)
    for name, u := range c.Users {
        if msg.Recipients == nil || msg.Recipients[name] {
            u.Out <- msg.Text
        }
    }
}

