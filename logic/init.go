package logic

func Init() {

	// Define regex
	regex.Username = `[a-zA-Z0-9_]$`
	regex.Email = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

}
