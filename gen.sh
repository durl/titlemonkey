#!/bin/sh

if [ "$#" -ne 1 ]; then
	echo "usage: $0 <filename>";
	exit 2;
fi

# append titles to out file
touch -a $1;
tmpfile=$(mktemp /tmp/titlemonkey.XXXXXX);
cp $1 $tmpfile
./titlemonkey gen 10000 >> $tmpfile;

# sort new output, remove duplicates
tmpfile2=$(mktemp /tmp/titlemonkey.XXXXXX);
sort $tmpfile | uniq > $tmpfile2;
mv $tmpfile2 $tmpfile

# compare with previous outputfile
comm -13 $1 $tmpfile;
mv $tmpfile $1;
