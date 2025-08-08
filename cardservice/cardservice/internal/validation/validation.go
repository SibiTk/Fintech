package validation


import (
   // "errors"
    "strconv"
    "strings"
    "time"
)


type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return e.Field + ": " + e.Message
}


// // ValidateAccountID validates account ID
// func ValidateAccountID(accountID int64) error {
//     if accountID <= 0 {
//         return ValidationError{
//             Field:   "account_id",
//             Message: "account ID must be positive",
//         }
//     }
//     return nil
// }



//ValidateCardType validates card type


func ValidateCardType(cardType string) error {
    validTypes := map[string]bool{
        "credit":  true,
        "debit":   true,
        "prepaid": true,
    }

    if cardType == "" {
        return ValidationError{
            Field:   "card_type",
            Message: "card type is required",
        }
    }

    if !validTypes[strings.ToLower(cardType)] {
        return ValidationError{
            Field:   "card_type",
            Message: "invalid card type. Must be credit, debit, prepaid",
        }
    }

    return nil
}


// ValidateExpiryDate validates expiry date format without regex
func ValidateExpiryDate(expiryDate string) error {
    if expiryDate == "" {
        return ValidationError{
            Field:   "expiry_date",
            Message: "expiry date is required",
        }
    }
    
    // Check if contains exactly one slash
    parts := strings.Split(expiryDate, "/")
    if len(parts) != 2 {
        return ValidationError{
            Field:   "expiry_date",
            Message: "expiry date must be in MM/YY or MM/YYYY format",
        }
    }
    
    monthStr := parts[0]
    yearStr := parts[1]
    
    // Validate month part
    if len(monthStr) != 2 {
        return ValidationError{
            Field:   "expiry_date",
            Message: "month must be 2 digits (MM)",
        }
    }
    
    // Check if month contains only digits
    for _, char := range monthStr {
        if char < '0' || char > '9' {
            return ValidationError{
                Field:   "expiry_date",
                Message: "month must contain only digits",
            }
        }
    }
    
    month, err := strconv.Atoi(monthStr)
    if err != nil || month < 1 || month > 12 {
        return ValidationError{
            Field:   "expiry_date",
            Message: "month must be between 01 and 12",
        }
    }
    
    // Validate year part
    if len(yearStr) != 2 && len(yearStr) != 4 {
        return ValidationError{
            Field:   "expiry_date",
            Message: "year must be 2 digits (YY) or 4 digits (YYYY)",
        }
    }
    
    // Check if year contains only digits
    for _, char := range yearStr {
        if char < '0' || char > '9' {
            return ValidationError{
                Field:   "expiry_date",
                Message: "year must contain only digits",
            }
        }
    }
    
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        return ValidationError{
            Field:   "expiry_date",
            Message: "invalid year format",
        }
    }
    
    // Convert 2-digit year to 4-digit
    if year < 100 {
        if year < 50 {
            year += 2000
        } else {
            year += 1900
        }
    }
    
    // Validate year range
    currentYear := time.Now().Year()
    if year < currentYear || year > currentYear+50 {
        return ValidationError{
            Field:   "expiry_date",
            Message: "year must be reasonable (current year to 50 years in future)",
        }
    }
    
    // Check if date is in the future
    expiryTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    currentTime := time.Now()
    currentMonthStart := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, time.UTC)
    
    if expiryTime.Before(currentMonthStart) {
        return ValidationError{
            Field:   "expiry_date",
            Message: "expiry date must be in the future",
        }
    }
    
    return nil
}

// ValidateLimits validates daily and monthly limits
func ValidateLimits(dailyLimit, monthlyLimit float64) error {
    if dailyLimit < 0 {
        return ValidationError{
            Field:   "daily_limit",
            Message: "daily limit cannot be negative",
        }
    }
    
    if monthlyLimit < 0 {
        return ValidationError{
            Field:   "monthly_limit",
            Message: "monthly limit cannot be negative",
        }
    }
    
    if dailyLimit > monthlyLimit {
        return ValidationError{
            Field:   "limits",
            Message: "daily limit cannot exceed monthly limit",
        }
    }
    
    // Additional business rule validations
    if dailyLimit > 50000 {
        return ValidationError{
            Field:   "daily_limit",
            Message: "daily limit cannot exceed 50,000",
        }
    }
    
    if monthlyLimit > 1000000 {
        return ValidationError{
            Field:   "monthly_limit",
            Message: "monthly limit cannot exceed 1,000,000",
        }
    }
    
    return nil
}

// ValidatePinAttempts validates PIN attempts
func ValidatePinAttempts(attempts int) error {
    if attempts < 0 || attempts > 5 {
        return ValidationError{
            Field:   "pin_attempts",
            Message: "PIN attempts must be between 0 and 5",
        }
    }
    return nil
}

// ValidateCreateCardRequest validates all fields for card creation
func ValidateCreateCardRequest( accountID int64, cardType, expiryDate string, 
    dailyLimit, monthlyLimit float64, pinAttempts int) error {
    
   
    // if err := ValidateAccountID(accountID); err != nil {
    //     return err
    // }
    
   
    
    if err := ValidateCardType(cardType); err != nil {
        return err
    }
    
    if err := ValidateExpiryDate(expiryDate); err != nil {
        return err
    }
    
    if err := ValidateLimits(dailyLimit, monthlyLimit); err != nil {
        return err
    }
    
    if err := ValidatePinAttempts(pinAttempts); err != nil {
        return err
    }
    
    return nil
}

// ValidateUpdateCardRequest validates fields for card update (no expiry date change)
func ValidateUpdateCardRequest(accountID int64, cardNumber, cardType string, 
    dailyLimit, monthlyLimit float64) error {
    
   
    
    // if err := ValidateAccountID(accountID); err != nil {
    //     return err
    // }
   
    // if err := ValidateCardType(cardType); err != nil {
    //     return err
    // }
    
    if err := ValidateLimits(dailyLimit, monthlyLimit); err != nil {
        return err
    }
    
    return nil
}
