package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/answerdev/answer/plugin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"meilisearch/i18n"
	"strings"
)

const primaryKey = "objectID"
const defaultIndexName = "answer_post"

type Search struct {
	Config *SearchConfig
	Client *meilisearch.Client
}

type SearchConfig struct {
	Host      string `json:"host"`
	ApiKey    string `json:"api_key"`
	IndexName string `json:"index_name"`
	Async     bool   `json:"async"`
}

func init() {
	plugin.Register(&Search{
		Config: &SearchConfig{},
	})
}

func (s *Search) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "meilisearch",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "sivdead",
		Version:     "0.0.1",
		Link:        "https://github.com/answerdev/plugins/tree/main/storage/aliyunoss",
	}
}

func (s *Search) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "host",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigHostTitle),
			Description: plugin.MakeTranslator(i18n.ConfigHostDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Host,
		},
		{
			Name:        "api_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigApiKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigApiKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.ApiKey,
		},
		{
			Name:        "index_name",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigIndexTitle),
			Description: plugin.MakeTranslator(i18n.ConfigIndexDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.IndexName,
		},
		{
			Name:        "async",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigAsyncTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAsyncDescription),
			Required:    false,
			Value:       s.Config.Async,
		},
	}
}

func (s *Search) ConfigReceiver(config []byte) error {
	conf := &SearchConfig{}
	_ = json.Unmarshal(config, conf)

	// if index name is empty, use default index name
	if conf.IndexName == "" {
		conf.IndexName = defaultIndexName
	}

	s.Config = conf

	log.Debugf("try to init meilisearch client: %s", conf.Host)

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   conf.Host,
		APIKey: conf.ApiKey,
	})
	s.Client = client
	resp, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        conf.IndexName,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		log.Errorf("create index error: %s", err.Error())
	}
	err = waitForTask(client, resp)
	if err != nil {
		log.Errorf("create index error: %s", err.Error())
		//ignore index exists error
	}
	index := client.Index(conf.IndexName)
	_, err = index.UpdateSearchableAttributes(&[]string{"title", "content", "tags", "status", "answers", "type", "questionID", "userID", "views", "created", "active", "score", "hasAccepted"})
	if err != nil {
		log.Errorf("update searchable attributes error: %s", err.Error())
		return err
	}
	_, err = index.UpdateSortableAttributes(&[]string{"active", "created", "active", "score"})
	if err != nil {
		log.Errorf("update sortable attributes error: %s", err.Error())
		return err
	}
	_, err = index.UpdateDisplayedAttributes(&[]string{"title", "content", "objectID", "type"})
	if err != nil {
		log.Errorf("update displayed attributes error: %s", err.Error())
		return err
	}
	return nil
}

