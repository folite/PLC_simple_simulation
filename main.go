package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var DM = make(map[string]int)

func Init() {
	for i := 1001; i <= 1002; i++ {
		DM[strconv.Itoa(i)] = i
	}
	for i := 2001; i <= 2007; i++ {
		DM[strconv.Itoa(i)] = i
	}
}

func main() {
	Init()
	li, err := net.Listen("tcp", ":8501")
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
		conn.Write(res)

	}
}
func ParseCMD(cmd string) string {
	reg_rd, _ := regexp.Compile(`^RD `)
	if reg_rd.MatchString(cmd) {
		a := strings.Split(cmd, " ")
		log.Printf("%v+", a)
	}
	reg_rds1001, _ := regexp.Compile(`^RDS DM1001.H 2`)
	reg_rds2001, _ := regexp.Compile(`^RDS DM2001.H 7`)
	if reg_rds1001.MatchString(cmd) || reg_rds2001.MatchString(cmd) {
		split := strings.Split(cmd, " ")
		id := split[1][2:6]
		step, _ := strconv.Atoi(string(split[2][0]))
		s := ""
		for i := 0; i < step; i++ {
			nid, _ := strconv.Atoi(id)
			s += " " + fmt.Sprintf("%04X", DM[strconv.Itoa(nid+i)])
		}
		return strings.TrimSpace(s) + "\r\n"
	}
	reg_wr, _ := regexp.Compile(`^WR `)
	if reg_wr.MatchString(cmd) {
		strings.Split(cmd, " ")
	}

	// if err == nil {
	// log.Print(reg)
	// }

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
