package components

import (
	"encoding/hex"
	"fmt"
	"goXdagWallet/xdago/cryptography"
	"goXdagWallet/xdago/secp256k1"
	xdagoUtils "goXdagWallet/xdago/utils"
	"testing"
)

const walletAddress = "mO88ml4B++TmUVMicswt4pmFWIHZeDQ9"

const rStr = "9af60175f278b5066fb236703ddfa340d75584825f5428e496c3367831cf93ab"
const sStr = "57226c8e903f7a7066c65e13458a3e1686d00256f9b1e40b6e620ad35dffe10e"

// const rStr = "1fbfcecd45a2ff8b75274cdf22610aa211acc84e768a4e8ab756acb3d4807455"
// const sStr = "4fef06c991f3bf385d77db2d2b4936726a00b189f3555c8d158dfa62dd08e0b1"

func TestAddress(t *testing.T) {
	var privByte1 = [32]byte{0x70, 0x77, 0x5e, 0x94, 0xd8, 0xd8, 0x3e, 0x34, 0xce, 0xbc, 0x9b, 0xcf, 0xe1, 0xda, 0x5a, 0x4d, 0x3c, 0x5a, 0xbe, 0x25, 0xce, 0xf1, 0x12, 0x8c, 0xdc, 0xb7, 0xa9, 0x34, 0xf6, 0x4b, 0x30, 0x1a}
	//var privByte1 = [32]byte{0x6d, 0x0e, 0x3d, 0x5b, 0xce, 0x66, 0x56, 0xe7, 0x34, 0x08, 0x34, 0x58, 0x1e, 0x73, 0x00, 0xd8, 0x45, 0xc4, 0x0b, 0x6f, 0x57, 0x90, 0x12, 0xbf, 0x05, 0x2b, 0x32, 0x1f, 0xbc, 0xbc, 0xdb, 0xf4}

	var timestamp = [8]byte{0xca, 0x55, 0x72, 0xfd, 0x84, 0x01, 0x00, 0x00}
	var fieldTypes = [8]byte{0x51, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var block0 = [1024]byte{}

	copy(block0[8:16], fieldTypes[:])
	copy(block0[16:24], timestamp[:])

	var privKey = secp256k1.PrivKeyFromBytes(privByte1[:])
	pubKey := privKey.PubKey()
	copy(block0[512:], pubKey.SerializeCompressed())
	fmt.Println(hex.EncodeToString(block0[:545]))

	signHash := cryptography.HashTwice(block0[:545])

	fmt.Println(hex.EncodeToString(signHash[:]))

	r, s := cryptography.EcdsaSign(privKey, signHash[:])
	copy(block0[32:], r[:])
	copy(block0[64:], s[:])

	fmt.Println(hex.EncodeToString(block0[:512]))

	hash := cryptography.HashTwice(block0[:512])
	fmt.Println(hex.EncodeToString(hash[:]))

	address := xdagoUtils.Hash2Address(hash)
	fmt.Println(address)

}

func TestBlock(t *testing.T) {
	var timestamp = [8]byte{0x8a, 0x3c, 0x6e, 0xfd, 0x84, 0x01, 0x00, 0x00}
	var fieldTypes = [8]byte{0x51, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r, _ := hex.DecodeString(rStr)
	s, _ := hex.DecodeString(sStr)

	var block0 = [1024]byte{}

	copy(block0[8:16], fieldTypes[:])
	copy(block0[16:24], timestamp[:])

	copy(block0[32:], r[:])
	copy(block0[64:], s[:])

	fmt.Println(hex.EncodeToString(block0[:512]))

	hash := cryptography.HashTwice(block0[:512])
	fmt.Println(hex.EncodeToString(hash[:]))

	address := xdagoUtils.Hash2Address(hash)
	fmt.Println(address)

}

func TestVerify(t *testing.T) {
	// var privByte1 = [32]byte{0x70, 0x77, 0x5e, 0x94, 0xd8, 0xd8, 0x3e, 0x34, 0xce, 0xbc, 0x9b, 0xcf, 0xe1, 0xda, 0x5a, 0x4d, 0x3c, 0x5a, 0xbe, 0x25, 0xce, 0xf1, 0x12, 0x8c, 0xdc, 0xb7, 0xa9, 0x34, 0xf6, 0x4b, 0x30, 0x1a}
	var privByte1 = [32]byte{0xeb, 0x7b, 0xc3, 0x28, 0x33, 0xbc, 0x7f, 0x9e, 0x80, 0xb7, 0x43, 0x03, 0xe9, 0x78, 0x46, 0xd8, 0x02, 0xa8, 0x2c, 0xd0, 0xe9, 0xb1, 0xbc, 0xc8, 0x1a, 0x9b, 0xb2, 0x57, 0xac, 0x6e, 0x11, 0x8c}

	// var timestamp = [8]byte{0xca, 0x55, 0x72, 0xfd, 0x84, 0x01, 0x00, 0x00}
	var timestamp = [8]byte{0x8a, 0x3c, 0x6e, 0xfd, 0x84, 0x01, 0x00, 0x00}
	var fieldTypes = [8]byte{0x51, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r, _ := hex.DecodeString(rStr)
	s, _ := hex.DecodeString(sStr)

	var block0 = [1024]byte{}

	copy(block0[8:16], fieldTypes[:])
	copy(block0[16:24], timestamp[:])

	var privKey = secp256k1.PrivKeyFromBytes(privByte1[:])
	pubKey := privKey.PubKey()
	copy(block0[512:], pubKey.SerializeCompressed())
	fmt.Println(hex.EncodeToString(block0[:545]))

	signHash := cryptography.HashTwice(block0[:545])

	fmt.Println(hex.EncodeToString(signHash[:]))

	verify := cryptography.EcdsaVerify(pubKey, signHash[:], r, s)

	fmt.Println(verify)

}
