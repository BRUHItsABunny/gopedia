# gopedia
Wikipedia wrapper in Go

## Example
```
client := gopedia.GetClient()
result, err, err2 := client.SearchBasic("Rabbit")
if result != nil {
    res, _ := json.Marshal(result)
    fmt.Println(string(res))
    result2, err, err2 := client.GetPage(result[0].Title)
    if result2 != nil {
        res2, _ := json.Marshal(result2)
    	fmt.Println(string(res2))
    }
    if err != nil {
    	fmt.Println(err)
    }
    if err2 != nil {
    	fmt.Println(err2)
    }
}
if err != nil {
	fmt.Println(err)
}
if err2 != nil {
	fmt.Println(err2)
}
```
