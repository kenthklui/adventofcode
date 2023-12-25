# 2023 Day 24
Reinventing the wheel and writing hundreds of lines for solving a system of linear equations with 6-10 variables seemed like a bad idea. Doing it by hand also seemed painful. Since Golang doesn't support Z3 directly, I wrote part 2 in python instead.

```bash
pip install z3-solver importlib_resources
python3 task2.py < input.txt
```
