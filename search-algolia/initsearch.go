package algolia

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// initSettings update algolia search settings
func (s *SearchAlgolia) initSettings() (err error) {
	var (
		settings = search.Settings{}
	)
	err = settings.UnmarshalJSON(AlgoliaSearchServerConfig)
	if err != nil {
		return
	}

	// point virtual index to sort
	settings.Replicas = opt.Replicas(
		"virtual("+s.getIndexName(NewestIndex)+")",
		"virtual("+s.getIndexName(ActiveIndex)+")",
		"virtual("+s.getIndexName(ScoreIndex)+")",
	)

	_, err = s.getIndex("").SetSettings(settings, opt.ForwardToReplicas(true))
	if err != nil {
		return
	}
	err = s.initVirtualReplicaSetting()
	return
}

// initVirtualReplicaSetting init virtual index replica setting
func (s *SearchAlgolia) initVirtualReplicaSetting() (err error) {

	_, err = s.getIndex(NewestIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(created)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}

	_, err = s.getIndex(ActiveIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(active)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}

	_, err = s.getIndex(ScoreIndex).SetSettings(search.Settings{
		CustomRanking: opt.CustomRanking(
			"desc(score)",
			"desc(content)",
			"desc(title)"),
	})
	if err != nil {
		return
	}
	return
}
