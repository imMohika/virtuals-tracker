package populate

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"virtuals-tracker/cmd/global"
	"virtuals-tracker/database"
	"virtuals-tracker/virtuals"
)

type Cmd struct {
}

func (c *Cmd) Run(_ *global.Flags, queries *database.Queries) error {
	ctx := context.Background()

	// get virtuals price
	vprice, err := virtuals.GetPrice()
	if err != nil {
		return err
	}
	fmt.Printf("Virtuals Price: %f\n", vprice)

	minMCap := int(math.Ceil(1_000_000/vprice) - vprice)
	fmt.Printf("For being >1M mcap should be atleast %d\n", minMCap)

	// paginate
	hasNext := true
	currPage := 1
	totalAdded := 0
	for hasNext {
		res, err := virtuals.GetPage(currPage, minMCap)
		if err != nil {
			return err
		}

		fmt.Printf("===%d===\n", res.Page)
		addedAgents := 0
		for _, data := range res.Data {
			if data.McapInVirtual*vprice < 1_000_000 {
				fmt.Println(data)
				continue
			}

			_, err := queries.CreateAgent(ctx, database.CreateAgentParams{
				Uid:      data.UID,
				Name:     data.Name,
				Status:   data.Status,
				Category: data.Category,
				Mcap:     strconv.FormatFloat(data.McapInVirtual, 'f', -1, 64),
				Notified: true, // don't notify when populating
			})
			if err != nil {
				return err
			}

			//fmt.Println("Added agent:", agent)
			addedAgents++
		}
		fmt.Printf("Added %d agents\n", addedAgents)
		totalAdded += addedAgents

		hasNext = res.HasNextPage
		currPage++
	}

	fmt.Println("\nFinished")
	fmt.Printf("Added a total of %d agents\n", totalAdded)

	return nil
}
