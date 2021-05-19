package util

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(num int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	 b := make([]rune, num)
	for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomOwner() string {
	names := []string{"Peter", "Hudson", "Lillian", "Thomson", "Nicola", "Robertson", "Oliver", "May", "Matt", "Roberts", "Warren", "Forsyth", "Sophie", "James", "Sam", "Johnston", "Amy", "Poole", "Adam", "Duncan", "Zoe", "Alsop", "Tracey", "Welch", "Keith", "Gill", "Emma", "Arnold", "Evan", "Miller", "Joseph", "Greene", "Sean", "Reid", "Ruth", "Duncan", "Dorothy", "Poole", "Cameron", "Morgan", "Jake", "Ellison", "Luke", "Cornish", "Sophie", "Chapman", "Jane", "Grant", "Piers", "Burgess", "Ryan", "Bower", "Dan", "MacLeod", "Megan", "Miller", "Isaac", "Butler", "Peter", "Watson", "Owen", "Lyman", "Mary", "Davidson", "Gavin", "Knox", "Dan", "Wilkins", "Owen", "White", "Paul", "Graham", "Andrea", "Dickens", "Leonard", "Dickens", "Carolyn", "Piper", "Caroline", "Sanderson", "David", "Arnold", "Anthony", "Henderson", "Charles", "Blake", "Joseph", "Springer", "Cameron", "Bower", "Liam", "Walker", "Ruth", "Burgess", "Katherine", "MacDonald", "Adam", "Davies", "Maria", "Morrison", "Christopher", "Scott", "Deirdre", "Peake", "Robert", "Edmunds", "Gavin", "Rutherford", "Amelia", "Peters", "Audrey", "Morrison", "Evan", "Smith", "Benjamin", "Martin", "Alan", "Mitchell", "Sam", "Young", "Dan", "Hunter", "Yvonne", "Miller", "Robert", "Payne", "Sophie", "Alsop", "Rachel", "Allan", "Vanessa", "Rutherford", "Amelia", "Sanderson", "Lauren", "Cornish", "Carolyn", "Churchill", "Molly", "Powell", "Sean", "Wright", "Keith", "Mitchell", "Andrea", "Dyer", "Jane", "Langdon", "Dylan", "Forsyth", "Diane", "Langdon"}
	return names[rand.Intn(len(names))]
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, KRW}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	name := RandomOwner()
	number := RandomInt(1, 10000)
	return fmt.Sprintf("%s%d@email.com", name, number)
}
