package main

import (
	"github.com/Djoulzy/mmu"
)

var (
	nbPage   uint = 256
	pageSize uint = mmu.PAGE_SIZE

	RAM *mmu.RAM = mmu.NewRAM("RAM", pageSize*nbPage, false)
	AUX *mmu.RAM = mmu.NewRAM("AUX", pageSize*nbPage, false)
	ROM *mmu.ROM = mmu.NewROM("ROM_D", 4096, "./cmd/main/assets/D.bin")

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
