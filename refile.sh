#!/bin/sh
 
now="$(date +'%m%d')"

yesterday="$(date -d yesterday +'%m%d')"
 
rename "s/0111$yesterday/0111$now/" *.TXT.*

# sed -i "s/0111$yesterday/0111$now/g" *.TXT*
