# The session files
In order to make data entry as easy as possible, rather than implementing a clunky web form (which we've done but it's not part of this work), we've developed a concise text file format for storing all the games played in a given session.

## Session file name
Each session file is named <tt><i>yyyymmdd</i>.wsv</tt>, where <i>yyyymmdd</i> indicates the date of the session, and the <tt>.wsv</tt> extension indicates whitespace-separated values, e.g., <tt>20230522.wsv</tt>.

## What's stored
There are ten pieces of data for each game: The two players' handles, scores, and bingos, the two blanks, the round number within the session, and the lexicon or lexicon family.

## Handles
Handles are short lower-case strings, usually the player's first name. When a handle is already in use, a new player's handle will have a suffix disabmiguating it. For example, the handle of the founder of our club is <tt>scott</tt>. If another Scott starts playing, e.g., Scott Jones, their handle might be <tt>scott-j</tt>.

## Lexica
The default lexicon is the North American lexicon in effect on the date of the game; it may be explicitly specified as <tt>nwl</tt> or <tt>$</tt>. If the Collns lexicon is used, it may be specified as <tt>csw</tt> or <tt>#</tt>. A specific lexicon may also be used; the available lexica are:

<tt>twl98 owl1 owl2 owl2_1 twl14 twl16 nwl18 nwl20 csw07 csw12 csw15 csw19 csw21 volost</tt>

## Session file format
In a session file, each line represents one game. Historical data up to mid March, 2023 only contained the winner, loser, their scores, their bingos (bonuses) and those scores, and whether the game was played using the CSW lexicon (implying the 5-point challenge rule). For these games, their lines in the session file have the following space-separated fields, which are all on one physical line:

- Winner's handle
- Winner's score
- Winner's bingos (described below)
- Loser's handle
- Loser's score
- Loser's bingos
- The lexicon or lexicon family, if not the North American lexicon in effect that day

From March 20, 2023, the format of each line in a session file is:

- The round number
- First player's handle (not necessarily the winner)
- First player's score
- First player's bingos (described below)
- Second player's handle
- Second player's score
- Second player's bingos
- The first blank
- The second blank
- The lexicon or lexicon family, if not the North American lexicon in effect that day

If a score is unknown, it is listed as <tt>0</tt>. If a blank is known to be unplayed, it is listed as <tt>-</tt>. If it is unknown whether a blank was played, it is listed as <tt>?</tt>.

Blank lines and lines beginning with <tt>#</tt> are ignored.

## Bingo field format
Each bingo is specified in its own field, consisting of the word in lower case, a character indicating its validity, and the score of the play. The character is one of:

<table border=0>
  <tr><td style="width:2rem"><tt>.</tt><td>The word is acceptable.
  <tr><td><tt>/</tt><td>The word is acceptable and was unsuccessfully challenged.
  <tr><td><tt>*</tt><td>The word is unacceptable and not in any other word list.
  <tr><td><tt>#</tt><td>The word is unacceptable but is acceptable in CSW (NWL games only).
  <tr><td><tt>$</tt><td>The word is unacceptable but is acceptable in NWL (CSW games only).
  <tr><td><tt>+</tt><td>The word was unacceptable at the time but has since been added.
  <tr><td><tt>_</tt><td>The word was acceptable at the time but has since been removed. (A dash is usually used for this but that would cause an ambiguity between bingos and handles, which can also contain dashes.)
  <tr><td><tt>!</tt><td>The word is not a bingo, but an otherwise interesting high-scoring acceptable play.
  <tr><td><tt>^</tt><td>The word is a triple-triple (nine-timer).
</table>

(Q-bombs etc. are not interesting.)

There may be more than one of these characters, as long as they don't conflict, e.g., <tt>#</tt> and <tt>$</tt> and are not redundant, e.g., <tt>.</tt> and <tt>/</tt>.

##  Examples
In the format for older records, a line might look like this:

`nigel 635 sinters.91 starnie/75 goopily*83 infested^194 ed 223 csw`

meaning that Nigel beat Ed, 635-223, Nigel played SINTERS for 91, STARNIE for 80, which Ed unsuccessfully challenged (the score of 80 reflects a play score of 75 plus the 5-point challenge bonus), the phoney GOOPILY* for 83, and the triple-triple INFESTED for 194. Ed had no bingos, and the game used the CSW lexicon in effect on that date.

An equivalent line in the format for current records might look like this:

`3 nigel 635 sinters.91 starNie/80 goopily*83 infested^194 ed 223 n - csw`

In addition to the above, this line indicates that this game occurred in round 3 of the session, the first blank was an N, and the second blank was unplayed.

In bingos, we note blanks in upper-case (if we know), the opposite of the conventional notation like STARnIE. This is again for ease of data entry.