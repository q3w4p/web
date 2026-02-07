package main

import (
        "fmt"
        "os"
        "os/exec"
        "strconv"
        "sync"
)

type ProcessManager struct {
        mu sync.Mutex
}

var pm = &ProcessManager{}

func (m *ProcessManager) StartBot(token string, discordId string) (int, error) {
        // Log starting message as requested
        fmt.Printf("user started pid: %d\n", 0) // Placeholder for real PID

        // Use PM2 to start main.go
        // We use discordId as the PM2 process name
        cmd := exec.Command("pm2", "start", "go run main.go", "--name", discordId, "--interpreter", "none")
        err := cmd.Start()
        if err != nil {
                return 0, err
        }
        
        // In a real scenario, we might want to wait or get the actual PID from PM2
        // For now, we return the cmd process PID or a simulated one
        return cmd.Process.Pid, nil
}

func (m *ProcessManager) KillBot(pid int) error {
        proc, err := os.FindProcess(pid)
        if err != nil {
                return err
        }
        return proc.Kill()
}

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Usage: host <token> or kill <pid>")
                return
        }

        cmd := os.Args[1]
        switch cmd {
        case "host":
                if len(os.Args) < 3 {
                        fmt.Println("Token required")
                        return
                }
                // In this environment, we execute from the root, so we need to ensure main.go is found in go-bot/
                cmd := exec.Command("pm2", "start", "go run main.go", "--name", os.Args[3], "--interpreter", "none", "--cwd", "./go-bot")
                err := cmd.Start()
                if err != nil {
                        fmt.Printf("Error: %v\n", err)
                        return
                }
                fmt.Printf("PID: %d\n", cmd.Process.Pid)
        case "kill":
                if len(os.Args) < 3 {
                        fmt.Println("PID required")
                        return
                }
                pid, _ := strconv.Atoi(os.Args[2])
                err := pm.KillBot(pid)
                if err != nil {
                        fmt.Printf("Error: %v\n", err)
                        return
                }
                fmt.Println("Killed")
        default:
                fmt.Println("Unknown command")
        }
}
