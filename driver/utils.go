package driver

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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
			fmt.Println("Download:", downloadSpeed)
			fmt.Println("Units:", units)
			fmt.Println("---------------")
			if units != "Mbits/sec" {
				return 0, err
			}
			return downloadSpeed, nil
		}
		lineNum = lineNum + 1
	}
}
