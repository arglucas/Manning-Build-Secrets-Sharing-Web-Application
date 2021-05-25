package handlers

import "testing"

func TestMD5(t *testing.T) {
	gotOutput := md5hex("My super secret123")

	expectedOutput := "c616584ac64a93aafe1c16b6620f5bcd"
	if gotOutput != expectedOutput {
		t.Fatalf("Expected: %s, Got: %s\n", expectedOutput, gotOutput)
	}
}