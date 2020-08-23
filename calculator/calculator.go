package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/fatih/structs"
	// "github.com/Nishant173/statcalc"
)


// Struct to store absolute tabular statistics
type StatsAbs struct {
	Rank int
	Team string
	GamesPlayed int
	Points int
	GoalDifference int
	Wins int
	Losses int
	Draws int
	GoalsScored int
	GoalsAllowed int
	CleanSheets int
	CleanSheetsAgainst int
	BigWins int
	BigLosses int
}


// Struct to store normalized tabular statistics i.e; StatAbs / GamesPlayed
type StatsNorm struct {
	Rank int
	Team string
	GamesPlayed int
	PPG float64
	GDPG float64
	WinPct float64
	LossPct float64
	DrawPct float64
	GSPG float64
	GAPG float64
	CsPct float64
	CsaPct float64
	BigWinPct float64
	BigLossPct float64
}


/*
Method that gets slice of stringified elements of `StatsNorm` struct (by record).
Used as helper function in storing data of `StatsNorm` struct to CSV file.
*/
func (obj StatsNorm) ListStringifiedValues() []string {
	var values []string
	values = append(values, strconv.Itoa(obj.Rank))
	values = append(values, obj.Team)
	values = append(values, strconv.Itoa(obj.GamesPlayed))
	values = append(values, fmt.Sprintf("%g", obj.PPG))
	values = append(values, fmt.Sprintf("%g", obj.GDPG))
	values = append(values, fmt.Sprintf("%g", obj.WinPct))
	values = append(values, fmt.Sprintf("%g", obj.LossPct))
	values = append(values, fmt.Sprintf("%g", obj.DrawPct))
	values = append(values, fmt.Sprintf("%g", obj.GSPG))
	values = append(values, fmt.Sprintf("%g", obj.GAPG))
	values = append(values, fmt.Sprintf("%g", obj.CsPct))
	values = append(values, fmt.Sprintf("%g", obj.CsaPct))
	values = append(values, fmt.Sprintf("%g", obj.BigWinPct))
	values = append(values, fmt.Sprintf("%g", obj.BigLossPct))
	return values
}


/*
Method that gets slice of stringified elements of `StatsAbs` struct (by record).
Used as helper function in storing data of `StatsAbs` struct to CSV file.
*/
func (obj StatsAbs) ListStringifiedValues() []string {
	var values []string
	values = append(values, strconv.Itoa(obj.Rank))
	values = append(values, obj.Team)
	values = append(values, strconv.Itoa(obj.GamesPlayed))
	values = append(values, strconv.Itoa(obj.Points))
	values = append(values, strconv.Itoa(obj.GoalDifference))
	values = append(values, strconv.Itoa(obj.Wins))
	values = append(values, strconv.Itoa(obj.Losses))
	values = append(values, strconv.Itoa(obj.Draws))
	values = append(values, strconv.Itoa(obj.GoalsScored))
	values = append(values, strconv.Itoa(obj.GoalsAllowed))
	values = append(values, strconv.Itoa(obj.CleanSheets))
	values = append(values, strconv.Itoa(obj.CleanSheetsAgainst))
	values = append(values, strconv.Itoa(obj.BigWins))
	values = append(values, strconv.Itoa(obj.BigLosses))
	return values
}


func integerify(num float64) int {
    return int(num + math.Copysign(0.5, num))
}


func round(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(integerify(num * output)) / output
}


func getGamesPlayedCount(records []RawData, team string) int {
	count := 0
	for _, record := range records {
		if record.HomeTeam == team {
			count ++
		} else if record.AwayTeam == team {
			count ++
		}
	}
	return count
}


func getWinCount(records []RawData, team string) int {
	count := 0
	for _, record := range records {
		if record.HomeTeam == team && record.HomeGoals > record.AwayGoals {
			count ++
		} else if record.AwayTeam == team && record.AwayGoals > record.HomeGoals {
			count ++
		}
	}
	return count
}


func getLossCount(records []RawData, team string) int {
	count := 0
	for _, record := range records {
		if record.HomeTeam == team && record.HomeGoals < record.AwayGoals {
			count ++
		} else if record.AwayTeam == team && record.AwayGoals < record.HomeGoals {
			count ++
		}
	}
	return count
}


func getDrawCount(records []RawData, team string) int {
	count := 0
	for _, record := range records {
		if record.HomeTeam == team && record.HomeGoals == record.AwayGoals {
			count ++
		} else if record.AwayTeam == team && record.AwayGoals == record.HomeGoals {
			count ++
		}
	}
	return count
}


