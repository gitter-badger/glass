package glass

import (
    "io"
    "encoding/binary"
    "bytes"
)

func write_uint16(w io.Writer, n uint16) error {
  return binary.Write(w, binary.BigEndian, n)
}

func write_uint32(w io.Writer, n uint32) error {
  return binary.Write(w, binary.BigEndian, n)
}

func read_uint16(r string) (n uint16) {
    buf := bytes.NewReader([]byte{r[0], r[1]})
    binary.Read(buf, binary.BigEndian, &n)
    return
}
func read_uint32(r string) (n uint32) {
    buf := bytes.NewReader([]byte{r[0], r[1], r[2], r[3]})
    binary.Read(buf, binary.BigEndian, &n)
    return
}
