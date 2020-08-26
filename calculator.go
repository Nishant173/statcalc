package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)


// Constants - Paths to data (source) and results (destination) folders
const (
	pathDataFolder = "data"
	pathResultsFolder = "results"
)


// Struct to store raw data
type RawData struct {
	HomeTeam  string
	HomeGoals int
	AwayGoals int
	AwayTeam  string
}


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


// Struct to store latest form (decided by latest PPG)
type LatestForm struct {
	Rank int
	Team string
	Form string // WLD (Wins, Losses, Draws) representation of latest form
	LatestPPG float64
	NumGamesConsidered int
}


/*
Method that gets slice of stringified elements of `StatsAbs` struct (by record).
NOTE: Elements of the slice returned must be in same order as the attributes defined in the struct.
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


/*
Method that gets slice of stringified elements of `StatsNorm` struct (by record).
NOTE: Elements of the slice returned must be in same order as the attributes defined in the struct.
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
Method that gets slice of stringified elements of `LatestForm` struct (by record).
NOTE: Elements of the slice returned must be in same order as the attributes defined in the struct.
Used as helper function in storing data of `LatestForm` struct to CSV file.
*/
func (obj LatestForm) ListStringifiedValues() []string {
	var values []string
	values = append(values, strconv.Itoa(obj.Rank))
	values = append(values, obj.Team)
	values = append(values, obj.Form)
	values = append(values, fmt.Sprintf("%g", obj.LatestPPG))
	values = append(values, strconv.Itoa(obj.NumGamesConsidered))
	return values
}


// Read CSV file having columns "HomeTeam, HomeGoals, AwayGoals, AwayTeam" in that order
func readRawRecordsFromCsv(filepath string) []RawData {
	csvfile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the CSV file", err)
	}
	r := csv.NewReader(csvfile)
	lineCount := 0
	var records []RawData
	for {
		lineCount ++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if lineCount != 1 {
			homeGoals, strConvErrHome := strconv.Atoi(record[1])
			awayGoals, strConvErrAway := strconv.Atoi(record[2])
			if strConvErrHome != nil {
				log.Fatalln("Error while converting HomeGoals to int", strConvErrHome)
			}
			if strConvErrAway != nil {
				log.Fatalln("Error while converting AwayGoals to int", strConvErrAway)
			}
			record := RawData{
				HomeTeam: record[0],
				HomeGoals: homeGoals,
				AwayGoals: awayGoals,
				AwayTeam: record[3],
			}
			records = append(records, record)
		}
	}
	return records
}


func removeExtension(filenameWithExt string) string {
	return strings.TrimSuffix(filenameWithExt, path.Ext(filenameWithExt))
}


func filenameContains2v2(filenameWithExt string) bool {
	return strings.Contains(strings.ToLower(filenameWithExt), "2v2")
}


func reverseRecordsOrder(records []RawData) []RawData {
	sliceRecordsReversed := []RawData{}
	for i := len(records) - 1; i >= 0; i-- {
		sliceRecordsReversed = append(sliceRecordsReversed, records[i])
	}
	return sliceRecordsReversed
}


// Get unique team names from slice of records of `RawData`
func getUniqueTeamNames(records []RawData) []string {
	var uniqueTeamNames []string
	for _, record := range records {
		homeTeam := record.HomeTeam
		awayTeam := record.AwayTeam
		if !stringInSlice(homeTeam, uniqueTeamNames) {
			uniqueTeamNames = append(uniqueTeamNames, homeTeam)
		}
		if !stringInSlice(awayTeam, uniqueTeamNames) {
			uniqueTeamNames = append(uniqueTeamNames, awayTeam)
		}
	}
	sort.Strings(uniqueTeamNames)
	return uniqueTeamNames
}


// Get unique individual names from slice of records of `RawData`
func getUniqueIndividualNames(records []RawData) []string {
	var uniqueIndividualNames []string
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	uniqueTeams := getUniqueTeamNames(records)
	for _, team := range uniqueTeams {
		matchedIndividuals := re.FindAllString(team, -1)
		uniqueIndividualNames = extendUniqueElements(uniqueIndividualNames, matchedIndividuals)
	}
	sort.Strings(uniqueIndividualNames)
	return uniqueIndividualNames
}


func stringInSlice(str string, slice []string) bool {
    for _, element := range slice {
        if element == str {
            return true
        }
    }
    return false
}


