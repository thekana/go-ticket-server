package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

func Contains(arr []string, s string) bool {
	for _, n := range arr {
		if s == n {
			return true
		}
	}
	return false
}

func EnsureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return errors.Wrap(err, "Could not create directory: "+dir)
		}
	}
	return nil
}

func GenerateRSAKeyPair(keyID string, keyDirPath string, keyLength int) error {
	if err := EnsureDir(keyDirPath, 0700); err != nil {
		return errors.Wrap(err, "Could not create key directory")
	}

	// Generate key pair
	reader := rand.Reader
	bitSize := keyLength

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return err
	}

	// Write private key to file
	privateKeyFile, err := os.Create(keyDirPath + "/" + keyID + ".pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		return err
	}
	//

	// Write public key to file
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}

	var publicKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	publicKeyFile, err := os.Create(keyDirPath + "/" + keyID + "_pub.pem")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKey)
	if err != nil {
		return err
	}

	return err
}

func ReadPrivateKeyFromFile(keyID string, keyDirPath string) ([]byte, error) {
	filePath := keyDirPath + "/" + keyID + ".pem"
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func ReadPublicKeyFromFile(keyID string, keyDirPath string) ([]byte, error) {
	filePath := keyDirPath + "/" + keyID + "_pub.pem"
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil
	}
	return key
}

func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil
	}
	return key
}

func EncryptWithPublicKey(data []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, data, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EncryptAESGCM(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func DecryptAESGCM(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:aesgcm.NonceSize()], data[aesgcm.NonceSize():]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
