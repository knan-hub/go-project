package service

import (
	"context"
	"fmt"
	"go-project/consts"
	"go-project/logger"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type GitRepoParams struct {
	GitUrl       string `json:"gitUrl"`
	LastCommitId string `json:"lastCommitId"`
	NewCommitId  string `json:"newCommitId"`
}

type GitRepo struct {
	Url string `json:"url"`
	Dir string `json:"dir"`
}

func NewGitRepo(url string) *GitRepo {
	uuid := uuid.New()
	uuidStr := uuid.String()
	dir := filepath.Join(consts.GIT_TMP_PATH, uuidStr)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to create directory %s: %s", dir, err))
		return nil
	}
	return &GitRepo{
		Url: url,
		Dir: dir,
	}
}

func (repo *GitRepo) ProcessSite(c context.Context, lastCommitId string, newCommitId string) (err error) {
	err = repo.Download(c)
	if err != nil {
		return
	}

	defer repo.Clean(c)

	if lastCommitId == "" || newCommitId == "" {
		err = repo.UploadMdFiles(c)
	} else {
		updateFiles, deleteFiles, _err := repo.GetDiffFiles(c, lastCommitId, newCommitId)
		sitePage := NewSitePage()
		if _err != nil {
			return _err
		}

		for _, file := range updateFiles {
			username, sitename, dir := extractValuesFromPath(file)
			logger.Logger.Info(fmt.Sprintf("Uploading Diff file: %s", filepath.Join(repo.Dir, file)))
			err = sitePage.UploadFile(c, filepath.Join(repo.Dir, file), sitename, username, dir, file)
			if err != nil {
				return
			}
		}

		if len(deleteFiles) > 0 {
			err = sitePage.DeleteFiles(c, deleteFiles)
			if err != nil {
				return
			}
		}
	}

	return
}

func (repo *GitRepo) Download(c context.Context) (err error) {
	logger.Logger.Info(fmt.Sprintf("Cloning %s repository to %s", repo.Url, repo.Dir))
	cmd := exec.Command("git", "clone", repo.Url, repo.Dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to clone %s repository: %s, output: %s", repo.Url, err, out))
		return
	}
	logger.Logger.Info("Repository cloned successfully")
	return nil
}

func (repo *GitRepo) Clean(c context.Context) (err error) {
	err = os.RemoveAll(repo.Dir)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to remove %s repository: %s", repo.Dir, err))
		return
	}
	logger.Logger.Info("Repository removed successfully")
	return nil
}

func (repo *GitRepo) UploadMdFiles(c context.Context) (err error) {
	sitePage := NewSitePage()
	files, err := repo.WalkMdFiles(c, "")
	if err != nil {
		return err
	}

	for _, file := range files {
		// 上传文件到向量数据库
		username, sitename, dir := extractValuesFromPath(file)
		logger.Logger.Info(fmt.Sprintf("Uploading file: %s", filepath.Join(repo.Dir, file)))
		err = sitePage.UploadFile(c, filepath.Join(repo.Dir, file), sitename, username, dir, file)
		if err != nil {
			return err
		}
	}

	return
}

func (repo *GitRepo) WalkMdFiles(ctx context.Context, subPath string) (files []string, err error) {
	mdPattern := regexp.MustCompile(`\.md$`)
	// 指定要排除的目录名称
	excludeDirs := []string{".git", "_config"}

	fullPath := filepath.Join(repo.Dir, subPath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		if mdPattern.MatchString(filepath.Base(fullPath)) {
			relativePath, err := filepath.Rel(repo.Dir, fullPath)
			if err != nil {
				return nil, err
			}
			files = append(files, relativePath)
		}
		return files, nil
	}

	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过指定的目录
		if info.IsDir() && isExcludedDir(path, excludeDirs) {
			logger.Logger.Info(fmt.Sprintf("Skipping directory: %s", path))
			return filepath.SkipDir
		}

		if !info.IsDir() && mdPattern.MatchString(filepath.Base(path)) {
			relativePath, err := filepath.Rel(repo.Dir, path)
			if err != nil {
				return err
			}
			files = append(files, relativePath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// 检查给定的路径是否包含在要排除的目录列表中
func isExcludedDir(path string, excludeDirs []string) bool {
	for _, dir := range excludeDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

// 从给定的路径中提取username，sitename和第一级目录
func extractValuesFromPath(path string) (username string, sitename string, firstLevelDir string) {
	mdPattern := regexp.MustCompile(`\.md$`)
	// 使用 ilepath.SplitList处理路径，以适应不同操作系统中的路径分隔符
	parts := strings.Split(path, string(filepath.Separator))
	// 提取username
	if len(parts) > 0 {
		username = parts[0]
	}

	// 提取sitename
	if len(parts) > 1 {
		sitename = parts[1]
	}

	// 提取第一级目录
	if len(parts) > 2 {
		if mdPattern.MatchString(parts[2]) {
			firstLevelDir = "/"
		} else {
			firstLevelDir = parts[2]
		}
	}

	return
}

func (repo *GitRepo) GetDiffFiles(c context.Context, lastCommitId string, newCommitId string) (updateFiles []string, deleteFiles []string, err error) {
	cmd := exec.Command("git", "diff", "--name-status", lastCommitId, newCommitId)
	cmd.Dir = repo.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to get diff files: %s, output: %s", err, out))
		return
	}

	// 判断是删除还是修改，并且排除exclude dir
	mdPattern := regexp.MustCompile(`\.md$`)
	excludeDirs := []string{".git", "_config"}

	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}

		status := line[0:1]
		file := line[1:]
		// 移除两边空格
		file = strings.TrimSpace(file)

		if isExcludedDir(file, excludeDirs) || !mdPattern.MatchString(file) {
			continue
		}

		if status == "D" {
			deleteFiles = append(deleteFiles, file)
		} else {
			updateFiles = append(updateFiles, file)
		}
	}

	return
}

func (repo *GitRepo) ProcessKnowledgeBase(c context.Context, uploadParams KnowledgeSiteUploadParams) (err error) {
	// 先删除knowledgeBaseItem
	knowledgeBase := NewKnowledgeBase()
	err = knowledgeBase.DeleteFileByKnowledgeBaseItemId(c, uploadParams.KnowledgeBaseItemId)
	if err != nil {
		return
	}

	// TODO 考虑缓存下来git文件，不用每次都download
	err = repo.Download(c)
	if err != nil {
		return
	}

	defer repo.Clean(c)

	subPath := filepath.Join(uploadParams.Site, uploadParams.Path)
	files, err := repo.WalkMdFiles(c, subPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 上传文件到向量数据库
		knowledgeBase := NewKnowledgeBase()
		res, err := knowledgeBase.UploadFile(c, filepath.Join(repo.Dir, file), KnowledgeBaseIndex{
			KnowledgeBaseId:     uploadParams.KnowledgeBaseId,
			KnowledgeBaseItemId: uploadParams.KnowledgeBaseItemId,
			Path:                uploadParams.Path,
			Site:                uploadParams.Site,
			RelativePath:        file,
		})
		if err != nil {
			return err
		}
		logger.Logger.Info(fmt.Sprintf("ProcessKnowledgeBase Uploaded file: %+v", res))
	}

	return

}
