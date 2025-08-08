package handler

import (
    "crypto/rand"
    "math/big"
    "strconv"
)

func GenerateAccountNumber() (int64, error) {
    var result [16]byte

    // First digit: 1–9
    d, err := rand.Int(rand.Reader, big.NewInt(9))
    if err != nil {
        return 0, err
    }
    result[0] = byte(d.Int64() + 1 + '0')

    // Remaining 15 digits: 0–9
    for i := 1; i < 16; i++ {
        d, err = rand.Int(rand.Reader, big.NewInt(10))
        if err != nil {
            return 0, err
        }
        result[i] = byte(d.Int64() + '0')
    }

    // Convert the byte slice to a string
    accountNumberStr := string(result[:])

    // Convert the string to int64
    accountNumber, err := strconv.ParseInt(accountNumberStr, 10, 64)
    if err != nil {
        return 0, err
    }

    return accountNumber, nil
}

