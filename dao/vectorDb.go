package dao

import (
	"context"
	"go-project/consts"
	"go-project/setting"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

var (
	Client *tcvectordb.RpcClient
	Db     *tcvectordb.AIDatabase
)

func Init(cfg *setting.VectorDbConfig) (err error) {
	var defaultOption = &tcvectordb.ClientOption{
		Timeout:            30 * time.Second,
		MaxIdldConnPerHost: 2,
		IdleConnTimeout:    0,
		ReadConsistency:    tcvectordb.EventualConsistency,
	}

	Client, err = tcvectordb.NewRpcClient(cfg.Url, cfg.Username, cfg.Key, defaultOption)
	if err != nil {
		return
	}

	if err = initAIDb(cfg.DatabaseName); err != nil {
		return
	}

	if err = initSiteCollectionView(cfg.DatabaseName); err != nil {
		return
	}

	if err = initKnowledgeBaseCollectionView(cfg.DatabaseName); err != nil {
		return
	}

	return
}

func initAIDb(dbName string) (err error) {
	ctx := context.Background()
	dbs, err := Client.ListDatabase(ctx)
	if err != nil {
		return
	}

	var db tcvectordb.AIDatabase
	for _, d := range dbs.AIDatabases {
		if d.DatabaseName == dbName {
			db = d
			break
		}
	}

	if db.DatabaseName != dbName {
		res, _err := Client.CreateAIDatabase(ctx, dbName)
		if _err != nil {
			return _err
		}
		db = res.AIDatabase
	}

	Db = &db

	return
}

func initSiteCollectionView(dbName string) (err error) {
	ctx := context.Background()
	db := Client.AIDatabase(dbName)
	collectionViewName := consts.SiteCollectionView

	collections, err := db.ListCollectionViews(ctx)
	if err != nil {
		return
	}

	var collectionView *tcvectordb.AICollectionView
	for _, c := range collections.CollectionViews {
		if c.CollectionViewName == collectionViewName {
			collectionView = c
			break
		}
	}

	if collectionView == nil {
		filterIndexes := []tcvectordb.FilterIndex{
			{FieldName: "username", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "sitename", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "dir", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "fullname", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
		}

		index := tcvectordb.Indexes{}
		index.FilterIndex = filterIndexes
		db.WithTimeout(time.Second * 30)

		enableWordsEmbedding := true
		appendTitleToChunk := true
		appendKeywordsToChunk := false

		_, err := db.CreateCollectionView(ctx, collectionViewName, tcvectordb.CreateCollectionViewParams{
			Description: "网站网页数据",
			Indexes:     index,
			Embedding: &collection_view.DocumentEmbedding{
				Language:             string(tcvectordb.LanguageChinese),
				EnableWordsEmbedding: &enableWordsEmbedding,
			},
			SplitterPreprocess: &collection_view.SplitterPreprocess{
				AppendTitleToChunk:    &appendTitleToChunk,
				AppendKeywordsToChunk: &appendKeywordsToChunk,
			},
		})

		if err != nil {
			return err
		}
	}

	return
}

func initKnowledgeBaseCollectionView(dbName string) (err error) {
	ctx := context.Background()
	db := Client.AIDatabase(dbName)
	collectionViewName := consts.KnowledgeBaseCollectionView

	collections, err := db.ListCollectionViews(ctx)
	if err != nil {
		return
	}

	var collectionView *tcvectordb.AICollectionView
	for _, c := range collections.CollectionViews {
		if c.CollectionViewName == collectionViewName {
			collectionView = c
			break
		}
	}

	if collectionView == nil {
		filterIndexes := []tcvectordb.FilterIndex{
			{FieldName: "knowledgeBaseId", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "knowledgeBaseItemId", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "filename", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "site", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "path", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
		}

		index := tcvectordb.Indexes{}
		index.FilterIndex = filterIndexes
		db.WithTimeout(time.Second * 30)

		enableWordsEmbedding := true
		appendTitleToChunk := true
		appendKeywordsToChunk := false

		_, err := db.CreateCollectionView(ctx, collectionViewName, tcvectordb.CreateCollectionViewParams{
			Description: "知识库数据",
			Indexes:     index,
			Embedding: &collection_view.DocumentEmbedding{
				Language:             string(tcvectordb.LanguageChinese),
				EnableWordsEmbedding: &enableWordsEmbedding,
			},
			SplitterPreprocess: &collection_view.SplitterPreprocess{
				AppendTitleToChunk:    &appendTitleToChunk,
				AppendKeywordsToChunk: &appendKeywordsToChunk,
			},
		})

		if err != nil {
			return err
		}
	}

	return
}
