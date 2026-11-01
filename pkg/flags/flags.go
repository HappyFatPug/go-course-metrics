package flags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port uint16
}

func (n *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (n *NetAddress) Set(flagVal string) error {
	valSplit := strings.Split(flagVal, ":")
	if len(valSplit) != 2 {
		return errors.New("use value like `host:port`")
	}
	port, err := strconv.ParseUint(valSplit[1], 10, 16)
	if err != nil {
		return err
	}

	n.Host = valSplit[0]
	n.Port = uint16(port)
	return nil
}
