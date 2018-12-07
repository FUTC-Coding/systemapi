package main

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/mem", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"total":  getMemory(0),
			"used":   getMemory(1),
			"cached": getMemory(2),
			"free":   getMemory(3),
		})
	})

	r.GET("/cpu", func(c *gin.Context){
		c.JSON(200, gin.H{
			"user": getCpu(0),
			"system": getCpu(1),
			"idle": getCpu(2),
		})
	})

	r.GET("/net", func(c *gin.Context){
		c.JSON(200, gin.H{
			"network": getNetwork(1),
		})
	})

	r.GET("/uptime", func(c *gin.Context){
		c.JSON(200, gin.H{
			"uptime": getUptime(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func getMemory(i int) (uint64){
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	if i == 0 {
		return memory.Total/1024/1024
	} else if i == 1 {
		return memory.Used/1024/1024
	} else if i == 2 {
		return memory.Cached/1024/1024
	} else if i == 3 {
		return memory.Free/1024/1024
	}
	return 0
}

func getCpu(i int) (float64){
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	total := float64(after.Total - before.Total)
	if i == 0 {
		return float64(after.User-before.User)/total*100
	} else if i == 1 {
		return float64(after.System-before.System)/total*100
	} else if i == 2 {
		return float64(after.Idle-before.Idle)/total*100
	}
	return 0
}

func getNetwork(i int) (float64){
	before, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	fmt.Println(before)
	fmt.Println(after)
	return 0
}

func getUptime() (string){
	out, err := exec.Command("uptime", "-p").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	output := fmt.Sprintf("%s", out)
	output = strings.TrimSuffix(output, "\n")
	return output
}