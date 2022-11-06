package query

type fragment[RES interface{}, REPO interface{}] struct {
	RepositoriesQuery    string
	RetrieveRepositories func(res RES) REPO
}

func listRepositories[
	RES interface{},
	REPO interface{},
	FRAGMENT fragment[RES, REPO],
	DATA interface{},
](
	run func(
		fragment FRAGMENT,
		data []DATA,
	) ([]DATA, error),
	fragments ...FRAGMENT,
) ([]DATA, error) {
	repos := []DATA{}

	for _, v := range fragments {
		rs, err := run(
			v,
			[]DATA{},
		)
		if err != nil {
			return nil, err
		}
		repos = append(repos, rs...)
	}

	return repos, nil
}
