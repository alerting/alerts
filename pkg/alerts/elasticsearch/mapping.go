package elasticsearch

// ElasticSearch mapping definition.
const mapping = `{
	"settings": {
	  "analysis": {
		 "analyzer": {
			"folding": {
			  "tokenizer": "standard",
			  "filter": ["lowercase", "asciifolding"]
			}
		 },
		 "normalizer": {
			"keyword_normalizer": {
			  "type": "custom",
			  "filter": ["lowercase", "asciifolding"]
			}
		 }
	  }
	},
	"mappings": {
	  "doc": {
		 "dynamic": false,
		 "properties": {
			"identifier": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"sender": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"sent": { "type": "date" },
			"status": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"message_type": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"source": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"scope": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"restriction": { "type": "text" },
			"addresses": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"codes": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"note": { "type": "text", "analyzer": "folding" },
			"references": {
			  "type": "nested",
			  "dynamic": false,
			  "properties": {
				 "sender": { "type": "keyword", "normalizer": "keyword_normalizer" },
				 "sent": { "type": "date" },
				 "indentifier": { "type": "keyword", "normalizer": "keyword_normalizer" },
				 "id": { "type": "keyword", "normalizer": "keyword_normalizer" }
			  }
			},
			"incidents": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"superseded": { "type": "boolean" },

			"language": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"categories": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"event": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"response_types": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"urgency": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"severity": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"certainty": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"audience": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"event_codes": { "type": "object" },
			"effective": { "type": "date" },
			"onset": { "type": "date" },
			"expires": { "type": "date" },
			"sender_name": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"headline": { "type": "text", "analyzer": "folding" },
			"description": { "type": "text", "analyzer": "folding" },
			"instruction": { "type": "text", "analyzer": "folding" },
			"web": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"contact": { "type": "keyword", "normalizer": "keyword_normalizer" },
			"parameters": { "type": "object" },
			"resources": {
			  "type": "nested",
			  "dynamic": false,
			  "properties": {
				 "description": { "type": "text", "analyzer": "folding" },
				 "mime_type": { "type": "keyword", "normalizer": "keyword_normalizer" },
				 "size": { "type": "integer" },
				 "uri": { "type": "keyword", "normalizer": "keyword_normalizer" },
				 "derefUri": { "type": "binary" },
				 "digest": { "type": "keyword", "normalizer": "keyword_normalizer" }
			  }
			},
			"areas": {
			  "type": "nested",
			  "dynamic": false,
			  "properties": {
				 "description": { "type": "text", "analyzer": "folding" },
				 "polygons": { "type": "geo_shape", "ignore_malformed": true },
				 "circles": { "type": "geo_shape", "ignore_malformed": true },
				 "geocodes": { "type": "object" },
				 "altitude": { "type": "float" },
				 "ceiling": { "type": "float" }
			  }
			},

			"_object": {
			  "type": "join",
			  "relations": {
				 "alert": "info"
			  }
			}
		 }
	  }
	}
 }`
