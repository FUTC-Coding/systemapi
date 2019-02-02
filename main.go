package main

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/disk"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"strconv"
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
			"model": getCpuModel(),
			"cores": getCpuCores(),
			"user": getCpu(0),
			"system": getCpu(1),
			"idle": getCpu(2),
		})
	})

	r.GET("/net", func(c *gin.Context){
		c.JSON(200, gin.H{
			"RxBytes": getNetwork(0),
			"TxBytes": getNetwork(1),
		})
	})

	r.GET("/uptime", func(c *gin.Context){
		c.JSON(200, gin.H{
			"uptime": getUptime(),
		})
	})

	r.GET("/disk", func(c *gin.Context){
		c.JSON(200, gin.H{
			"ReadsCompleted": getDisk(0),
			"WritesCompleted": getDisk(1),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
func getCpuModel() (string) {
	out, err := exec.Command("bash", "-c", "cat /proc/cpuinfo | grep \"model name\"").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	output := fmt.Sprintf("%s", out)
	output = strings.TrimSuffix(output, "\n")
	fmt.Println(output[13:])
	return output[13:]
}

func getCpuCores() (int) {
	out, err := exec.Command("bash", "-c", "cat /proc/cpuinfo | grep \"cpu cores\"").Output()
	if err != nil {
                fmt.Fprintf(os.Stderr, "%s\n", err)
                return 0
        }
	output := fmt.Sprintf("%s", out)
        output = strings.TrimSuffix(output, "\n")
	output = output[12:]
	fmt.Println(output)
	i, err := strconv.Atoi(output)
	if err != nil {
                fmt.Fprintf(os.Stderr, "%s\n", err)
                return 0
        }
	return i
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

func getNetwork(i int) (uint64){
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
	if i == 0 {
		return after[0].RxBytes - before[0].RxBytes
	}
	if i == 1 {
		return after[0].TxBytes - before[0].TxBytes
	}

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
	re := regexp.MustCompile("[0-9]+")
	nums := re.FindAllString(output, -1)
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 { //reverse the array
		nums[i], nums[j] = nums[j], nums[i]
	}
	s := ""
	for i,n := range nums {
		if i == 0 {
			s = n + "m"
		}
		if i == 1 {
			s = n + "h" + s
		}
		if i == 2 {
			s = n + "d" + s
		}
	}
	return s
}

func getDisk(i int) (uint64){
	before, err := disk.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := disk.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	if i == 0 {
		return after[0].ReadsCompleted - before[0].ReadsCompleted
	}
	if i == 1 {
		return after[0].WritesCompleted - before[0].WritesCompleted
	}
	return 0
}
