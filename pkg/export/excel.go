package export

import (
	"go_gin_base/pkg/setting"
)

const ExcelEXT = ".xlsx"

func GetExcelPath() string {
	return setting.App.ExportExcelSavePath
}

//link path
func GetExcelFullUrl(name string) string {
	return setting.App.PrefixUrl + "/" + GetExcelPath() + name
}

//file path
func GetExcelFullPath() string {
	return setting.App.RuntimeRootPath + GetExcelPath()
}