func getGoalsScored(records []RawData, team string) int {
	goalsScored := 0
	for _, record := range records {
		if record.HomeTeam == team {
			goalsScored += record.HomeGoals
		} else if record.AwayTeam == team {
			goalsScored += record.AwayGoals
		}
	}
	return goalsScored
}


func getGoalsAllowed(records []RawData, team string) int {
	goalsAllowed := 0
	for _, record := range records {
		if record.HomeTeam == team {
			goalsAllowed += record.AwayGoals
		} else if record.AwayTeam == team {
			goalsAllowed += record.HomeGoals
		}
	}
	return goalsAllowed
}


func getCleanSheets(records []RawData, team string) int {
	cleanSheetCount := 0
	for _, record := range records {
		if record.HomeTeam == team && record.AwayGoals == 0 {
			cleanSheetCount ++
		} else if record.AwayTeam == team && record.HomeGoals == 0 {
			cleanSheetCount ++
		}
	}
	return cleanSheetCount
}


func getCleanSheetsAgainst(records []RawData, team string) int {
	cleanSheetAgainstCount := 0
	for _, record := range records {
		if record.HomeTeam == team && record.HomeGoals == 0 {
			cleanSheetAgainstCount ++
		} else if record.AwayTeam == team && record.AwayGoals == 0 {
			cleanSheetAgainstCount ++
		}
	}
	return cleanSheetAgainstCount
}


func getBigWinCount(records []RawData, team string, margin int) int {
	bigWinCount := 0
	for _, record := range records {
		hg := record.HomeGoals
		ag := record.AwayGoals
		goalMargin := int(math.Abs(float64(hg - ag)))
		if record.HomeTeam == team && hg > ag && goalMargin >= margin {
			bigWinCount ++
		} else if record.AwayTeam == team && ag > hg && goalMargin >= margin {
			bigWinCount ++
		}
	}
	return bigWinCount
}


func getBigLossCount(records []RawData, team string, margin int) int {
	bigLossCount := 0
	for _, record := range records {
		hg := record.HomeGoals
		ag := record.AwayGoals
		goalMargin := int(math.Abs(float64(hg - ag)))
		if record.HomeTeam == team && hg < ag && goalMargin >= margin {
			bigLossCount ++
		} else if record.AwayTeam == team && ag < hg && goalMargin >= margin {
			bigLossCount ++
		}
	}
	return bigLossCount
}


/*
Gets slice of absolute stats from raw records.
Returns slice wherein each element of the slice is an object of the struct `StatsAbs`
*/
func getAbsoluteStats(records []RawData) []StatsAbs {
	teams := getUniqueTeamNames(records)
	bigResultGoalMargin := 3 // Will be considered as big result if GoalDifference >= this number
	var sliceAbsoluteStats []StatsAbs
	for _, team := range teams {
		wins := getWinCount(records, team)
		draws := getDrawCount(records, team)
		gs := getGoalsScored(records, team)
		ga := getGoalsAllowed(records, team)
		gd := gs - ga
		points := 3 * wins + draws
		tempAbsoluteStats := StatsAbs{
			Team: team,
			GamesPlayed: getGamesPlayedCount(records, team),
			Points: points,
			GoalDifference: gd,
			Wins: wins,
			Losses: getLossCount(records, team),
			Draws: draws,
			GoalsScored: gs,
			GoalsAllowed: ga,
			CleanSheets: getCleanSheets(records, team),
			CleanSheetsAgainst: getCleanSheetsAgainst(records, team),
			BigWins: getBigWinCount(records, team, bigResultGoalMargin),
			BigLosses: getBigLossCount(records, team, bigResultGoalMargin),
		}
		sliceAbsoluteStats = append(sliceAbsoluteStats, tempAbsoluteStats)
	}
	return sliceAbsoluteStats
}


