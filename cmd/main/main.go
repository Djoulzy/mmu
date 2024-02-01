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

	AX_ZPS = mmu.NewRAM("AX_ZPS", 0x0200)
	AX___1 = mmu.NewRAM("AX___1", 0x0200)
	AX_TXT = mmu.NewRAM("AX_TXT", 0x0400)
	AX___2 = mmu.NewRAM("AX___2", 0x1800)
	AX_HGR = mmu.NewRAM("AX_HGR", 0x2000)
	AX___3 = mmu.NewRAM("AX___3", 0x9000)
	AX_SLT = mmu.NewRAM("AX_SLT", 0x0800)

	AX_BK1 = mmu.NewRAM("AX_BK1", 0x1000)
	AX_BK2 = mmu.NewRAM("AX_BK2", 0x1000)
	AX___4 = mmu.NewRAM("AX___4", 0x2000)

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

	myMMU.Attach(AX_ZPS, 0x00)
	myMMU.Attach(AX___1, 0x02)
	myMMU.Attach(AX_TXT, 0x04)
	myMMU.Attach(AX___2, 0x08)
	myMMU.Attach(AX_HGR, 0x20)
	myMMU.Attach(AX___3, 0x40)
	myMMU.Attach(AX_SLT, 0xC8)
	myMMU.Attach(AX_BK1, 0xD0)
	myMMU.Attach(AX_BK2, 0xD0)
	myMMU.Attach(AX___4, 0xE0)

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

	myMMU.CheckMapIntegrity()

	// // $C080
	// myMMU.Mount("MN_BK2", "")
	// myMMU.Mount("MN___4", "")

	// // // $C081
	// myMMU.Mount("ROM_D", "MN_BK2")
	// myMMU.Mount("ROM_EF", "MN___4")

	// // $C082
	// myMMU.Mount("ROM_D", "")
	// myMMU.Mount("ROM_EF", "")

	// // $C083
	myMMU.Mount("MN_BK2", "MN_BK2")
	myMMU.Mount("MN___4", "MN___4")

	// // $C088
	// myMMU.Mount("MN_BK1", "")
	// myMMU.Mount("MN___4", "")

	// // $C089
	// myMMU.Mount("ROM_D", "MN_BK1")
	// myMMU.Mount("ROM_EF", "MN___4")

	// // $C08A
	// myMMU.Mount("ROM_D", "")
	// myMMU.Mount("ROM_EF", "")

	// // $C08B
	// myMMU.Mount("MN_BK1", "MN_BK1")
	// myMMU.Mount("MN___4", "MN___4")

	// // $C006
	// for i := 1; i < 8; i++ {
	// 	myMMU.Mount("SLOT_"+strconv.Itoa(i), "")
	// }

	// // $C007
	// myMMU.Mount("ROM_C", "")
	// myMMU.Mount("IO", "IO")

	// // $C00B
	// myMMU.Mount("SLOT_3", "")

	// // $C002
	// myMMU.MountReader("MN_ZPS")
	// myMMU.MountReader("MN___1")
	// myMMU.MountReader("MN_TXT")
	// myMMU.MountReader("MN___2")
	// myMMU.MountReader("MN_HGR")
	// myMMU.MountReader("MN___3")

	// // $C003
	// myMMU.MountReader("AX_ZPS")
	// myMMU.MountReader("AX___1")
	// myMMU.MountReader("AX_TXT")
	// myMMU.MountReader("AX___2")
	// myMMU.MountReader("AX_HGR")
	// myMMU.MountReader("AX___3")

	// // $C004
	// myMMU.MountWriter("MN_ZPS")
	// myMMU.MountWriter("MN___1")
	// myMMU.MountWriter("MN_TXT")
	// myMMU.MountWriter("MN___2")
	// myMMU.MountWriter("MN_HGR")
	// myMMU.MountWriter("MN___3")

	// // $C005
	// myMMU.MountWriter("AX_ZPS")
	// myMMU.MountWriter("AX___1")
	// myMMU.MountWriter("AX_TXT")
	// myMMU.MountWriter("AX___2")
	// myMMU.MountWriter("AX_HGR")
	// myMMU.MountWriter("AX___3")

	myMMU.SwapChip("SLOT_3", "ROM_C")

	myMMU.SwapChip("MN_BK1", "AX_BK1")
	myMMU.SwapChip("MN_BK2", "AX_BK2")
	myMMU.SwapChip("MN___4", "AX___4")

	myMMU.DumpMap()
}

func main() {
	Init()

	myMMU.Write(0x8000, 0xAA)
	myMMU.Write(0xC0FF, 0xBB)
	myMMU.Write(0xD099, 0xCC)
	myMMU.Read(0xD099)
}
