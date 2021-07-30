package gpio

func (g *GPIO) Sync(pin int, val bool) {
	g.Set(pin, val)
}

func (g *GPIO) ReverseSync(pin int, val bool) {
	if val {
		g.Set(pin, false)
	} else {
		g.Set(pin, true)
	}
}
