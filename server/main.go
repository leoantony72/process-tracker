package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Cmd   string `json:"cmd"`
	Pname string `json:"pid_name"`
}

// type Process struct {
// 	Pname    string    `json:"pid_name"`
// 	TimeUsed time.Time `json:"time"`
// }
type Process struct {
	Data map[string]time.Time `json:"data"`
}

var process_data Process = Process{}
var broadcast chan Message = make(chan Message, 4)

func main() {
	// broadcast = make(chan Message,4)
	r := gin.Default()
	r.LoadHTMLFiles("static/index.html")
	r.Static("/css", "static/css")
	r.Static("/js", "static/js")

	r.GET("/", Home)
	r.GET("/events", sse)
	r.GET("/getdata", getData)
	r.POST("/send-command", parse_command)
	r.POST("/process", process)

	r.Run("0.0.0.0:900")
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func getData(c *gin.Context) {
	type Data struct {
		Process []string `json:"data"`
	}
	data := Data{}

	for key := range process_data.Data {
		cleanedKey := strings.Trim(key, "\"")
		data.Process = append(data.Process, cleanedKey)
	}
	// b, _ := json.Marshal(data)
	fmt.Println(data.Process[0], data.Process[1])
	c.JSON(200, gin.H{"data": data})

}

func process(c *gin.Context) {

	decoder := json.NewDecoder(c.Request.Body)
	for k := range process_data.Data {
		delete(process_data.Data, k)
	}
	
	// &process_data.Data = nil
	err := decoder.Decode(&process_data.Data)
	if err != nil {
		c.JSON(400, gin.H{"message": "something went wrong"})
		return
	}
	c.JSON(200, gin.H{"message": "done"})
}

func parse_command(c *gin.Context) {
	command := c.Query("cmd")
	process_name := c.Query("pid_name")
	if command == "" || process_name == "" {

		c.JSON(400, gin.H{"message": "provide command and pid"})
		return
	}
	msg := Message{Cmd: command, Pname: process_name}
	broadcast <- msg
	c.JSON(200, gin.H{"message": "command given"})
}

func sse(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		select {
		// case <-clientGone:
		// 	return false
		case message := <-broadcast:
			c.SSEvent("message", message)
			return true
		}
	})
}
