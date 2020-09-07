# statcalc
Stat calculator in Golang - Reads raw data from CSV files and transforms it into valuable insights

## Usage
- Drop CSV data files into the `data` folder. It **must** have the columns `HomeTeam HomeGoals AwayGoals AwayTeam` in this particular order (but the column names can be different).
- Install dependencies with `go get github.com/fatih/structs`
- Run the code with `go run calculator.go`
- View results in the `results` folder

## Naming conventions
- Filenames with 2v2 data i.e; `data/FIFA19-2v2.csv` must contain the string "2v2" in their filename (not case sensitive).
- The naming format for 2v2 teams must be the same as shown in the mentioned file i.e; unique-name of both individuals (one after the other) with first letter of each individual's unique-name capitalized.