package internal

import (
	"github.com/miekg/dns"
)

// 使用 DNS 解析器查询指定的域名
func QueryDNS(resolver, domain string) (dns.Msg, error) {
	client := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA) // 查询 A 记录
	msg.SetEdns0(4096, true)
	msg.RecursionDesired = true

	resp, _, err := client.Exchange(msg, resolver+":53")
	if err != nil {
		return dns.Msg{}, err
	}
	return *resp, err
}
func PrepareCSVData(resolver string, response *dns.Msg) []string {
	ad := ""
	if response.MsgHdr.AuthenticatedData {
		ad = "1"
		// fmt.Println("-1-1-1")
	} else {
		ad = "0"
	}
	// if response.MsgHdr.Truncated {
	// 	fmt.Println(000)
	// }
	// if response.MsgHdr.Authoritative {
	// 	fmt.Println(111)
	// }
	// if response.MsgHdr.RecursionAvailable {
	// 	fmt.Println(222)
	// }
	// if response.MsgHdr.RecursionDesired {
	// 	fmt.Println(333)
	// }
	// var answers []string
	// for _, answer := range response.Answer {
	// 	// fmt.Println(answer.String())
	// 	answers = append(answers, answer.String())
	// }
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
		// strings.Join(answers, "||"),
		// strings.Join(authority, "||"),
		// strings.Join(additional, "||"),
	}
}
