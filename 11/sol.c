#include <stdio.h>
#include <string.h>

#define MAX 101
#define SIZE 512

long long int dp[MAX][MAX];
#define p(k, n) dp[k - 1][n - 1]

long long int dp2[MAX][MAX];
#define q(k, n) dp2[k - 1][n - 1]

long long int dp3[MAX][MAX];
#define x(k, n) dp3[k - 1][n - 1]

long long int pf[MAX][MAX];

long long int lldabs(long long int a) {
    return a < 0 ? -a : a;
}

long long int min(long long int a, long long int b) {
    return a < b ? a : b;
}

int main(int argc, char **argv)
{
    for (int i = 1; i <= MAX; i++)
    {
        for (int j = i; j > 0; j--)
        {
            if (j > i)
                p(j, i) = 0;
            else if (j == i)
                p(j, i) = 1;
            else
                p(j, i) = p(j, i - j) + p(j + 1, i);
        }
    }

    for (int n = 1; n <= MAX; n++)
    {
        for (int k = n; k > 0; k--)
        {
            if (k == 1 || (n == MAX && n == k))
                q(k, n) = 1;
            else if (k > n)
                q(k, n) = 0;
            else
                q(k, n) = q(k-1, n-1) + q(k, n-k);
        }
    }

    for (int n = 1; n <= MAX; n++)
    {
        for (int k = 1; k <= MAX; k++)
        {
            if (k >= n)
                x(k, n) = 0;
            else
                x(k, n) = p(1, n - k);
        }
    }

    int t, n, m, tmp;
    int data[MAX];
    long long int sol;

    int ret_code;
    char line[SIZE], *val;
    char delims[] = " \t\r\n";

    scanf("%d\n", &t);

    for (int i = 1; i <= t; i++)
    {
        if (fgets(line, SIZE, stdin) == NULL)
        {
            fprintf(stderr, "No input\n");
            return 1;
        }

        if (line[strlen(line) - 1] != '\n')
        {
            fprintf(stderr, "Line too long\n");
            return 2;
        }

        val = strtok(line, delims);
        ret_code = sscanf(val, "%d", &n);

        sol = p(1, n) - 1;
        memset(data, 0, sizeof(int) * n);

        val = strtok(NULL, delims);
        ret_code = (val == NULL) ? 0 : sscanf(val, "%d", &tmp);

        while (ret_code > 0)
        {
            data[tmp] = 1;
            val = strtok(NULL, delims);
            ret_code = (val == NULL) ? 0 : sscanf(val, "%d", &tmp);
        }

        printf("n = %d\n", n);
        long long int left, right, sl, sr;
        int mid, sm;

        // for (int i = 1; i < n; i++)
        // {
        //     sl = 0; 
        //     sm = 0; 
        //     sr = 0; 

        //     for (int j = i; j < n; j++)
        //     {
        //         mid   = n % j == 0 ? 1 : 0;
        //         right = - q(j, n) + x(j, n);
        //         left  = x(j, n) - mid - right;
        //         printf("(%lld | %d | %lld) %lld ", left, mid, right, q(j+1, n) - p(j, n));

        //         // if (j == i) continue;

        //         sl += left; 
        //         sm += mid; 
        //         sr += right;
        //     }

        //     printf(" -- %ld | %d | %lld\n", sl, sm, sr);
        // }

        for (int i = 0; i < n; i++)
            memset(pf, 0, sizeof(int) * n);

        int xd1;

        for (int i = 2; i < n; i++)
        {
            long long int x1, x2, tt;

            xd1 = n % i == 0 ? 1 : 0;
            x1 = x(1, n) - 1;
            x2 = x(i, n) - xd1;

            x1 = min(x1, x2);
            tt = lldabs(p(i, n) - q(i, n));

            printf("(1, %d) %lld %lld %lld\n", i, x1, tt, x1-tt);
        }
        
        printf("Case #%d: %lld\n", i, sol > 0 ? sol : 0);
    }

    int test_n = 8;

    printf("\n");

    for (int i = 1; i <= test_n; i++)
    {
        for (int j = 1; j <= i; j++) 
        {
            printf("p(%d,%d)=%lld  ", j, i, p(j, i));
        }
        printf("\n");
    }

    printf("\n");

    for (int i = 1; i <= test_n; i++)
    {
        for (int j = 1; j <= i; j++) 
        {
            printf("q(%d,%d)=%lld  ", j, i, q(j, i));
        }
        printf("\n");
    }
    

    printf("\n");

    for (int i = 1; i <= test_n; i++)
    {
        for (int j = 1; j <= i; j++) 
        {
            printf("x(%d,%d)=%lld  ", j, i, x(j, i));
        }
        printf("\n");
    }
}