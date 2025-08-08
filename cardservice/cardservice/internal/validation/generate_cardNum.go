package validation




import (
    "crypto/rand"
    "math/big"
)

func GenerateCardNumber() (string, error) {
    var result [16]byte

    // First digit: 1–9 (no leading zero)
    d, err := rand.Int(rand.Reader, big.NewInt(9))
    if err != nil {
        return "", err
    }
    result[0] = byte(d.Int64() + 1 + '0')

    // Remaining 15 digits: 0–9
    for i := 1; i < 16; i++ {
        d, err = rand.Int(rand.Reader, big.NewInt(10))
        if err != nil {
            return "", err
        }
        result[i] = byte(d.Int64() + '0')
    }

    return string(result[:]), nil
}
