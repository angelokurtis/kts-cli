package git

import (
	"testing"
)

func TestNewLocalDir(t *testing.T) {
	dir, err := NewLocalDir("https://github.com/angelokurtis/hellognome")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	address := dir.SSHAddress()
	if address != "git@github.com:angelokurtis/hellognome.git" {
		t.Errorf("Expect that `SSHAddress()` return correctly but got %s", address)
	}
}

func TestNewLocalDir2(t *testing.T) {
	dir, err := NewLocalDir("https://www.github.com/angelokurtis/hellognome")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	address := dir.SSHAddress()
	if address != "git@github.com:angelokurtis/hellognome.git" {
		t.Errorf("Expect that `SSHAddress()` return correctly but got %s", address)
	}
}
