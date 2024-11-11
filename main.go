package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	utils "file-streamer/utils"
	"fmt"
	"io"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) start() {
	ln := utils.Must(net.Listen("tcp", ":3000"))

	for {
		conn := utils.Must(ln.Accept())

		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		nBytes := utils.Must(io.CopyN(buf, conn, size))

		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over network \n", nBytes)
	}
}

// simple function to mock file streaming over network, almost like from client side
func mockSendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	nBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes over the network\n", nBytes)

	return nil
}

func main() {
	// wait 4 seconds and send huge file
	go func() {
		time.Sleep(4 * time.Second)
		mockSendFile(200000)
	}()

	server := &FileServer{}
	server.start()
}
