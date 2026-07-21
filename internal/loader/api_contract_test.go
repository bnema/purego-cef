package loader

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestValidateAPIHash(t *testing.T) {
	const expectedLinuxHash = "210767725a6feb2e4becd3956b648cab6a006712"
	if err := validateAPIHash(expectedLinuxHash); err != nil {
		t.Fatal(err)
	}

	hash := append([]byte(expectedLinuxHash), 0)
	if err := validateAPIHashPointer(uintptr(unsafe.Pointer(&hash[0]))); err != nil {
		t.Fatalf("valid hash pointer rejected: %v", err)
	}
	runtime.KeepAlive(hash)

	if err := validateAPIHashPointer(0); err == nil {
		t.Fatal("nil hash accepted")
	}
	if err := validateAPIHash("wrong"); err == nil {
		t.Fatal("mismatched hash accepted")
	}
}

func TestRuntimeContractIsCEF150(t *testing.T) {
	if defaultCEFVersion != 150 || cefAPIVersion != 15000 {
		t.Fatalf("major=%d api=%d", defaultCEFVersion, cefAPIVersion)
	}
}
