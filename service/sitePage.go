package service

import (
	"context"
	"fmt"
	"go-project/consts"
	"go-project/dao"
	"go-project/logger"
	"strings"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type SitePage struct {
	Coll *tcvectordb.AICollectionView
}

type SiteSearchParams struct {
	Sitename string `json:"sitename"`
	Username string `json:"username"`
	Dir      string `json:"dir"`
	Query    string `json:"query" binding:"required"`
	Limit    int64  `json:"limit"`
}

func NewSitePage() *SitePage {
	coll := dao.Db.CollectionView(consts.SiteCollectionView)
	return &SitePage{
		Coll: coll,
	}
}

func (sitePage *SitePage) UploadFile(ctx context.Context, filepath string, sitename string, username string, dir string, fullname string) (err error) {
	documentSetName := generateDocumentSetName(fullname)
	res, err := sitePage.Coll.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		LocalFilePath:   filepath,
		DocumentSetName: documentSetName,
		MetaData: map[string]interface{}{
			"username": username,
			"sitename": sitename,
			"dir":      dir,
			"fullname": fullname,
		},
	})
	logger.Logger.Info(fmt.Sprintf("UploadFileResult: %+v", res))
	return
}

func generateDocumentSetName(fullname string, args ...string) string {
	// 取出文件的后缀，判断出类型，后缀可能不同
	tmp := strings.Split(fullname, ".")
	fileType := tmp[len(tmp)-1]
	return strings.Join(args, "-") + "." + fileType
}

func (sitePage *SitePage) DeleteFiles(ctx context.Context, files []string) (err error) {
	documentSetNames := make([]string, len(files))
	for i, file := range files {
		documentSetNames[i] = generateDocumentSetName(file)
	}
	res, err := sitePage.Coll.DeleteByNames(ctx, documentSetNames...)
	logger.Logger.Info(fmt.Sprintf("DeleteFilesResult: %+v", res.AffectedCount))
	return
}

func (sitePage *SitePage) Search(ctx context.Context, params SiteSearchParams) (results []tcvectordb.AISearchDocumentSet, err error) {
	if params.Limit == 0 {
		params.Limit = 3
	}

	var filters []string
	if params.Sitename != "" {
		filters = append(filters, fmt.Sprintf(`sitename = "%s"`, params.Sitename))
	}

	if params.Username != "" {
		filters = append(filters, fmt.Sprintf(`username = "%s"`, params.Username))
	}

	if params.Dir != "" {
		filters = append(filters, fmt.Sprintf(`dir = "%s"`, params.Dir))
	}

	filterStr := strings.Join(filters, " and ")
	res, err := sitePage.Coll.Search(ctx, tcvectordb.SearchAIDocumentSetsParams{
		Content: params.Query,
		Limit:   params.Limit,
		Filter:  tcvectordb.NewFilter(filterStr),
	})

	if err != nil {
		return
	}

	results = res.Documents

	return
}
