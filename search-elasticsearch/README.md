# Elasticsearch Engine (preview)
> The default Answer uses a built-in database such as MySQL as its search engine. 
> However, when dealing with large amounts of data, the speed and accuracy of searches can be affected. 
> Therefore, we provide a plugin that uses Elasticsearch as the search engine, which greatly improves search speed and accuracy.

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/search-elasticsearch
```

### Configuration
- `Endpoints` - Elasticsearch connection address, such as http://127.0.0.1:9200 or multiple addresses separated by ','
- `Username` - Elasticsearch username
- `Password` - Elasticsearch password

## Note
- Only support Elasticsearch 7.x
- Index name is `answer_post`. It will create automatically if not exists. 
- You also can create index manually if you want to specify `search_analyzer` or other settings(replicas and shards).
