package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())

}

func TestSkalMappeToTabeller(t *testing.T) {
	nøkkelKilde = "2,3,4"
	nøkkelMål = "0,1,2"

	kildeData := `
			KOL1, KOL2, KOL3, KOL4, KOL5, KOL6
			1, 2, 3, 4, 5, 6
			2, 3, 4, 5, 6, 7
			3, 4, 5, 6, 7, 8
			`

	målData := `
			KOL3, KOL4, KOL5, KOL6
			3, 4, 5, 6
			4, 5, 6, 7
			5, 6, 7, 8
			`

	diffTabell, err := sammellignCsv(nøkkelKilde, strings.NewReader(kildeData), nøkkelMål, strings.NewReader(målData))

	if err != nil {
		t.Error(err.Error())
	}

	if len(diffTabell) != 1 {
		t.Errorf("Forventet ingen diff")
	}

	if strings.TrimSpace(diffTabell[0][0]) != "KOL1" {
		t.Errorf("Forventet KOL1, fikk %s", diffTabell[0][0])
	}
}

func TestSkalFinne2RaderSomIkkeMappes(t *testing.T) {
	nøkkelKilde = "2,3,4"
	nøkkelMål = "0,1,2"

	kildeData := `
			KOL1, KOL2, KOL3, KOL4, KOL5, KOL6
			1, 2, 3, 4, 5, 6
			2, 3, 4, 5, 6, 7
			3, 4, 5, 6, 7, 8
			`

	målData := `
			KOL3, KOL4, KOL5, KOL6
			3, 4, 5, 6			
			`

	diffTabell, err := sammellignCsv(nøkkelKilde, strings.NewReader(kildeData), nøkkelMål, strings.NewReader(målData))

	if err != nil {
		t.Error(err.Error())
	}

	if len(diffTabell) != 3 {
		t.Errorf("Forventet 3 rader i tabell, men fikk %d rad.", len(diffTabell))
		t.Fail()
	}

	kol1FørsteMangel := strings.TrimSpace(diffTabell[1][0])
	kol1AndreMangel := strings.TrimSpace(diffTabell[2][0])

	if kol1FørsteMangel != "3" {
		t.Errorf("Forventet %s, fant %s", "3", kol1FørsteMangel)
	}

	if kol1AndreMangel != "2" {
		t.Errorf("Forventet %s, fant %s", "2", kol1AndreMangel)
	}
}
