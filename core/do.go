package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type clockDay struct {
	Do  []string
	Not []string
}

//go build -o clockinbin main.go
func Clockinbin() {

	var (
		err   error
		env   string
		force bool
	)

	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	if len(os.Args) > 2 {
		force = os.Args[2] == "true"
	}

	if !force {
		//非强制打卡，走这两个流程
		if !shouldClock() {
			log.Println("今日不打卡")
			return
		}

		if env == "pro" {
			randSleep()
		}
	}

	if err = Run(env); err != nil {
		log.Println(err)
	}
}

func loadConfig() (res clockDay) {
	var (
		c   []byte
		err error
	)
	if c, err = ioutil.ReadFile("data/clockday.json"); err != nil {
		log.Println(err)
		return
	}

	if len(c) > 0 {
		if err = json.Unmarshal(c, &res); err != nil {
			log.Println(err)
		}
	}

	return

}

func in(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

//判定是否应该打卡
func shouldClock() bool {
	var (
		today string
		day   clockDay
	)

	today = time.Now().Format("2006-01-02")

	day = loadConfig()

	if in(today, day.Do) {
		return true
	}

	weekDay := time.Now().Weekday()
	if weekDay != time.Sunday && weekDay != time.Saturday {
		//当前是周一到周五
		if !in(today, day.Not) {
			return true
		}
	}

	return false
}

//随机睡眠
func randSleep() {
	rand.Seed(time.Now().Unix())
	randM := 0
	morning := time.Now().Hour() < 12

	if morning {
		randM = rand.Intn(10)
	} else {
		randM = rand.Intn(1)
	}

	log.Println("随机睡眠时间:", randM)

	time.Sleep(time.Duration(randM) * time.Minute)
}
