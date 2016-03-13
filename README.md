Trending application
====================

**tl;dr** - Found a simple algorithm for scoring keywords: TF-IDF - https://en.wikipedia.org/wiki/Tf-idf, which I then tried to tweak to find “trending” keywords within a target group consisting of all job postings within the last X days.

### Compiling:

    go get github.com/rhinoman/muse_trending
    go build
    
### Run unit tests:
    
    cd trend_analysis
    go test -v

### Usage:
The executable takes the following arguments:
- numPages: Number of pages to pull from The Muse API for analysis. Defaults to 10
- flexible: include jobs with flexible location. Defaults to true
- location (optional): limit jobs to a specific location. Defaults to none
- days: time period to identify trending terms (i.e., analyze trending terms in jobs during the last X days).  Defaults to 30
- stopWords (optional): fully qualified path to a text file containing language "stop words".  Default uses the small included file of english stop words

Example Usage:
```./muse_trending -numPages=1 -numResults=10 -days=30 -location="San Francisco Bay Area" ```

Example Output:

    Finished Processing

    **** 0 Errors occurred during processing ****
    ==== Trending Terms ====

    Location: San Francisco Bay Area
    Displaying the top 10 trending terms
    0)   Term: sales        Score: 11.3261
    1)   Term: mulesoft     Score: 8.1020
    2)   Term: partner      Score: 6.7706
    3)   Term: payments     Score: 6.4816
    4)   Term: marketing    Score: 6.4622
    5)   Term: quantcast    Score: 5.6714
    6)   Term: kernel       Score: 4.8612
    7)   Term: thumbtack    Score: 4.8612
    8)   Term: customer     Score: 4.8382
    9)   Term: company      Score: 4.6340
    10)  Term: customers    Score: 4.5137

### Basic algorithm:
- Job postings are divided into two “sets”: our target set(within day range and (optionally) location) and everything else(our ‘control’ set)
- Compute the Term Frequency for each term in our target set(within day range and optionally location)
- Compute the Inverse Document Frequency for each target set term.
- The target term frequency is multiplied by the IDF to provide a trending score.

### Notes
The API is not particularly helpful here:

1.	It would be very useful if the results returned from the jobs endpoint were in time order (by publication_date).  They’re not in any order I could ascertain.

2.	My first thought was to simply look at the “tags” field for keywords.  However, it appears the tags field isn’t used frequently enough in job records to be very useful, so going through the content is necessary.

3.	The contents field contains HTML.  So all of the HTML tags must be stripped out or otherwise dealt with.  I’ve decided to use regex as a quick-and-dirty solution for this, despite the dire warnings in this SO post: http://stackoverflow.com/questions/1732348/regex-match-open-tags-except-xhtml-self-contained-tags/1732454#1732454

