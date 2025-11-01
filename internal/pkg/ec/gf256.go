package ec

var (
	// A table which represents the following structure
	// by using one-dimensional slice.
	// {0, 0}
	// {1, 0}, {1, 1}
	// ...
	// {254, 0}, {254, 1}, ... , {254, 254}
	// {255, 0}, {255, 0}, ... , {255, 254}, {255, 255}
	//
	// Here, {a, b} represents the result of a*b.
	mulTable []gf256

	invTable [256]gf256
)

// Irreducible polynomial for GF(2^8): x^8 + x^4 + x^3 + x + 1
const gfPolynomial = 0x11B

type gf256 byte

func init() {
	// Initialize mulTable.
	const n = 256
	size := n * (n + 1) / 2
	mulTable = make([]gf256, size)

	for a := range n {
		for b := range a + 1 {
			mulTable[offset(a)+b] = gf256(a).calcMul(gf256(b))
		}
	}

	// Initialize invTable
	// Since he inverse of 0 does not exist,
	// this value will not be used.
	invTable[0] = 0

	for x := 1; x < 256; x++ {
		invTable[x] = gf256(x).calcInv()
	}
}

func offset(a int) int {
	return (a * (a + 1)) / 2
}

// ref. https://en.wikipedia.org/wiki/Ancient_Egyptian_multiplication#Russian_peasant_multiplication
func (g gf256) calcMul(x gf256) gf256 {
	var a uint16 = uint16(g)
	var b uint16 = uint16(x)
	var p uint16 = 0
	for range 8 {
		if (b & 1) != 0 {
			p ^= a
		}
		b >>= 1
		carry := (a & 0x80) != 0
		a <<= 1
		if carry {
			a ^= gfPolynomial
		}
	}
	return gf256(p & 0xFF)
}

func (g gf256) calcInv() gf256 {
	if g == 0 {
		panic("the inverse of 0 does not exist")
	}
	result := gf256(1)
	// Since the multiplicative group of a finite field
	// is a cyclic group, x^{-1} = x^{254} holds.
	for range 254 {
		result = result.calcMul(g)
	}
	return result
}

func Gf256Generator() gf256 {
	return gf256(2)
}

func (g gf256) Add(x gf256) gf256 {
	return g ^ x
}

func (g gf256) Sub(x gf256) gf256 {
	return g ^ x
}

func (g gf256) Mul(x gf256) gf256 {
	if g >= x {
		return mulTable[offset(int(g))+int(x)]
	}
	return mulTable[offset(int(x))+int(g)]
}

func (g gf256) Div(x gf256) gf256 {
	if x == 0 {
		panic("divide by zero")
	}
	return g.calcMul(invTable[x])
}

func (g gf256) Pow(i int) gf256 {
	if i < 0 {
		return 0
	}

	result := gf256(1)
	for range i {
		result = result.Mul(g)
	}
	return result
}

func (g gf256) Inv() gf256 {
	return invTable[g]
}
