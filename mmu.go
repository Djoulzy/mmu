package mmu

import (
	"fmt"
	"os"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/emutools/charset"
)

const (
	PAGE_SIZE = 256
	READONLY  = 1
	WRITEONLY = 2
	READWRITE = 3
)

type ChipAccess interface {
	SetMMU(*MMU)
	GetSize() uint
	GetName() string
	GetBuff() []byte
	Read(uint16) byte
	Dump(uint16) byte
	Write(uint16, byte)
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
	chips        map[string]*chipInfos
	reader       []*chipInfos
	writter      []*chipInfos
}

func Init(pageSize uint, nbPages uint) *MMU {
	tmp := MMU{
		NbPage:       nbPages,
		AddressRange: pageSize * nbPages,
		reader:       make([]*chipInfos, nbPages),
		writter:      make([]*chipInfos, nbPages),
		chips:        make(map[string]*chipInfos),
	}
	return &tmp
}

func (m *MMU) GetSize() uint {
	return m.AddressRange
}

func (m *MMU) Attach(chip ChipAccess, startPage uint) {
	// var i uint

	nbPages := chip.GetSize() / PAGE_SIZE
	if chip.GetSize()+startPage > m.AddressRange {
		fmt.Printf("%s: Size Error\n", chip.GetName())
		os.Exit(0)
	}
	m.chips[chip.GetName()] = &chipInfos{
		startPage: startPage,
		nbPages:   nbPages,
		baseAddr:  uint16(startPage) << 8,
		access:    chip,
	}

	fmt.Printf("Attach %s page $%02X to $%02X\n", chip.GetName(), startPage, startPage+nbPages-1)
	// m.Mount(chip.GetName(), mode)

	chip.SetMMU(m)
}

func (m *MMU) GetChipMem(name string) []byte {
	return m.chips[name].access.GetBuff()
}

func (m *MMU) Mount(reader string, writer string) {
	var i uint
	var writerChip *chipInfos = nil

	readerChip := m.chips[reader]
	if writer != "" {
		writerChip = m.chips[writer]
	}
	for i = readerChip.startPage; i < (readerChip.startPage + readerChip.nbPages); i++ {
		m.reader[i] = readerChip
		m.writter[i] = writerChip
	}
}

func (m *MMU) SwapRom(from string, to string) {
	var i uint
	fromChip := m.chips[from]
	toChip := m.chips[to]
	for i = fromChip.startPage; i < (fromChip.startPage + fromChip.nbPages); i++ {
		m.reader[i] = toChip
	}
}

func (m *MMU) Read(addr uint16) byte {
	chipInfo := m.reader[addr>>8]
	return chipInfo.access.Read(addr - chipInfo.baseAddr)
}

func (m *MMU) DirectRead(addr uint16) byte {
	chipInfo := m.reader[addr>>8]
	return chipInfo.access.Dump(addr - chipInfo.baseAddr)
}

func (m *MMU) Write(addr uint16, data byte) {
	chipInfo := m.writter[addr>>8]
	// fmt.Printf("Found %s at %02X\n", chipInfo.access.GetName(), addr>>8)
	if chipInfo != nil {
		chipInfo.access.Write(addr-chipInfo.baseAddr, data)
	}
}

//////////////////////////////////////////////////////////////////////////
// Dumpers
//////////////////////////////////////////////////////////////////////////

func (m *MMU) getWriter(i uint) string {
	if m.writter[i] != nil {
		return m.writter[i].access.GetName()
	} else {
		return " - "
	}
}

func (m *MMU) DumpMap() {
	var i uint

	step := m.NbPage / 4
	for i = 0; i < step; i++ {
		fmt.Printf("%02X - %8s / %8s\t", i, m.reader[i].access.GetName(), m.getWriter(i))
		fmt.Printf("%02X - %8s / %8s\t", i+step, m.reader[i+step].access.GetName(), m.getWriter(i+step))
		fmt.Printf("%02X - %8s / %8s\t", i+(step*2), m.reader[i+(step*2)].access.GetName(), m.getWriter(i+(step*2)))
		fmt.Printf("%02X - %8s / %8s\n", i+(step*3), m.reader[i+(step*3)].access.GetName(), m.getWriter(i+(step*3)))
	}
}

func (m *MMU) CheckMapIntegrity() {
	var i uint
	for i = 0; i < m.NbPage; i++ {
		if m.reader[i].access == nil {
			fmt.Printf("Page %02X not allocated for reading !\n", i)
			os.Exit(0)
		}
		// if m.writter[i].access == nil {
		// 	fmt.Printf("Page %02X not allocated for writting !\n", i)
		// 	os.Exit(0)
		// }
	}
}

func (m *MMU) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = m.DirectRead(cpt)
			if val != 0x00 && val != 0xFF {
				line = line + clog.CSprintf("white", "black", "%02X", val) + " "
			} else {
				line = fmt.Sprintf("%s%02X ", line, val)
			}
			if _, ok := charset.PETSCII[val]; ok {
				ascii += string(charset.PETSCII[val])
			} else {
				ascii += "."
			}
			if i == 7 {
				line = fmt.Sprintf("%s  ", line)
			}
			cpt++
		}
		fmt.Printf("%s - %s\n", line, ascii)
	}
}
