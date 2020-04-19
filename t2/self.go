package main

import "syscall"

var buf [100]uint8
var r [100]uint8

func Itoa(ival int) string {
	var next int
	var right int
	var ix int = 0
	if ival == 0 {
		return "0"
	}
	var minus bool
	minus = false
	for ix = 0; ival != 0; ix = ix + 1 {
		if ival < 0 {
			ival = -1 * ival
			minus = true
			r[0] = '-'
		} else {
			next = ival / 10
			right = ival - next*10
			ival = next
			buf[ix] = uint8('0' + right)
		}
	}

	var j int
	var c uint8
	for j = 0; j < ix; j = j + 1 {
		c = buf[ix-j-1]
		if minus {
			r[j+1] = c
		} else {
			r[j] = c
		}
	}

	return string(r[0:ix])
}

func fmtPrintf(s string) {
	var slc []uint8 = []uint8(s)
	syscall.Write(1, slc)
}

func semanticAnalyze(s string) string {
	globalFuncs = make([]*Func, 2, 2)
	var fnc *Func = new(Func)
	fnc.name = "main"
	globalFuncs[0] = fnc
	var fnc2 *Func = new(Func)
	fnc2.name = "foo"
	globalFuncs[1] = fnc2

	stringLiterals = make([]string, 2, 2)
	stringLiterals[0] = "hello0"
	stringLiterals[1] = "hello1"

	return s
}


func emitData(pkgName string) {
	fmtPrintf(".data\n")
	var i int = 0
	for i=0;i<len(stringLiterals);i++ {
		fmtPrintf("# string literals\n")
		var seq string = Itoa(i)
		fmtPrintf("." + pkgName + ".S" + seq + ":\n")
		fmtPrintf("  .string " + stringLiterals[i] + "\n")
	}
	fmtPrintf("# ===== Global Variables =====\n")
	fmtPrintf("# ==============================\n")
}

func emitFuncDecl(pkgPrefix string, fn *Func) {
	var fname string = fn.name
	fmtPrintf(pkgPrefix + "." + fname + ":\n")
	fmtPrintf("  ret\n")
}

func emitText(pkgName string) {
	fmtPrintf(".text\n")
	var i int
	for i = 0; i<len(globalFuncs); i++ {
		var fnc *Func = globalFuncs[i]
		emitFuncDecl(pkgName, fnc)
	}
}

func generateCode(pkgName string) {
	emitData(pkgName)
	emitText(pkgName)
}


type astValueSpec struct {
}

type Func struct {
//	decl      *ast.FuncDecl
//	localvars []*ast.ValueSpec
	localarea int
	argsarea  int
	name string
}

var stringLiterals []string
var stringIndex int
var globalVars []*astValueSpec
var globalFuncs []*Func

var sourceFiles [1]string
var _garbage string
func main() {
	sourceFiles[0] = "main"
	var i int
	for i=0;i<len(sourceFiles); i++ {
		globalVars = nil
		globalFuncs = nil
		stringLiterals = nil
		stringIndex = 0
		_garbage = sourceFiles[i]
		var pkgName string = semanticAnalyze(sourceFiles[0])
		generateCode(pkgName)
	}
}
