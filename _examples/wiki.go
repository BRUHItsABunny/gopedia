package main

import (
	"encoding/json"
	"fmt"
	"gopedia"
)

func main() {
	// Init client
	client := gopedia.GetClient()

	// Search topic
	result, err := client.SearchBasic("Rabbit")
	if err == nil {
		// Print result
		res, _ := json.Marshal(result)
		fmt.Println(string(res))

		// Get full page
		result2, err := client.GetPage(result[0].Title)
		if err == nil {

			// Print result2
			res2, _ := json.Marshal(result2)
			fmt.Println(string(res2))
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
