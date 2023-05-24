# club-records
This repository contains the PostgreSQL database schema and related tools and documents for storage of Scrabble club records and reports based on that data. No actual data from any Scrabble club is included here, nor are any word lists, even though every North American word list since TWL98 and every Collins list since CSW07 are stored in our local database.

The early goal is to make CLI-based data entry as simple as possible. We have over 20 years of handwritten records to transcribe, along with new data from club sessions going forward. 

We have three types of game records to consider:

1. Data through February of 2023 consists of each game's:
   * Date
   * Winner, their score, and their bingos (bonuses);
   * Loser, their scores, and their bingos;
   * The lexicon used if not the TWL/NWL in effect at the time;
   * For each bingo, its score and whether it was acceptable.
2. Beginning in March, 2023, the following additional information is stored:
   * Which player went first;
   * The blanks that were played.
3. Moving forward, we also want to store:
   * Interesting non-bingo plays;
   * Challenges and their results;
   * Whether bingos were naturals or played with blanks;
   * Triple-triples (nine-timers).

Each phoney bingo is checked as to whether it was acceptable in its CSW or NWL counterpart, and whether it was added later. Acceptable bingos are checked to see if the word was subsequently removed from the lexicon, e.g., the recent removal of slurs.

We store data in the following tables:
1. In the database, each player is referred to by a unique handle. The `players` table stores that handle along with the player's name and contact information.
2. The `sessions` table stores every club session's date, whether it was virtual or live, and the location (for live sessions) or web site (e.g., Woogles) forvirtual sessions.
3. The `words` table contains the union of all the word lists mentioned above, and for each word, which word lists it belongs to.
4. The `games` table stores the information listed above about each game, except for bingos.
5. The `bingos` table stores each bingo played, the game and player, and the per-bingo information listed above.
6. The `lexica` table lists all lexica, e.g., TWL98, the family they belong to, NWL or CSW, if applicable (VOLOST is neither), and dates between which the lexicon was in effect.

Scrabble is a registered trademark of [Hasbro, Inc.](https://shop.hasbro.com/scrabble) in the USA and Canada, and of [Mattel Inc.](https://shopping.mattel.com/en-gb/collections/scrabble) elsewhere.

This work was originated by members of the Bridgewater Scrabble Club, which is also [NASPA Club #580](https://www.bridgewaterscrabble.org). The Club is solely repsonsible for, and is the sole owner of, this Github repository and its contents. It is neither affiliated with, nor endorsed by Hasbro, Mattel, [NASPA Games](https://scrabbleplayers.org) or any other sanctioning body for competitive Scrabble play.
