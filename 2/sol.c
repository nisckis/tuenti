#include <stdio.h>
#include <stdlib.h>

void solve(int i) {
    int l, bestp, bests = 0, *m;
    scanf("%d", &l);
    m = (int *)calloc(l*2, sizeof(int));

    int a, b, c;

    for (int j = 0; j < l; j++) {
        scanf("%d %d %d", &a, &b, &c);
        
        if (c == 1) {
            m[a-1]++;

            if (m[a-1] > bests) {
                bests = m[a-1];
                bestp = a;
            }
        } else {
            m[b-1]++;

            if (m[b-1] > bests) {
                bests = m[b-1];
                bestp = b;
            }
        }
    }

    free(m);
    printf("Case #%d: %d\n", i, bestp);
}

int main() {
	int n;
	scanf("%d", &n);

	for (int i = 0; i < n; i++) {
		solve(i+1);
	}

	return 0;
}