/*
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
	for _, obj := range sliceAbsStats {
		gamesPlayed := float64(obj.GamesPlayed)
		tempNormalizedStats := StatsNorm{
			Team: obj.Team,
			GamesPlayed: obj.GamesPlayed,
			PPG: round(float64(obj.Points) / gamesPlayed, 4),
			GDPG: round(float64(obj.GoalDifference) / gamesPlayed, 3),
			WinPct: round(float64(obj.Wins) * hundred / gamesPlayed, 2),
			LossPct: round(float64(obj.Losses) * hundred / gamesPlayed, 2),
			DrawPct: round(float64(obj.Draws) * hundred / gamesPlayed, 2),
			GSPG: round(float64(obj.GoalsScored) / gamesPlayed, 3),
			GAPG: round(float64(obj.GoalsAllowed) / gamesPlayed, 3),
			CsPct: round(float64(obj.CleanSheets) * hundred / gamesPlayed, 2),
			CsaPct: round(float64(obj.CleanSheetsAgainst) * hundred / gamesPlayed, 2),
			BigWinPct: round(float64(obj.BigWins) * hundred / gamesPlayed, 2),
			BigLossPct: round(float64(obj.BigLosses) * hundred / gamesPlayed, 2),
		}
		sliceNormalizedStats = append(sliceNormalizedStats, tempNormalizedStats)
	}
	return sliceNormalizedStats
}


// Sorts absolute stats based on certain metric/s
func sortAbsStatsByMetric(sliceAbsoluteStats []StatsAbs) []StatsAbs {
	sort.SliceStable(sliceAbsoluteStats, func(i, j int) bool {
		pointsOfI := 3 * sliceAbsoluteStats[i].Wins + sliceAbsoluteStats[i].Draws
		pointsOfJ := 3 * sliceAbsoluteStats[j].Wins + sliceAbsoluteStats[j].Draws
		ppgOfI := float64(pointsOfI) / float64(sliceAbsoluteStats[i].GamesPlayed)
		ppgOfJ := float64(pointsOfJ) / float64(sliceAbsoluteStats[j].GamesPlayed)
		return ppgOfI > ppgOfJ
	})
	return sliceAbsoluteStats
}


// Sorts normalized stats based on certain metric/s
func sortNormStatsByMetric(sliceNormalizedStats []StatsNorm) []StatsNorm {
	sort.SliceStable(sliceNormalizedStats, func(i, j int) bool {
		return sliceNormalizedStats[i].PPG > sliceNormalizedStats[j].PPG
	})
	return sliceNormalizedStats
}


// Sorts latest form based on certain metric/s
func sortLatestFormByMetric(sliceLatestForm []LatestForm) []LatestForm {
	sort.SliceStable(sliceLatestForm, func(i, j int) bool {
		return sliceLatestForm[i].LatestPPG > sliceLatestForm[j].LatestPPG
	})
	return sliceLatestForm
}


// NOTE: Only attaches incremental ranking, since the slice is already sorted by ranking metric/s
// Attach ranking AFTER slice of `StatsAbs` objects is sorted based on ranking metric/s
func attachRankingToAbsStats(sliceAbsoluteStats []StatsAbs) []StatsAbs {
	var sliceAbsoluteStatsRanked []StatsAbs
	for idx, tempStats := range sliceAbsoluteStats {
		tempStats.Rank = idx + 1
		sliceAbsoluteStatsRanked = append(sliceAbsoluteStatsRanked, tempStats)
	}
	return sliceAbsoluteStatsRanked
}


// Attach ranking AFTER slice of `StatsNorm` objects is sorted based on ranking metric/s
func attachRankingToNormStats(sliceNormalizedStats []StatsNorm) []StatsNorm {
	var sliceNormalizedStatsRanked []StatsNorm
	for idx, tempStats := range sliceNormalizedStats {
		tempStats.Rank = idx + 1
		sliceNormalizedStatsRanked = append(sliceNormalizedStatsRanked, tempStats)
	}
	return sliceNormalizedStatsRanked
}


// Attach ranking AFTER slice of `LatestForm` objects is sorted based on ranking metric/s
func attachRankingToLatestForm(sliceLatestForm []LatestForm) []LatestForm {
	var sliceLatestFormRanked []LatestForm
	for idx, tempStats := range sliceLatestForm {
		tempStats.Rank = idx + 1
		sliceLatestFormRanked = append(sliceLatestFormRanked, tempStats)
	}
	return sliceLatestFormRanked
}


/*
[Helper function]
Extracts individuals' absolute stats from slice of teams' absolute stats.
Returns slice wherein each element of the slice is an object of the struct `StatsAbs`
*/
func extractIndividualStatsFromTeamStats(individual string, sliceTeamAbsStats []StatsAbs) map[string]int {
	mapStatsByIndividual := map[string]int{}
	for _, objStat := range sliceTeamAbsStats {
		if individualInTeam(individual, objStat.Team) {
			mapStatsByIndividual["GamesPlayed"] += objStat.GamesPlayed
			mapStatsByIndividual["Points"] += objStat.Points
			mapStatsByIndividual["GoalDifference"] += objStat.GoalDifference
			mapStatsByIndividual["Wins"] += objStat.Wins
			mapStatsByIndividual["Losses"] += objStat.Losses
			mapStatsByIndividual["Draws"] += objStat.Draws
			mapStatsByIndividual["GoalsScored"] += objStat.GoalsScored
			mapStatsByIndividual["GoalsAllowed"] += objStat.GoalsAllowed
			mapStatsByIndividual["CleanSheets"] += objStat.CleanSheets
			mapStatsByIndividual["CleanSheetsAgainst"] += objStat.CleanSheetsAgainst
			mapStatsByIndividual["BigWins"] += objStat.BigWins
			mapStatsByIndividual["BigLosses"] += objStat.BigLosses
		}
	}
	return mapStatsByIndividual
}


