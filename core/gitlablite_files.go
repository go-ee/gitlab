package core

func NewGitlabLiteMemJson(fileOrFolderJson string, filePattern string) (ret *GitlabLiteMem, err error) {
	ret = NewGitlabLiteMem()
	var fileLoader *FileLoaderJson
	if fileLoader, err = NewFileLoaderJson(filePattern, ret.AddGroup); err == nil {
		err = fileLoader.LoadFileOrFolder(fileOrFolderJson)
	}
	return
}
