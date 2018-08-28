package common

import (
    "github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
    "fmt"
)

// 解析Path
func ParsePath(data *[][]interface{})  {
    var paths = [][]string{}
    if len(*data) == 0{
        return
    }
    for _, d := range *data{
        // 记录路径
        var aPath = []string{}
        if len(d) == 0 {
            continue
        }
        dPath := d[0].(graph.Path)
        // 解析点
        dNodes := dPath.Nodes
        if len(dNodes) == 0{
            continue
        }
        for _, v := range dNodes{
            if len(v.Labels) == 0{
                continue
            }
            vProperties := v.Properties
            fmt.Println("vProperties ", vProperties)
        }
        paths = append(paths, aPath)
        // 解析关系
        dRelationships := dPath.Relationships
        dSequence := dPath.Sequence
        dSequenceFlag := append([]int{0}, dSequence...)
        for sIdx := 0; sIdx < len(dSequenceFlag) / 2; sIdx++ {
            rIdx := dSequenceFlag[2*sIdx+1]
            rIdxAbs := 0
            linkId, target, source := "", "", ""
            fmt.Println("linkId", linkId)
            // todo
            if rIdx > 0 {
                rIdxAbs = rIdx
                source = aPath[dSequenceFlag[2*sIdx]]
                target = aPath[dSequenceFlag[2*(sIdx+1)]]
            } else if rIdx < 0 {
                rIdxAbs = rIdx * -1
                source = aPath[dSequenceFlag[2*(sIdx+1)]]
                target = aPath[dSequenceFlag[2*sIdx]]
            } else {
                continue
            }
            v := dRelationships[rIdxAbs-1]
            fmt.Println("v ", v)
            // todo
            linkId = source + target
        }
    }
}
