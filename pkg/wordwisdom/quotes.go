package wordwisdom

import (
	"math/rand"
)

// quotes contains a list of quotes of wisdom.
var quotes = []string{
	"It may be that that was not the final word of wisdom on the matter.",
	"But surely that is not the last word of wisdom even in war, and certainly not in politics.",
	"Others see in the word of wisdom a teaching function.",
}

// GetRandQuote returns a random quote.
// quotes must not be empty otherwise panic will be generated.
func GetRandQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
