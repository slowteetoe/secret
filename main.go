package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

var (
	passphrase = "6368616e676520746869732070617373"
)

type Vault struct {
	secret   string
	filename string
	data     map[string]string
}

func NewVault(passphrase, filename string) (*Vault, error) {

	f, err := os.Open(filename)
	defer f.Close()
	var data map[string]string
	if err != nil {
		// ignore the vault doesn't have any data yet
		data = make(map[string]string)
	} else {
		// decrypt and read map

		block, err := aes.NewCipher([]byte(passphrase))
		if err != nil {
			return nil, errors.Wrap(err, "failed getting cipher")
		}

		// If the key is unique for each ciphertext, then it's ok to use a zero
		// IV.
		var iv [aes.BlockSize]byte
		stream := cipher.NewOFB(block, iv[:])

		reader := &cipher.StreamReader{S: stream, R: f}

		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&data)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("vault data is: %+v", data)

	return &Vault{secret: passphrase, filename: filename, data: data}, nil
}

func (v *Vault) save() error {

	f, err := os.Create(v.filename)
	if err != nil {
		return errors.Wrapf(err, "failed to write encrypted data to %s", v.filename)
	}
	defer f.Close()

	block, err := aes.NewCipher([]byte(v.secret))
	if err != nil {
		return errors.Wrap(err, "failed to obtain cipher")
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: f}

	err = json.NewEncoder(writer).Encode(v.data)

	if err != nil {
		return errors.Wrap(err, "failed to write marshalled json")
	}
	return nil
}

func (v *Vault) Set(key, value string) error {
	v.data[key] = value
	return v.save()
}

func (v *Vault) Get(key string) (string, error) {
	val, ok := v.data[key]
	if !ok {
		return "", fmt.Errorf("no value for key: %s", key)
	}
	return val, nil
}

func main() {

	v, err := NewVault(passphrase, "vault.txt")

	if err != nil {
		log.Fatalf("failed to open vault: %v", err)
	}

	key := "some-key"

	err = v.Set(key, "some-value")

	if err != nil {
		log.Fatalf("failed to set key: %s", key)
	}

	val, err := v.Get(key)

	if err != nil {
		log.Fatalf("failed to find key: %s", key)
	}

	log.Printf("value is: %s", val)
}
