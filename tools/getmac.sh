#!/bin/bash

arping -c 1 $1 | grep -o '[0-9A-F]\{2\}\(:[0-9A-F]\{2\}\)\{5\}'