/*
Gets slice of absolute stats of individuals from `RawData` records, `StatsAbs` of teams.
Returns slice wherein each element of the slice is an object of the struct `StatsAbs`
*/
func getAbsoluteStatsByIndividual(records []RawData, sliceAbsoluteTeamStats []StatsAbs) []StatsAbs {
	individuals := getUniqueIndividualNames(records)
	var sliceStatsAllIndividuals []StatsAbs
	for _, individual := range individuals {
		mapIndividualStats := extractIndividualStatsFromTeamStats(individual, sliceAbsoluteTeamStats)
		objStatsByIndividual := StatsAbs{
			Rank: 0,
			Team: individual,
			GamesPlayed: mapIndividualStats["GamesPlayed"],
			Points: mapIndividualStats["Points"],
			GoalDifference: mapIndividualStats["GoalDifference"],
			Wins: mapIndividualStats["Wins"],
			Losses: mapIndividualStats["Losses"],
			Draws: mapIndividualStats["Draws"],
			GoalsScored: mapIndividualStats["GoalsScored"],
			GoalsAllowed: mapIndividualStats["GoalsAllowed"],
			CleanSheets: mapIndividualStats["CleanSheets"],
			CleanSheetsAgainst: mapIndividualStats["CleanSheetsAgainst"],
			BigWins: mapIndividualStats["BigWins"],
			BigLosses: mapIndividualStats["BigLosses"],
		}
		sliceStatsAllIndividuals = append(sliceStatsAllIndividuals, objStatsByIndividual)
	}
	return sliceStatsAllIndividuals
}


// Gets string of WLD (Wins, Losses, Draws) representation of `LatestForm` for Teams
func representLatestForm(records []RawData, team string, nLatestGames int) string {
	recordsReversed := reverseRecordsOrder(records)
	representationLatestForm := ""
	numGamesConsidered := 0
		for _, match := range recordsReversed {
			if team == match.HomeTeam {
				if match.HomeGoals > match.AwayGoals {
					representationLatestForm += "W"
				} else if match.HomeGoals == match.AwayGoals {
					representationLatestForm += "D"
				} else if match.HomeGoals < match.AwayGoals {
					representationLatestForm += "L"
				}
				numGamesConsidered ++
			} else if team == match.AwayTeam {
				if match.AwayGoals > match.HomeGoals {
					representationLatestForm += "W"
				} else if match.AwayGoals == match.HomeGoals {
					representationLatestForm += "D"
				} else if match.AwayGoals < match.HomeGoals {
					representationLatestForm += "L"
				}
				numGamesConsidered ++
			}
			if numGamesConsidered == nLatestGames {
				break
			}
		}
	return representationLatestForm
}


// Gets string of WLD (Wins, Losses, Draws) representation of `LatestForm` for Individuals
func representLatestFormSolo(records []RawData, individual string, nLatestGames int) string {
	recordsReversed := reverseRecordsOrder(records)
	representationLatestForm := ""
	numGamesConsidered := 0
		for _, match := range recordsReversed {
			if individualInTeam(individual, match.HomeTeam) {
				if match.HomeGoals > match.AwayGoals {
					representationLatestForm += "W"
				} else if match.HomeGoals == match.AwayGoals {
					representationLatestForm += "D"
				} else if match.HomeGoals < match.AwayGoals {
					representationLatestForm += "L"
				}
				numGamesConsidered ++
			} else if individualInTeam(individual, match.AwayTeam) {
				if match.AwayGoals > match.HomeGoals {
					representationLatestForm += "W"
				} else if match.AwayGoals == match.HomeGoals {
					representationLatestForm += "D"
				} else if match.AwayGoals < match.HomeGoals {
					representationLatestForm += "L"
				}
				numGamesConsidered ++
			}
			if numGamesConsidered == nLatestGames {
				break
			}
		}
	return representationLatestForm
}


