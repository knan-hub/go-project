package service

import (
	"context"
	"fmt"
	"go-project/consts"
	"go-project/dao"
	"go-project/logger"
	"os"
	"strings"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type KnowledgeBaseIndex struct {
	// 使用string搜索时可以用in
	KnowledgeBaseId     string `json:"knowledgeBaseId"`
	KnowledgeBaseItemId string `json:"knowledgeBaseItemId"`
	Site                string `json:"site"`
	Path                string `json:"path"`
	Filename            string `json:"filename"`
	RelativePath        string `json:"relativePath"`
}

type KnowledgeBase struct {
	Coll *tcvectordb.AICollectionView
}

type KnowledgeSiteUploadParams struct {
	Site                string `json:"site"`
	GitUrl              string `json:"gitUrl"`
	Path                string `json:"path"`
	KnowledgeBaseId     string `json:"knowledgeBaseId"`
	KnowledgeBaseItemId string `json:"knowledgeBaseItemId"`
}

type KnowledgeSearchParams struct {
	KnowledgeBaseIds []string `json:"knowledgeBaseIds"`
	Query            string   `json:"query" binding:"required"`
	Limit            int64    `json:"limit"`
}

func NewKnowledgeBase() *KnowledgeBase {
	coll := dao.Db.CollectionView(consts.KnowledgeBaseCollectionView)
	return &KnowledgeBase{
		Coll: coll,
	}
}

func (base *KnowledgeBase) UploadFile(ctx context.Context, filepath string, baseIndex KnowledgeBaseIndex) (res *tcvectordb.LoadAndSplitTextResult, err error) {
	// 判断下文件内容是否为空
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() == 0 {
		logger.Logger.Info(fmt.Sprintf("uploadFile warning: empty file %+v", baseIndex))
		return nil, nil
	}

	addTitleToChunk := true
	documentSetName := generateDocumentSetName(filepath, baseIndex.KnowledgeBaseId, baseIndex.KnowledgeBaseItemId, baseIndex.Filename)

	res, err = base.Coll.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		LocalFilePath:   filepath,
		DocumentSetName: documentSetName,
		MetaData: map[string]interface{}{
			"knowledgeBaseId":     baseIndex.KnowledgeBaseId,
			"knowledgeBaseItemId": baseIndex.KnowledgeBaseItemId,
			"site":                baseIndex.Site,
			"path":                baseIndex.Path,
			"filename":            baseIndex.Filename,
			"relativePath":        baseIndex.RelativePath,
		},
		SplitterPreprocess: ai_document_set.DocumentSplitterPreprocess{
			AppendTitleToChunk: &addTitleToChunk,
		},
	})

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("UploadFile failed, err: %+v, %+v", err, baseIndex))
		return
	}

	logger.Logger.Info(fmt.Sprintf("UploadFileResult: %+v", res))

	return
}

func (base *KnowledgeBase) DeleteFileByKnowledgeBaseItemId(ctx context.Context, knowledgeBaseItemId string) (err error) {
	res, err := base.Coll.Delete(ctx, tcvectordb.DeleteAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(fmt.Sprintf(`knowledgeBaseItemId = "%s"`, knowledgeBaseItemId)),
	})
	logger.Logger.Info(fmt.Sprintf("DeleteFilesResult: %+v", res.AffectedCount))
	return
}

func (base *KnowledgeBase) DeleteFileByKnowledgeBaseId(ctx context.Context, knowledgeBaseId string) (err error) {
	res, err := base.Coll.Delete(ctx, tcvectordb.DeleteAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(fmt.Sprintf(`knowledgeBaseId = "%s"`, knowledgeBaseId)),
	})
	logger.Logger.Info(fmt.Sprintf("DeleteFilesResult: %+v", res.AffectedCount))
	return
}

func (base *KnowledgeBase) Search(ctx context.Context, params KnowledgeSearchParams) (results []tcvectordb.AISearchDocumentSet, err error) {
	if params.Limit == 0 {
		params.Limit = 3
	}

	var filters []string
	if len(params.KnowledgeBaseIds) > 0 {
		var res string
		for i, id := range params.KnowledgeBaseIds {
			res += "\"" + id + "\""
			if i < len(params.KnowledgeBaseIds)-1 {
				res += ", "
			}
		}
		filters = append(filters, fmt.Sprintf(`knowledgeBaseId in (%s)`, res))

	}

	filterStr := strings.Join(filters, " and ")
	enableRerank := true

	res, err := base.Coll.Search(ctx, tcvectordb.SearchAIDocumentSetsParams{
		Content: params.Query,
		Limit:   params.Limit,
		Filter:  tcvectordb.NewFilter(filterStr),
		RerankOption: &ai_document_set.RerankOption{
			Enable:                &enableRerank,
			ExpectRecallMultiples: 10,
		},
	})
	if err != nil {
		return
	}

	results = res.Documents

	return
}
