package version

import "testing"

func TestPrint(t *testing.T) {
	Version = "testVersion"
	Branch = "testBranch"
	Revision = "testRevision"
	BuildUser = "testUser"
	BuildDate = "testDate"

	expected := "bpt, version testVersion (branch: testBranch, revision: testRevision)\n  build user:       testUser\n  build date:       testDate\n  go version:       go1.11.3"
	result := Print()
	if result != expected {
		t.Errorf("Print: %q, want %q", result, expected)
	}
}

func TestNewInfo(t *testing.T) {
	Version = "testVersion"
	Branch = "testBranch"
	Revision = "testRevision"
	BuildUser = "testUser"
	BuildDate = "testDate"

	info := NewInfo()
	if info.BuildDate != BuildDate {
		t.Errorf("BuildDate: %q, want %q", info.BuildDate, BuildDate)
	}
	if info.BuildUser != BuildUser {
		t.Errorf("BuildUser: %q, want %q", info.BuildUser, BuildUser)
	}
	if info.Revision != Revision {
		t.Errorf("Revision: %q, want %q", info.Revision, Revision)
	}
	if info.Branch != Branch {
		t.Errorf("Branch: %q, want %q", info.Branch, Branch)
	}
	if info.Version != Version {
		t.Errorf("Version: %q, want %q", info.Version, Version)
	}

}
