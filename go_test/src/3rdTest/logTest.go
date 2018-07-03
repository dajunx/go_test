package officialLibTest

import "log"

func TestSystemLog() {
	log.Printf("message %d\n", 111)
	log.Fatalln("message")
}