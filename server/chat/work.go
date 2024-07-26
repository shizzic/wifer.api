package chat

import "fmt"

// Юзер покидает соединение полностью
func quit(id int) {
	fmt.Print("\n", "Quit")
	clients.Delete(id)
}
