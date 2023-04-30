package sti2023

import (
	"testing"
)

func TestMail(t *testing.T) {
	mailConfigDir = "test-data/"
	mailConfigFile = "mail"	

	if got := Mail("michal.kukla@tul.cz", "test"); got {
		t.Errorf("Expected '%t' but, got '%t'", false, got)
	}
}
