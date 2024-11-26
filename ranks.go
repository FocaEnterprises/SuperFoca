package main

import (
	"log"
)

type Rank struct {
	Id      uint    `sql:"id"`
	UserId  string  `sql:"user_id"`
	GuildId string  `sql:"guild_id"`
	IQ      float64 `sql:"iq"`
	TitleId int     `sql:"title_id"`
}

func createRank(userId, guildId string) *Rank {
	createQuery, err := db.Prepare("INSERT INTO ranks (user_id, guild_id) VALUES ($1, $2)")

	if err != nil {
		log.Printf("couldn't prepare query: %s", err)
		return nil
	}

	defer createQuery.Close()

	_, err = createQuery.Exec(userId, guildId)

	if err != nil {
		log.Printf("failed to insert: %s", err)
		return nil
	}

	rank := findRank(userId, guildId)

	return rank
}

func updateRank(rank Rank, iq float64) {
	updateQuery, err := db.Prepare("UPDATE ranks SET iq = $1 WHERE id = $2")

	if err != nil {
		log.Printf("failed to update rank: %s\n", err)
		return
	}

	defer updateQuery.Close()

	_, err = updateQuery.Exec(iq, rank.Id)

	if err != nil {
		log.Printf("failed executing query: %s", err)
		return
	}
}

func findRank(userId, guildId string) *Rank {
	readQuery, err := db.Query("SELECT * FROM ranks WHERE user_id = $1 AND guild_id = $2", userId, guildId)

	if err != nil {
		log.Printf("couldn't query %s and %s: %s", userId, guildId, err)
		return nil
	}

	defer readQuery.Close()

	if !readQuery.Next() {
		return nil
	}

	var rank Rank

	err = readQuery.Scan(&rank.Id, &rank.UserId, &rank.GuildId, &rank.IQ, &rank.TitleId)

	if err != nil {
		log.Printf("failed scanning query: %s", err)
		return nil
	}

	return &rank
}

func getRanking(page int, limit int) []*Rank {
	readQuery, err := db.Query("SELECT * FROM ranks ORDER BY iq DESC LIMIT $1 OFFSET $2", limit, page*limit)

	if err != nil {
		log.Printf("failed querying table rank: %s", err)
		return nil
	}

	defer readQuery.Close()

	ranking := []*Rank{}

	for readQuery.Next() {
		if err := readQuery.Err(); err != nil {
			log.Printf("failed iterating rows: %s", err)
			return nil
		}

		rank := &Rank{}

		err := readQuery.Scan(&rank.Id, &rank.UserId, &rank.GuildId, &rank.IQ)

		if err != nil {
			log.Printf("failed scanning row: %s", err)
			return nil
		}

		ranking = append(ranking, rank)
	}

	return ranking
}
