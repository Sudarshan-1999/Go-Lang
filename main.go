package main

func main() {
	a := App{}
	app := a
	app.Initialize(DbUser, DbPass, DbHost, DbPort, DbName)
	app.Run("localhost:8080")

}
