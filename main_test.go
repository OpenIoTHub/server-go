package main

import (
	"github.com/OpenIoTHub/utils/models"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	token, err := models.DecodeToken("123abc", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSdW5JZCI6IjBiYmM1NWRkLWNhNzAtNGRhNC1iNDAzLTRmODQ3MThkNWEyNCIsIkhvc3QiOiIzNi42My4zOS4xODUiLCJUY3BQb3J0IjozNDMyMCwiS2NwUG9ydCI6MzQzMjAsIlRsc1BvcnQiOjM0MzIxLCJHcnBjUG9ydCI6MzQzMjIsIlVEUEFwaVBvcnQiOjM0MzIxLCJLQ1BBcGlQb3J0IjozNDMyMiwiUGVybWlzc2lvbiI6MSwiZXhwIjoyMDE1OTU3NzU5OTIsIm5iZiI6MTU5NTc0NzE5Mn0.COUJwk3x6RoHZ-ajeOGqTsek9BinwLxjAlgRDriI_Wc")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(token)
}
