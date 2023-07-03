test-refactor:
	echo "y" | tzap init
	tzap refactor \
	  --filein e2e_examples/utils.py \
	  --task "Fix the implementation of the function 'is_even' to make the tests work. Only answer with code and do not explain your solution." \
	  --outputformat python \
	  --temperature 0.0 \
	  --plan "Make sure to only write python code without any further explanation since the results are direclry stored to the file "e2e_examples/utils.py". Just answer with the plain code. Also do not wrap the code in backticks." \
	  --automode=true
	poetry run pytest
