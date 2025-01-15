package check

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"virtuals-tracker/cmd/global"
	"virtuals-tracker/database"
	"virtuals-tracker/telegram"
	"virtuals-tracker/virtuals"
)

type Cmd struct {
}

func (cmd *Cmd) Run(_ *global.Flags, queries *database.Queries) error {
	ctx := context.Background()

	// get virtuals price
	vprice, err := virtuals.GetPrice()
	if err != nil {
		return err
	}
	fmt.Printf("Virtuals Price: %f\n", vprice)

	minMCap := int(math.Ceil(1_000_000/vprice) - vprice)
	fmt.Printf("For being >1M mcap should be atleast %d\n", minMCap)

	var sendNotif []virtuals.Data

	// paginate
	hasNext := true
	currPage := 1
	totalNew := 0
	for hasNext {
		res, err := virtuals.GetPage(currPage, minMCap)
		if err != nil {
			return err
		}

		fmt.Printf("===%d===\n", res.Page)
		newAgents := 0
		for _, data := range res.Data {
			if data.McapInVirtual*vprice < 1_000_000 {
				fmt.Println(data)
				continue
			}

			// check if exists in db
			exists, err := queries.ExistsAgentByUID(ctx, data.UID)
			if err != nil {
				return err
			}
			if exists > 0 {
				continue
			}

			sendNotif = append(sendNotif, data)
			fmt.Println("New agent:", data)
			newAgents++
		}
		totalNew += newAgents

		hasNext = res.HasNextPage
		currPage++
	}

	if len(sendNotif) == 0 {
		fmt.Println("No new agents found")
		return nil
	}

	bot, err := telegram.GetBot()
	if err != nil {
		return err
	}
	defer bot.Close()

	fmt.Printf("Sending notifs for %d agents\n", len(sendNotif))
	for _, data := range sendNotif {
		err := telegram.SendNotif(data)
		if err != nil {
			return err
		}

		_, err = queries.CreateAgent(ctx, database.CreateAgentParams{
			Uid:      data.UID,
			Name:     data.Name,
			Status:   data.Status,
			Category: data.Category,
			Mcap:     strconv.FormatFloat(data.McapInVirtual, 'f', -1, 64),
			Notified: true,
		})
		if err != nil {
			return err
		}
	}

	fmt.Println("\nFinished")
	return nil
}
