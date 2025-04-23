package consts

import "go-project/setting"

var (
	SiteCollectionView          = ""
	KnowledgeBaseCollectionView = ""
	GIT_TMP_PATH                = "/git_temp"
	KNOWLEDGE_BASE_TMP_PATH     = "/knowledgeBase_temp"
)

func Init(cfg *setting.VerctorDBConfig) {
	SiteCollectionView = cfg.SiteCollectionView
	KnowledgeBaseCollectionView = cfg.KnowledgeBaseCollectionView
}
