#!/usr/bin/env bash

echo yskj1234 | passwd root --stdin

function ys_adduser() 
{
    useradd $1 -u $2; passwd -d $1; echo $1"123" | passwd $1 --stdin
}

# yskj:x:1000:1000::/home/yskj:/bin/bash
# user1:x:1001:1001::/home/user1:/bin/bash
# user2:x:1002:1002::/home/user2:/bin/bash
# user3:x:1003:1003::/home/user3:/bin/bash
# user4:x:1004:1004::/home/user4:/bin/bash
# user5:x:1005:1005::/home/user5:/bin/bash
ys_adduser yskj 1000
for i in {1..5}; do ys_adduser user$i 100$i; done