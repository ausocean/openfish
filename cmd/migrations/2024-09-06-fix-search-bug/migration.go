package main

import (
	"fmt"

	"github.com/ausocean/openfish/cmd/openfish/ds_client"
	"github.com/ausocean/openfish/cmd/openfish/services"
)

func main() {
	ds_client.Init(false)

	for i := 0; ; i += 100 {
		fmt.Printf("Updated %d species\n", i)

		_, ids, _ := services.GetRecommendedSpecies(100, i, nil, nil, nil)

		for _, id := range ids {
			services.UpdateSpecies(id, nil, nil, nil, nil)
		}

		if len(ids) == 0 {
			break
		}
	}
}
