###check_logs is a Nagios plugin that can parse logs according to defined time period and regular expression.
Usage of ./check_logs:

  -b int

    	Amount of bytes to read from the end of a file (default 4096)

  -h int

    	Time interval for monitoring in hours

   Log's date format must be yyyy[/-]mm[/-]dd hh:mm:ss

   or yyyy[/-]mm[/-]ddThh:mm:ss at the beginning of a line; if --h=0 then it checks a current hour only (default 1)

  -r string

    	A regular expression for seeking (default "[Ee]xception")
