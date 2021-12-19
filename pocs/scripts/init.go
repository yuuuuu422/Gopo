package scripts

import (
	"fmt"
)

var threads int

type scriptFunc func(target string)

var ScriptMap = map[string]scriptFunc{}

func ScriptInit(scriptName string,num int) scriptFunc{
	threads=num
	 script,ok := ScriptMap[scriptName]
	 if ok{
		 return script
	 }
	 return nil
}

func scriptRegister(scriptName string,scriptFunc scriptFunc){
	ScriptMap[scriptName]=scriptFunc
}

func ShowRegister(){
	for k,_:=range ScriptMap{
		fmt.Println("   ",k)
	}
}