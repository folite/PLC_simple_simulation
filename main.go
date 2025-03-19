package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var DM = make(map[string]int)

func Init() {
	DM["1001"] = 0
	DM["1002"] = 0
	DM["2001"] = 2
	DM["2002"] = 0
	DM["2003"] = 0
	DM["2004"] = 2
	DM["2005"] = 0
	DM["2006"] = 0
	DM["2007"] = 0
	DM["2008"] = 0
	DM["2009"] = 0
	DM["2010"] = 0
	return
	// for i := 1001; i <= 1002; i++ {
	// 	DM[strconv.Itoa(i)] = i
	// }
	// for i := 2001; i <= 2007; i++ {
	// 	DM[strconv.Itoa(i)] = i
	// }
}

func main() {
	Init()
	li, err := net.Listen("tcp", "127.0.0.1:8501")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Accepted a new TCP connection.")

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	fmt.Println(conn.RemoteAddr())
	defer conn.Close()

	bytes := make([]byte, 1024)
	for {
		n, err := conn.Read(bytes)
		if err != nil {
			log.Println(err)
			return
		}
		cmd := string(bytes[:n])
		log.Printf("Received : %q\n", cmd)
		res := []byte((ParseCMD(cmd)))
		if DM["2001"] == 1 {
			go Mock()
		}
		conn.Write(res)

	}
}
func ParseCMD(cmd string) string {
	reg_init, _ := regexp.Compile(`^init`)
	if reg_init.MatchString(cmd) {
		Init()
		return "OK\r\n"
	}
	reg_rd, _ := regexp.Compile(`^RD `)
	if reg_rd.MatchString(cmd) {
		split := strings.Split(cmd, " ")
		res, err := RD(split[1][2:6])
		if err == nil {
			return strconv.Itoa(res) + "\r\n"
		}
	}
	reg_rds1001, _ := regexp.Compile(`^RDS DM1001.H 2`)
	reg_rds2001, _ := regexp.Compile(`^RDS DM2001.H 10`)
	if reg_rds1001.MatchString(cmd) || reg_rds2001.MatchString(cmd) {
		split := strings.Split(cmd, " ")
		id := split[1][2:6]
		step, _ := strconv.Atoi(strings.Split(split[2], "\r")[0])
		s := ""
		for i := 0; i < step; i++ {
			nid, _ := strconv.Atoi(id)
			s += " " + fmt.Sprintf("%04X", DM[strconv.Itoa(nid+i)])
		}
		return strings.TrimSpace(s) + "\r\n"
	}
	reg_wr, _ := regexp.Compile(`^WR `)
	if reg_wr.MatchString(cmd) {
		split := strings.Split(cmd, " ")
		a := strings.Split(split[2], "\r")
		data, _ := strconv.Atoi(a[0])
		_, err := WR(split[1][2:6], data)
		if err == nil {
			return "OK\r\n"
		}
	}
	return "E1\r\n"
}

func DataSpaceExist(index string) bool {
	_, ok := DM[index]
	if ok {
		return true
	} else {
		return false
	}
}

func RD(index string) (int, error) {
	if !DataSpaceExist(index) {
		return 0, errors.New(fmt.Sprintf("DM%s is nothing.\n", index))
	}
	return DM[index], nil
}

func WR(index string, data int) (int, error) {
	DM[index] = data
	return data, nil
}

func Mock() {
	DM["2007"] += rand.Intn(100)
	DM["2010"] += rand.Intn(100)
	if DM["2007"] > 5000 {
		DM["2001"] = 2
	}
	if DM["2001"] != 1 {
		return
	}
}
