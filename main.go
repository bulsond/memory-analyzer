package main

func main() {
	app, err := NewApp(3, 10)
	if err != nil {
		panic(err)
	}

	app.Run()
}
