package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"time"
)

func Root(a *big.Float, n uint64) *big.Float {
	limit := Exp(New(2), 1024)
	n1 := n - 1
	n1f, rn := New(float64(n1)), Div(New(1.0), New(float64(n)))
	x, x0 := New(1.0), Zero()
	_ = x0
	for {
		potx, t2 := Div(New(1.0), x), a
		for b := n1; b > 0; b >>= 1 {
			if b&1 == 1 {
				t2 = Mul(t2, potx)
			}
			potx = Mul(potx, potx)
		}
		x0, x = x, Mul(rn, Add(Mul(n1f, x), t2))
		if Lesser(Mul(Abs(Sub(x, x0)), limit), x) {
			break
		}
	}
	return x
}

func Abs(a *big.Float) *big.Float {
	return Zero().Abs(a)
}

func Exp(a *big.Float, e uint64) *big.Float {
	result := Zero().Copy(a)
	for i := uint64(0); i < e-1; i++ {
		result = Mul(result, a)
	}
	return result
}

func New(f float64) *big.Float {
	r := big.NewFloat(f)
	r.SetPrec(1024)
	return r
}

func Div(a, b *big.Float) *big.Float {
	return Zero().Quo(a, b)
}

func Zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(256)
	return r
}

func Mul(a, b *big.Float) *big.Float {
	return Zero().Mul(a, b)
}

func Add(a, b *big.Float) *big.Float {
	return Zero().Add(a, b)
}

func Sub(a, b *big.Float) *big.Float {
	return Zero().Sub(a, b)
}

func Lesser(x, y *big.Float) bool {
	return x.Cmp(y) == -1
}

type Pub struct {
	E int64
	N *big.Int
}

func (pub *Pub) Size() int {
	return (pub.N.BitLen() + 7) / 8
}

