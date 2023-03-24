package main

import (
	"os"
	"testing"

	"github.com/DQNEO/babygo/lib/fmt"
	"github.com/DQNEO/babygo/lib/path"
)

func TestExampleCompile(t *testing.T) {
	prjSrcPath = "./src"
	workdir := "./tmp"
	initAsm, err := os.Create(workdir + "/a.s")
	if err != nil {
		panic(err)
	}
	fout = initAsm
	logf("Build start\n")

	debugFrontEnd = true
	debugCodeGen = true
	inputFiles := []string{"./example/map.go"}

	paths := collectAllPackages(inputFiles)
	var packagesToBuild []*PkgContainer
	// 针对本地编译项目所导入的包，获取其包的路径和包下所有的文件
	for _, _path := range paths {
		files := collectSourceFiles(getPackageDir(_path))
		packagesToBuild = append(packagesToBuild, &PkgContainer{
			name:  path.Base(_path),
			path:  _path,
			files: files,
		})
	}

	// 额外添加main包，以及main包中需要编译的文件
	packagesToBuild = append(packagesToBuild, &PkgContainer{
		name:  "main",
		files: inputFiles,
	})

	var universe = createUniverse()
	// 将每个包都编译成一个汇编文件
	for _, _pkg := range packagesToBuild {
		currentPkg = _pkg
		if _pkg.name == "" {
			panic("empty pkg name")
		}
		pgkAsm, err := os.Create(fmt.Sprintf("%s/%s.s", workdir, _pkg.name))
		if err != nil {
			panic(err)
		}
		fout = pgkAsm
		buildPackage(_pkg, universe)
		pgkAsm.Close()
	}
	initAsm.Close()
}
