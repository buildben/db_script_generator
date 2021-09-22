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
	versionEnv            = "BUILDBEN_VERSION"
	indexFileName         = "index.json"
	initFolderPathPart    = "init"
	schemasFolderPathPart = "schemas"
	defaultFileSuffix     = "_initial.sql"
)

var (
	orderArray    = []string{"types", "tables", "data", "views", "routines"}
	scriptVersion string
)

func main() {
	outFileName := ""
	outFileFlag := parseFlags()
	wd, err := os.Getwd()
	if err != nil {
		logrus.Panic(err)
	}

	indexFileBytes, err := os.ReadFile(path.Join(wd, indexFileName))
	if err != nil {
		logrus.Panic(err)
	}

	indexStruct := model.IndexFile{}
	err = jsoniter.Unmarshal(indexFileBytes, &indexStruct)
	if err != nil {
		logrus.Panic(err)
	}

	versionFromEnv := os.Getenv(versionEnv)
	if len(versionFromEnv) == 0 {
		scriptVersion = indexStruct.Version
	} else {
		scriptVersion = versionFromEnv
	}

	if len(outFileFlag) > 0 {
		outFileName = outFileFlag
	} else {
		outFileName = scriptVersion + defaultFileSuffix
	}

	f, err := os.Create(outFileName)
	if err != nil {
		logrus.Panic(err)
	}

	defer utils.CloseFile(f)

	for i := range indexStruct.Init {
		err = utils.FileAppendToFile(f, path.Join(wd, initFolderPathPart, indexStruct.Init[i]))
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
