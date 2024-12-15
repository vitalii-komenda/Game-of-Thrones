package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/vitalii-komenda/got/entities"
	"github.com/vitalii-komenda/got/utils"
)

func SendElasticSearchRequest(term string) ([]entities.CharacterEntryElastic, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     term,
				"fields":    []string{"character_name^2", "actor_name", "siblings"},
				"fuzziness": "AUTO",
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", utils.MustGetEnvOrPanic("ELASTICSEARCH_HOST")+"/character_details/_search?pretty", bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	var value []entities.CharacterEntryElastic
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceJSON, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var characterEntry entities.CharacterEntryElastic
		err = json.Unmarshal(sourceJSON, &characterEntry)
		if err != nil {
			return nil, err
		}
		value = append(value, characterEntry)
	}

	return value, nil
}
