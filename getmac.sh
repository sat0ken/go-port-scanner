#!/bin/bash

arping -c 1 $1 | grep reply | cut -d" " -f5 | sed -e "s/\[//" -e "s/\]//"