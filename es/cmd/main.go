package main

import (
	// 标准库
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	// 三方包
	"github.com/elastic/go-elasticsearch/v8"
	// 本地包
	"es/config"
	. "es/pkg/es"
)

func main() {
	config.Init()

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	// 配置es地址
	Init(config.Conf.Elasticsearch.Host, config.Conf.Elasticsearch.Username, config.Conf.Elasticsearch.Password)
	es := NewESClient()

	// 1. 查询集群信息
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	list := []map[string]interface{}{
		{
			"title":   "test one",
			"desc":    "test one desc",
			"url":     "https://www.baiud.com",
			"context": "test one context",
			"click":   "1",
		},
		{
			"title":   "test two",
			"desc":    "test two desc",
			"url":     "https://www.baiud.com",
			"context": "test two context",
			"click":   "1",
		},
	}

	// 2. 创建索引
	for i, article := range list {
		wg.Add(1)

		go func(i int, article map[string]interface{}) {
			defer wg.Done()

			// 序列化
			newArticle, err := json.Marshal(article)

			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(string(newArticle)),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))

					log.Printf("Status: %s", res.Status())
					log.Printf("Result: %s", r["result"])
					log.Printf("Version: %v", int(r["_version"].(float64)))
					log.Printf("res: %v", r)
				}
			}

		}(i, article)
	}
	wg.Wait()

	/*




	 */

	log.Println(strings.Repeat("-", 37))
	log.Println("查询全部")
	// 3. 查询信息
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.执行查询
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	/*




	 */

	log.Println(strings.Repeat("=", 37))

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("更新文档")
		var updateBuf bytes.Buffer
		// 5. 更新文档
		// Build the request body.
		updateData := map[string]interface{}{
			"doc": map[string]interface{}{
				"click": "5",
			},
		}

		if err := json.NewEncoder(&updateBuf).Encode(updateData); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}

		res, err = es.Update(
			"test",
			"2",
			&updateBuf,
			es.Update.WithContext(context.Background()),
		)

		defer res.Body.Close()

		log.Println(fmt.Sprintf("Document updated, status code: %d", res.StatusCode))
	}()
	wg.Wait()

	// 防止下面精确查询获取不到修改的数据
	time.Sleep(time.Second)
	/*




	 */

	log.Println(strings.Repeat("=", 37))
	log.Println("精确查询")
	// 4. 精确查询
	// Build the request body.
	query = map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"_id": map[string]interface{}{
					"value": "2",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.执行查询
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	// Print the ID and document source for each hit.
	log.Printf(" * ID=%v, %v",
		r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_id"],
		r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_source"],
	)

	/*




	 */

	log.Println(strings.Repeat("=", 37))
	log.Println("删除一条")
	// 6. 删除一条
	// Build the request body.
	res, err = es.Delete(
		"test",
		"1",
		es.Delete.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	log.Println(fmt.Sprintf("Document deleted, status code: %d", res.StatusCode))

	/*




	 */

	log.Println(strings.Repeat("=", 37))
	log.Println("7. 删除全部")
	// 7. 删除全部
	//Build the request body.
	res, err = es.Indices.Delete(
		[]string{"test"},
		es.Indices.Delete.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	log.Println(fmt.Sprintf("Document deleted all, status code: %d", res.StatusCode))
}
