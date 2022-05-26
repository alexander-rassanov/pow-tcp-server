package wordwisdom

import (
	"math/rand"
)

var quotes = []string{
	"It may be that that was not the final word of wisdom on the matter.",
	"But surely that is not the last word of wisdom even in war, and certainly not in politics.",
	"Others see in the word of wisdom a teaching function.",
}

func getRandQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
