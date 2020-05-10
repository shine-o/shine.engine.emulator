package service

func (z *zone) mapQueries() {
	for {
		select {
		case eq := <-z.queries[queryMap]:
			log.Info(eq)
		}
	}
}

func (z *zone) playerQueries() {
	for {
		select {
		case eq := <-z.queries[queryPlayer]:
			log.Info(eq)
		}
	}
}

func (z *zone) monsterQueries() {
	for {
		select {
		case eq := <-z.queries[queryMonster]:
			log.Info(eq)
		}
	}
}
