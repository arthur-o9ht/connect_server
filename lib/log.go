package lib

import "log"

func Warn(msg string) {
	log.Println("[Warn]:" + msg)
}

func Error(msg string) {
	log.Fatalln("[Error]" + msg)
}

func Info(msg string) {
	log.Println("[Info]:" + msg)
}
