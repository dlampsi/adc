echo "Init tests data..."

samba-tool user add testuser1 --random-password
samba-tool user add testuser2 --random-password

samba-tool group add testgroup1
samba-tool group addmembers testgroup1 testuser1

samba-tool group add testgroup2
samba-tool group addmembers testgroup2 testuser2
