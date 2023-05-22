package mmu

import (
	"fmt"
	"os"
)

const PAGE_SIZE = 256

type ChipAccess interface {
	GetName() string
	Read(uint16) byte
	Write(uint16, byte)
	ReadOnly() bool
}

type chipInfos struct {
	startPage uint
	nbPages   uint
	baseAddr  uint16
	access    ChipAccess
}

type MMU struct {
	NbPage       uint
	AddressRange uint
	chips        map[string]chipInfos
	reader       []chipInfos
	writter      []chipInfos
}

func Init(pageSize uint, nbPages uint) *MMU {
	tmp := MMU{
		NbPage:       nbPages,
		AddressRange: pageSize * nbPages,
		reader:       make([]chipInfos, nbPages),
		writter:      make([]chipInfos, nbPages),
		chips:        make(map[string]chipInfos),
	}
	return &tmp
}

func (m *MMU) GetSize() uint {
	return m.AddressRange
}

func (m *MMU) Attach(chip ChipAccess, startPage uint, nbPages uint) {
	var i uint

	if nbPages*PAGE_SIZE+startPage > m.AddressRange {
		fmt.Printf("%s: Size Error\n", chip.GetName())
		os.Exit(0)
	}
	m.chips[chip.GetName()] = chipInfos{
		startPage: startPage,
		nbPages:   nbPages,
		baseAddr:  uint16(startPage) << 8,
		access:    chip,
	}

	fmt.Printf("Attach %s page $%02X to $%02X\n", chip.GetName(), startPage, startPage+nbPages-1)
	for i = startPage; i < (startPage + nbPages); i++ {
		m.reader[i] = m.chips[chip.GetName()]
		if !chip.ReadOnly() {
			m.writter[i] = m.chips[chip.GetName()]
		}
	}
}

func (m *MMU) SwitchZoneTo(name string, startPage uint, nbPages uint) {
	var i uint
	chip := m.chips[name]
	for i = startPage; i < (startPage + nbPages); i++ {
		m.reader[i] = chip
		if !chip.access.ReadOnly() {
			m.writter[i] = chip
		}
	}
}

func (m *MMU) SwitchFullTo(name string) {
	var i uint
	chip := m.chips[name]
	for i = chip.startPage; i < (chip.startPage + chip.nbPages); i++ {
		m.reader[i] = chip
		if !chip.access.ReadOnly() {
			m.writter[i] = chip
		}
	}
}

func (m *MMU) Read(addr uint16) byte {
	chipInfo := m.reader[addr>>8]
	return chipInfo.access.Read(addr - chipInfo.baseAddr)
}

func (m *MMU) Write(addr uint16, data byte) {
	chipInfo := m.writter[addr>>8]
	// fmt.Printf("Found %s at %02X\n", chipInfo.access.GetName(), addr>>8)
	chipInfo.access.Write(addr-chipInfo.baseAddr, data)
}

func (m *MMU) DumpMap() {
	for page, chip := range m.reader {
		fmt.Printf("%02X - %s\n", page, chip.access.GetName())
	}
}