func (s *Search) SearchContents(_ context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {

	//search and parse
	client := s.Client
	index := client.Index(s.Config.IndexName)

	query, searchRequest := s.buildQuery(cond)
	// convert searchRequest.Filter to []string

	filter := s.buildFilter(cond)
	searchRequest.Filter = filter
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("search error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) SearchQuestions(_ context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {

	//search and parse
	client := s.Client
	index := client.Index(s.Config.IndexName)

	query, searchRequest := s.buildQuery(cond)
	// convert searchRequest.Filter to []string

	filter := s.buildFilter(cond)
	filter = append(filter, "type = question")
	searchRequest.Filter = filter
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("search error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) SearchAnswers(_ context.Context, cond *plugin.SearchBasicCond) (res []plugin.SearchResult, total int64, err error) {

	//search and parse
	client := s.Client
	index := client.Index(s.Config.IndexName)

	query, searchRequest := s.buildQuery(cond)
	// convert searchRequest.Filter to []string

	filter := s.buildFilter(cond)
	filter = append(filter, "type = answer")
	searchRequest.Filter = filter
	searchResult, err := index.Search(query, searchRequest)
	if err != nil {
		log.Errorf("search error: %s", err.Error())
		return nil, 0, err
	}
	return s.warpResult(searchResult)
}

func (s *Search) warpResult(resp *meilisearch.SearchResponse) ([]plugin.SearchResult, int64, error) {
	res := make([]plugin.SearchResult, 0)
	for _, hit := range resp.Hits {

		var content plugin.SearchContent
		bytes, err := json.Marshal(hit)
		if err != nil {
			log.Errorf("marshal hit error: %s", err.Error())
			return nil, 0, err
		}
		err = json.Unmarshal(bytes, &content)
		if err != nil {
			log.Errorf("unmarshal hit error: %s", err.Error())
			return nil, 0, err
		}

		res = append(res, plugin.SearchResult{
			ID:   content.ObjectID,
			Type: content.Type,
		})
	}
	log.Debugf("search result: %d", len(res))
	return res, resp.TotalHits, nil
}

func (s *Search) UpdateContent(_ context.Context, _ string, content *plugin.SearchContent) error {

	client := s.Client
	index := client.Index(s.Config.IndexName)
	if s.Config.Async {
		_, err := index.AddDocuments([]*plugin.SearchContent{content}, primaryKey)
		return err
	} else {
		resp, err := index.AddDocuments([]*plugin.SearchContent{content}, primaryKey)
		if err != nil {
			return err
		}
		return waitForTask(client, resp)
	}
}

func (s *Search) DeleteContent(_ context.Context, contentID string) error {

	client := s.Client
	index := client.Index(s.Config.IndexName)

	if s.Config.Async {
		_, err := index.DeleteDocument(contentID)
		return err
	} else {
		resp, err := index.DeleteDocument(contentID)
		err = waitForTask(client, resp)
		return err
	}
}

func (s *Search) buildQuery(cond *plugin.SearchBasicCond) (string, *meilisearch.SearchRequest) {

	searchRequest := meilisearch.SearchRequest{}

	// page
	if cond.Page > 0 {
		searchRequest.Page = int64(cond.Page)
	}
	if cond.PageSize > 0 {
		searchRequest.HitsPerPage = int64(cond.PageSize)
	}

	// order
	switch cond.Order {
	case plugin.SearchNewestOrder:
		searchRequest.Sort = []string{"created:desc"}
	case plugin.SearchActiveOrder:
		searchRequest.Sort = []string{"created:desc"}
	case plugin.SearchScoreOrder:
		searchRequest.Sort = []string{"score:desc"}
	}

	var query string
	if cond.Words != nil && len(cond.Words) > 0 {
		query = strings.Join(cond.Words, " ")
	}
	return query, &searchRequest
}

func (s *Search) buildFilter(cond *plugin.SearchBasicCond) []string {
	var filter []string
	if cond.TagIDs != nil && len(cond.TagIDs) > 0 {
		filter = append(filter, fmt.Sprintf("tags IN (%s)", strings.Join(cond.TagIDs, ",")))
	}
	if cond.UserID != "" {
		filter = append(filter, fmt.Sprintf("userID = %s", cond.UserID))
	}
	// QuestionAccepted
	if cond.QuestionAccepted == plugin.AcceptedCondTrue {
		filter = append(filter, "hasAccepted = true")
	} else if cond.QuestionAccepted == plugin.AcceptedCondFalse {
		filter = append(filter, "hasAccepted = false")
	}

	// AnswerAccepted
	if cond.AnswerAccepted == plugin.AcceptedCondTrue {
		filter = append(filter, "hasAccepted = true")
	} else if cond.AnswerAccepted == plugin.AcceptedCondFalse {
		filter = append(filter, "hasAccepted = false")
	}

	// QuestionID
	if cond.QuestionID != "" {
		filter = append(filter, fmt.Sprintf("questionID = %s", cond.QuestionID))
	}

	// VoteAmount
	if cond.VoteAmount > 0 {
		filter = append(filter, fmt.Sprintf("voteAmount >= %d", cond.VoteAmount))
	}

	// ViewAmount
	if cond.ViewAmount > 0 {
		filter = append(filter, fmt.Sprintf("viewAmount >= %d", cond.ViewAmount))
	}

	// AnswerAmount
	if cond.AnswerAmount > 0 {
		filter = append(filter, fmt.Sprintf("answerAmount >= %d", cond.AnswerAmount))
	}
	return filter
}

func waitForTask(client *meilisearch.Client, resp *meilisearch.TaskInfo) error {
	task, err := client.WaitForTask(resp.TaskUID)
	if err != nil {
		return err
	}
	if task.Status != meilisearch.TaskStatusSucceeded {
		err = errors.InternalServer(task.Error.Code).WithMsg("invoke meili failed:" + task.Error.Message).WithStack()
		return err
	}
	return nil
}
