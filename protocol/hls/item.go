package hls

import "time"

type TSItem struct {
	Name     string
	SeqNum   int
	Duration int
	Create   time.Time
	Data     []byte
}

func NewTSItem(name string, duration, seqNum int, b []byte) TSItem {
	var item TSItem
	item.Name = name
	item.SeqNum = seqNum
	item.Duration = duration
	item.Create = time.Now()
	item.Data = make([]byte, len(b))
	copy(item.Data, b)
	return item
}