/*
Get latest form of team/individual in last `nLatestGames` games. Metric used is PPG (Points per game).
NOTE: Assumes that the records are sorted in ascending order of time of occurence of matches.
*/
func getLatestForm(records []RawData, nLatestGames int) []LatestForm {
	recordsReversed := reverseRecordsOrder(records)
	teams := getUniqueTeamNames(records)
	sliceLatestFormData := []LatestForm{}
	for _, team := range teams {
		wins, draws, gamesPlayed := 0, 0, 0
		for _, match := range recordsReversed {
			if team == match.HomeTeam {
				if match.HomeGoals > match.AwayGoals {
					wins ++
				} else if match.HomeGoals == match.AwayGoals {
					draws ++
				}
				gamesPlayed ++
			} else if team == match.AwayTeam {
				if match.AwayGoals > match.HomeGoals {
					wins ++
				} else if match.HomeGoals == match.AwayGoals {
					draws ++
				}
				gamesPlayed ++
			}
			if gamesPlayed == nLatestGames {
				break
			}
		}
		latestPPG := float64(3 * wins + draws) / float64(gamesPlayed)
		tempObj := LatestForm{
			Rank: 0,
			Team: team,
			Form: representLatestForm(records, team, nLatestGames),
			LatestPPG: round(latestPPG, 4),
			NumGamesConsidered: gamesPlayed,
		}
		sliceLatestFormData = append(sliceLatestFormData, tempObj)
	}
	return sliceLatestFormData
}


func getLatestFormSolo(records []RawData, nLatestGames int) []LatestForm {
	recordsReversed := reverseRecordsOrder(records)
	individuals := getUniqueIndividualNames(records)
	sliceLatestFormData := []LatestForm{}
	for _, individual := range individuals {
		wins, draws, gamesPlayed := 0, 0, 0
		for _, match := range recordsReversed {
			if individualInTeam(individual, match.HomeTeam) {
				if match.HomeGoals > match.AwayGoals {
					wins ++
				} else if match.HomeGoals == match.AwayGoals {
					draws ++
				}
				gamesPlayed ++
			} else if individualInTeam(individual, match.AwayTeam) {
				if match.AwayGoals > match.HomeGoals {
					wins ++
				} else if match.HomeGoals == match.AwayGoals {
					draws ++
				}
				gamesPlayed ++
			}
			if gamesPlayed == nLatestGames {
				break
			}
		}
		latestPPG := float64(3 * wins + draws) / float64(gamesPlayed)
		tempObj := LatestForm{
			Rank: 0,
			Team: individual,
			Form: representLatestFormSolo(records, individual, nLatestGames),
			LatestPPG: round(latestPPG, 4),
			NumGamesConsidered: gamesPlayed,
		}
		sliceLatestFormData = append(sliceLatestFormData, tempObj)
	}
	return sliceLatestFormData
}


// Saves slice having objects of `StatsAbs` struct to CSV file
func saveAbsToCsv(sliceData []StatsAbs, filepath string) {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
    if err != nil {
        os.Exit(1)
	}
	var sliceStringifiedRecords [][]string // Slice of slice of strings, where each sub-slice represents a record
	statFields := structs.Names(&StatsAbs{})
	sliceStringifiedRecords = append(sliceStringifiedRecords, statFields)
	for _, obj := range sliceData {
		record := obj.ListStringifiedValues()
		sliceStringifiedRecords = append(sliceStringifiedRecords, record)
	}
    csvWriter := csv.NewWriter(file)
    csvWriter.WriteAll(sliceStringifiedRecords)
    csvWriter.Flush()
}


// Saves slice having objects of `StatsNorm` struct to CSV file
func saveNormToCsv(sliceData []StatsNorm, filepath string) {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
    if err != nil {
        os.Exit(1)
	}
	var sliceStringifiedRecords [][]string // Slice of slice of strings, where each sub-slice represents a record
	statFields := structs.Names(&StatsNorm{})
	sliceStringifiedRecords = append(sliceStringifiedRecords, statFields)
	for _, obj := range sliceData {
		record := obj.ListStringifiedValues()
		sliceStringifiedRecords = append(sliceStringifiedRecords, record)
	}
    csvWriter := csv.NewWriter(file)
    csvWriter.WriteAll(sliceStringifiedRecords)
    csvWriter.Flush()
}


