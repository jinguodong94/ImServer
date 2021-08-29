package mq

import "fmt"

//mq消息接收
func Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}
