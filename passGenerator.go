package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func IsValid(length int, number, upperCase, lowerCase, specialCase bool) (bool, error) {
	if length <= 0 {
		return false, fmt.Errorf("length of password is zero or less")
	}
	if length > 500 {
		return false, fmt.Errorf("password is too long")
	}
	if number == false && upperCase == false && lowerCase == false && specialCase == false {
		return false, fmt.Errorf("password params if false")
	}
	return true, nil
}

// checker: return symbol "Check" if true, "Cross" if false
func checker(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

// GeneratePassword : Генератор паролей; length - длина генерируемого пароля
// numbers - наличие цифр в пароле, upperCase - наличие букв верхнего регистра в пароле,
// lowerCase - наличие букв нижнего регистра в пароле, specialCase - спец.символов в пароле
func GeneratePassword(length int, number, upperCase, lowerCase, specialCase bool) (string, error) {
	if ok, err := IsValid(length, number, upperCase, lowerCase, specialCase); !ok {
		return "", err
	}
	kit := ""
	if number {
		kit += "0123456789"
	}
	if lowerCase {
		kit += "abcdefghijklmnopqrstuvwxyz"
	}
	if upperCase {
		kit += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if specialCase {
		kit += ",.%*()?@#$~"
	}
	res := make([]byte, length)
	for i := range res {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(kit))))
		if err != nil {
			log.Panic(err)
		}
		res[i] = kit[r.Int64()]
	}
	return string(res), nil
}
