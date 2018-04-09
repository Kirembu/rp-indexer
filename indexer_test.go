package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
)

const elasticURL = "http://localhost:9200"
const indexName = "rp_elastic_test"

func setup(t *testing.T) (*sql.DB, *elastic.Client) {
	testDB, err := ioutil.ReadFile("testdb.sql")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", "postgres://localhost/elastic_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(string(testDB))
	assert.NoError(t, err)

	client, err := elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	assert.NoError(t, err)

	existing := FindPhysicalIndexes(elasticURL, indexName)
	for _, idx := range existing {
		_, err = client.DeleteIndex(idx).Do(context.Background())
		assert.NoError(t, err)
	}

	return db, client
}

func assertQuery(t *testing.T, client *elastic.Client, index string, query elastic.Query, hits []int64) {
	results, err := client.Search().Index(index).Query(query).Sort("id", true).Pretty(true).Do(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(len(hits)), results.Hits.TotalHits)

	if int64(len(hits)) == results.Hits.TotalHits {
		for i, hit := range results.Hits.Hits {
			assert.Equal(t, fmt.Sprintf("%d", hits[i]), hit.Id)
		}
	}
}

func TestIndexing(t *testing.T) {
	batchSize = 4
	db, client := setup(t)

	physicalName, err := CreateNewIndex(elasticURL, indexName)
	assert.NoError(t, err)

	added, deleted, err := IndexContacts(db, elasticURL, physicalName, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, 9, added)
	assert.Equal(t, 0, deleted)

	time.Sleep(2 * time.Second)

	assertQuery(t, client, physicalName, elastic.NewMatchQuery("name", "john"), []int64{5})
	assertQuery(t, client, physicalName, elastic.NewTermQuery("name.keyword", "john doe"), []int64{5})

	assertQuery(t, client, physicalName, elastic.NewMatchQuery("language", "eng"), []int64{1})

	// test contact, not indexed
	assertQuery(t, client, physicalName, elastic.NewMatchQuery("language", "fra"), []int64{})

	assertQuery(t, client, physicalName, elastic.NewMatchQuery("is_blocked", "true"), []int64{4})
	assertQuery(t, client, physicalName, elastic.NewMatchQuery("is_stopped", "true"), []int64{3})

	assertQuery(t, client, physicalName, elastic.NewMatchQuery("org_id", "1"), []int64{1, 3, 4, 5})

	// urn query
	query := elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("urns.scheme", "facebook"),
		elastic.NewMatchQuery("urns.path.keyword", "1000001")))
	assertQuery(t, client, physicalName, query, []int64{9})

	// urn substring query
	query = elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("urns.scheme", "tel"),
		elastic.NewMatchQuery("urns.path", "779")))
	assertQuery(t, client, physicalName, query, []int64{1, 3, 4, 5, 6, 7})

	// urn substring query with more characters (77911)
	query = elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("urns.scheme", "tel"),
		elastic.NewMatchQuery("urns.path", "779"),
		elastic.NewMatchQuery("urns.path", "791"),
		elastic.NewMatchQuery("urns.path", "911")))
	assertQuery(t, client, physicalName, query, []int64{1})

	// match a contact with multiple tel urns
	query = elastic.NewNestedQuery("urns", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("urns.scheme", "tel"),
		elastic.NewMatchQuery("urns.path", "222")))
	assertQuery(t, client, physicalName, query, []int64{1})

	// text query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "17103bb1-1b48-4b70-92f7-1f6b73bd3488"),
		elastic.NewMatchQuery("fields.text", "the rock")))
	assertQuery(t, client, physicalName, query, []int64{1})

	// no tokenizing of field text
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "17103bb1-1b48-4b70-92f7-1f6b73bd3488"),
		elastic.NewMatchQuery("fields.text", "rock")))
	assertQuery(t, client, physicalName, query, []int64{})

	// decimal field range query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "05bca1cd-e322-4837-9595-86d0d85e5adb"),
		elastic.NewRangeQuery("fields.decimal").Gt(10)))
	assertQuery(t, client, physicalName, query, []int64{3})

	// datetime field range query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "e0eac267-463a-4c00-9732-cab62df07b16"),
		elastic.NewRangeQuery("fields.datetime").Lt(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))))
	assertQuery(t, client, physicalName, query, []int64{4})

	// state query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "22d11697-edba-4186-b084-793e3b876379"),
		elastic.NewMatchPhraseQuery("fields.state", "washington")))
	assertQuery(t, client, physicalName, query, []int64{6})

	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "22d11697-edba-4186-b084-793e3b876379"),
		elastic.NewMatchPhraseQuery("fields.state.keyword", "usa > washington")))
	assertQuery(t, client, physicalName, query, []int64{6})

	// doesn't include country
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "22d11697-edba-4186-b084-793e3b876379"),
		elastic.NewMatchPhraseQuery("fields.state", "usa")))
	assertQuery(t, client, physicalName, query, []int64{})

	// district query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "fcab2439-861c-4832-aa54-0c97f38f24ab"),
		elastic.NewMatchPhraseQuery("fields.district", "king")))
	assertQuery(t, client, physicalName, query, []int64{8, 10})

	// phrase matches all
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "fcab2439-861c-4832-aa54-0c97f38f24ab"),
		elastic.NewMatchPhraseQuery("fields.district", "king county")))
	assertQuery(t, client, physicalName, query, []int64{8})

	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "fcab2439-861c-4832-aa54-0c97f38f24ab"),
		elastic.NewMatchPhraseQuery("fields.district.keyword", "usa > washington > king county")))
	assertQuery(t, client, physicalName, query, []int64{8})

	// ward query
	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "a551ade4-e5a0-4d83-b185-53b515ad2f2a"),
		elastic.NewMatchPhraseQuery("fields.ward", "district")))
	assertQuery(t, client, physicalName, query, []int64{9})

	query = elastic.NewNestedQuery("fields", elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("fields.field", "a551ade4-e5a0-4d83-b185-53b515ad2f2a"),
		elastic.NewMatchPhraseQuery("fields.ward.keyword", "usa > washington > king county > central district")))
	assertQuery(t, client, physicalName, query, []int64{9})

	lastModified, err := GetLastModified(elasticURL, physicalName)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2017, 11, 10, 21, 11, 59, 890662000, time.UTC), lastModified.In(time.UTC))

	// map our index over
	err = MapIndexAlias(elasticURL, indexName, physicalName)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	// try a test query to check it worked
	assertQuery(t, client, indexName, elastic.NewMatchQuery("name", "john"), []int64{5})

	// look up our mapping
	physical := FindPhysicalIndexes(elasticURL, indexName)
	assert.Equal(t, physicalName, physical[0])

	// rebuild again
	newIndex, err := CreateNewIndex(elasticURL, indexName)
	assert.NoError(t, err)

	added, deleted, err = IndexContacts(db, elasticURL, newIndex, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, 9, added)
	assert.Equal(t, 0, deleted)

	// remap again
	err = MapIndexAlias(elasticURL, indexName, newIndex)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	assertQuery(t, client, newIndex, elastic.NewMatchQuery("name", "john"), []int64{5})

	// update our database, removing one contact, updating another
	dbUpdate, err := ioutil.ReadFile("testdb_update.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(dbUpdate))
	assert.NoError(t, err)

	added, deleted, err = IndexContacts(db, elasticURL, indexName, lastModified)
	assert.NoError(t, err)
	assert.Equal(t, 1, added)
	assert.Equal(t, 1, deleted)

	time.Sleep(2 * time.Second)

	// should only match new john, old john is gone
	assertQuery(t, client, indexName, elastic.NewMatchQuery("name", "john"), []int64{3})
}
