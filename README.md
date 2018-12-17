# gopedia
Wikipedia wrapper in Go

## Example
```
client := gopedia.GetClient()
result, err, err2 := client.SearchBasic("Rabbit")
if result != nil {
	res, _ := json.Marshal(result)
	fmt.Println(string(res))
}
if err != nil {
	fmt.Println(err)
}
if err2 != nil {
	fmt.Println(err2)
}
```
