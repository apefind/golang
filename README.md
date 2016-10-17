
EDL
===
EDL = Edit Decision List, extract information from edit decision list.

Usage
-----

	edltool convert -fps [24|30] -format [xlsx|csv] -output <file or directory> ...

		Convert edit decision list into xlsx or csv

	  -format string
	    	csv or xlsx (default "xlsx")
	  -fps int
	    	frames per second (default 30)
	  -o string
	    	short for -output
	  -output string
	    	csv output file or standard output

	edltool report

		Create usage report from edit decision list

	  -fps int
	    	frames per second (default 30)
	  -i string
	    	short for -input
	  -input string
	    	converted edl
	  -o string
	    	short for -output
	  -output string
	    	usage report

	edltool diff

		Compare two edit decision list

	  -columns string
	    	columns to define identical rows (default "Reel,Source In,Source Out,Notes")
	  -i string
	    	short for -input
	  -input string
	    	usage reports
	  -o string
	    	short for -output
	  -output string
	    	usage report


Episode Guide
=============
Rename series episodes according to tvrage.

Usage
-----

	episode_guide [rename|info]

  		-d=false:
		-dry_run=false: Just print, do not actually rename
		-m="tvrage":
		-method="tvrage": Only tvrage for now
		-n=false:
		-no_title=false: Ignore the title, just use S01E01, S01E02, ...


Shell Util
==========

Usage
-----

	shell_util [prompt|timeout|timeit]

	shell_util prompt -l <length> -r <ratio>

		print a length limited prompt

	  -l=32:
	  -length=32: maximum length of the prompt
	  -r=0.75:
	  -ratio=0.75: ratio between head/tail of the path

	shell_util timeout -d <duration> command [args ...]

		run a command under time limitation

	  -d=0:
	  -duration=0: kill the command after given duration

	shell_util timeit -n <repetitions> command [args ...]

	  	measure execution time of a command

	    -n=1: number of repetitions
	    -q=false: quiet run
