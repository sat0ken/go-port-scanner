#!/bin/bash

arping -c 1 $1 | grep -oi '[0-9A-F]\{2\}\(:[0-9A-F]\{2\}\)\{5\}'
