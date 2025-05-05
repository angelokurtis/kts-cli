package gpg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewSecretKeys(t *testing.T) {
	const gpgOutput = `
sec   rsa4096/ABC123456789DEF0 2020-01-01 [SC] [expires: 2022-01-01]
      1234ABCD5678EFGH1234IJKL5678MNOP1234QRST
uid           [ultimate] John Doe <john@example.com>
ssb   rsa4096/1234567890ABCDEF 2020-01-01 [E]

sec   rsa2048/DEF0987654321CBA 2021-02-02 [SC]
      9876ZYXW5432VUTS9876RQPO5432NMKL9876JIHG
uid           [unknown] Jane Unknown <unknown@example.com>
ssb   rsa2048/0FEDCBA987654321 2021-02-02 [E]
`

	expires := "2022-01-01"

	want := SecretKeys{
		{
			KeyType: "rsa4096",
			KeyID:   "ABC123456789DEF0",
			Created: "2020-01-01",
			Expires: &expires,
			Uids: []*Uid{
				{Name: "John Doe", Email: "john@example.com", Trust: "ultimate"},
			},
			Subkeys: []*Subkey{
				{KeyType: "rsa4096", KeyID: "1234567890ABCDEF", Created: "2020-01-01", Usage: []string{"E"}},
			},
		},
		{
			KeyType: "rsa2048",
			KeyID:   "DEF0987654321CBA",
			Created: "2021-02-02",
			Expires: nil,
			Uids: []*Uid{
				{Name: "Jane Unknown", Email: "unknown@example.com", Trust: "unknown"},
			},
			Subkeys: []*Subkey{
				{KeyType: "rsa2048", KeyID: "0FEDCBA987654321", Created: "2021-02-02", Usage: []string{"E"}},
			},
		},
	}

	got, err := NewSecretKeys([]byte(gpgOutput))
	if err != nil {
		t.Fatalf("NewSecretKeys() error = %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewSecretKeys() mismatch (-want +got):\n%s", diff)
	}
}
