package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"

	"os"
	"strconv"
	"strings"
)

var (
	kildeFil    string
	målFil      string
	nøkkelKilde string
	nøkkelMål   string
	verbose     bool

	tabellHode []string
)

func init() {

	flag.StringVar(&kildeFil, "s", "", "CSV kildefil for mapping.")
	flag.StringVar(&nøkkelKilde, "k1", "", "Kolonnernummer som utgjør nøkkel i Map (-k1=1,2,3)")
	flag.StringVar(&målFil, "t", "", "CSV målfil for mapping.")
	flag.StringVar(&nøkkelMål, "k2", "", "Kolonnernummer som utgjør nøkkel i Map(-k1=0,1,2)")
	flag.BoolVar(&verbose, "v", false, "Skriv ut debug info.")

	flag.Parse()
}

func main() {

	if validateFlags() == false {
		flag.Usage()
		os.Exit(0)
	}

	åpenKildeFil, _ := os.Open(kildeFil)
	åpenMålFil, _ := os.Open(målFil)
	defer åpenKildeFil.Close()
	defer åpenMålFil.Close()

	if verbose == true {
		fmt.Println("Startert sammenligning av")
		fmt.Println("kildefil:", kildeFil, ", nøkkelkolonner", nøkkelKilde)
		fmt.Println("målfil:", målFil, ", nøkkelkolonner", nøkkelMål)
	}

	diff, err := sammellignCsv(nøkkelKilde, åpenKildeFil, nøkkelMål, åpenMålFil)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if verbose == true {
		for i, rad := range diff {
			fmt.Printf("#%d %s finnes bare i %s\n", i+1, rad, kildeFil)
		}
	} else {
		w := csv.NewWriter(os.Stdout)
		w.Comma = ';'
		w.WriteAll(diff)
	}
}

func sammellignCsv(k1 string, kildeReader io.Reader, k2 string, målReader io.Reader) ([][]string, error) {

	nøkkelKildeInt, err := toIntArray(k1)
	if err != nil {
		return nil, err
	}

	nøkkelMålInt, err := toIntArray(k2)
	if err != nil {
		return nil, err
	}

	buffretKildeReader := bufio.NewReader(kildeReader)
	buffretMålReader := bufio.NewReader(målReader)

	tabellHode, err := hentTabellHode(buffretKildeReader)

	if err != nil {
		return nil, err
	}

	kildeTabell := hentCsvAsMap(kildeFil, nøkkelKildeInt, buffretKildeReader)
	destinasjonsTabell := hentCsvAsMap(målFil, nøkkelMålInt, buffretMålReader)

	diff := finnForskjellene(kildeTabell, destinasjonsTabell)

	diffTabell := append(diff, tabellHode)

	return reverse(diffTabell), nil
}

func toIntArray(nøkkelKollonner string) ([]int, error) {
	n := strings.Split(nøkkelKollonner, ",")
	var nøkkler []int

	for _, v := range n {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		nøkkler = append(nøkkler, num)
	}

	return nøkkler, nil
}

func hentCsvAsMap(file string, nøkkelKol []int, r io.Reader) map[string][]string {
	csvMap := make(map[string][]string)

	csvReader := csv.NewReader(r)

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		nøkkel := mapNøkkel(nøkkelKol, record)

		csvMap[nøkkel] = record
	}

	return csvMap
}

func mapNøkkel(nøkkelKol []int, record []string) string {
	var nøkkel bytes.Buffer

	for n := range nøkkelKol {
		nøkkel.WriteString(strings.TrimSpace(record[nøkkelKol[n]]))
	}

	return nøkkel.String()
}

func finnForskjellene(kildeTabell map[string][]string, destinasjonsTabell map[string][]string) [][]string {
	var diff [][]string

	for key, rad := range kildeTabell {
		if _, ok := destinasjonsTabell[key]; !ok {
			diff = append(diff, rad)
		}
	}

	return diff
}

func hentTabellHode(r io.Reader) ([]string, error) {
	csvReader := csv.NewReader(r)
	record, err := csvReader.Read()

	if err != nil {
		return nil, err
	}

	return record, nil
}

func reverse(t [][]string) [][]string {

	if len(t) > 0 {
		return append(reverse(t[1:]), t[0])
	}

	return t
}

func validateFlags() bool {

	if _, err := os.Stat(kildeFil); err != nil {
		fmt.Printf("Fant ikke kildefil %s\n", err.Error())
		return false
	}

	if _, err := os.Stat(målFil); err != nil {
		fmt.Printf("Fant ikke målfil %s\n", målFil)
		return false
	}

	if validerNøkler(nøkkelKilde, nøkkelMål) != true {
		fmt.Printf("Kildenøkkel %s eller målnøkkel %s er ikke gyldig\n", nøkkelKilde, nøkkelMål)
		return false
	}

	return true
}

func validerNøkler(k string, m string) bool {
	if len(m) == 0 {
		return false
	}

	if len(k) == 0 {
		return false
	}

	if len(k) != len(m) {
		return false
	}

	return true
}
