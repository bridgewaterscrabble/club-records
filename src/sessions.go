// enter-games file...
//
// Add games to the club database.

package main

import (
    "database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	//"slices"
	"strconv"
	"strings"
	"time"
    _ "github.com/lib/pq" // for Postgres
)

// For using functions that return multiple values inside Printf()
func wrap(vs ...interface{}) []interface{} {
    return vs
}
/*func Info(level int, a ...interface{}) {
    if level <= logLevel {
        log.Printf(a)
    }
    }*/

const clubSessions = "/nas/ebh/Scrabble/Sessions/2006"
const logLevel = 1

func main() {
    connStr := "postgresql://ebh@localhost/scrabble?sslmode=disable"
    // Connect to database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

	//////
	//
	//  This is not a general data entry app. In particular, the players
	//  listed in every game must have already been inserted into the
	//  players table. Likewise, the session date must already have been
	//  inserted into the sessions table. The total number of each is
	//  small enough that we can store copies of them in memory. Also,
	//  grab the dates for each NWL and CSW lexica.
	//
	//////
	var handle string
	var handles []string
	handles_db, err := db.Query("SELECT handle FROM players ORDER BY handle")
    defer handles_db.Close()
    if err != nil {
        log.Fatalln(err)
    }
	for handles_db.Next() {
		if err := handles_db.Scan(&handle); err != nil {
			log.Fatalln(err)
		}
		handles = append(handles, handle)
		log.Printf("Handle: %s\n", handle)
	}

	var date string
	var dates []string
    sessions_db, err := db.Query("SELECT date FROM sessions ORDER BY date")
    defer sessions_db.Close()
    if err != nil {
        log.Fatalln(err)
    }
    for sessions_db.Next() {
        if err := sessions_db.Scan(&date); err != nil {
			log.Fatalln(err)
		}
		// Date is in the form yyyy-mm-ddT00:00:00Z
		date = strings.TrimSuffix(date, "T00:00:00Z")
		date = strings.ReplaceAll(date, "-", "")
		dates = append(dates, date)
		//dt, _ := time.Parse("2006-01-02T15:04:05", strings.TrimSuffix(date, "Z"))
		log.Printf("Session date: %s\n", date/*dt.Format("2006-01-02")*/)
    }

	type Lexicon struct {
		name string
		bdate time.Time
		edate time.Time
	}
	var lex_bdate, lex_edate string

	var nwl_lexicon Lexicon
	var nwl_lexica []Lexicon
    nwl_lexica_db, err := db.Query(`SELECT lexicon, 
                                      TO_CHAR(start_date, 'YYYY-MM-DD'),
                                      TO_CHAR(end_date, 'YYYY-MM-DD')
                                    FROM lexica
                                    WHERE lex_family = 'nwl'
                                    ORDER BY end_date`)
    defer nwl_lexica_db.Close()
    if err != nil {
        log.Fatalln(err)
    }
    for nwl_lexica_db.Next() {
        if err := nwl_lexica_db.Scan(&nwl_lexicon.name, &lex_bdate, &lex_edate); err != nil {
			log.Fatalln(err)
		}
		nwl_lexicon.bdate, _ = time.Parse("2006-01-02", lex_bdate)
		nwl_lexicon.edate, _ = time.Parse("2006-01-02", lex_edate)
		nwl_lexica = append(nwl_lexica, nwl_lexicon)

		//yb, mb, db := nwl_lexicon.bdate.Date()
		//ye, me, de := nwl_lexicon.edate.Date()
		//log.Printf("NWL Lexicon: %s, from %d/%02d/%02d to %d/%02d/%02d\n", nwl_lexicon.name, yb, mb, db, ye, me, de)
    }

	var csw_lexicon Lexicon
	var csw_lexica []Lexicon
    csw_lexica_db, err := db.Query(`SELECT lexicon,
                                      TO_CHAR(start_date, 'YYYY-MM-DD'),
                                      TO_CHAR(end_date, 'YYYY-MM-DD')
                                    FROM lexica
                                    WHERE lex_family = 'csw'
                                    ORDER BY end_date`)
    defer csw_lexica_db.Close()
    if err != nil {
        log.Fatalln(err)
    }
    for csw_lexica_db.Next() {
        if err := csw_lexica_db.Scan(&csw_lexicon.name, &lex_bdate, &lex_edate); err != nil {
			log.Fatalln(err)
		}
		csw_lexicon.bdate, _ = time.Parse("2006-01-02", lex_bdate)
		csw_lexicon.edate, _ = time.Parse("2006-01-02", lex_edate)
		csw_lexica = append(csw_lexica, csw_lexicon)

		//yb, mb, db := csw_lexicon.bdate.Date()
		//ye, me, de := csw_lexicon.edate.Date()
		//log.Printf("CSW Lexicon: %s, from %d/%02d/%02d to %d/%02d/%02d\n", csw_lexicon.name, yb, mb, db, ye, me, de)
    }


	//////
	//
	//  Now go through all the session files at or below the directory
	//  specified by clubSessions. [TODO: comand-line!]
	//
	//////
	fsys := os.DirFS(clubSessions)
	fs.WalkDir(fsys, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		// WalkDir doesn't have a "find . -type f" equivalent so
		// we have to manually skip directories
		fmt.Printf("Path: %s\n", fpath)
		if d.IsDir() || fpath == "." {
			//fmt.Printf("Skipping directory %s\n", fpath)
			return nil
		}

		// Skip any file not ending in .csv or .wsv
		if !strings.HasSuffix(fpath, ".csv") &&
		   !strings.HasSuffix(fpath, ".wsv") {
			return nil
		}

		/////
		//  Skip this session file if the session has already been recorded. As a
		//  shortcut, we only query the sessions table, not the games themselves.
		/////
		sessionFilename := path.Base(fpath) // yymmdd.{csv,wsv}
		sessionDate := strings.TrimSuffix(sessionFilename, path.Ext(sessionFilename)) // yymmdd
		if Contains(dates, sessionDate) {
			log.Printf("Session %s already recorded", sessionDate)
			//return nil
		} else {
			log.Printf("Session %s not yet recorded", sessionDate)
		}

		//////
		//  We need to read the session file and add it into the database.
		//  The session files are short enough to read in all at once.
		//////
		bytes, err := os.ReadFile(clubSessions + "/" + fpath)
		if err != nil {
			log.Fatalln(err)
		}

		//////
		//
		//  There are ten pieces of data for each game: The two players'
		//  names, scores, and bingoes, the two blanks, the round number
		//  and the lexicon.
		//
		//  In a session file, each line represents one game. Prior to
		//  20 February 2023, rounds were not indicated, and the winner
		//  was listed first, not the player who went first.
		//
		//  Prior to 20 February 2023, lines of the forms
		//    1) winner w_score [w_bingo[.*#]score ...] loser l_score [l_bingo[.*#]score ...]
		//  From 20 February 2023, lines are of the form
		//    2) p1,p2,p1_score,p2_score,p1_bingo:score[:...],p2_bingo:score[:...],blank_1,blank_2,lexicon
		//  We determine which is which heuristically.
		//
		//  Lexicon defaults to "nwl", the TWL/NWL lexicon in effect
		//  that day. It can be specified as "csw", the CSW lexicon in
		//  effect that day, or it can be specified explicitly, e.g.,
		//  "volost". We'll change "nwl" or "csw" into the specific one
		//  in effect. It must be one of the columns in the words table
		//  or phoneys can't be checked.
		//
		//  The player fields must already exist as handles in the
		//  players table. (This is for checking typos.)
		//
		//  A player's bingos are stored in the form
		//    word[phoney]:score[:...]
		//  phoney is one of:
		//    Nothing: Acceptable and not challenged
		//    !: Acceptable but unsuccessully challenged
		//    *: Phoney in either lexicon
		//    #: (NWL Games) Phoney in NWL but acceptable in CSW
		//    $: (CSW games) Phoney in CSW but acceptable in NWL
		//    +: Unacceptable but added later
		//    -: Acceptable but removed later
		//    ^: A non-bingo but otherwise interesting high play
		//       (must be acceptable)
		//  When playing lexicons outside NWL or CSW, e.g., volost,
		//  # and $ are moot. ^ should not be used for common high
		//  plays like Q-bombs.
		//
		// ISSUE: What about cases where a word is added and later
		// removed or vice versa?
		//
		//  From 15 May 2023, blanks used in bingos are indicated
		//  by the lower-case notational convention.
		//
		//////
		lines := strings.Split(string(bytes[:]), "\n")
		for i, line := range lines {
			// Skip blank lines and lines starting with #
			if line == "" || strings.HasPrefix(line, "#") {
				//log.Printf("File: %s Line: %d skipping comment: %s", fpath, i, line)
				continue
			}

			// In the session files, empty trailing fields don't have to
			// have the extra commas. Add them if needed.
			var fields []string
			commas := strings.Count(line, ",")
			if commas > 0 {
				// .csv
				line += strings.Repeat(",", 9-commas)
				fields = strings.Split(line, ",")
			} else {
				// .wsv
				fields = strings.Split(line, " ")
			}
			//log.Printf("File: %s Line: %d Content: %s", fpath, i, line)

			//////
			//
			//  Determine which of the four forms a line in a session
			//  file has, and parse the line accordingly to get the
			//  specific game data.
			//
			//////
			var handle1, handle2, bingos1, bingos2, blank1, blank2, lexicon string
			var score1, score2 int
			var round = 0
			if commas == 0 {
				// Form 1
				// The first two fields are the winning handle and score.
				// The next fields containing non-alphanumeric are the
				// winner's bingos. The next two fields are the losing
				// handle and score. The remaining fields are the loser's
				// bingos.
				const (
					START = iota
					HAVE_HANDLE1
					HAVE_SCORE1
					GETTING_BINGOS1
					HAVE_HANDLE2
					HAVE_SCORE2
					GETTING_BINGOS2
				)
				var state int = START
				for _, v := range fields {
					switch state {
					case START:
						//log.Printf("File: %s Line: %d State: START Field: %s\n", fpath, i, v)
						handle1 = v
						state = HAVE_HANDLE1
					case HAVE_HANDLE1:
						//log.Printf("File: %s Line: %d State: HAVE_HANDLE1 = %s Field: %s\n", fpath, i, handle1, v)
						if score1, err = strconv.Atoi(v); err != nil {
							log.Fatalf("File: %s Line: %d: Form 4 score1 not numeric\n", fpath, i)
						}
						state = HAVE_SCORE1
					case HAVE_SCORE1, GETTING_BINGOS1:
						//if state == HAVE_SCORE1 {
						//	log.Printf("File: %s Line: %d State: HAVE_SCORE1 = %d Field: %s\n", fpath, i, score1, v)
						//} else {
						//	log.Printf("File: %s Line: %d State: GETTING_BINGOS1 = %s Field: %s\n", fpath, i, bingos1, v)
						//}
						if strings.Contains(v, ":") {
							log.Fatalf("File: %s Line: %d Form 4 lines don't use colons: %s\n", fpath, i, v)
						}
						if strings.ContainsAny(v, ".*#$+") ||
						  (strings.Contains(v, "-") && strings.ContainsAny(v, "0123456789")) {
							if bingos1 == "" {
								bingos1 = v
							} else {
								bingos1 += " " + v
							}
							state = GETTING_BINGOS1
						} else {
							handle2 = v
							state = HAVE_HANDLE2
						}
					case HAVE_HANDLE2:
						//log.Printf("File: %s Line: %d State: HAVE_HANDLE2 = %s Field: %s\n", fpath, i, handle2, v)
						if score2, err = strconv.Atoi(v); err != nil {
							log.Fatalf("File: %s Line: %d: Form 4 score2 not numeric\n", fpath, i)
						}
						state = HAVE_SCORE2
					case HAVE_SCORE2, GETTING_BINGOS2:
						//if state == HAVE_SCORE1 {
						//	log.Printf("File: %s Line: %d State: HAVE_SCORE2 Field: %s\n", fpath, i, v)
						//} else {
						//	log.Printf("File: %s Line: %d State: GETTING_BINGOS2 Field: %s\n", fpath, i, v)
						//}
						if strings.Contains(v, ":") {
							log.Fatalf("File: %s Line: %d Form 4 lines don't use colons: %s\n", fpath, i, v)
						}
						if strings.ContainsAny(v, ".*#$+") {
							if bingos2 == "" {
								bingos2 = v
							} else {
								bingos2 += " " + v
							}
							state = GETTING_BINGOS2
						} else {
							log.Fatalf("File: %s Line: %d: Extra fields after loser's bingos.\n", fpath, i)
						}
					}
				}
			} else {
				// Form 2 -- current format
				fmt.Println("Form 1")
				if round, err = strconv.Atoi(fields[0]); err == nil {
					log.Fatal(err)
				}
				handle1 = fields[1]
				handle2 = fields[2]
				if score1, err = strconv.Atoi(fields[3]); err != nil {
					log.Fatalf("File: %s Line: %d: Form 1 score1 not numeric\n", fpath, i)
				}
				if score2, err = strconv.Atoi(fields[4]); err != nil {
					log.Fatalf("File: %s Line: %d: Form 1 score2 not numeric\n", fpath, i)
				}
				bingos1 = fields[5]
				bingos2 = fields[6]
				blank1 = fields[7] // May be a letter, - for unplayed, ? for unknown
				blank1 = fields[8] // May be a letter, - for unplayed, ? for unknown
				lexicon = fields[9]
			} // end of line parsing

			//////
			//
			//  We have all the data about the game (other than bingos,
			//  which we'll deal with later because they're stored in
			//  a different table) stuffed into all the right variables.
			//  Now look for inconsistencies and obvious mistakes.
			//
			//////
			if round == 0 {
				// Handle 1 is the winner and who went first is unknown.
				// Make sure the loser's score isn't higher than the
				// winner's. Ties are OK.
				fmt.Printf("%s %s %d %s %s %d %s\n", fpath, handle1, score1, bingos1, handle2, score2, bingos2)
				if (score2 > score1) {
					log.Fatalln("Loser has higher score than winner")
				}
				if !Contains(handles, handle1) {
					log.Fatalf("Unknown winner: %s", handle1)
				}
				if !Contains(handles, handle2) {
					log.Fatalf("Unknown loser: %s", handle2)
				}
			} else {
				// Handle 1 played first; calculate winner or tie
				log.Printf("First: %s %d %s; Second: %s %d %s\n", handle1, score1, bingos1, handle2, score2, bingos2)
				if !Contains(handles, handle1) {
					log.Fatalf("Unknown handle 1: %s", handle1)
				}
				if !Contains(handles, handle2) {
					log.Fatalf("Unknown handle 2: %s", handle2)
				}
			} // end of game analysis

			_ = score1; _ = score2; _ = blank1; _ = blank2; _ = round; _ = lexicon
			//log.Printf("Handle 1: %s, bingos %s, Handle 2: %s  bingos %s\n", handle1, bingos1, handle2, bingos2)
	} // end of lines in file

		return nil
	})
//  ^^
//  ||
//	|+-- end of WalkDir argument list
//	+--- end of WalkDir's anonymous function
}
