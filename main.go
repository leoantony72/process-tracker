package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/r3labs/sse"
	"golang.org/x/sys/windows"
)

var process map[string]time.Time

type Message struct {
	Cmd   string `json:"cmd"`
	Pname string `json:"pid_name"`
}

const (
	serverURL = "http://192.168.18.11:900"
)

func main() {
	go func() {
		for {
			process = make(map[string]time.Time)
			getProcess()
			time.Sleep(time.Second * 30)
		}
	}()

	client := sse.NewClient(serverURL + "/events")
	client.Subscribe("message", func(msg *sse.Event) {
		serverCommand := Message{}
		err := json.Unmarshal(msg.Data, &serverCommand)
		if err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			return
		}

		switch serverCommand.Cmd {
		case "kill":
			killP(serverCommand.Pname)
		case "get":
			for k := range process {
				delete(process, k)
			}
			getProcess()
			sendData()
		}
	})
}

func killP(p string) {
	processName := p
	pids, err := findProcessIDByName(processName)
	if err != nil {
		fmt.Printf("Error finding process: %v\n", err)
		return
	}

	for _, pid := range pids {
		err = terminateProcess(pid)
		if err != nil {
			fmt.Printf("Error terminating process with PID %d: %v\n", pid, err)
		} else {
			fmt.Printf("Successfully terminated process with PID %d\n", pid)
		}
	}
}

func getProcess() {
	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		idx := strings.Index(line, ".exe")
		if idx != -1 {
			processName := line[:idx+len(".exe")]
			_, ok := process[processName]
			if !ok {
				process[processName] = time.Now()
			}
		}
	}
	sendData()
}

func sendData() {
	jsonData, err := json.Marshal(process)
	if err != nil {
		fmt.Println(err)
		return
	}

	retry(func() error {
		return postJSON(serverURL+"/process", jsonData)
	}, 30, time.Second*10)
}

func postJSON(url string, data []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if isNetworkError(err) {
			return err
		}
		return nil
	}
	defer resp.Body.Close()

	return nil
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(net.Error)
	return ok
}

func retry(operation func() error, attempts int, sleep time.Duration) {
	for i := 0; i < attempts; i++ {
		err := operation()
		if err == nil {
			return
		}

		fmt.Printf("Operation failed: %v. Retrying in %v...\n", err, sleep)
		time.Sleep(sleep)
	}
	fmt.Println("Operation failed after all attempts.")
}

func findProcessIDByName(processName string) ([]uint32, error) {
	var pids []uint32

	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	err = windows.Process32First(snapshot, &processEntry)
	for err == nil {
		name := windows.UTF16ToString(processEntry.ExeFile[:])
		if strings.EqualFold(name, processName) {
			pids = append(pids, processEntry.ProcessID)
		}
		err = windows.Process32Next(snapshot, &processEntry)
	}
	if len(pids) == 0 {
		return nil, fmt.Errorf("no process found with name: %s", processName)
	}
	return pids, nil
}

func terminateProcess(pid uint32) error {
	handle, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	err = windows.TerminateProcess(handle, 0)
	if err != nil {
		return err
	}
	return nil
}
