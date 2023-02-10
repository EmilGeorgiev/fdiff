package rollinghash_test

import (
	"fmt"
	"github.com/EmilGeorgiev/fdiff/rollinghash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	fmt.Println(rollinghash.Pow(10, 97))
}

func TestOoff(t *testing.T) {
	data := []byte("LLorem ipsum dolor sit amet,Lore consectetur Lore adipiscingLore Loreelit, sedLore do eiusmodLore ")

	h := rollinghash.NewRabinFingerprint(data[:4])

	// 7336
	for i := 4; i < len(data); i++ {
		if h.Next(data[i]) == 7336 {
			fmt.Println("Index is breakpoint: ", i)
		}
	}

	//expected1 := rollinghash.NewRabinFingerprint(data[1:49])
	//expected2 := rollinghash.NewRabinFingerprint(data[2:50])
	//expected3 := rollinghash.NewRabinFingerprint(data[3:51])
	////expected4 := rollinghash.NewRabinFingerprint(data[4:44])
	////expected5 := rollinghash.NewRabinFingerprint(data[5:45])
	////expected6 := rollinghash.NewRabinFingerprint(data[6:46])
	//
	//assert.EqualValues(t, expected1.Value(), h.Next(data[48]))
	//assert.EqualValues(t, expected2.Value(), h.Next(data[49]))
	//assert.EqualValues(t, expected3.Value(), h.Next(data[50]))
	////assert.EqualValues(t, expected4.Value(), h.Next(data[43]))
	////assert.EqualValues(t, expected5.Value(), h.Next(data[44]))
	////assert.EqualValues(t, expected6.Value(), h.Next(data[45]))
	//
	//fmt.Println(expected1.Value()) // 152   - 3292
	//fmt.Println(expected2.Value()) // 5824  - 152
	//fmt.Println(expected3.Value()) //  x    - 5824
	//fmt.Println(expected4.Value()) //   -
	//fmt.Println(expected5.Value()) //   -
	//fmt.Println(expected6.Value()) // -   -
}

func TestOffSet257(t *testing.T) {

	fmt.Println(6384828 % 113)

	a := ((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 +
		((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 + 1
	fmt.Println("Hash of 'a': ", a%8191)

	b := ((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 +
		((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 + 257
	fmt.Println("Hash of 'b': ", b%8191)

	fmt.Println("ab: ", (a+b)%8191)

	c := ((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 +
		((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 + 257 + 1
	fmt.Println("Hash of 'c': ", c%8191)

	d := ((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 +
		((257%8191)*(257%8191)*(257%8191)*(257%8191)*(257%8191))%8191 + 257*257
	fmt.Println("Hash of 'd': ", d%8191)

	// 97 =  01100001 = (257^6 + 257^5 +257^0) % 8191 = 737
	// 98 =  01100010 = (257^6 + 257^5 +257^1) % 8191 = 993
	// 99 =  01100011 = (257^6 + 257^5 +257^1 + 257^0) % 8191 = 994
	// 100 = 01100100 = (257^6 + 257^5 +257^2) % 8191 = 1257
	// Hash of 'a':  737
	// Hash of 'b':  993
	// Hash of 'c':  994
	// Hash of 'd':  1257

	// Hash ab = (737 + 993) % 8191 = 1730

	data := []byte("abcd")

	fmt.Println("C: ", rollinghash.NewRabinFingerprint(data[:2]).Value())

	h := rollinghash.NewRabinFingerprint(data[:2])
	fmt.Println(h.Value())
	fmt.Println(h.Next(data[2]))
	fmt.Println(h.Next(data[3]))
}

func TestCalculateRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	// Vestibulum neque massa, scelerisque sit amet ligula eu, congue molestie mi. Praesent ut varius sem. Nullam at porttitor arcu, nec lacinia nisi. Ut ac dolor vitae odio interdum condimentum. Vivamus dapibus sodales ex, vitae malesuada ipsum cursus convallis. Maecenas sed egestas nulla, ac condimentum orci. Mauris diam felis, vulputate ac suscipit et, iaculis non est. Curabitur semper arcu ac ligula semper, nec luctus nisl blandit. Integer lacinia ante ac libero lobortis imperdiet. Nullam mollis convallis ipsum, ac accumsan nunc vehicula vitae. Nulla eget justo in felis tristique fringilla. Morbi sit amet tortor quis risus auctor condimentum. Morbi in ullamcorper elit. Nulla iaculis tellus sit amet mauris tempus fringilla.
	oldB := []byte("abcd")
	NewB := []byte("aabcd")
	//fmt.Println(len(p))
	// Action
	h := rollinghash.NewRabinFingerprint(oldB[:2])
	actual := h.Value()

	old := []uint64{actual}
	for _, b := range oldB {
		v := h.Next(b)
		old = append(old, v)
	}
	fmt.Println(old)

	h2 := rollinghash.NewRabinFingerprint(NewB[:2])
	newS := []uint64{h2.Value()}
	for _, b := range NewB {
		v := h.Next(b)
		newS = append(newS, v)
	}
	fmt.Println(newS)

	// Assert
	expected := 3873
	assert.EqualValues(t, expected, actual)
}

func TestCalculateRabinFingerPrintFor4Bytess(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("bcde")
	// Action
	h := rollinghash.NewRabinFingerprint(p)
	actual := h.Value()

	// Assert
	expected := 7293
	assert.EqualValues(t, expected, actual)
}

func TestCalculateNextRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("abcd")
	// Action
	h := rollinghash.NewRabinFingerprint(p)
	b := []byte("e")
	actual := h.Next(b[0])

	// Assert
	expected := 7293
	assert.EqualValues(t, expected, actual)
}

//
//func TestCalculateRabinFingerPrintFor48Bytes(t *testing.T) {
//	// Set up
//	p := "This is a program that calculates the difference"
//	// Action
//	h := rabin.NewHash(p)
//	actual := h.Value()
//
//	// Assert
//	expected := 1477
//	assert.EqualValues(t, expected, actual)
//}
//
//func TestCalculateRabinFingerPrintOnNonASCIICharacters(t *testing.T) {
//	// Set up
//	p := "абвгྠྡ ྡ"
//
//	// Action
//	h := rabin.NewHash(p)
//	actual := h.Value()
//
//	// Assert
//	expected := 5378
//	assert.EqualValues(t, expected, actual)
//}
//
//func TestNextRollingHash(t *testing.T) {
//	// Set up
//	rh := rabin.NewHash("abcd")
//
//	// Action
//	actual := rh.Next('e')
//
//	// Assert
//	expected := 15
//	assert.EqualValues(t, expected, actual)
//}
