# statcalc
Stat calculator in Golang

## Usage
- Drop CSV data files into the `data` folder. It **must** have the columns `HomeTeam HomeGoals AwayGoals AwayTeam` in this particular order (but the column names can be different).
- Install dependencies with `go get github.com/fatih/structs`
- Run the code with `go run calculator.go`
- View results in the `results` folder