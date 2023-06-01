package mmu

type RAM struct {
	IC
}

func NewRAM(name string, size uint) *RAM {
	var tmp RAM

	tmp.Size = size
	tmp.Name = name
	tmp.Buff = make([]byte, size)

	// fmt.Printf("Create new RAM %4s with size %d\n", name, size)
	return &tmp
}

func (c *RAM) Read(addr uint16) byte {
	data := c.Buff[addr]
	// fmt.Printf("Reading $%02X from %4s at $%04X\n", data, c.Name, addr)
	return data
}

func (c *RAM) Write(addr uint16, data byte) {
	// fmt.Printf("Writting $%02X to %4s at $%04X\n", data, c.Name, addr)
	c.Buff[addr] = data
}

func (c *RAM) Clear(interval int, startWith byte) {
	// interval: 0x40 pour C64
	//           0x1000 pour Apple
	// startWith: 0x00 pour C64
	//            0xFF pour Apple
	cpt := 0
	fill := byte(startWith)
	for i := range c.Buff {
		c.Buff[i] = fill
		cpt++
		if cpt == interval {
			fill = ^fill
			cpt = 0
		}
	}
}
