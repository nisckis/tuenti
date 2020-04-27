#include <stdio.h>
#include <stdint.h>

int main() {
    int c;
    uint64_t n;

    scanf("%d", &c);

    for (int i = 1; i <= c; i++) {
        scanf("%ld", &n);
        
        uint64_t res = n % 20, div = n / 20;

        if (div > 0) {
            if ((res <= 9) || (div >= 2 && res < 19) || (div >= 3)) {
                printf("Case #%d: %ld\n", i, div);
                continue;
            }
        }

        printf("Case #%d: IMPOSSIBLE\n", i);
    }
}