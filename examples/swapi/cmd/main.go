package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"os"

	"github.com/google/uuid"
	"github.com/lumina-tech/gooq/examples/swapi/model"
	"github.com/lumina-tech/gooq/examples/swapi/table"
	"github.com/lumina-tech/gooq/pkg/database"
	"github.com/lumina-tech/gooq/pkg/gooq"
)

func main() {
	dockerDB := database.NewDockerizedDB(&database.DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		Username:     "postgres",
		Password:     "password",
		DatabaseName: "swapi",
		SSLMode:      "disable",
	}, "11.4-alpine")
	defer dockerDB.Close()

	ctx := context.Background()
	database.MigrateDatabase(dockerDB.DB.DB, "migrations")

	speciesStmt := gooq.InsertInto(table.Species).
		Set(table.Species.ID, uuid.New()).
		Set(table.Species.Name, "Human").
		Set(table.Species.Classification, "Mammal").
		Set(table.Species.AverageHeight, 160.5).
		Set(table.Species.AverageLifespan, 1000000000).
		Set(table.Species.HairColor, model.ColorBlack).
		Set(table.Species.SkinColor, model.ColorOrange).
		Set(table.Species.EyeColor, model.ColorBrown).
		Set(table.Species.HomeWorld, "Earth").
		Set(table.Species.Language, "English").
		Set(table.Species.Hash, "0xfoobar").
		Returning(table.Species.Asterisk)
	species, err := table.Species.ScanRowWithContext(ctx, dockerDB.DB, speciesStmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	personStmt := gooq.InsertInto(table.Person).
		Set(table.Person.ID, uuid.New()).
		Set(table.Person.Name, "Frank").
		Set(table.Person.Height, 170.3).
		Set(table.Person.Mass, 150.5).
		Set(table.Person.BirthYear, 1998).
		Set(table.Person.HomeWorld, "Runescape").
		Set(table.Person.Gender, model.GenderMale).
		Set(table.Person.EyeColor, model.ColorBrown).
		Set(table.Person.HairColor, model.ColorBlack).
		Set(table.Person.SkinColor, model.ColorOrange).
		Set(table.Person.SpeciesID, species.ID).
		Returning(table.Person.Asterisk)
	frank, err := table.Person.ScanRowWithContext(ctx, dockerDB.DB, personStmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	personStmtUpdate := gooq.InsertInto(table.Person).
		Set(table.Person.ID, uuid.New()).
		Set(table.Person.Name, "Frank").
		Set(table.Person.Height, 170.3).
		Set(table.Person.Mass, 150.5).
		Set(table.Person.BirthYear, 1998).
		Set(table.Person.HomeWorld, "Runescape").
		Set(table.Person.Gender, model.GenderMale).
		Set(table.Person.EyeColor, model.ColorBrown).
		Set(table.Person.HairColor, model.ColorBlue).
		Set(table.Person.SkinColor, model.ColorOrange).
		Set(table.Person.SpeciesID, species.ID).
		OnConflictDoUpdate(&table.Person.Constraints.NameBirthyearConstraint).
		SetUpdateColumns(table.Person.HairColor).
		Returning(table.Person.Asterisk)
	frankUpdated, err := table.Person.ScanRowWithContext(ctx, dockerDB.DB, personStmtUpdate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Fprintf(os.Stderr, "frank updated haircolor: %s to %s\n", frank.HairColor, frankUpdated.HairColor)
	fmt.Fprintf(os.Stderr, "frank did not update eyecolor: %s to %s\n", frank.EyeColor, frankUpdated.EyeColor)

	type PersonWithSpecies struct {
		model.Person
		Species *model.Species `db:"species"`
	}

	{
		stmt := gooq.Select().From(table.Species).
			Where(table.Species.Hash.IsEq(pgtype.UndecodedBytes("0xfoobar")))

		var results []model.Species
		if err := gooq.ScanRowsWithContext(ctx, dockerDB.DB, stmt, &results); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}

		printQuery(stmt, results)
	}

	{
		speciesWithAlias := table.Species.As("species_alias")
		stmt := gooq.Select(
			table.Person.Asterisk,
			speciesWithAlias.Name.As("species.name"),
			speciesWithAlias.Classification.As("species.name"),
			speciesWithAlias.AverageHeight.As("species.average_height"),
			speciesWithAlias.AverageLifespan.As("species.average_lifespan"),
			speciesWithAlias.HairColor.As("species.hair_color"),
			speciesWithAlias.SkinColor.As("species.skin_color"),
			speciesWithAlias.EyeColor.As("species.eye_color"),
			speciesWithAlias.HomeWorld.As("species.home_world"),
			speciesWithAlias.Language.As("species.language"),
			speciesWithAlias.Hash.As("species.hash"),
		).From(table.Person).
			Join(speciesWithAlias).
			On(table.Person.SpeciesID.Eq(speciesWithAlias.ID))

		var results []PersonWithSpecies
		if err := gooq.ScanRowsWithContext(ctx, dockerDB.DB, stmt, &results); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}

		printQuery(stmt, results)
	}

	// same as above but we don't have to manually enumerate all the column in species
	// inside the projection
	{
		selection := []gooq.Selectable{table.Person.Asterisk}
		selection = append(selection,
			getColumnsWithPrefix("species", table.Species.GetColumns())...)
		stmt := gooq.Select(selection...).From(table.Person).
			Join(table.Species).
			On(table.Person.SpeciesID.Eq(table.Species.ID))

		var results []PersonWithSpecies
		if err := gooq.ScanRowsWithContext(ctx, dockerDB.DB, stmt, &results); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}

		printQuery(stmt, results)
	}
}

func printQuery(stmt gooq.Selectable, results interface{}) {
	builder := &gooq.Builder{}
	stmt.Render(builder)
	bytes, _ := json.Marshal(results)

	fmt.Println("### query ###")
	fmt.Println()
	fmt.Println(builder.String())
	fmt.Println()
	fmt.Println(string(bytes))
	fmt.Println()
}

func getColumnsWithPrefix(
	prefix string, expressions []gooq.Expression,
) []gooq.Selectable {
	results := make([]gooq.Selectable, 0)
	for _, exp := range expressions {
		if field, ok := exp.(gooq.Field); ok {
			alias := fmt.Sprintf("%s.%s", prefix, field.GetName())
			results = append(results, exp.As(alias))
		}
	}
	return results
}
