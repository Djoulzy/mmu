package main

import (
	"strconv"

	"github.com/Djoulzy/mmu"
)

const (
	romSize      = 4096
	softSwitches = 256
	slot_roms    = 256
)

var (
	nbPage   uint = 256
	pageSize uint = mmu.PAGE_SIZE

	MN_ZPS = mmu.NewRAM("MN_ZPS", 0x0200)
	MN___1 = mmu.NewRAM("MN___1", 0x0200)
	MN_TXT = mmu.NewRAM("MN_TXT", 0x0400)
	MN___2 = mmu.NewRAM("MN___2", 0x1800)
	MN_HGR = mmu.NewRAM("MN_HGR", 0x2000)
	MN___3 = mmu.NewRAM("MN___3", 0x9000)
	MN_SLT = mmu.NewRAM("MN_SLT", 0x0800)

	MN_BK1 = mmu.NewRAM("MN_BK1", 0x1000)
	MN_BK2 = mmu.NewRAM("MN_BK2", 0x1000)
	MN___4 = mmu.NewRAM("MN___4", 0x2000)

	// AUX_ZP = mmu.NewRAM("AX_ZP", zpStack)
	// AUX_LO = mmu.NewRAM("AX_LO", lowRamSize)
	// AUX_B1 = mmu.NewRAM("AX_B1", bankSize)
	// AUX_B2 = mmu.NewRAM("AX_B2", bankSize)
	// AUX_HI = mmu.NewRAM("AX_HI", hiRamSize)

	SLOTS  [8]*mmu.ROM
	IO     = mmu.NewRAM("IO", softSwitches)
	ROM_C  = mmu.NewROM("ROM_C", romSize, "./cmd/main/assets/C.bin")
	ROM_D  = mmu.NewROM("ROM_D", romSize, "./cmd/main/assets/D.bin")
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "./cmd/main/assets/EF.bin")

	myMMU *mmu.MMU = mmu.Init(pageSize, nbPage)
)

func Init() {
	myMMU.Attach(MN_ZPS, 0x00)
	myMMU.Attach(MN___1, 0x02)
	myMMU.Attach(MN_TXT, 0x04)
	myMMU.Attach(MN___2, 0x08)
	myMMU.Attach(MN_HGR, 0x20)
	myMMU.Attach(MN___3, 0x40)
	myMMU.Attach(MN_SLT, 0xC8)
	myMMU.Attach(MN_BK1, 0xD0)
	myMMU.Attach(MN_BK2, 0xD0)
	myMMU.Attach(MN___4, 0xE0)

	// myMMU.Attach(AUX_ZP, 0x00)
	// myMMU.Attach(AUX_LO, 0x02)
	// myMMU.Attach(AUX_B1, 0xD0)
	// myMMU.Attach(AUX_B2, 0xD0)
	// myMMU.Attach(AUX_HI, 0xE0)

	myMMU.Attach(IO, 0xC0)
	myMMU.Attach(ROM_C, 0xC0)
	myMMU.Attach(ROM_D, 0xD0)
	myMMU.Attach(ROM_EF, 0xE0)

	myMMU.Mount("MN_ZPS", "MN_ZPS")
	myMMU.Mount("MN___1", "MN___1")
	myMMU.Mount("MN_TXT", "MN_TXT")
	myMMU.Mount("MN___2", "MN___2")
	myMMU.Mount("MN_HGR", "MN_HGR")
	myMMU.Mount("MN___3", "MN___3")
	myMMU.Mount("MN_SLT", "MN_SLT")

	myMMU.Mount("ROM_C", "MN_SLT")
	myMMU.Mount("ROM_D", "MN_BK1")
	myMMU.Mount("ROM_EF", "MN___4")
	myMMU.Mount("IO", "IO")

	for i := 1; i < 8; i++ {
		SLOTS[i] = mmu.NewROM("SLOT_"+strconv.Itoa(i), slot_roms, "")
		myMMU.Attach(SLOTS[i], 0xC0+uint(i))
		myMMU.Mount("SLOT_"+strconv.Itoa(i), "")
	}

	// myMMU.CheckMapIntegrity()

	// // $C080
	// myMMU.Mount("MAIN_B2", "")
	// myMMU.Mount("MAIN_HI", "")

	// // $C081
	// myMMU.Mount("ROM_D", "MAIN_B2")
	// myMMU.Mount("ROM_EF", "MAIN_HI")

	// // $C082
	// myMMU.Mount("ROM_D", "")
	// myMMU.Mount("ROM_EF", "")

	// // $C083
	// myMMU.Mount("MAIN_B2", "MAIN_B2")
	// myMMU.Mount("MAIN_HI", "MAIN_HI")

	// // $C088
	// myMMU.Mount("MAIN_B1", "")
	// myMMU.Mount("MAIN_HI", "")

	// // $C089
	// myMMU.Mount("ROM_D", "MAIN_B1")
	// myMMU.Mount("ROM_EF", "MAIN_HI")

	// // $C08A
	// myMMU.Mount("ROM_D", "")
	// myMMU.Mount("ROM_EF", "")

	// // $C08B
	// myMMU.Mount("MAIN_B2", "MAIN_B1")
	// myMMU.Mount("MAIN_HI", "MAIN_HI")

	// // $C006
	// for i := 1; i < 8; i++ {
	// 	myMMU.Mount("SLOT_"+strconv.Itoa(i), "")
	// }

	// // $C007
	// myMMU.Mount("ROM_C", "")
	// myMMU.Mount("IO", "IO")

	// $C00A
	// myMMU.SwapRom("SLOT_3", "ROM_C")

	// // $C00B
	// myMMU.Mount("SLOT_3", "")

	myMMU.DumpMap()
}

func main() {
	Init()

	myMMU.Write(0x8000, 0xAA)
	myMMU.Write(0xC0FF, 0xBB)
	myMMU.Write(0xD099, 0xCC)
	myMMU.Read(0xD099)
}
