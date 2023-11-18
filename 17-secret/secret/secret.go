package secret

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aboelkassem/gophercises/secret/encrypt"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Vault struct {
	encodingKey string
	keysPath    string
}

func FileVault(encodingKey, keysPath string) *Vault {
	return &Vault{encodingKey: encodingKey, keysPath: keysPath}
}

func (v *Vault) Set(keyName, KeyValue string) error {
	// 1. read keys file
	encryptedData, err := v.readFile()
	kv := map[string]string{}
	if err != nil {
		return err
	}

	// if the file not empty
	if encryptedData != "" {
		// 2. decrypt
		data, err := v.decryptFile(encryptedData)
		if err != nil {
			return err
		}
		kv = v.decodeFile(data)
	}

	// handle delete
	if KeyValue == "" {
		delete(kv, keyName)
	} else {
		// 3. set key
		kv[keyName] = KeyValue
	}
	data := v.encodeFile(kv)

	// 4. encrypt
	encryptedData, err = v.encryptFile(data)
	if err != nil {
		return err
	}

	// 5. write file
	return v.writeFile(encryptedData)
}

func (v *Vault) Get(keyName string) (string, error) {
	// 1. read keys file
	encryptedData, err := v.readFile()
	if err != nil {
		return "", err
	}

	if encryptedData == "" {
		return "", ErrKeyNotFound
	}

	// 2. decrypt
	data, err := v.decryptFile(encryptedData)
	if err != nil {
		return "", err
	}
	// 3. find key
	kv := v.decodeFile(data)
	// 4. return value
	value, ok := kv[keyName]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (v *Vault) List() ([]string, error) {
	// 1. read keys file
	encryptedData, err := v.readFile()
	if err != nil {
		return nil, err
	}

	if encryptedData == "" {
		return nil, nil
	}

	// 2. decrypt
	data, err := v.decryptFile(encryptedData)
	if err != nil {
		return nil, err
	}
	// 3. populate keys
	kv := v.decodeFile(data)
	var keys []string

	// loop over map, by default k is the key, if want value use _, v := range kv
	for k := range kv {
		keys = append(keys, k)
	}

	return keys, nil
}

func (v *Vault) Delete(keyName string) error {
	if _, err := v.Get(keyName); err != nil {
		return err
	}
	return v.Set(keyName, "")
}

// ReadCloser is mean to be closed
func (v *Vault) readFile() (string, error) {
	// create file if not exist
	f, err := os.OpenFile(v.keysPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (v *Vault) writeFile(data string) error {
	return ioutil.WriteFile(v.keysPath, []byte(data), 0600)
}

func (v *Vault) decryptFile(data string) (string, error) {
	return encrypt.Decrypt(v.encodingKey, data)
}

// if with different encoding key, its override
func (v *Vault) encryptFile(data string) (string, error) {
	return encrypt.Encrypt(v.encodingKey, data)
}

// read vault file as map of key value store
func (v *Vault) decodeFile(data string) map[string]string {
	kv := map[string]string{}

	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		line := strings.Split(sc.Text(), "=")
		if len(line) < 2 {
			continue
		}
		keyName, valueName := line[0], line[1]
		kv[keyName] = valueName
	}
	return kv
}

func (v *Vault) encodeFile(kv map[string]string) string {
	var data []string
	for key, value := range kv {
		data = append(data, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(data, "\n")
}
