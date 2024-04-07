all:
	CGO_ENABLED=0 go build

man:
	mmark -man dnsfmt.1.md > dnsfmt.1
