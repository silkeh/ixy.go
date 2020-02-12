package driver

import (
	"fmt"
	"syscall"

	"github.com/ixy-languages/ixy.go/pci"
)

const maxQueues = 64

//IxyInterface is the interface that has to be implemented for all substrates such as the ixgbe or virtio
type IxyInterface interface {
	RxBatch(uint16, []*PktBuf) uint32
	TxBatch(uint16, []*PktBuf) uint32
	ReadStats(*DeviceStats)
	setPromisc(bool)
	getLinkSpeed() uint32
	getIxyDev() IxyDevice
}

//IxyDevice contains information common across all substrates
type IxyDevice struct {
	PCI         *pci.Device
	DriverName  string
	NumRxQueues uint16
	NumTxQueues uint16
}

//IxyInit initializes the driver and hands back the interface
func IxyInit(pciAddr string, rxQueues, txQueues uint16) (IxyInterface, error) {
	if syscall.Getuid() != 0 {
		fmt.Println("Not running as root, this will probably fail")
	}

	dev, err := pci.NewDevice(pciAddr)
	if err != nil {
		return nil, err
	}

	// Read Device configuration
	config, err := dev.Config()
	if err != nil {
		return nil, err
	}
	defer config.Close()

	// Check if the device is an ethernet controller
	if base, sub, iface := config.ClassCode(); base != 2 || sub != 0 || iface != 0 {
		return nil, fmt.Errorf("device %v is not an ethernet controller", pciAddr)
	}

	// Check if device is an unsupported virtio device
	if config.VendorID() == 0x1af4 && config.DeviceID() >= 0x1000 {
		return nil, fmt.Errorf("virtio not supported")
	}

	// Probably an ixgbe device
	return NewIxgbe(dev, rxQueues, txQueues)
}

// IxyTxBatchBusyWait calls dev.TxBatch until all packets are queued with busy waiting
func IxyTxBatchBusyWait(dev IxyInterface, queueID uint16, bufs []*PktBuf) {
	numBufs := uint32(len(bufs))
	for numSent := uint32(0); numSent != numBufs; numSent += dev.TxBatch(queueID, bufs[numSent:]) {
		//busy wait
	}
}
