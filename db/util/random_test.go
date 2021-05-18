package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(10)
	max := int64(10000)
	for i := 0; i < 100; i++ {
		randInt := RandomInt(min, max)
		require.GreaterOrEqual(t, randInt, min)
		require.LessOrEqual(t, randInt, max)
	}
}

func TestRandomOwner(t *testing.T) {
	names := []string{"Peter", "Hudson", "Lillian", "Thomson", "Nicola", "Robertson", "Oliver", "May", "Matt", "Roberts", "Warren", "Forsyth", "Sophie", "James", "Sam", "Johnston", "Amy", "Poole", "Adam", "Duncan", "Zoe", "Alsop", "Tracey", "Welch", "Keith", "Gill", "Emma", "Arnold", "Evan", "Miller", "Joseph", "Greene", "Sean", "Reid", "Ruth", "Duncan", "Dorothy", "Poole", "Cameron", "Morgan", "Jake", "Ellison", "Luke", "Cornish", "Sophie", "Chapman", "Jane", "Grant", "Piers", "Burgess", "Ryan", "Bower", "Dan", "MacLeod", "Megan", "Miller", "Isaac", "Butler", "Peter", "Watson", "Owen", "Lyman", "Mary", "Davidson", "Gavin", "Knox", "Dan", "Wilkins", "Owen", "White", "Paul", "Graham", "Andrea", "Dickens", "Leonard", "Dickens", "Carolyn", "Piper", "Caroline", "Sanderson", "David", "Arnold", "Anthony", "Henderson", "Charles", "Blake", "Joseph", "Springer", "Cameron", "Bower", "Liam", "Walker", "Ruth", "Burgess", "Katherine", "MacDonald", "Adam", "Davies", "Maria", "Morrison", "Christopher", "Scott", "Deirdre", "Peake", "Robert", "Edmunds", "Gavin", "Rutherford", "Amelia", "Peters", "Audrey", "Morrison", "Evan", "Smith", "Benjamin", "Martin", "Alan", "Mitchell", "Sam", "Young", "Dan", "Hunter", "Yvonne", "Miller", "Robert", "Payne", "Sophie", "Alsop", "Rachel", "Allan", "Vanessa", "Rutherford", "Amelia", "Sanderson", "Lauren", "Cornish", "Carolyn", "Churchill", "Molly", "Powell", "Sean", "Wright", "Keith", "Mitchell", "Andrea", "Dyer", "Jane", "Langdon", "Dylan", "Forsyth", "Diane", "Langdon"}
	for i := 0; i < 100; i++ {
		randName := RandomOwner()
		require.Contains(t, names, randName)
	}
}