/*
Gets slice of normalized stats from slice of absolute stats.
Returns slice wherein each element of the slice is an object of the struct `StatsNorm`
*/
func getNormalizedStats(sliceAbsStats []StatsAbs) []StatsNorm {
	hundred := 100.0
	var sliceNormalizedStats []StatsNorm
	for _, absStatSlice := range sliceAbsStats {
		gamesPlayed := float64(absStatSlice.GamesPlayed)
		tempNormalizedStats := StatsNorm{
			Team: absStatSlice.Team,
			GamesPlayed: absStatSlice.GamesPlayed,
			PPG: round(float64(absStatSlice.Points) / gamesPlayed, 4),
			GDPG: round(float64(absStatSlice.GoalDifference) / gamesPlayed, 3),
			WinPct: round(float64(absStatSlice.Wins) * hundred / gamesPlayed, 2),
			LossPct: round(float64(absStatSlice.Losses) * hundred / gamesPlayed, 2),
			DrawPct: round(float64(absStatSlice.Draws) * hundred / gamesPlayed, 2),
			GSPG: round(float64(absStatSlice.GoalsScored) / gamesPlayed, 3),
			GAPG: round(float64(absStatSlice.GoalsAllowed) / gamesPlayed, 3),
			CsPct: round(float64(absStatSlice.CleanSheets) * hundred / gamesPlayed, 2),
			CsaPct: round(float64(absStatSlice.CleanSheetsAgainst) * hundred / gamesPlayed, 2),
			BigWinPct: round(float64(absStatSlice.BigWins) * hundred / gamesPlayed, 2),
			BigLossPct: round(float64(absStatSlice.BigLosses) * hundred / gamesPlayed, 2),
		}
		sliceNormalizedStats = append(sliceNormalizedStats, tempNormalizedStats)
	}
	return sliceNormalizedStats
}


// Sorts normalized stats based on certain metric
func sortNormalizedStats(sliceNormalizedStats []StatsNorm) []StatsNorm {
	sort.SliceStable(sliceNormalizedStats, func(i, j int) bool {
		return sliceNormalizedStats[i].PPG > sliceNormalizedStats[j].PPG
	})
	return sliceNormalizedStats
}


// Sorts absolute stats based on certain metric
func sortAbsoluteStats(sliceAbsoluteStats []StatsAbs) []StatsAbs {
	sort.SliceStable(sliceAbsoluteStats, func(i, j int) bool {
		pointsOfI := 3 * sliceAbsoluteStats[i].Wins + sliceAbsoluteStats[i].Draws
		pointsOfJ := 3 * sliceAbsoluteStats[j].Wins + sliceAbsoluteStats[j].Draws
		ppgOfI := float64(pointsOfI) / float64(sliceAbsoluteStats[i].GamesPlayed)
		ppgOfJ := float64(pointsOfJ) / float64(sliceAbsoluteStats[j].GamesPlayed)
		return ppgOfI > ppgOfJ
	})
	return sliceAbsoluteStats
}


// Generate ranking AFTER slice of `StatsNorm` objects is sorted based on ranking metric/s
func rankNormalizedStats(sliceNormalizedStats []StatsNorm) []StatsNorm {
	var sliceNormalizedStatsRanked []StatsNorm
	for idx, tempStats := range sliceNormalizedStats {
		tempStats.Rank = idx + 1
		sliceNormalizedStatsRanked = append(sliceNormalizedStatsRanked, tempStats)
	}
	return sliceNormalizedStatsRanked
}


// Generate ranking AFTER slice of `StatsAbs` objects is sorted based on ranking metric/s
func rankAbsoluteStats(sliceAbsoluteStats []StatsAbs) []StatsAbs {
	var sliceAbsoluteStatsRanked []StatsAbs
	for idx, tempStats := range sliceAbsoluteStats {
		tempStats.Rank = idx + 1
		sliceAbsoluteStatsRanked = append(sliceAbsoluteStatsRanked, tempStats)
	}
	return sliceAbsoluteStatsRanked
}


/*
[Helper function]
Extends a slice of strings with another slice of strings, returning
only the unique elements among the two given slices.
*/
func extendUniqueElements(slice1 []string, slice2 []string) []string {
	for _, element := range slice2 {
		if !stringInSlice(element, slice1) {
			slice1 = append(slice1, element)
		}
	}
	return slice1
}


func getUniqueIndividualNames(records []RawData) []string {
	var uniqueIndividuals []string
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	uniqueTeams := getUniqueTeamNames(records)
	for _, team := range uniqueTeams {
		matchedIndividuals := re.FindAllString(team, -1)
		uniqueIndividuals = extendUniqueElements(uniqueIndividuals, matchedIndividuals)
	}
	sort.Strings(uniqueIndividuals)
	return uniqueIndividuals
}


func individualInTeam(individual string, team string) bool {
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	teamMembers := re.FindAllString(team, -1) // Slice of team members
	for _, teamMember := range teamMembers {
		if teamMember == individual {
			return true
		}
	}
	return false
}


