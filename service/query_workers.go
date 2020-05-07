package service

func (z *zone) mapQueryWorkers() {
	for {
		select {
		case eq := <- z.queries[queryMap]:
			log.Info(eq)
		}
	}
}

func (z *zone) playerQueryWorkers() {
	for {
		select {
		case eq := <- z.queries[queryPlayer]:
			log.Info(eq)
		}
	}
}

func (z *zone) monsterQueryWorkers() {
	for {
		select {
		case eq := <- z.queries[queryMonster]:
			log.Info(eq)
		}
	}
}