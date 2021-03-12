module github.com/FerdinandWittmann/coli_crawler

go 1.15

require (
	github.com/kelvins/geocoder v0.0.0-20200113010004-f579500e9e27
	github.com/neo4j/neo4j-go-driver v1.8.4-0.20210129114204-8adb2b3a394f
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
)

replace github.com/FerdinandWittmann/coli_crawler/neo4j_extended => /home/workerferd/go/src/github.com/FerdinandWittmann/neo4j_extended/neo4j_extended
