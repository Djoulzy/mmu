package main

import (
	"github.com/Djoulzy/mmu"
)

const (
	lowRamSize   = 53248
	hiRamSize    = 8192
	bankSize     = 4096
	romSize      = 4096
	softSwitches = 256
	chargenSize  = 2048
	keyboardSize = 2048
	blanckSize   = 12288
	slot_roms    = 256
)

var (
	nbPage   uint = 256
	pageSize uint = mmu.PAGE_SIZE

	MAIN_LOW = mmu.NewRAM("MAIN_LOW", lowRamSize)
	MAIN_B1  = mmu.NewRAM("MAIN_B1", bankSize)
	MAIN_B2  = mmu.NewRAM("MAIN_B2", bankSize)
	MAIN_HI  = mmu.NewRAM("MAIN_HI", hiRamSize)

	IO     = mmu.NewRAM("IO", softSwitches)
	ROM_C  = mmu.NewROM("ROM_C", romSize, "./cmd/main/assets/C.bin")
	ROM_D  = mmu.NewROM("ROM_D", romSize, "./cmd/main/assets/D.bin")
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "./cmd/main/assets/EF.bin")

	myMMU *mmu.MMU = mmu.Init(pageSize, nbPage)
)

func Init() {
	myMMU.Attach(MAIN_LOW, 0x00)
	myMMU.Attach(MAIN_B1, 0xD0)
	myMMU.Attach(MAIN_B2, 0xD0)
	myMMU.Attach(MAIN_HI, 0xE0)
	myMMU.Attach(IO, 0xC0)
	myMMU.Attach(ROM_C, 0xC0)
	myMMU.Attach(ROM_D, 0xD0)
	myMMU.Attach(ROM_EF, 0xE0)

	myMMU.Mount("MAIN_LOW", "MAIN_LOW")
	myMMU.Mount("ROM_C", "MAIN_LOW")
	myMMU.Mount("ROM_D", "MAIN_B1")
	myMMU.Mount("ROM_EF", "MAIN_HI")
	myMMU.Mount("IO", "IO")

	myMMU.CheckMapIntegrity()

	// $C080
	myMMU.Mount("MAIN_B2", "")
	myMMU.Mount("MAIN_HI", "")

	// $C081
	myMMU.Mount("ROM_D", "MAIN_B2")
	myMMU.Mount("ROM_EF", "MAIN_HI")

	// $C082
	myMMU.Mount("ROM_D", "")
	myMMU.Mount("ROM_EF", "")

	// $C083
	myMMU.Mount("MAIN_B2", "MAIN_B2")
	myMMU.Mount("MAIN_HI", "MAIN_HI")

	// $C088
	myMMU.Mount("MAIN_B1", "")
	myMMU.Mount("MAIN_HI", "")

	// $C089
	myMMU.Mount("ROM_D", "MAIN_B1")
	myMMU.Mount("ROM_EF", "MAIN_HI")

	// $C08A
	myMMU.Mount("ROM_D", "")
	myMMU.Mount("ROM_EF", "")

	// $C08B
	myMMU.Mount("MAIN_B2", "MAIN_B1")
	myMMU.Mount("MAIN_HI", "MAIN_HI")

	myMMU.DumpMap()
}

func main() {
	Init()

	myMMU.Write(0x8000, 0xAA)
	myMMU.Write(0xC0FF, 0xBB)
	myMMU.Write(0xD099, 0xCC)
	myMMU.Read(0xD099)
}
