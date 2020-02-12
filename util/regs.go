package util

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

//map C functions to Go

func SetReg32(addr []byte, reg int, value uint32) {
	atomic.StoreUint32((*uint32)(unsafe.Pointer(&addr[reg])), value)
}

func GetReg32(addr []byte, reg int) uint32 {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
}

func SetFlags32(addr []byte, reg int, flags uint32) {
	SetReg32(addr, reg, GetReg32(addr, reg)|flags)
}

func ClearFlags32(addr []byte, reg int, flags uint32) {
	SetReg32(addr, reg, GetReg32(addr, reg)&^flags)
}

func WaitClearReg32(addr []byte, reg int, mask uint32) {
	cur := atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	for (cur & mask) != 0 {
		fmt.Printf("waiting for flags %#x in register %#x to clear, current value %#x", mask, reg, cur)
		time.Sleep(10000 * time.Microsecond)
		cur = atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	}
}

func WaitSetReg32(addr []byte, reg int, mask uint32) {
	cur := atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	for (cur & mask) != mask {
		fmt.Printf("waiting for flags %#x in register %#x to clear, current value %#x", mask, reg, cur)
		time.Sleep(10000 * time.Microsecond)
		cur = atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	}
}

//getter for pci io port resources
func ReadIo32(fd *os.File, offset uint) uint32 {
	fd.Sync()
	b := make([]byte, 4)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	return HostOrder.Uint32(b)
}

func ReadIo16(fd *os.File, offset uint) uint16 {
	fd.Sync()
	b := make([]byte, 2)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	return HostOrder.Uint16(b)
}

func ReadIo8(fd *os.File, offset uint) uint8 {
	fd.Sync()
	b := make([]byte, 1)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	return uint8(b[0])
}

//setter for pci io port resources
func writeIo32(fd *os.File, value uint32, offset uint) {
	b := make([]byte, 4)
	HostOrder.PutUint32(b, value)
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}

func WriteIo16(fd *os.File, value uint16, offset uint) {
	b := make([]byte, 2)
	HostOrder.PutUint16(b, value)
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}

func WriteIo8(fd *os.File, value uint8, offset uint) {
	b := make([]byte, 1)
	b[0] = byte(value)
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}
