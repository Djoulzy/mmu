package main

import (
	"github.com/Djoulzy/chip"
	"github.com/Djoulzy/mmu"
)

var (
	nbPage   uint = 256
	pageSize uint = mmu.PAGE_SIZE

	RAM *chip.RAM = chip.NewRAM("RAM", pageSize*nbPage, false)
	AUX *chip.RAM = chip.NewRAM("AUX", pageSize*nbPage, false)
	ROM *chip.ROM = chip.NewROM("ROM", 0x2FFF)
	IO  *chip.ROM = chip.NewROM("IO", 0x0FFF)

	myMMU *mmu.MMU = mmu.Init(pageSize, nbPage)
)

func Init() {
	myMMU.Attach(RAM, 0, 256)
	myMMU.Attach(AUX, 0, 256)
	// myMMU.Attach(IO, 0xC0, 1)
	myMMU.Attach(ROM, 0xD0, 1)

	// myMMU.SwitchFullTo("RAM")
	// // myMMU.SwitchFullTo("IO")
	// myMMU.SwitchFullTo("ROM")

	myMMU.DumpMap()
}

func main() {
	Init()

	myMMU.Write(0x8000, 0xAA)
	myMMU.Write(0xC0FF, 0xBB)
	myMMU.Write(0xD099, 0xCC)
	myMMU.Read(0xD099)
}