/*
Gets slice of absolute stats of individuals from raw records and slice of absolute stats of teams.
Returns slice wherein each element of the slice is an object of the struct `StatsAbs`
*/
func getAbsoluteStatsByIndividual(records []RawData, sliceAbsoluteStats []StatsAbs) []StatsAbs {
	individuals := getUniqueIndividualNames(records)
	var sliceStatsAllIndividuals []StatsAbs // Slice of all individuals' stats
	for _, individual := range individuals {
		var sliceStatsByIndividual []StatsAbs // Slice of stats by individual
		var gamesPlayed int
		var points int
		var goalDifference int
		var wins int
		var losses int
		var draws int
		var goalsScored int
		var goalsAllowed int
		var cleanSheets int
		var cleanSheetsAgainst int
		var bigWins int
		var bigLosses int
		for _, obj := range sliceAbsoluteStats {
			if individualInTeam(individual, obj.Team) {
				sliceStatsByIndividual = append(sliceStatsByIndividual, obj)
				gamesPlayed += obj.GamesPlayed
				points += obj.Points
				goalDifference += obj.GoalDifference
				wins += obj.Wins
				losses += obj.Losses
				draws += obj.Draws
				goalsScored += obj.GoalsScored
				goalsAllowed += obj.GoalsAllowed
				cleanSheets += obj.CleanSheets
				cleanSheetsAgainst += obj.CleanSheetsAgainst
				bigWins += obj.BigWins
				bigLosses += obj.BigLosses
			}
		}
		// Construct struct `StatsAbs` by individual
		tempObj := StatsAbs{
			Rank: 0,
			Team: individual,
			GamesPlayed: gamesPlayed,
			Points: points,
			GoalDifference: goalDifference,
			Wins: wins,
			Losses: losses,
			Draws: draws,
			GoalsScored: goalsScored,
			GoalsAllowed: goalsAllowed,
			CleanSheets: cleanSheets,
			CleanSheetsAgainst: cleanSheetsAgainst,
			BigWins: bigWins,
			BigLosses: bigLosses,
		}
		sliceStatsAllIndividuals = append(sliceStatsAllIndividuals, tempObj)
	}
	return sliceStatsAllIndividuals
}


// Saves slice having objects of `StatsNorm` struct to CSV file
func saveNormToCsv(sliceData []StatsNorm, filepath string) {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
    if err != nil {
        os.Exit(1)
	}
	var strWrite [][]string // Slice of slice of strings, where each sub-slice represents a record
	statFields := structs.Names(&StatsNorm{})
	strWrite = append(strWrite, statFields)
	for _, obj := range sliceData {
		record := obj.ListStringifiedValues()
		strWrite = append(strWrite, record)
	}
    csvWriter := csv.NewWriter(file)
    csvWriter.WriteAll(strWrite)
    csvWriter.Flush()
}


// Saves slice having objects of `StatsAbs` struct to CSV file
func saveAbsToCsv(sliceData []StatsAbs, filepath string) {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
    if err != nil {
        os.Exit(1)
	}
	var strWrite [][]string // Slice of slice of strings, where each sub-slice represents a record
	statFields := structs.Names(&StatsAbs{})
	strWrite = append(strWrite, statFields)
	for _, obj := range sliceData {
		record := obj.ListStringifiedValues()
		strWrite = append(strWrite, record)
	}
    csvWriter := csv.NewWriter(file)
    csvWriter.WriteAll(strWrite)
    csvWriter.Flush()
}


func main() {
	pathRawData := "../data/FIFA19-2v2.csv"
	rawRecords := readRawRecordsFromCsv(pathRawData)

	// Teams stats
	sliceAbsStats := getAbsoluteStats(rawRecords)
	sliceNormStats := getNormalizedStats(sliceAbsStats)

	sliceAbsStats = sortAbsoluteStats(sliceAbsStats)
	sliceAbsStats = rankAbsoluteStats(sliceAbsStats)
	sliceNormStats = sortNormalizedStats(sliceNormStats)
	sliceNormStats = rankNormalizedStats(sliceNormStats)
	
	// Individuals' stats
	sliceAbsStatsSolo := getAbsoluteStatsByIndividual(rawRecords, sliceAbsStats)
	sliceNormStatsSolo := getNormalizedStats(sliceAbsStatsSolo)

	sliceAbsStatsSolo = sortAbsoluteStats(sliceAbsStatsSolo)
	sliceAbsStatsSolo = rankAbsoluteStats(sliceAbsStatsSolo)
	sliceNormStatsSolo = sortNormalizedStats(sliceNormStatsSolo)
	sliceNormStatsSolo = rankNormalizedStats(sliceNormStatsSolo)

	// Store results
	saveAbsToCsv(sliceAbsStats, "../results/Teams - Absolute Stats.csv")
	saveNormToCsv(sliceNormStats, "../results/Teams - Normalized Stats.csv")
	saveAbsToCsv(sliceAbsStatsSolo, "../results/Individuals - Absolute Stats.csv")
	saveNormToCsv(sliceNormStatsSolo, "../results/Individuals - Normalized Stats.csv")
	fmt.Println("Done!")
}