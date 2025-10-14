# Elasticsearch Setup

Docker Compose setup untuk Elasticsearch standalone (tanpa integrasi dengan analytic-service).

## Fitur

- **Elasticsearch 8.11.3** - Latest stable version
- **Single-node setup** - Untuk development
- **Kibana (Optional)** - Web UI untuk Elasticsearch
- **Health checks** - Auto-monitoring container health
- **Persistent storage** - Data tersimpan di Docker volume

## Quick Start

### 1. Start Elasticsearch Saja

```bash
docker-compose -f docker-compose.elasticsearch.yml up -d elasticsearch
```

**Akses:**
- Elasticsearch: http://localhost:9200
- Health Check: http://localhost:9200/_cluster/health

### 2. Start Elasticsearch + Kibana

```bash
docker-compose -f docker-compose.elasticsearch.yml --profile kibana up -d
```

**Akses:**
- Elasticsearch: http://localhost:9200
- Kibana: http://localhost:5601

## Docker Compose Configuration

### Services

#### Elasticsearch
- **Image**: `docker.elastic.co/elasticsearch/elasticsearch:8.11.3`
- **Ports**:
  - 9200 (HTTP API)
  - 9300 (Transport)
- **Memory**: 512MB (min & max)
- **Security**: Disabled untuk development
- **Data**: Persistent volume `elasticsearch-data`

#### Kibana (Optional)
- **Image**: `docker.elastic.co/kibana/kibana:8.11.3`
- **Port**: 5601
- **Depends on**: Elasticsearch
- **Profile**: `kibana` (run dengan `--profile kibana`)

## Commands

### Start Services

```bash
# Elasticsearch only
docker-compose -f docker-compose.elasticsearch.yml up -d elasticsearch

# Elasticsearch + Kibana
docker-compose -f docker-compose.elasticsearch.yml --profile kibana up -d

# With logs
docker-compose -f docker-compose.elasticsearch.yml up elasticsearch
```

### Stop Services

```bash
# Stop all
docker-compose -f docker-compose.elasticsearch.yml down

# Stop and remove volumes (⚠️ akan menghapus data)
docker-compose -f docker-compose.elasticsearch.yml down -v
```

### View Logs

```bash
# All services
docker-compose -f docker-compose.elasticsearch.yml logs -f

# Elasticsearch only
docker-compose -f docker-compose.elasticsearch.yml logs -f elasticsearch

# Kibana only
docker-compose -f docker-compose.elasticsearch.yml logs -f kibana
```

### Check Status

```bash
# Container status
docker-compose -f docker-compose.elasticsearch.yml ps

# Elasticsearch health
curl http://localhost:9200/_cluster/health?pretty
```

## Testing Elasticsearch

### 1. Check Cluster Health

```bash
curl http://localhost:9200/_cluster/health?pretty
```

**Response:**
```json
{
  "cluster_name" : "repath-cluster",
  "status" : "green",
  "number_of_nodes" : 1,
  "number_of_data_nodes" : 1
}
```

### 2. Create Index

```bash
curl -X PUT "localhost:9200/test-index" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  }
}'
```

### 3. Index Document

```bash
curl -X POST "localhost:9200/test-index/_doc" -H 'Content-Type: application/json' -d'
{
  "user": "john_doe",
  "message": "Hello Elasticsearch!",
  "timestamp": "2025-10-14T10:30:00Z"
}'
```

### 4. Search Documents

```bash
curl "localhost:9200/test-index/_search?pretty"
```

### 5. Delete Index

```bash
curl -X DELETE "localhost:9200/test-index"
```

## Using with Python

### Install Client

```bash
pip install elasticsearch
```

### Example Code

