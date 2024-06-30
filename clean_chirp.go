package main

import "strings"

func cleanChirp(chirp string) string {
	var cleanedChirp strings.Builder

	words := strings.Split(chirp, " ")

	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax": // Profone words
			cleanedChirp.WriteString("****")
		default:
			cleanedChirp.WriteString(word)
		}

		if isLastIndex := i == len(words)-1; !isLastIndex {
			cleanedChirp.WriteByte(' ')
		}
	}

	return cleanedChirp.String()
}
