package glassbox

import (
  "io"
)

func writeLength(w io.Writer, int16 length) error {
  return binary.Write(w, binary.BigEndian, length)
}
