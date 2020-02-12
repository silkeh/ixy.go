package pci

import (
	"github.com/ixy-languages/ixy.go/util"
	"os"
)

const (
	offsetVendor     = 0
	offsetDeviceID   = 2
	offsetCommand    = 4
	offsetStatus     = 6
	offsetRevisionID = 8
	offsetClassCode  = 9
)

type Config struct {
	fd *os.File
}

func NewPCIConfig(path string) (*Config, error) {
	fd, err := os.OpenFile(path, os.O_RDWR, 0700)
	return &Config{fd: fd}, err
}

func (c *Config) Close() error {
	return c.fd.Close()
}

func (c *Config) VendorID() uint16 {
	return util.ReadIo16(c.fd, offsetVendor)
}

func (c *Config) DeviceID() uint16 {
	return util.ReadIo16(c.fd, offsetDeviceID)
}

func (c *Config) Command() uint16 {
	return util.ReadIo16(c.fd, offsetCommand)
}

func (c *Config) Status() uint16 {
	return util.ReadIo16(c.fd, offsetStatus)
}

func (c *Config) RevisionID() uint8 {
	return util.ReadIo8(c.fd, offsetRevisionID)
}

func (c *Config) ClassCode() (base, sub, iface uint8) {
	b := make([]byte, 3)
	_, err := c.fd.ReadAt(b, offsetClassCode)
	if err != nil {
		panic(err)
	}
	return b[0], b[1], b[2]
}

func (c *Config) SetCommand(cmd uint16) error {
	util.WriteIo16(c.fd, cmd, offsetCommand)
	return nil
}

func (c *Config) SetStatus(status uint16) error {
	util.WriteIo16(c.fd, status, offsetStatus)
	return nil
}
