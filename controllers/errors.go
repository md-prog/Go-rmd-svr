package Controllers

import log "github.com/Sirupsen/logrus"

func handleErr(err error) {
	if err != nil {
		log.Warn(err)
	}
	return
}
