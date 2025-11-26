package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const FileProtocolID = "/anonbox/file/1.0.0"

// SendFile sends a file to a peer
func (p *P2PManager) SendFile(peerID peer.ID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	s, err := p.Host.NewStream(p.Ctx, peerID, FileProtocolID)
	if err != nil {
		return err
	}
	defer s.Close()

	// Send metadata (filename length, filename, file size)
	filename := filepath.Base(filePath)
	filenameLen := int32(len(filename))
	fileSize := fileInfo.Size()

	// Write filename length
	if err := binary.Write(s, binary.LittleEndian, filenameLen); err != nil {
		return err
	}
	// Write filename
	if _, err := s.Write([]byte(filename)); err != nil {
		return err
	}
	// Write file size
	if err := binary.Write(s, binary.LittleEndian, fileSize); err != nil {
		return err
	}

	// Send file content in chunks
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if _, err := s.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

// HandleFileStream handles incoming file streams
func (p *P2PManager) HandleFileStream(s network.Stream) {
	defer s.Close()

	// Read metadata
	var filenameLen int32
	if err := binary.Read(s, binary.LittleEndian, &filenameLen); err != nil {
		fmt.Printf("Error reading filename length: %v\n", err)
		return
	}

	filenameBuf := make([]byte, filenameLen)
	if _, err := io.ReadFull(s, filenameBuf); err != nil {
		fmt.Printf("Error reading filename: %v\n", err)
		return
	}
	filename := string(filenameBuf)

	var fileSize int64
	if err := binary.Read(s, binary.LittleEndian, &fileSize); err != nil {
		fmt.Printf("Error reading file size: %v\n", err)
		return
	}

	fmt.Printf("Receiving file: %s (%d bytes)\n", filename, fileSize)

	// Create file (in current directory for now, or a 'downloads' folder)
	// TODO: Make download path configurable
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer outFile.Close()

	// Copy stream to file
	copied, err := io.CopyN(outFile, s, fileSize)
	if err != nil {
		fmt.Printf("Error writing file content: %v\n", err)
		return
	}

	fmt.Printf("Received file: %s (%d bytes)\n", filename, copied)
	p.MsgChan <- Message{
		Sender:  s.Conn().RemotePeer().String(),
		Content: fmt.Sprintf("Received file: %s", filename),
		IsFile:  true,
	}
}
