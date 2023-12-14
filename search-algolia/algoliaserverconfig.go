package algolia

var AlgoliaSearchServerConfig = []byte(`
{
  "minWordSizefor1Typo": 4,
  "minWordSizefor2Typos": 8,
  "hitsPerPage": 20,
  "maxValuesPerFacet": 100,
  "searchableAttributes": [
    "content",
    "title"
  ],
  "numericAttributesToIndex": null,
  "attributesToRetrieve": null,
  "unretrievableAttributes": null,
  "optionalWords": null,
  "attributesForFaceting": [
    "status",
    "tags",
    "type",
    "user_id"
  ],
  "attributesToSnippet": null,
  "attributesToHighlight": null,
  "paginationLimitedTo": 1000,
  "attributeForDistinct": null,
  "exactOnSingleWordQuery": "attribute",
  "ranking": [
    "desc(active)",
    "desc(score)",
    "desc(created)",
    "typo",
    "geo",
    "words",
    "filters",
    "proximity",
    "attribute",
    "exact",
    "custom"
  ],
  "customRanking": [
    "desc(title)",
    "desc(content)"
  ],
  "separatorsToIndex": "",
  "removeWordsIfNoResults": "none",
  "queryType": "prefixLast",
  "highlightPreTag": "<em>",
  "highlightPostTag": "</em>",
  "alternativesAsExact": [
    "ignorePlurals",
    "singleWordSynonym"
  ]
}
`)
