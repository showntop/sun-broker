package main

import (
	"fmt"
)

type Animal interface {
	Say()
}

type Dog struct {
	name string
}

func (d *Dog) Say() {
	fmt.Println("i'm a dog.....%s", d.name)
}

func main() {
	animals := make(map[string]Animal)
	animals["xiaohua"] = &Dog{"xiaohua"}
	animals["xiaojiji"] = &Dog{"xiaojiji"}
	animals["xiaoheihei"] = &Dog{"xiaoheihei"}
	// animals := []Animal{&Dog{"xiaoming"}, &Dog{"xiaohua"}}
	for _, value := range animals {
		value.Say()
	}
}



1.系统结构

2.http调用接口
create index
delete index
get index
list index

index doc
get doc
delete doc
list fields

search


type SearchRequest struct {
	Query            query.Query       `json:"query"`
	Size             int               `json:"size"`
	From             int               `json:"from"`
	Highlight        *HighlightRequest `json:"highlight"`
	Fields           []string          `json:"fields"`
	Facets           FacetsRequest     `json:"facets"`
	Explain          bool              `json:"explain"`
	Sort             search.SortOrder  `json:"sort"`
	IncludeLocations bool              `json:"includeLocations"`
}

3.支持search类型
exact term match
boolean字段查询
复合boolean查询
合取范式查询
析取范式查询
时间段查询、数字区间查询、词汇段查询
指定docid组查询
模糊词查询
前缀匹配查询
正则表达式查询
术语匹配查询（可指定分析器作为分词依据）match query-match phrase query
区域查询
地理维度查询
Boosting：可以影响评分（description:water name:water^5）
ps:支持查询dsl

4.术语
Analyzer
Character Filter
Term
Token
Tokenizer
Token Filter
Token Stream
Facets
highlight


5.示例

6.模块&特性
Index-mapping 会根据文档类型map到本地结构，包含documentmapping和feildmapping
analizer[char-filter + tokenizer + token-filter]
upsidedown-index
collector
store

分页查询 limit/offset
支持索引所有go数据结构以及json格式数据（扩展其它格式也容易）
高度可定制
支持text/number/date字段类型
支持多种查询方式
Tf-idf 打分
查询结果高亮
方面聚合

Index alias可以切换底层的物理index，还可以实现多个index并发查询，并且merge结果？
提供了一系列命令行工具


索引结构


检索阶段使用的主要结构term_freq_row

7.性能

Bolt stores its keys in byte-sorted order within a bucket. This makes sequential iteration over these keys extremely fast. To iterate over keys we'll use a Cursor:


8.优缺
https://groups.google.com/forum/#!topic/bleve/HUwKS-GDh7o
Release 0.x.x，作者表示api可能会变动
https://medium.com/developers-writing/full-text-search-and-indexing-with-bleve-part-1-bd73599d82ef
基于GO，简单，可嵌入程序内
https://groups.google.com/forum/#!topic/bleve/jbv61ACjYNA
个人提交
几乎每一个阶段都可以定制

https://github.com/blevesearch/bleve/issues/616


https://github.com/couchbase/cbft在bleve之上的分布式检索引擎

文档少，社区不活跃

ES:分布式、更成熟、社区更大更活跃

单机，分布式需手动实现（参考mongo分片、ES分片、kafka副本机制）
