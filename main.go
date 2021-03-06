package main

import (
	"db_script_generator/internal/model"
	"db_script_generator/utils"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

const (
	basePathEnv           = "DB_SOURCES_BASE_PATH"
	indexFileName         = "index.json"
	initFolderPathPart    = "init"
	schemasFolderPathPart = "schemas"
	outputFileName        = "initial_script.sql"
)

var (
	orderArray    = []string{"types", "tables", "data", "views", "routines"}
	scriptVersion int64
)

func main() {
	basePath := os.Getenv(basePathEnv)
	indexFileBytes, err := os.ReadFile(path.Join(basePath, indexFileName))
	if err != nil {
		logrus.Panic(err)
	}

	indexStruct := model.IndexFile{}
	err = jsoniter.Unmarshal(indexFileBytes, &indexStruct)
	if err != nil {
		logrus.Panic(err)
	}

	scriptVersion = indexStruct.Version
	f, err := os.Create(outputFileName)
	if err != nil {
		logrus.Panic(err)
	}

	defer utils.CloseFile(f)

	for i := range indexStruct.Init {
		err = utils.FileAppendToFile(f, path.Join(basePath, initFolderPathPart, indexStruct.Init[i]))
		if err != nil {
			logrus.Panic(err)
		}
	}

	err = initSchemas(f, indexStruct.Schemas)
	if err != nil {
		logrus.Panic(err)
	}
}

func initSchemas(f *os.File, schemas []model.Schema) error {
	const schemaSeparatorTemplate = "\n---------- %s ----------\n"
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for i := range schemas {
		_, err = f.WriteString(fmt.Sprintf(schemaSeparatorTemplate, schemas[i].Name))
		if err != nil {
			return err
		}

		err = schemasObjectsByOrder(f, wd, schemas[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func schemasObjectsByOrder(f *os.File, wd string, schema model.Schema) error {
	pathToSchemeDir := path.Join(wd, schemasFolderPathPart, schema.Name)
	schemeDir, err := os.ReadDir(pathToSchemeDir)
	if err != nil {
		return err
	}

	err = utils.FileAppendToFile(f, path.Join(pathToSchemeDir, schema.Name+".sql"))
	if err != nil {
		return err
	}

	for j := range orderArray {
		for z := range schemeDir {
			if !schemeDir[z].IsDir() || schemeDir[z].Name() != orderArray[j] {
				continue
			}

			switch orderArray[j] {
			case "types":
				err = utils.AllFilesFromDirAppendToFile(f, path.Join(pathToSchemeDir, schemeDir[z].Name()), schema.Types)
			case "tables":
				err = utils.AllFilesFromDirAppendToFile(f, path.Join(pathToSchemeDir, schemeDir[z].Name()), schema.Tables)
			case "data":
				err = utils.AllFilesFromDirAppendToFile(f, path.Join(pathToSchemeDir, schemeDir[z].Name()), schema.Data)
			case "views":
				err = utils.AllFilesFromDirAppendToFile(f, path.Join(pathToSchemeDir, schemeDir[z].Name()), schema.Views)
			case "routines":
				err = utils.AllFilesFromDirAppendToFile(f, path.Join(pathToSchemeDir, schemeDir[z].Name()), schema.Routines)
			}

			if err != nil {
				return err
			}
		}
	}

	if schema.Name == "vsn" {
		callVsnStoredProc := fmt.Sprintf("\n SELECT * FROM vsn.versions_create('{\"version_id\": %v}'); \n", scriptVersion)
		_, err = f.WriteString(callVsnStoredProc)
	}

	if err != nil {
		return err
	}

	return nil
}
