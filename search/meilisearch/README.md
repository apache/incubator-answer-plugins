# Meilisearch (preview)
> Meilisearch is A lightning-fast search engine that fits effortlessly into your apps, websites, and workflow 🔍

## How to use

### Build
```bash
./answer build --with github.com/answerdev/plugins/search/meilisearch
```

### Configuration
- `Host` - Meilisearch connection address, such as http://127.0.0.1:7700
- `ApiKey` - Meilisearch api key
- `IndexName` - The index answer will use. Default is `answer_post`
- `Async` - Should answer use async mode to send data to Meilisearch. Default is `false`. use Async means you will not get any error message if Meilisearch task failed. 
