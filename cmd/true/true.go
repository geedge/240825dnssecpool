package main

import (
	"dnssecpool/internal"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	domainToQuery = "ftp.isc.org"            // 替换为你想查询的域名
	outputFile    = "files/resolvers_ad.csv" // 输出的 CSV 文件
	numWorkers    = 50                       // 并发的 goroutine 数量
)

func main() {
	start := time.Now()

	resolvers, err1 := internal.ReadResolversTXT("db/stable_open_resolvers_240726.txt")
	if err1 != nil {
		log.Fatalf("Failed to read resolvers: %v", err1)
	}

	// 创建 CSV 文件
	csvFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	var wg sync.WaitGroup
	jobs := make(chan string, len(resolvers))
	results := make(chan []string, len(resolvers))

	// 启动 worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for resolver := range jobs {
				response, err := internal.QueryDNS(resolver, domainToQuery)
				if err != nil {
					// results <- [][]string{[][]string, err.Error()}
					// results <- []string{resolver, "err", "err", "err"}

					// results <- []string{resolver, "err", "err"}
				} else {
					csvData := internal.PrepareCSVData(resolver, &response)
					if csvData[1] == "1" {
						results <- csvData
					}
				}
			}
		}()
	}

	// 发送解析器到 jobs 通道
	for _, resolver := range resolvers {
		jobs <- resolver
	}
	close(jobs)

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(results)
	}()
	// header := []string{"Resolver", "AD", "Response Code"}
	// // header := []string{"Resolver", "AD", "Response Code", "Answers"}
	// if err := csvWriter.Write(header); err != nil {
	// 	log.Printf("Failed to write to CSV: %v", err)
	// }
	// 处理结果并写入 CSV 文件
	for result := range results {
		if err := csvWriter.Write(result); err != nil {
			log.Printf("Failed to write to CSV: %v", err)
		}
	}

	// 记录程序结束时间
	end := time.Now()
	// 计算和输出执行时间
	elapsed := end.Sub(start)
	fmt.Printf("程序执行时间: %s\n", elapsed)
}
