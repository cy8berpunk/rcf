/*
Package rcfutil implements basic parsing and de-encoding for rcf_node & rcf_node_client
*/
package rcfUtil

import (
	"net"
	"io"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"math/rand"
)

// serializable protocol msg
type Smsg struct {
	Type string
	Name string
	Id int
	Operation string
	Payload []byte
	MultiplePayload [][]byte
}

func EncodeMsg(msg *Smsg) ([]byte, error) {
	serializedMsg, err := json.Marshal(&msg)
	if err != nil {
		return []byte{}, err
	}
	return serializedMsg, nil
}


func DecodeMsg(msg *Smsg, data []byte) error {
	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}


// CompareSlice compares two slices for equality
// slices must be of same length
// returns false if slices are not equal
func CompareSlice(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

// TopicsContainTopic checks if the topics map contains a certain topic(name)
// returns false if topic(name) is not included in the list
func TopicsContainTopic(imap map[string][][]byte, key string) bool {
	if _, ok := imap[key]; ok {
		return true
	}
	return false
}

// GenRandomIntID generates random id
// returns generated random id
func GenRandomIntID() int {
	pullReqID := rand.Intn(1000000000)
	if pullReqID == 0 || pullReqID == 2 {
		pullReqID = rand.Intn(100000000)
	}
	return pullReqID
}

func WriteFrame(writer *bufio.Writer, data []byte) error {
	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(dataLen, uint64(len(data)))
	_, err := writer.Write(dataLen)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func readFixedSize(r io.Reader, len int) (b []byte, err error) {
    b = make([]byte, len)
    _, err = io.ReadFull(r, b)
    if err != nil {
        return
    }
    return
}

func ReadFrame(conn net.Conn) ([]byte, error) {
	lenData, err := readFixedSize(conn, 8)
	if err != nil {
		return []byte{}, err
	}
	dataLen := binary.LittleEndian.Uint64(lenData)

	dataBuffer, err := readFixedSize(conn, int(dataLen))
	if err != nil {
		return []byte{}, err
	}
	return dataBuffer, nil
}
