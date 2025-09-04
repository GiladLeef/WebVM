package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type VMController struct{}

type startvmResponse struct {
	VMID  string `json:"VMID"`
	ContainerID string `json:"containerId,omitempty"`
	Running     bool   `json:"running"`
	Message     string `json:"message"`
	Logs        string `json:"logs,omitempty"`
}

// Startvm runs a container based on platform-vm:latest exposing VNC on a random host port.
func (vc *VMController) Startvm(c *gin.Context) {
	VMID := uuid.New().String()
	containerName := fmt.Sprintf("vm-%s", VMID)

	// Run container without publishing ports; we'll connect over bridge network IP
	cmd := exec.Command("docker", "run",
		"-d",
		"--name", containerName,
		"--network", "platform_default",
		"platform-vm:latest",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("docker run failed: %v, output: %s", err, string(out))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to start vm: %v", err), "output": string(out)})
		return
	}
	containerID := strings.TrimSpace(string(out))

	// Give it a moment to start
	time.Sleep(300 * time.Millisecond)

	// Verify running state
	runCmd := exec.Command("bash", "-lc", fmt.Sprintf("docker inspect -f '{{.State.Running}}' %s", containerName))
	runOut, runErr := runCmd.CombinedOutput()
	running := strings.TrimSpace(string(runOut)) == "true" && runErr == nil

	var logsText string
	if !running {
		logsCmd := exec.Command("docker", "logs", containerName)
		lOut, _ := logsCmd.CombinedOutput()
		logsText = string(lOut)
		log.Printf("vm container %s not running. logs:\n%s", containerName, logsText)
		// cleanup failed container
		_ = exec.Command("docker", "rm", "-f", containerName).Run()
	}

	c.JSON(http.StatusOK, startvmResponse{VMID: VMID, ContainerID: containerID, Running: running, Message: "vm started", Logs: logsText})
}

// Stopvm stops and removes the container for the given VMID if present.
func (vc *VMController) Stopvm(c *gin.Context) {
	VMID := c.Param("id")
	if VMID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing vm id"})
		return
	}
	containerName := fmt.Sprintf("vm-%s", VMID)

	_ = exec.Command("docker", "stop", containerName).Run()
	_ = exec.Command("docker", "rm", "-f", containerName).Run()

	c.JSON(http.StatusOK, gin.H{"message": "vm stopped"})
}

// upgrader for WS without origin checks (public demo)
var vmUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32,
	WriteBufferSize: 1024 * 32,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Streamvm proxies raw TCP from the vm's VNC server to the browser over WebSocket.
func (vc *VMController) Streamvm(c *gin.Context) {
	VMID := c.Param("id")
	if VMID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing vm id"})
		return
	}
	containerName := fmt.Sprintf("vm-%s", VMID)

	// Discover the container IP on bridge network
	ipCmd := exec.Command("bash", "-lc", fmt.Sprintf("docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' %s", containerName))
	ipOut, err := ipCmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to resolve vm IP: %v", err), "output": string(ipOut)})
		return
	}
	ip := strings.TrimSpace(string(ipOut))
	if ip == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "vm IP not available"})
		return
	}

	// Wait for VNC to be ready (up to ~20s)
	var tcpConn net.Conn
	deadline := time.Now().Add(20 * time.Second)
	for {
		conn, dErr := net.DialTimeout("tcp", fmt.Sprintf("%s:5901", ip), 1*time.Second)
		if dErr == nil {
			tcpConn = conn
			break
		}
		if time.Now().After(deadline) {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("VNC not ready: %v", dErr)})
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	defer tcpConn.Close()

	ws, err := vmUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	done := make(chan struct{}, 2)

	// WS -> TCP
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			mt, data, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if mt == websocket.BinaryMessage {
				if _, err := tcpConn.Write(data); err != nil {
					return
				}
			}
		}
	}()

	// TCP -> WS
	go func() {
		defer func() { done <- struct{}{} }()
		buf := make([]byte, 8192)
		for {
			n, err := tcpConn.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}
	}()

	<-done

	// On disconnect, clean up the vm container
	go func() {
		if err := exec.Command("docker", "rm", "-f", containerName).Run(); err != nil {
			log.Printf("failed to remove vm container %s: %v", containerName, err)
		}
	}()
}

// For debugging: simple info endpoint
func (vc *VMController) Info(c *gin.Context) {
	VMID := c.Param("id")
	containerName := fmt.Sprintf("vm-%s", VMID)
	inspect := exec.Command("bash", "-lc", fmt.Sprintf("docker inspect %s", containerName))
	out, _ := inspect.CombinedOutput()
	var parsed []map[string]interface{}
	_ = json.Unmarshal(out, &parsed)
	c.JSON(http.StatusOK, gin.H{"VMID": VMID, "container": parsed})
}


