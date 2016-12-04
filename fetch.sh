#!/bin/sh

sources="https://www.heute.at/storage/rss/rss/heute.xml
https://feeds.feedburner.com/oe24?format=xml
http://www.krone.at/nachrichten/rss.html"

if [ "$#" -ne 1 ]; then
	echo "usage: $0 <filename>";
	exit 2;
fi

touch -a $1;
before=$(wc -l < $1);

for src in $sources; do
	echo "fetching from: $src";
	./titlemonkey fetch $src >> $1;
done

tmpfile=$(mktemp /tmp/titlemonkey.XXXXXX);
sort $1 | uniq > $tmpfile;
mv $tmpfile $1;

after=$(wc -l < $1);
new=$(expr $after - $before);
echo "fetched $new new titles";
