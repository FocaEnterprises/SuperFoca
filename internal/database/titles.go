package database

import "log"

type Title struct {
	Id    int     `sql:"id"`
	Title string  `sql:"title"`
	IQ    float64 `sql:"float"`
}

func ReadTitleFromRank(rank Rank) *Title {
	query, err := db.Query("SELECT * FROM titles WHERE id = $1", rank.TitleId)

	defer query.Close()

	if err != nil {
		log.Printf("failed querying title: %s", err)
		return nil
	}

	if !query.Next() {
		log.Printf("failed preparing scan: %s", query.Err())
		return nil
	}

	title := &Title{}

	query.Scan(&title.Id, &title.Title, &title.IQ)

	return title
}
