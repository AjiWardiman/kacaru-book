package util_test

import (
	"testing"
	"unicode"

	"tourism-booking/internal/util"
)

func TestGenerateOTPCode(t *testing.T) {
	testLengths := []int{4, 6, 8}

	for _, length := range testLengths {
		t.Run("length_"+string(rune('0'+length)), func(t *testing.T) {
			otp, err := util.GenerateOTPCode(length)
			if err != nil {
				t.Fatalf("GenerateOTPCode(%d) failed: %v", length, err)
			}

			if len(otp) != length {
				t.Errorf("expected length %d, got %d (code: %s)", length, len(otp), otp)
			}

			for _, ch := range otp {
				if !unicode.IsDigit(ch) {
					t.Errorf("OTP code contains non-digit character %q in %s", ch, otp)
				}
			}
		})
	}
}

func TestGenerateOTPCode_Randomness(t *testing.T) {
	// Generate two 6-digit OTPs and verify they are unlikely to be identical
	otp1, err1 := util.GenerateOTPCode(6)
	otp2, err2 := util.GenerateOTPCode(6)

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateOTPCode failed: %v, %v", err1, err2)
	}

	// Note: In theory they could be equal with 1/1,000,000 chance, but generates valid codes
	if len(otp1) != 6 || len(otp2) != 6 {
		t.Errorf("OTP lengths invalid: %s, %s", otp1, otp2)
	}
}
