package pci

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

type Device struct {
	Addr string
}

func NewDevice(addr string) (*Device, error) {
	pci := &Device{Addr: addr}
	if _, err := os.Stat(pci.path("")); os.IsNotExist(err) {
		return nil, fmt.Errorf("device %s does not exist", addr)
	}
	return pci, nil
}

func (dev *Device) Config() (*Config, error) {
	return NewPCIConfig(dev.path("config"))
}

func (dev *Device) RemoveDriver() error {
	return ioutil.WriteFile(dev.path("driver/unbind"), []byte(dev.Addr), 0700)
}

func (dev *Device) EnableDMA() error {
	c, err := dev.Config()
	if err != nil {
		return err
	}
	defer c.Close()

	// Set 'Bus Master Enable' (PCIe 3.0 specification section 7.5.1.1)
	return c.SetCommand(1 << 2)
}

func (dev *Device) MapResource() ([]byte, error) {
	err := dev.RemoveDriver()
	if err != nil {
		return nil, err
	}

	err = dev.EnableDMA()
	if err != nil {
		return nil, err
	}

	fd, err := dev.open("resource0")
	if err != nil {
		return nil, err
	}

	stat, _ := fd.Stat()
	mmap, err := syscall.Mmap(int(fd.Fd()), 0, int(stat.Size()),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	return mmap, err
}

func (dev *Device) open(file string) (*os.File, error) {
	return os.OpenFile(dev.path(file), os.O_RDWR, 0700)
}

func (dev *Device) path(file string) string {
	return fmt.Sprintf("/sys/bus/dev/devices/%v/%v", dev.Addr, file)
}