// Saves slice having objects of `LatestForm` struct to CSV file
func saveLatestFormToCsv(sliceData []LatestForm, filepath string) {
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
    if err != nil {
        os.Exit(1)
	}
	var sliceStringifiedRecords [][]string // Slice of slice of strings, where each sub-slice represents a record
	statFields := structs.Names(&LatestForm{})
	sliceStringifiedRecords = append(sliceStringifiedRecords, statFields)
	for _, obj := range sliceData {
		record := obj.ListStringifiedValues()
		sliceStringifiedRecords = append(sliceStringifiedRecords, record)
	}
    csvWriter := csv.NewWriter(file)
    csvWriter.WriteAll(sliceStringifiedRecords)
    csvWriter.Flush()
}


// Gets slice of all filenames from data source
func getListOfDataFilenames() []string {
    f, err := os.Open(pathDataFolder)
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Fatal(err)
	}
	filenamesDesired := []string{}
    for _, file := range files {
		filenamesDesired = append(filenamesDesired, file.Name())
	}
	return filenamesDesired
}


// Executes ETL pipeline for a raw data file, and stores results appropriately
func executePipeline(filename string) {
	filenameWithoutExt := removeExtension(filename)
	pathRawData := pathDataFolder + "/" + filename
	rawRecords := readRawRecordsFromCsv(pathRawData)
	nLatestGames := 10 // Number of latest games to consider for LatestForm

	// ########## Teams stats ##########
	sliceAbsStats := getAbsoluteStats(rawRecords)
	sliceNormStats := getNormalizedStats(sliceAbsStats)
	sliceAbsStats = sortAbsStatsByMetric(sliceAbsStats)
	sliceAbsStats = attachRankingToAbsStats(sliceAbsStats)
	sliceNormStats = sortNormStatsByMetric(sliceNormStats)
	sliceNormStats = attachRankingToNormStats(sliceNormStats)
	// LatestForm
	sliceLatestForm := getLatestForm(rawRecords, nLatestGames)
	sliceLatestForm = sortLatestFormByMetric(sliceLatestForm)
	sliceLatestForm = attachRankingToLatestForm(sliceLatestForm)
	// Save results
	saveAbsToCsv(sliceAbsStats, pathResultsFolder + "/" + filenameWithoutExt + " - Teams - Absolute Stats.csv")
	saveNormToCsv(sliceNormStats, pathResultsFolder +  "/" + filenameWithoutExt + " - Teams - Normalized Stats.csv")
	saveLatestFormToCsv(sliceLatestForm, pathResultsFolder +  "/" + filenameWithoutExt + " - Teams - Latest Form.csv")
	
	// ########## Individuals' stats ##########
	if filenameContains2v2(filename) {
		sliceAbsStatsSolo := getAbsoluteStatsByIndividual(rawRecords, sliceAbsStats)
		sliceNormStatsSolo := getNormalizedStats(sliceAbsStatsSolo)
		sliceAbsStatsSolo = sortAbsStatsByMetric(sliceAbsStatsSolo)
		sliceAbsStatsSolo = attachRankingToAbsStats(sliceAbsStatsSolo)
		sliceNormStatsSolo = sortNormStatsByMetric(sliceNormStatsSolo)
		sliceNormStatsSolo = attachRankingToNormStats(sliceNormStatsSolo)
		// LatestForm
		sliceLatestFormSolo := getLatestFormSolo(rawRecords, nLatestGames)
		sliceLatestFormSolo = sortLatestFormByMetric(sliceLatestFormSolo)
		sliceLatestFormSolo = attachRankingToLatestForm(sliceLatestFormSolo)
		// Save results
		saveAbsToCsv(sliceAbsStatsSolo, pathResultsFolder +  "/" + filenameWithoutExt + " - Individuals - Absolute Stats.csv")
		saveNormToCsv(sliceNormStatsSolo, pathResultsFolder +  "/" + filenameWithoutExt + " - Individuals - Normalized Stats.csv")
		saveLatestFormToCsv(sliceLatestFormSolo, pathResultsFolder +  "/" + filenameWithoutExt + " - Individuals - Latest Form.csv")
	}
	fmt.Println("Computed stats for " + filename)
}


func main() {
	filenames := getListOfDataFilenames()
	for _, filename := range filenames {
		executePipeline(filename)
	}
	fmt.Println("\nDone!")
}