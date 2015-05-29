
EDL
===
EDL = Edit Decision List, extract information from edit decision list into csv.

Usage
-----
::

	edltool extract -i <edl file> -o <csv file> -fps [24|30] [-auto]

		Extract information from edit decision list into csv

		-a=false:
		-auto=false: automatically choose input files and output names
		-fps=30:
		-frames-per-second=30: frames per second, usually 24 or 30
		-i="":
		-input="": edit decision list or standard input
		-o="":
		-output="": csv output file or standard output


Episode Guide
=============
Rename series episodes according to tvrage.

Usage
-----
::

	episode_guide [rename|info]

  		-d=false:
		-dry_run=false: Just print, do not actually rename
		-m="tvrage":
		-method="tvrage": Only tvrage for now
		-n=false:
		-no_title=false: Ignore the title, just use S01E01, S01E02, ...