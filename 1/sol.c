#include <stdio.h>

int main() {
	int n;
	scanf("%d", &n);

	for (int i = 0; i < n; i++) {
		char a, b;
		scanf(" %c %c", &a, &b);

		if (a == b) {
			printf("Case #%d: -\n", i+1);
		} else if (a == 'R' && b == 'P') {
			printf("Case #%d: P\n", i+1);
		} else if (a == 'P' && b == 'S') {
			printf("Case #%d: S\n", i+1);
		} else if (a == 'S' && b == 'R') {
			printf("Case #%d: R\n", i+1);
		} else {
			printf("Case #%d: %c\n", i+1, a);
		}
	}

	return 0;
}
