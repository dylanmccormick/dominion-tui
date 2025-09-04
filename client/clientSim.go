package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/dylanmccormick/dominion-tui/internal/utils"
	"github.com/dylanmccormick/dominion-tui/server"
)

func SimulateClient() {
	conn, err := net.Dial("tcp", "localhost:42069")
	if err != nil {
		log.Fatalf("Unable to connect to server")
	}
	defer conn.Close()

	conn.Write([]byte("Test"))
	time.Sleep(5 * time.Second)
	conn.Write([]byte("lobby"))
	time.Sleep(5 * time.Second)

	fmt.Println("Simulating client")
	go PrintResponses(&conn)

	for range 100 {
		msg, err := json.Marshal(CreateRandomChat())
		if err != nil {
			continue
		}
		conn.Write(append(msg, []byte("\r\n")...))
	}

	time.Sleep(100 * time.Second)
}

func PrintResponses(conn *net.Conn) {
	count := 0
	for {
		scanner := bufio.NewReader(*conn)
		buffer := make([]byte, 4096)
		_, err := scanner.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of connection closed gracefully")
				break
			} else {
				fmt.Printf("Unknown error from TCP request: %s", err)
				break
			}
		}
		data := bytes.SplitSeq(buffer, []byte("\r\n"))
		for dat := range data {
			clean := utils.ClearZeros(dat)
			if len(clean) == 0 {
				continue
			}
			fmt.Printf("\t\t%d %s\n", count, string(clean))
			count += 1
		}
	}
}


func CreateRandomChat() server.Message {
		answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}

	names := []string {
		"Dylan",
		"Oliver",
		"ClapTrap",
		"Meeble",
	}

	return server.Message{
		Requester: names[rand.Intn(len(names))],
		Typ: "chat",
		Body: map[string]any{
			"message": fmt.Sprintf("Magic 8-Ball says: %s", answers[rand.Intn(len(answers))]),
		},

	}

}

