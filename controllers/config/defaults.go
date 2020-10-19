package config

const defaultTopic = "platform.inventory.events"
const defaultConnectCluster = "my-connect-cluster"

const defaultConnectorTemplate = `{
	"connector.class": "io.confluent.connect.jdbc.JdbcSinkConnector",
	"tasks.max": "{{.TasksMax}}",
	"topics": "{{.Topic}}",
	"key.converter": "org.apache.kafka.connect.storage.StringConverter",
	"value.converter": "org.apache.kafka.connect.json.JsonConverter",
	"value.converter.schemas.enable": false,
	"connection.url": "jdbc:postgresql://{{.DBHostname}}:{{.DBPort}}/{{.DBName}}",
	"connection.user": "{{.DBUser}}",
	"connection.password": "{{.DBPassword}}",
	"dialect.name": "EnhancedPostgreSqlDatabaseDialect",
	"auto.create": false,
	"insert.mode": "upsert",
	"delete.enabled": true,
	"batch.size": "{{.BatchSize}}",
	"table.name.format": "inventory.{{.TableName}}",
	"pk.mode": "record_key",
	"pk.fields": "id",
	"fields.whitelist": "account,display_name,tags,updated,created,stale_timestamp,system_profile",

	{{ if eq .InsightsOnly "true" }}
	"transforms": "timestampFilterShort,insightsFilter,deleteToTombstone,extractHost,systemProfileFilter,systemProfileToJson,tagsToJson,injectSchemaKey,injectSchemaValue",
	"transforms.insightsFilter.type":"com.redhat.insights.kafka.connect.transforms.Filter",
	"transforms.insightsFilter.predicate": "!!record.headers().lastWithName('insights_id').value()",
	{{ else  }}
	"transforms": "timestampFilterLong,deleteToTombstone,extractHost,systemProfileFilter,systemProfileToJson,tagsToJson,injectSchemaKey,injectSchemaValue",
	{{ end }}

	"transforms.timestampFilterLong.type":"com.redhat.insights.kafka.connect.transforms.Filter",
	"transforms.timestampFilterLong.predicate": "(Date.now() - record.timestamp()) < 45 * 24 * 60 * 60 * 1000",
	"transforms.timestampFilterShort.type":"com.redhat.insights.kafka.connect.transforms.Filter",
	"transforms.timestampFilterShort.predicate": "(Date.now() - record.timestamp()) < 21 * 24 * 60 * 60 * 1000",
	"transforms.deleteToTombstone.type":"com.redhat.insights.kafka.connect.transforms.DropIf$Value",
	"transforms.deleteToTombstone.predicate": "'delete'.equals(record.headers().lastWithName('event_type').value())",
	"transforms.extractHost.type":"org.apache.kafka.connect.transforms.ExtractField$Value",
	"transforms.extractHost.field":"host",
	"transforms.systemProfileFilter.type": "com.redhat.insights.kafka.connect.transforms.FilterFields$Value",
	"transforms.systemProfileFilter.field": "system_profile",
	"transforms.systemProfileFilter.allowlist": "sap_system,sap_sids",
	"transforms.systemProfileToJson.type": "com.redhat.insights.kafka.connect.transforms.FieldToJson$Value",
	"transforms.systemProfileToJson.originalField": "system_profile",
	"transforms.systemProfileToJson.destinationField": "system_profile",
	"transforms.tagsToJson.type": "com.redhat.insights.kafka.connect.transforms.FieldToJson$Value",
	"transforms.tagsToJson.originalField": "tags",
	"transforms.tagsToJson.destinationField": "tags",
	"transforms.injectSchemaKey.type": "com.redhat.insights.kafka.connect.transforms.InjectSchema$Key",
	"transforms.injectSchemaKey.schema": "{\"type\":\"string\",\"optional\":false, \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=uuid\"}",
	"transforms.injectSchemaValue.type": "com.redhat.insights.kafka.connect.transforms.InjectSchema$Value",
	"transforms.injectSchemaValue.schema": "{\"type\":\"struct\",\"fields\":[{\"type\":\"string\",\"optional\":false,\"field\":\"account\"},{\"type\":\"string\",\"optional\":false,\"field\":\"display_name\"},{\"type\":\"string\",\"optional\":false,\"field\":\"tags\", \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=jsonb\"},{\"type\":\"string\",\"optional\":false,\"field\":\"updated\", \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=timestamptz\"},{\"type\":\"string\",\"optional\":false,\"field\":\"created\", \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=timestamptz\"},{\"type\":\"string\",\"optional\":false,\"field\":\"stale_timestamp\", \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=timestamptz\"},{\"type\":\"string\",\"optional\":false,\"field\":\"system_profile\", \"name\": \"com.redhat.cloud.inventory.syndication.pgtype=jsonb\"}],\"optional\":false}",

	"errors.tolerance": "all",
	"errors.deadletterqueue.topic.name": "platform.cyndi.dlq",
	"errors.deadletterqueue.topic.replication.factor": 1,
	"errors.deadletterqueue.context.headers.enable":true,
	"errors.retry.delay.max.ms": 60000,
	"errors.retry.timeout": 600000,
	"errors.log.enable":true,
	"errors.log.include.messages":true
}`

const defaultConnectorTasksMax int64 = 16
const defaultConnectorBatchSize int64 = 100
const defaultConnectorMaxAge int64 = 45

const defaultDBTableInitScript = `
CREATE TABLE inventory.{{.TableName}} (
	id uuid PRIMARY KEY,
	account character varying(10) NOT NULL,
	display_name character varying(200) NOT NULL,
	tags jsonb NOT NULL,
	updated timestamp with time zone NOT NULL,
	created timestamp with time zone NOT NULL,
	stale_timestamp timestamp with time zone NOT NULL,
	system_profile jsonb NOT NULL
);

CREATE INDEX {{.TableName}}_account_index ON inventory.{{.TableName}}
(account);

CREATE INDEX {{.TableName}}_display_name_index ON inventory.{{.TableName}}
(display_name);

CREATE INDEX {{.TableName}}_tags_index ON inventory.{{.TableName}} USING GIN
(tags JSONB_PATH_OPS);

CREATE INDEX {{.TableName}}_stale_timestamp_index ON
inventory.{{.TableName}} (stale_timestamp);

CREATE INDEX {{.TableName}}_system_profile_index ON inventory.{{.TableName}}
USING GIN (system_profile JSONB_PATH_OPS);
`

var defaultValidationConfig = ValidationConfiguration{
	Interval:            60,
	AttemptsThreshold:   5,
	PercentageThreshold: 5,
}

var defaultValidationConfigInit = ValidationConfiguration{
	Interval:            15,
	AttemptsThreshold:   10,
	PercentageThreshold: 5,
}