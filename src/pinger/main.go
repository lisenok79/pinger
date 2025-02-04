package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// pingIP выполняет ping для заданного IP и выводит результат.
func pingIP(ip string, wg *sync.WaitGroup) {
	defer wg.Done()

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Printf("Ошибка создания пингера для %s: %v\n", ip, err)
		return
	}

	// Устанавливаем режим привилегированного пинга (может потребоваться root-права)
	pinger.SetPrivileged(true)
	pinger.Count = 3                 // Количество ICMP-запросов
	pinger.Timeout = 3 * time.Second // Общий таймаут пинга
	
	err = pinger.Run()
	if err != nil {
		log.Printf("Ошибка при пинге %s: %v\n", ip, err)
		return
	}
	
	timeStamp := time.Now().Format(time.RFC1123)
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		fmt.Printf("IP %s доступен (%d/%d пакетов получено)\n", ip, stats.PacketsRecv, stats.PacketsSent)
		log.Println(timeStamp)
	} else {
		fmt.Printf("IP %s недоступен\n", ip)
		log.Println(timeStamp)
	}
}

func main() {
	Ips := []string{}
	file, err := os.Open("ips.json")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &Ips)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	for _, ip := range Ips {
		wg.Add(1)
		go pingIP(ip, &wg)
	}

	wg.Wait()
}
