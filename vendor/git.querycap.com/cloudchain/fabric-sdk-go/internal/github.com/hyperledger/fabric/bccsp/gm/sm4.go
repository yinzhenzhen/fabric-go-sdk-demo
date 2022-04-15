package gm

import (
	"crypto/rand"
	"errors"
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
	"github.com/tjfoc/gmsm/sm4"
)

// GetRandomBytes returns len random looking bytes
func GetRandomBytes(len int) ([]byte, error) {
	if len < 0 {
		return nil, errors.New("Len must be larger than 0")
	}

	buffer := make([]byte, len)

	n, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n != len {
		return nil, fmt.Errorf("Buffer not filled. Requested [%d], got [%d]", len, n)
	}

	return buffer, nil
}

// AESCBCPKCS7Encrypt combines CBC encryption and PKCS7 padding
func SM4Encrypt(key, src []byte) ([]byte, error) {
	dst := make([]byte, len(src))
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipher.Encrypt(dst, src)
	//sm4.EncryptBlock(key, dst, src)
	return dst, nil
}

// AESCBCPKCS7Decrypt combines CBC decryption and PKCS7 unpadding
func SM4Decrypt(key, src []byte) ([]byte, error) {

	dst := make([]byte, len(src))
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipher.Decrypt(dst, src)
	//sm4.DecryptBlock(key, dst, src)
	return dst, nil
}

type gmsm4Encryptor struct{}

//实现 Encryptor 接口
func (*gmsm4Encryptor) Encrypt(k bccsp.Key, plaintext []byte, opts bccsp.EncrypterOpts) (ciphertext []byte, err error) {

	return SM4Encrypt(k.(*gmsm4PrivateKey).privKey, plaintext)
}

type gmsm4Decryptor struct{}

//实现 Decryptor 接口
func (*gmsm4Decryptor) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {

	return SM4Decrypt(k.(*gmsm4PrivateKey).privKey, ciphertext)
}
