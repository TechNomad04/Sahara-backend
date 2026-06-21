package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func GenerateEmbedding(text string) ([]float32, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("HF_TOKEN")).
		SetBody(map[string]string{
			"inputs": text,
		}).
		Post("https://router.huggingface.co/hf-inference/models/BAAI/bge-small-en-v1.5/pipeline/feature-extraction")

	if err != nil {
		return nil, err
	}

	var embedding []float32

	err = json.Unmarshal(resp.Body(), &embedding)
	if err != nil {
		return nil, err
	}

	return embedding, nil
}

func main() {
	embedding, err := GenerateEmbedding("What is machine learning?")
	if err != nil {
		panic(err)
	}

	fmt.Println("Dimension:", len(embedding))
}