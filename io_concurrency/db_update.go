package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var filename = flag.String("f", "", "file to read")
var database = flag.String("d", "", "database name")
var table = flag.String("t", "", "table name")
var column = flag.String("c", "", "column name")
var where = flag.String("w", "", "where condition")
var updateValue = flag.String("v", "", "column value")
var output = flag.String("o", "", "output filename")
var ch = make(chan string, 8)

func main() {
	flag.Parse()

	go scan()
	go generateSql()

	var a string
	fmt.Scanln(&a)
	//fmt.Println("hello world")
}

func scan() {
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		ch <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func generateSql() {
	count := 0
	defer close(ch)
	f, err := os.Create(*output)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	var sql string
	var sidsJoinStr string
	var sids []string
	for {
		count++
		s := <-ch
		if s != "" {
			sids = append(sids, s)
			//fmt.Printf("%+v\n", sids)
			if count%100 == 0 {
				sidsJoinStr = strings.Join(sids[:], ",")
				sql = fmt.Sprintf("update %s.%s set %s = '%s' where %s in(%s);\n", *database, *table, *column, *updateValue, *where, sidsJoinStr)
				sql = "select sleep(1);\n" + sql
				sids = []string{}
				_, err := w.WriteString(sql)
				if err != nil {
					panic(err)
				}
				w.Flush()
			}
		}
	}
}
