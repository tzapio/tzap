# this is a readme


tzap refactor --filein e2e_examples/utils.py \
--task "Fix the implementation of the function 'is_even' to make the tests work. Only answer with code and do not explain your solution." \
--outputformat python \
--temperature 0.0 \
--plan "Do not write any text because this file will be saved directly to e2e_examples/utils.py. Just answer with the plain code. Don not wrap the code in backticks." \
--automode=true
