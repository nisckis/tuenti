x = 1

while True:
  if x%2==1 and x%3==2 and x%4==3 and x%5==4 and x%6==5 and x%7==0:
    break

  x += 1

print(f"sol: {x}")
