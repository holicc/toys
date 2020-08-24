package spider

import "log"

func init() {
	log.SetPrefix("[::Spider::]")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}
