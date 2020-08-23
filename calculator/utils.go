package main

import "sort"


// Get unique team names from slice of records of RawData
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


func stringInSlice(str string, slice []string) bool {
    for _, element := range slice {
        if element == str {
            return true
        }
    }
    return false
}