func encrypt(c *big.Int, pub *Pub, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

func mgf1XOR(out []byte, hash hash.Hash, seed []byte) {
	var counter [4]byte
	var digest []byte

	done := 0
	for done < len(out) {
		hash.Write(seed)
		hash.Write(counter[0:4])
		digest = hash.Sum(digest[:0])
		hash.Reset()

		for i := 0; i < len(digest) && done < len(out); i++ {
			out[done] ^= digest[i]
			done++
		}
		incCounter(&counter)
	}
}

func incCounter(c *[4]byte) {
	if c[3]++; c[3] != 0 {
		return
	}
	if c[2]++; c[2] != 0 {
		return
	}
	if c[1]++; c[1] != 0 {
		return
	}
	c[0]++
}

func DefaultPolynomial(x, n *big.Int) *big.Int {
	one := big.NewInt(1)
	two := big.NewInt(2)
	x2 := new(big.Int).Exp(x, two, n) // x^2 mod n
	x2.Add(x2, one)                   // (x^2 mod n) + 1
	x2.Mod(x2, n)                     // (x^2 + 1) mod n
	return x2
}

func Rho(n *big.Int, g func(*big.Int, *big.Int) *big.Int) (*big.Int, error) {
	one := big.NewInt(1)
	x, y, d := big.NewInt(2), big.NewInt(2), big.NewInt(1)

	for d.Cmp(one) == 0 {
		x = g(x, n)
		y = g(g(y, n), n)
		sub := new(big.Int).Sub(x, y)
		d.GCD(nil, nil, sub.Abs(sub), n) // gcd(|x - y|, n)
	}

	if d.Cmp(n) == 0 {
		return nil, errors.New("algorithm failed with default parameters (x = y = 2)")
	}

	return d, nil
}

func messageInt(msg string, pub *Pub) *big.Int {
	label := []byte("")

	hash := sha256.New()

	hash.Reset()
	k := pub.Size()

	if len(msg) > k-2*hash.Size()-2 {
		panic("message too long")
	}

	hash.Write(label)
	lHash := hash.Sum(nil)
	hash.Reset()

	em := make([]byte, k)
	seed := em[1 : 1+hash.Size()]
	db := em[1+hash.Size():]

	copy(db[0:hash.Size()], lHash)
	db[len(db)-len(msg)-1] = 1
	copy(db[len(db)-len(msg):], msg)

	rng := rand.Reader

	_, err := io.ReadFull(rng, seed)
	if err != nil {
		panic(err)
	}

	mgf1XOR(db, hash, seed)
	mgf1XOR(seed, hash, db)

	m := new(big.Int)
	m.SetBytes(em)

	return m
}

const modulus string = "685418641534116524651241278167264524621421421421486214321451246214128495146217321657217621789217621324652142145000145794613021546487542151203124548881512000003164275464512130502163464972431612130316942769784234312064673421906423146164312198060306491313130360946434797816030219494342060302111114778895232211011"

var exponents = [...]int64{
	3,
	5,
	17,
	257,
	1337,
	65537,
}

func main() {
	n := new(big.Int)
	n.SetString(modulus, 10)

	for _, e := range exponents {

		pub := &Pub{
			e,
			n,
		}

		mf1, err := ioutil.ReadFile("testdata/plaintexts/test1.txt")
		if err != nil {
			panic(err)
		}
		m1 := messageInt(string(mf1), pub)
		mb1 := new(big.Int).SetBytes(mf1)
		// int1 := messageInt("First test", pub)
		// cb1 := encrypt(new(big.Int), pub, int1)
		cf1, err := ioutil.ReadFile("testdata/ciphered/test1.txt")
		if err != nil {
			panic(err)
		}
		cb1 := new(big.Int).SetBytes(cf1)

		mf2, err := ioutil.ReadFile("testdata/plaintexts/test2.txt")
		if err != nil {
			panic(err)
		}
		m2 := messageInt(string(mf2), pub)
		mb2 := new(big.Int).SetBytes(mf2)
		// int2 := messageInt("Second text", pub)
		// cb2 := encrypt(new(big.Int), pub, int2)
		cf2, err := ioutil.ReadFile("testdata/ciphered/test2.txt")
		if err != nil {
			panic(err)
		}
		cb2 := new(big.Int).SetBytes(cf2)

		fmt.Println("--------------------------------------------")
		fmt.Println("e:   ", pub.E)

		fmt.Println("m1:  ", m1.BitLen())
		fmt.Println("m2:  ", m2.BitLen())
		fmt.Println("cb1: ", cb1.BitLen())
		fmt.Println("cb2: ", cb2.BitLen())

		// f1 := new(big.Float)
		// f1.SetInt(cb1)

		// root1 := Root(f1, uint64(e))

		// if root1.IsInt() {
		// 	result := new(big.Int)
		// 	root1.Int(result)
		// 	fmt.Println("the e-th root of the cipher is a int!")
		// 	fmt.Println("c1 ^ (1 / e): ", result.BitLen())
		// 	// fmt.Println(result.Bytes())
		// 	fmt.Println("ib1: ", ib1.BitLen())
		// 	// fmt.Println(ib1.Bytes())
		// }

		// f2 := new(big.Float)
		// f2.SetInt(cb2)

		// root2 := Root(f2, uint64(e))

		// if root2.IsInt() {
		// 	result := new(big.Int)
		// 	root2.Int(result)
		// 	fmt.Println("the e-th root of the cipher is a int!")
		// 	fmt.Println("c2 ^ (1 / e): ", result.BitLen())
		// 	// fmt.Println(result.Bytes())
		// 	fmt.Println("ib2: ", ib2.BitLen())
		// 	// fmt.Println(ib2.Bytes())
		// }

		e1 := new(big.Int)
		e1.Exp(m1, big.NewInt(pub.E), nil)
		fmt.Println("m1 ^ e:  ", e1.BitLen())

		eb1 := new(big.Int)
		eb1.Exp(mb1, big.NewInt(pub.E), nil)
		fmt.Println("mb1 ^ e: ", eb1.BitLen())

		eb1.Sub(eb1, cb1)
		fmt.Println("x1 = mb1 ^ e - cb1: ", eb1.BitLen())

		// // r1 := e1.Mod(e1, pub.N)
		// // fmt.Println("m1 ^ e - cb1 (mod N):  ", r1)

		e2 := new(big.Int)
		e2.Exp(m2, big.NewInt(pub.E), nil)
		fmt.Println("m2 ^ e:  ", e2.BitLen())

		eb2 := new(big.Int)
		eb2.Exp(mb2, big.NewInt(pub.E), nil)
		fmt.Println("mb2 ^ e: ", eb2.BitLen())

		eb2.Sub(eb2, cb2)
		fmt.Println("x2 = mb2 ^ e - cb2: ", eb2.BitLen())

		// // r2 := e2.Mod(e2, pub.N)
		// // fmt.Println("m2 ^ e - cb2 (mod N): ", r2)

		// now we have to find the GCD of e1 and e2
		fmt.Printf("Begin GCD computation ...")
		t0 := time.Now()
		gcd := new(big.Int).GCD(nil, nil, eb1, eb2)
		fmt.Printf(" done in %v\n", time.Since(t0))

		fmt.Println("GCD(x1, x2): ", gcd)

	}

	// compute the factorization of e1 and e2
	// and look for common factors

	// f1, err := Rho(e1, DefaultPolynomial)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	fmt.Println(f1)

	// 	if f1.Cmp(big.NewInt(1)) == 0 {
	// 		break
	// 	}

	// 	e1 = e1.Div(e1, f1)

	// 	f1, err = Rho(e1, DefaultPolynomial)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// f2, err := Rho(e2, DefaultPolynomial)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	fmt.Println(f2)

	// 	if f2.Cmp(big.NewInt(1)) == 0 {
	// 		break
	// 	}

	// 	e2 = e2.Div(e2, f2)

	// 	f2, err = Rho(e2, DefaultPolynomial)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// e1 := new(big.Int)
	// e1.Exp(int1, big.NewInt(pub.E), nil)
	// fmt.Println("m1 ^ e: ", e1.BitLen())

	// e1.Mod(e1, pub.N)
	// fmt.Println("m1 ^ e (mod N): ", e1.BitLen())

	// r1 := e1.Cmp(c1)
	// fmt.Println("m1 ^ e === c1 (mod N): ", r1)

	// e2 := new(big.Int)
	// e2.Exp(int2, big.NewInt(pub.E), nil)
	// fmt.Println("m2 ^ e: ", e2.BitLen())

	// e2.Mod(e2, pub.N)
	// fmt.Println("m2 ^ e (mod N): ", e2.BitLen())

	// r2 := e2.Cmp(c2)
	// fmt.Println("m2 ^ e === c2 (mod N): ", r2)

	// fmt.Println("GEN cipher bytes: ", len(xd))

	// c := encrypt(new(big.Int), pub, m)
	// out := c.Bytes()

	// if len(out) < k {
	// 	// If the output is too small, we need to left-pad with zeros.
	// 	t := make([]byte, k)
	// 	copy(t[k-len(out):], out)
	// 	out = t
	// }

	// fmt.Println("GIVEN cipher bytes: ", len(out))
	// // fmt.Println("cipher: ", c)

	// e1 := new(big.Int)
	// e1.Exp(m, big.NewInt(pub.E), nil)
	// fmt.Println("e1: ", e1.BitLen())

	// e1.Mod(e1, pub.N)
	// fmt.Println("mod: ", e1.BitLen())

	// r := e1.Cmp(c)
	// fmt.Println("mod eq: ", r)

	// plainFile1 := "testdata/plaintexts/test1.txt"
	// cipherFile1 := "testdata/ciphered/test1.txt"

	// plainFile2 := "testdata/plaintexts/test2.txt"
	// cipherFile2 := "testdata/ciphered/test2.txt"

	// plainBytes1, err := ioutil.ReadFile(plainFile1)
	// if err != nil {
	// 	panic(err)
	// }

	// plainBytes2, err := ioutil.ReadFile(plainFile2)
	// if err != nil {
	// 	panic(err)
	// }

	// cipherBytes1, err := ioutil.ReadFile(cipherFile1)
	// if err != nil {
	// 	panic(err)
	// }

	// cipherBytes2, err := ioutil.ReadFile(cipherFile2)
	// if err != nil {
	// 	panic(err)
	// }

	// c1 := new(big.Int).SetBytes(cipherBytes1)
	// c2 := new(big.Int).SetBytes(cipherBytes2)

	// m1 := new(big.Int).SetBytes(plainBytes1)
	// m2 := new(big.Int).SetBytes(plainBytes2)

	// // rsae := big.NewInt(2731)
	// rsae := big.NewInt(3)

	// e1 := new(big.Int)
	// e1.Exp(m1, rsae, nil)

	// e2 := new(big.Int)
	// e2.Exp(m2, rsae, nil)

	// s1 := new(big.Int)
	// s1.Sub(e1, c1)

	// s2 := new(big.Int)
	// s2.Sub(e2, c2)

	// gcd := new(big.Int)
	// gcd.GCD(nil, nil, s1, s2)

	// fmt.Println("factorizing s1")

	// f1, err := Rho(s1, DefaultPolynomial)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	fmt.Println(f1)

	// 	if f1.Cmp(big.NewInt(1)) == 0 {
	// 		break
	// 	}

	// 	s1 = s1.Div(s1, f1)

	// 	f1, err = Rho(s1, DefaultPolynomial)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	break
	// }

	// fmt.Println("factorizing s2")

	// f2, err := Rho(s2, DefaultPolynomial)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	fmt.Println(f2)

	// 	if f2.Cmp(big.NewInt(1)) == 0 {
	// 		break
	// 	}

	// 	s2 = s2.Div(s2, f2)

	// 	f2, err = Rho(s2, DefaultPolynomial)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	break
	// }

}
