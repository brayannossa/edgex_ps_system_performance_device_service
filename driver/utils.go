package driver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CheckTemperature() (int64, string, error) {
	out, err := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp").Output()
	if err != nil {
		return 0, "", err
	}
	outInt, err := (strconv.ParseInt(string(out[:len(out)-1]), 10, 64))
	temp := int64(outInt / 1000)
	if err != nil {
		return 0, "", err
	}
	if temp > 60 {
		return temp, "Hi", err
	} else if temp < 58 {
		return temp, "Low", err
	} else {
		return 0, "", err
	}

}

func CheckStorage() (int64, error) {
	args := "-h"
	cmd := exec.Command("df", strings.Split(args, " ")...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	buf := bufio.NewReader(stdout)
	lineNum := 0
	for {
		if lineNum > 10 {
			return 0, nil
		}
		line, _, _ := buf.ReadLine()

		vector := strings.Fields(string(line))
		// vector[0] = s.file
		// vector[1] = overall size
		// vector[2] = used size
		// vector[3] = available size
		// vector[4] = use %
		// vector[5] = mounted on
		if vector[5] == "/" {
			usedString := strings.Replace(vector[4], "%", "", -1)
			usedInt, _ := strconv.ParseInt(usedString, 10, 64)
			fmt.Println("---------------")
			fmt.Println("Storage")
			fmt.Println("used:", vector[2])
			fmt.Println("available", vector[3])
			fmt.Println("%used", vector[4])
			fmt.Println("---------------")
			return usedInt, nil
		}

		lineNum += 1
	}
}

func CheckRAM() (float64, error) {
	cmd := exec.Command("free")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	buf := bufio.NewReader(stdout)
	lineNum := 0
	for {
		if lineNum > 5 {
			return 0, nil
		}
		line, _, _ := buf.ReadLine()
		if strings.Contains(string(line), "Mem") {
			vector := strings.Fields(string(line))
			used, err := strconv.ParseFloat(vector[2], 64)
			if err != nil {
				return 0, err
			}
			total, err := strconv.ParseFloat(vector[1], 64)
			if err != nil {
				return 0, err
			}
			ans := used / total * 100
			fmt.Println("---------------")
			fmt.Println("RAM")
			fmt.Println("used:", used)
			fmt.Println("total", total)
			fmt.Println("%used", ans)
			fmt.Println("---------------")
			return ans, nil

		}
		lineNum += 1
	}
}

func InternetSpeed() (float64, error) {
	args := "-c 54.227.54.139 -p 5002"
	cmd := exec.Command("iperf", strings.Split(args, " ")...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	buf := bufio.NewReader(stdout)
	lineNum := 0
	lineBandwidth := 10
	for {
		if lineNum > 10 {
			return 0, nil
		}

		line, _, _ := buf.ReadLine()
		if strings.Contains(string(line), "Bandwidth") {
			lineBandwidth = lineNum
		}
		if lineNum > lineBandwidth {
			vector := strings.Fields(string(line))
			downloadSpeed, err := strconv.ParseFloat(vector[6], 64)
			if err != nil {
				return 0, err
			}
			units := vector[7]

			fmt.Println("---------------")
			fmt.Println("Internet Speed")
			fmt.Println("download:", downloadSpeed)
			fmt.Println("units:", units)
			fmt.Println("---------------")
			if units != "Mbits/sec" {
				return 0, err
			}
			return downloadSpeed, nil
		}
		lineNum = lineNum + 1
	}
}

func RestartContainers() {
	fmt.Println("restarting device (containers)...")
	balenaAddress := os.Getenv("BALENA_SUPERVISOR_ADDRESS")
	balenaKey := os.Getenv("BALENA_SUPERVISOR_API_KEY")
	balenaAppID := os.Getenv("BALENA_APP_ID")
	appID := map[string]string{"appId": balenaAppID}
	jsonData, err := json.Marshal(appID)
	if err != nil {
		fmt.Println("error marshaling JSON data:", err)
	}
	url := balenaAddress + "/v1/restart?apikey=" + balenaKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error restarting device (ram):", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error restarting device (ram): status code", resp.StatusCode)
	}
}
