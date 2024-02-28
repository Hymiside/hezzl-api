package service

type repositoryClickhouse interface{}

type repositoryPostgres interface{}

type repositoryRedis interface{}

type queueNats interface{}

type Service struct {
	repoClickhouse repositoryClickhouse
	repoPostgres   repositoryPostgres
	repoRedis      repositoryRedis
	queueNats      queueNats
}

func NewService(
	repoClickhouse repositoryClickhouse,
	repoPostgres repositoryPostgres,
	repoRedis repositoryRedis,
	queueNats queueNats,
) *Service {
	return &Service{
		repoClickhouse: repoClickhouse,
		repoPostgres:   repoPostgres,
		repoRedis:      repoRedis,
		queueNats:      queueNats,
	}
}
