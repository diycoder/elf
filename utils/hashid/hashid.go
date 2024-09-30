package hashid

import (
	"errors"

	"github.com/speps/go-hashids/v2"
)

const (
	salt      = "~!@#*&^%$"                            // default salt
	alphabet  = "abcdefghijklmnopqrstuvwxyz0123456789" // default alphabet
	minLength = 8                                      // default min length
)

func Encode(number int64) (string, error) {
	var (
		err    error
		result string
		hashId *hashids.HashID
	)
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hd.Alphabet = alphabet

	hashId, err = hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	result, err = hashId.EncodeInt64([]int64{number})
	if err != nil {
		return "", err
	}
	return result, nil
}

func Decode(hash string) (int64, error) {
	if hash == "" {
		return 0, errors.New("empty hash")
	}
	var (
		err    error
		number int64
		result []int64
		hashId *hashids.HashID
	)
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hd.Alphabet = alphabet

	hashId, err = hashids.NewWithData(hd)
	if err != nil {
		return 0, err
	}
	result, err = hashId.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}
	if len(result) > 0 {
		number = result[0]
	}
	return number, nil
}

func EncodeBySalt(number int64, salt string) (string, error) {
	var (
		err    error
		result string
		hashId *hashids.HashID
	)
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hd.Alphabet = alphabet

	hashId, err = hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	result, err = hashId.EncodeInt64([]int64{number})
	if err != nil {
		return "", err
	}
	return result, nil
}

func DecodeBySalt(hash string, salt string) (int64, error) {
	if hash == "" {
		return 0, errors.New("empty hash")
	}
	var (
		err    error
		number int64
		result []int64
		hashId *hashids.HashID
	)
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hd.Alphabet = alphabet

	hashId, err = hashids.NewWithData(hd)
	if err != nil {
		return 0, err
	}
	result, err = hashId.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}
	if len(result) > 0 {
		number = result[0]
	}
	return number, nil
}
