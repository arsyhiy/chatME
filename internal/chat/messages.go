package chat

import "strings"

// Смайлы
var Emojis = map[string]string{
    ":smile:": "😄",
    ":heart:": "❤️",
    ":thumbs:": "👍",
    ":wink:": "😉",
}

func ApplyEmojis(text string) string {
    for k, v := range Emojis {
        text = strings.ReplaceAll(text, k, v)
    }
    return text
}