```python
from elasticsearch import Elasticsearch

# Connect to Elasticsearch
es = Elasticsearch(["http://localhost:9200"])

# Check connection
if es.ping():
    print("Connected to Elasticsearch!")

# Create index
es.indices.create(index="analytics", ignore=400)

# Index document
doc = {
    "user": "john_doe",
    "event": "page_view",
    "page": "/products",
    "timestamp": "2025-10-14T10:30:00Z"
}
es.index(index="analytics", document=doc)

# Search
result = es.search(index="analytics", query={"match_all": {}})
print(f"Found {result['hits']['total']['value']} documents")

for hit in result['hits']['hits']:
    print(hit['_source'])
```

## Configuration Details

### Environment Variables

```yaml
- node.name=elasticsearch              # Node name
- cluster.name=repath-cluster          # Cluster name
- discovery.type=single-node           # Single-node mode
- bootstrap.memory_lock=true           # Lock memory
- ES_JAVA_OPTS=-Xms512m -Xmx512m      # JVM heap size
- xpack.security.enabled=false         # Disable security (dev only)
```

### Memory Configuration

- **Minimum**: 512MB
- **Maximum**: 512MB
- **Recommended untuk production**: 2GB+

### Ulimits

```yaml
memlock:
  soft: -1
  hard: -1
nofile:
  soft: 65536
  hard: 65536
```

## Troubleshooting

### 1. Container Won't Start

**Check logs:**
```bash
docker-compose -f docker-compose.elasticsearch.yml logs elasticsearch
```

**Common issues:**
- Insufficient memory
- Port already in use
- Docker memory limit too low

**Solution:**
```bash
# Increase Docker memory to at least 4GB
# Docker Desktop > Settings > Resources > Memory
```

### 2. Cluster Health Yellow/Red

```bash
# Check cluster health
curl http://localhost:9200/_cluster/health?pretty

# Check indices
curl http://localhost:9200/_cat/indices?v
```

### 3. Port Already in Use

```bash
# Check what's using port 9200
lsof -i :9200

# Kill process if needed
kill -9 <PID>
```

### 4. Reset Everything

```bash
# Stop and remove all data
docker-compose -f docker-compose.elasticsearch.yml down -v

# Remove Docker volumes
docker volume rm elasticsearch-data

# Start fresh
docker-compose -f docker-compose.elasticsearch.yml up -d
```

## Performance Tips

### 1. Memory Settings

```yaml
# For development
ES_JAVA_OPTS=-Xms512m -Xmx512m

# For production
ES_JAVA_OPTS=-Xms2g -Xmx2g
```

### 2. Index Settings

```json
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "refresh_interval": "30s"
  }
}
```

### 3. Bulk Operations

Use bulk API untuk insert banyak documents:

```python
from elasticsearch import helpers

actions = [
    {
        "_index": "analytics",
        "_source": {"user": f"user_{i}", "event": "page_view"}
    }
    for i in range(1000)
]

helpers.bulk(es, actions)
```

## Security Notes

⚠️ **Development Setup Only**

Security disabled untuk kemudahan development:
- No authentication
- No SSL/TLS
- HTTP only

**Untuk Production:**
- Enable xpack.security
- Configure SSL/TLS
- Set up authentication
- Use firewall rules

## Useful Elasticsearch APIs

```bash
# Cluster info
curl localhost:9200

# Cluster health
curl localhost:9200/_cluster/health?pretty

# Node stats
curl localhost:9200/_nodes/stats?pretty

# All indices
curl localhost:9200/_cat/indices?v

# Index mapping
curl localhost:9200/your-index/_mapping?pretty

# Index settings
curl localhost:9200/your-index/_settings?pretty
```

## Resources

- [Elasticsearch Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [Elasticsearch Python Client](https://elasticsearch-py.readthedocs.io/)
- [Kibana Documentation](https://www.elastic.co/guide/en/kibana/current/index.html)
- [Docker Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)

## Next Steps

1. **Start Elasticsearch**: `docker-compose -f docker-compose.elasticsearch.yml up -d`
2. **Test connection**: `curl localhost:9200`
3. **Create your first index**: Follow examples above
4. **Optional: Start Kibana** untuk web interface
5. **Integrate with your application**
