#include <stdio.h>
#include <stdlib.h>

struct circle {
  int min, max, front;
  int * cups;

  circle(int min, int max, int front, int * cups):
    min(min), max(max), front(front), cups(cups) {}

  void print() {
    int cup = front;
    for (int i = 0; i < 9; i++) {
      printf("%d", cup);
      cup = cups[cup];
    }
    printf("\n");
  }

  void printSolution() {
    long long value1 = cups[1];
    long long value2 = cups[value1];
    long long product = value1 * value2;

    printf("%lld\n", product);
  }

  int nextDest(int current) {
    current--;
    if (current < min) {
      current = max;
    }
    return current;
  }

  void move() {
    int pickup[3];
    pickup[0] = cups[front];
    pickup[1] = cups[pickup[0]];
    pickup[2] = cups[pickup[1]];

    int dest = nextDest(front);
    while (dest == pickup[0] || dest == pickup[1] || dest == pickup[2]) {
      dest = nextDest(dest);
    }

    cups[front] = cups[pickup[2]];
    cups[pickup[2]] = cups[dest];
    cups[dest] = pickup[0];
    front = cups[front];
	}
};

circle readCircle(int realMax) {
	int * cups = (int*)malloc((realMax+1) * sizeof(int));

	int min = 9, max = 1, front = 0, prev = 0, cup;
  char str [20];
  int dummy = scanf("%s", str);
  for (int i = 0; str[i] != '\0'; i++) {
    cup = str[i] - '0';
    if (cup < min)
      min = cup;
    if (cup > max)
      max = cup;

    if (prev == 0) {
      front = cup;
    } else {
      cups[prev] = cup;
    }
    prev = cup;
  }

  if (realMax > max) {
    cups[prev] = max + 1;
    for (int i = max + 1; i < realMax; i++) {
      cups[i] = i + 1;
    }
    cups[realMax] = front;
  } else {
    cups[prev] = front;
  }


	return circle(min, realMax, front, cups);
}

int main() {
  int realMax = 1000000;
  int moves = 10000000;

  circle c = readCircle(realMax);
  for (int i = 0; i < moves; i++) {
    c.move();
  }

  c.printSolution();

  return 0;
}
