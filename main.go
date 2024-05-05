package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var process map[string]time.Time

func main() {
	process = make(map[string]time.Time)

	for {
		time.Sleep(time.Second * 5)
		getProcess()
	}

}

func getProcess() {
	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		// return
	}

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		idx := strings.Index(line, ".exe")
		if idx != -1 {
			processName := line[:idx+len(".exe")]
			_, ok := process[processName]
			if !ok {
				process[processName] = time.Now()
				v := process[processName]
				t := time.Since(v)
				fmt.Println(processName, t)
			}
			v := process[processName]
			t := time.Since(v)
			fmt.Println(processName, t)

			//send the map to the server
			// sendData()
		}
	}
}

// func sendData() {
// 	_, err := json.Marshal(process)
// 	if err != nil {
// 		fmt.Println(err)
// 		// return
// 	}

// 	// responseBody := bytes.NewBuffer(jsonData)
// 	// url := "www.example.com"
// 	// http.Post(url, "application/json", responseBody)

// }
