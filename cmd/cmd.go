package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alecthomas/kong"
	"os"
	"virtuals-tracker/cmd/check"
	"virtuals-tracker/cmd/global"
	"virtuals-tracker/cmd/populate"
	"virtuals-tracker/config"
	"virtuals-tracker/database"

	_ "github.com/mattn/go-sqlite3"
)

type CLI struct {
	Globals  global.Flags `embed:"" group:"Global Options"`
	Check    check.Cmd    `cmd:"" help:"Check virtuals for updated data and notify"`
	Populate populate.Cmd `cmd:"" help:"populate database with initial data"`
}

func Execute() {
	db, err := sql.Open("sqlite3", "file:./sqlite.db")
	if err != nil {
		fmt.Println("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	{
		ctx := context.Background()
		if _, err := db.ExecContext(ctx, database.DDL); err != nil {
			fmt.Println("Failed to create database tables", err)
			return
		}
	}
	queries := database.New(db)

	cli := new(CLI)
	ctx := kong.Parse(cli,
		kong.Name("virtuals-tracker"),
		kong.Description("CLI for tracking virtuals"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Bind(&cli.Globals),
		kong.Bind(queries),
		kong.Vars{
			"version": config.Version,
		})

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
