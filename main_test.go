package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCase2(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "2")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria T2-euston
T1-st_pancras T2-st_pancras`

	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase3(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "3")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria
T3-st_pancras`
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase4(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras`
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase5(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "100")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras T5-victoria T6-euston
T5-st_pancras T6-st_pancras T7-victoria T8-euston
T7-st_pancras T8-st_pancras T9-victoria T10-euston`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase6(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria
T1-st_pancras`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase7(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase8(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/8.map", "bond_square", "space_port", "4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-apple_avenue
T1-orange_junction T2-apple_avenue
T1-space_port T2-orange_junction T3-apple_avenue
T2-space_port T3-orange_junction T4-apple_avenue
T3-space_port T4-orange_junction
T4-space_port`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase9(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/9.map", "jungle", "desert", "10")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-grasslands T2-farms T3-green_belt
T1-suburbs T2-downtown T3-village T4-grasslands T5-farms T6-green_belt
T1-clouds T2-metropolis T3-mountain T4-suburbs T5-downtown T6-village T7-grasslands T8-farms
T1-wetlands T2-industrial T3-treetop T4-clouds T5-metropolis T6-mountain T7-suburbs T8-downtown T9-grasslands T10-farms
T1-desert T2-desert T3-desert T4-wetlands T5-industrial T6-treetop T7-clouds T8-metropolis T9-suburbs T10-downtown
T4-desert T5-desert T6-desert T7-wetlands T8-industrial T9-clouds T10-metropolis
T7-desert T8-desert T9-wetlands T10-industrial
T9-desert T10-desert`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase10(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/10.map", "beginning", "terminus", "20")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-terminus T2-near
T2-far T3-terminus T4-near
T2-terminus T4-far T5-terminus T6-near
T4-terminus T6-far T7-terminus T8-near
T6-terminus T8-far T9-terminus T10-near
T8-terminus T10-far T11-terminus T12-near
T10-terminus T12-far T13-terminus T14-near
T12-terminus T14-far T15-terminus T16-near
T14-terminus T16-far T17-terminus T18-near
T16-terminus T18-far T19-terminus
T18-terminus T20-terminus`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase11(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/11.map", "two", "four", "4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-three
T1-one T2-three
T1-four T2-one T3-three
T2-four T3-one T4-three
T3-four T4-one
T4-four`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase12(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/12.map", "beethoven", "part", "9")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-verdi T2-handel
T1-part T2-mozart T3-verdi T4-handel
T2-part T3-part T4-mozart T5-verdi T6-handel
T4-part T5-part T6-mozart T7-verdi T8-handel
T6-part T7-part T8-mozart T9-verdi
T8-part T9-part`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase13(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/13.map", "small", "large", "9")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	expected := `T1-10 T2-13 T3-00
T1-11 T2-14 T3-01 T4-10 T5-13
T1-12 T2-15 T3-02 T4-11 T5-14 T6-10 T7-13
T1-large T2-21 T3-03 T4-12 T5-15 T6-11 T7-14 T8-10
T2-22 T3-04 T4-large T5-21 T6-12 T7-15 T8-11 T9-10
T2-large T3-05 T5-22 T6-large T7-21 T8-12 T9-11
T3-large T5-large T7-22 T8-large T9-12
T7-large T9-large`
	
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected output not found. Got: %s", output)
	}
}

func TestCase14(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "a", "b")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: incorrect number of command line arguments"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase15(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "a", "b", "1", "hello")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: incorrect number of command line arguments"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase16(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/big.map", "station0", "station9999", "100")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Errorf("Expected some output, but got empty output")
	}
}

func TestCase17(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "hive", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: start station \"hive\" does not exist"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase18(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "hive", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: end station \"hive\" does not exist"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase19(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "waterloo", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: start and end station are the same"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase20(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/nopath.map", "waterloo", "central", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: no path exists between \"waterloo\" and \"central\""
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase21(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/dup.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: duplicate connection between \"euston\" and \"waterloo\" (line 15)"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase22(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/london.map", "waterloo", "st_pancras", "-2")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: number of trains (-2) is not a valid positive integer"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase23(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/mixed.map", "two", "four", "5")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: station \"four\" has a coordinate (9, -2) which is not a valid positive integer"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase24(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/dupcoor.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: stations \"waterloo\" and \"euston\" exist at the same coordinates"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase25(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/25.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: connection on line 11 refers to station \"hive\" which does not exist"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase26(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/26.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: duplicate station name \"euston\" (line 8)"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase27(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/27.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: invalid station name on line 7: \"euston-my-love\""
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase28(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/28.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: the map does not contain a \"stations:\" section"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase29(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/29.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: the map does not contain a \"connections:\" section"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}

func TestCase30(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "maps/30.map", "waterloo", "st_pancras", "4")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded with output: %s", output)
	}

	expectedError := "Error: the map contains more than 10000 stations"
	if !strings.Contains(string(output), expectedError) {
		t.Errorf("Expected error not found. Got: %s", output)
	}
}
