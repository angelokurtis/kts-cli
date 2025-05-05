package gpg

import (
	"testing"
)

func TestNewKeys(t *testing.T) {
	// Simulated GPG output for testing
	const gpgOutput = `
sec   rsa4096/ABC123456789DEF0 2020-01-01 [SC] [expires: 2022-01-01]
      1234ABCD5678EFGH1234IJKL5678MNOP1234QRST
uid           [ultimate] John Doe <john@example.com>
ssb   rsa4096/1234567890ABCDEF 2020-01-01 [E]

sec   rsa2048/DEF0987654321CBA 2021-02-02 [SC]
      9876ZYXW5432VUTS9876RQPO5432NMKL9876JIHG
uid           [ultimate] Jane Smith <jane@example.com>
ssb   rsa2048/0FEDCBA987654321 2021-02-02 [E]

sec   rsa2048/DEF0987654321CBA 2021-02-02 [SC]
      9876ZYXW5432VUTS9876RQPO5432NMKL9876JIHG
uid           [unknown] Jane Unknown <unknown@example.com>
ssb   rsa2048/0FEDCBA987654321 2021-02-02 [E]
`

	keys, err := NewKeys([]byte(gpgOutput))
	if err != nil {
		t.Fatalf("NewKeys() error = %v", err)
	}

	if keys == nil || len(keys.Items) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys.Items))
	}

	expected := []struct {
		Sec string
		UID string
	}{
		{"rsa4096/ABC123456789DEF0", "John Doe <john@example.com>"},
		{"rsa2048/DEF0987654321CBA", "Jane Smith <jane@example.com>"},
		{"rsa2048/DEF0987654321CBA", "Jane Unknown <unknown@example.com>"},
	}

	for i, key := range keys.Items {
		if key.Sec != expected[i].Sec {
			t.Errorf("key[%d] Sec = %q, want %q", i, key.Sec, expected[i].Sec)
		}

		if key.UID != expected[i].UID {
			t.Errorf("key[%d] UID = %q, want %q", i, key.UID, expected[i].UID)
		}
	}
}
