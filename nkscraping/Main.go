package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"bufio"
	"strings"
	"strconv"
	"bytes"
)

const (
	RaceInfoTable = ".race_table_01"
	Tr = "tr"
	Td = "td"
	H1 = "h1"
)

type RaceInfo struct {
	title     string
	orderInfo []OrderInfo
}

type OrderInfo struct {
	order        uint64
	postPosition uint64
	horseNumber  uint64
	horseName    string
}

func main() {
	doc, _ := getUrl2Doc("http://race.netkeiba.com/?pid=race&id=p201605021011")
	var raceInfo RaceInfo;
	parseRaceInfo(doc, &raceInfo)
	fmt.Println(raceInfo)
	fmt.Println(parseRaceOrder(doc))
}
func parseRaceInfo(doc *goquery.Document, raceinfo *RaceInfo) (err error) {
	text, err := scan(doc.Find(H1).Text())
	raceinfo.title = text
	return
}

func getUrl2Doc(url string) (doc *goquery.Document, err error) {
	return goquery.NewDocument(url)
}

func parseRaceOrder(doc *goquery.Document) (result []OrderInfo) {
	doc.Find(RaceInfoTable).Each(func(_ int, s *goquery.Selection) {
		s.Find(Tr).Each(func(trIndex int, s *goquery.Selection) {
			var info OrderInfo
			s.Find(Td).Each(func(tdIndex int, ss *goquery.Selection) {
				text, _ := scan(ss.Text())
				if trIndex != 0 {
					switch tdIndex {
					case 0:
						info.order, _ = strconv.ParseUint(text, 10, 0)
					case 1:
						info.postPosition, _ = strconv.ParseUint(text, 10, 0)
					case 2:
						info.horseNumber, _ = strconv.ParseUint(text, 10, 0)
					case 3:
						info.horseName = strings.Trim(text, " ")
						result = append(result, info)
					}
				}
			})
		})
	})
	return
}

func scan(text string) (string, error) {
	scanner := bufio.NewScanner(transform.NewReader(strings.NewReader(text), japanese.EUCJP.NewDecoder()))

	// TODO:byte
	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		//		result += scanner.Text()
	}

	return buffer.String(), scanner.Err()
}