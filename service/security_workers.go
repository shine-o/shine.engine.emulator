package service

func (z *zone) loadSecurityWorker() {
	for {
		select{
		case e := <- z.recv[clientSHN]:
			log.Info(e)
		}
	}
}