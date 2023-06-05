package mmu

import (
	"os"
)

type IC struct {
	Name string
	Size uint
	Buff []byte
	Mmu  *MMU
	ChipAccess
}

// func NewIC(name string, size uint, read func(uint16) byte, write func(uint16, byte)) *IC {
func NewIC(name string, size uint) *IC {
	var tmp IC

	tmp.Size = size
	tmp.Name = name
	tmp.Buff = make([]byte, size)

	// fmt.Printf("Create new RAM %4s with size %d\n", name, size)
	return &tmp
}

func (ic *IC) GetName() string {
	return ic.Name
}

func (ic *IC) GetSize() uint {
	return ic.Size
}

func (ic *IC) LoadData(file string, memStart uint16) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	for i, val := range content {
		// mem.Write(memStart+uint16(i), val)
		ic.Buff[memStart+uint16(i)] = val
	}
	return nil
}

func (ic *IC) SetMMU(ref *MMU) {
	ic.Mmu = ref
}

func (ic *IC) Dump(addr uint16) byte {
	return ic.Buff[addr]
}
