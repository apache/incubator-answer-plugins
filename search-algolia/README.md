# Algolia Search (preview)
> Algolia is a hosted search engine, offering full-text, numerical, and faceted search, capable of delivering real-time 
> results from the first keystroke. Algolia's powerful API lets you quickly and seamlessly implement search within your 
> websites and mobile applications. Our search API powers billions of queries for thousands of companies every month, 
> delivering relevant results in under 100ms anywhere in the world.

## How to use

### Algolia Official website
[Algolia](https://www.algolia.com/)

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/search-algolia
```

### Configuration
- `Application ID` - Your Algolia Application ID.
- `Search-only API Key` - Your Algolia Search-only API Key (public).
- `Admin API Key` - Your Algolia ADMIN API Key (kept private).
- `Index name prefix` - This prefix will be prepended to your index names.
- `Algolia logo` - Algolia requires that you keep the logo if you are using a free plan.

### Note
- If you have a large amount of data, it will be synchronized to algolia server auto when plugin configuration completed. If you need to know the specific progress, you need to check the console log information yourself.
