package main

import (
	"dnssecpool/internal"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// 在这个代码里完成以下操作：
// 发起对正确域名的查询并记录响应
func main() {
	//从文件中读取解析器列表
	resolvers, err1 := internal.ReadResolversTXT("db/resolver_test.txt")
	fmt.Println(resolvers)
	//对解析器列表进行多线程DNS查询，同时记录收到查询的结果，记录ad字段、rcode、三部分答案
	if err1 != nil {
		log.Fatalf("Failed to read resolvers: %v", err1)
	}
	domain := "ftp.isc.org"
	filePath := "files/dns_responses.csv"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Resolver", "AD", "Response Code", "Answers"}
	// header := []string{"Resolver", "AD", "Response Code", "Answers", "Authority", "Additional"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error writing header to CSV: %v", err)
	}

	// Create a DNS client
	client := dns.Client{}

	for _, resolver := range resolvers {
		// Prepare the DNS message
		resolver += ":53"
		msg := dns.Msg{}
		// msg.set
		msg.SetQuestion(domain+".", dns.TypeA)
		msg.SetEdns0(4096, true)
		msg.RecursionDesired = true
		// Send the DNS query to the resolver
		response, _, err := client.Exchange(&msg, resolver)
		if err != nil {
			log.Printf("Query error for resolver %s: %v", resolver, err)
			continue
		}

		// Prepare CSV data for this resolver's response
		csvData := prepareCSVData(resolver, response)
		fmt.Println(csvData)
		// Write to CSV file
		if err := writer.Write(csvData); err != nil {
			log.Printf("Error writing CSV data for resolver %s: %v", resolver, err)
		}
	}

	fmt.Println("DNS responses written to", filePath)
}
func prepareCSVData(resolver string, response *dns.Msg) []string {
	ad := ""
	if response.MsgHdr.AuthenticatedData {
		ad = "1"
		fmt.Println("-1-1-1")
	} else {
		ad = "0"
	}
	if response.MsgHdr.Truncated {
		fmt.Println(000)
	}
	if response.MsgHdr.Authoritative {
		fmt.Println(111)
	}
	if response.MsgHdr.RecursionAvailable {
		fmt.Println(222)
	}
	if response.MsgHdr.RecursionDesired {
		fmt.Println(333)
	}
	var answers []string
	for _, answer := range response.Answer {
		// fmt.Println(answer.String())
		answers = append(answers, answer.String())
	}
	// fmt.Println(answers)

	// var authority []string
	// for _, ns := range response.Ns {
	// 	authority = append(authority, ns.String())
	// }

	// var additional []string
	// for _, extra := range response.Extra {
	// 	additional = append(additional, extra.String())
	// }
	return []string{
		resolver,
		ad,
		dns.RcodeToString[response.Rcode],
		strings.Join(answers, "||"),
		// strings.Join(authority, "||"),
		// strings.Join(additional, "||"),
	}
}
