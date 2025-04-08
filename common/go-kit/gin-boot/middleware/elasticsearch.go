package middleware

import (
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

func (mw *Middleware) initElasticsearch() {
	mw.eslock = &sync.Mutex{}
	mw.elasticsearch = make(map[string]*elasticsearch.Client)
	for dbName, dbConfig := range mw.conf.App.Middleware.Elasticsearch {
		if dbConfig.Startup == true {
			conf := mw.conf.App.Middleware.Elasticsearch[dbName]
			elasticsearchConfig := mw.initElasticsearchConfigs(&conf)
			var err error
			mw.elasticsearch[dbName], err = elasticsearch.NewClient(*elasticsearchConfig)
			util.ChkErr(err)
		}
	}
}

func (mw *Middleware) initElasticsearchConfigs(opt *conf_type.Elasticsearch) *elasticsearch.Config {
	esOptions := &elasticsearch.Config{
		Addresses: opt.Addresses,
		Username:  opt.Username,
		Password:  opt.Password,
	}
	return esOptions
}

// DefaultElasticsearch DefaultElasticsearch
func (mw *Middleware) DefaultElasticsearch() *elasticsearch.Client {
	return mw.Elasticsearch(defaultName)
}

// Elasticsearch Elasticsearch
func (mw *Middleware) Elasticsearch(name string) *elasticsearch.Client {
	mw.eslock.Lock()
	defer mw.eslock.Unlock()
	if db, ok := mw.elasticsearch[name]; ok {
		return db
	}
	return nil
}
