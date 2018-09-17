package main

import (
    //"encoding/csv"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "io/ioutil"
    "strings"
    "gopkg.in/bufio.v1"
)

var MyS3Conn = S3Conn("cn-north-1")

// WorkDir : 当前工作目录(不一定是可执行文件所在目录)
var WorkDir = func () (string) {
    dir, _ := os.Getwd()

    return dir
}()


func readAndPutS3(aPath string){
    aPath = strings.TrimSpace(aPath)
    lineList := strings.Split(aPath, "\t")
    if len(lineList) != 2 {return }

    rawS3Path := lineList[1]
    //names := lineList[0]
    objDetail, err := MyS3Conn.SCGetObject("api-saic", rawS3Path)
    if err != nil || len(objDetail) == 0{
        fmt.Println(aPath)
        fmt.Println("err", err)
        return
    }
    // 忽略为空的年报
    if len(objDetail) < 20 {return }

    savePath := filepath.Join(WorkDir, "api-saic", rawS3Path)
    savePathDirNames := strings.Split(savePath, "/")
    if len(savePathDirNames) < 3 {return }

    savePathDir := savePathDirNames[:len(savePathDirNames)-1]
    localSavePath := strings.Join(savePathDir, "/")
    err = os.MkdirAll(localSavePath, 0755)
    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(savePath, objDetail, 0666)
    if err != nil {
        fmt.Println(aPath)
        fmt.Println("err", err)
    }

    //err1 := MyS3Conn.SCPutObject("analysis-saic", aPath, string(objDetail))
    //if err1 != nil {
    //    fmt.Println(aPath)
    //}
}


func ReadLine(filePth string, hookfn func(s string), task_num int) error {
    var lineNum = 0
    f, err := os.Open(filePth)
    if err != nil {
        fmt.Println(err)
        return err
    }

    defer f.Close()
    //reader := csv.NewReader(f)
    bfRd := bufio.NewReader(f)

    for {
        //record, err := reader.Read()
        line, err := bfRd.ReadBytes('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return err
        }

        lineNum += 1
        if task_num *100 *10000 > lineNum{
            continue
        }
        if lineNum > (task_num+1)*100  *10000 {
            break
        }
        if lineNum % 100 == 0{
            fmt.Println(lineNum)
        }

        //if len(line) != 1{
        //    continue
        //}
        hookfn(string(line)) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。

    }
    return nil
}

func main() {
    var jobs []func() interface{}
    i := 0
    allTaskNum := 9
    for i < allTaskNum {
        task_num := i
        jobs = append(jobs, func() interface{}{
            fmt.Println(task_num)
            ReadLine(filepath.Join(WorkDir, "file.txt"), readAndPutS3, task_num)
            //ReadLine(filepath.Join(WorkDir, "missName.path"), readAndPutS3, task_num)
            return nil
        })

        i += 1
    }
    var fg = ParallelDo(jobs)
    fg.WaitAll()
}
