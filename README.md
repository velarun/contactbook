# contactbook

Pre-requistie:

Golang should be installed. 

Version >= 1.10.1

MySQL Should be installed.

Version >= 8.0.15

Setup:

Description:

file_input.txt --> input file
parking_lot.go --> main file
parking_lot_test.go --> test file

Execution:


Run Test suite:

./parking_lot_testsuite


Sample:

Interactive Mode:

VelmuruganDhandapani-MacbookPro:ParkingLot velu$ ./parking_lot

$ 

$ create_parking_lot 6

Created parking lot with 6 slots

$ park KA-01-HH-1234 White
Allocated slot number: 1
$ park KA-01-HH-9999 White
Allocated slot number: 2
$ park KA-01-BB-0001 Black
Allocated slot number: 3
$ park KA-01-HH-7777 Red
Allocated slot number: 4
$ 
$ park KA-01-HH-2701 Blue
Allocated slot number: 5
$ park KA-01-HH-3141 Black
Allocated slot number: 6
$ 
$ 
$ 
$ lea
Invalid input
$ leave 4
Slot number 4 is free
$ status
Slot No.	Registration No			Colour
1		KA-01-HH-1234			White
2		KA-01-HH-9999			White
3		KA-01-BB-0001			Black
5		KA-01-HH-2701			Blue
6		KA-01-HH-3141			Black
$ park KA-01-P-333 White
Allocated slot number: 4
$ park DL-12-AA-9999 White
Sorry, parking lot is full
$ registration_numbers_for_cars_with_colour White
KA-01-HH-1234, KA-01-HH-9999, KA-01-P-333
$ slot_numbers_for_cars_with_colour White
1, 2, 4
$ slot_number_for_registration_number KA-01-HH-3141
6
$ slot_number_for_registration_number MH-04-AY-1111
Not Found
$ exit
VelmuruganDhandapani-MacbookPro:ParkingLot velu$ 

File mode:

VelmuruganDhandapani-MacbookPro:ParkingLot velu$ ./parking_lot file_input.txt 

Created parking lot with 6 slots
Allocated slot number: 1
Allocated slot number: 2
Allocated slot number: 3
Allocated slot number: 4
Allocated slot number: 5
Allocated slot number: 6
Slot number 4 is free
Slot No.	Registration No			Colour
1		KA-01-HH-1234			White
2		KA-01-HH-9999			White
3		KA-01-BB-0001			Black
5		KA-01-HH-2701			Blue
6		KA-01-HH-3141			Black
Allocated slot number: 4
Sorry, parking lot is full
KA-01-HH-1234, KA-01-HH-9999, KA-01-P-333
1, 2, 4
6
Not Found

Test Environment:

$ ./parking_lot_testsuite
=== RUN   TestCreatingParkingLot
--- PASS: TestCreatingParkingLot (0.00s)
=== RUN   TestParkingCarInParkingLot
--- PASS: TestParkingCarInParkingLot (0.00s)
=== RUN   TestLeaveCarFromParkingLot
--- PASS: TestLeaveCarFromParkingLot (0.00s)
=== RUN   TestGetRegNoByCarColor
--- PASS: TestGetRegNoByCarColor (0.00s)
=== RUN   TestGetSlotNoByCarColor
--- PASS: TestGetSlotNoByCarColor (0.00s)
=== RUN   TestGetSlotNoByCarRegNo
--- PASS: TestGetSlotNoByCarRegNo (0.00s)
=== RUN   TestParkingLotStatus
--- PASS: TestParkingLotStatus (0.00s)
=== RUN   TestRunCommandTask
--- PASS: TestRunCommandTask (0.00s)
PASS
ok  	command-line-arguments	(cached)
PASS
coverage: 79.5% of statements
ok  	_/Users/velu/workspace/Parking_Lot	0.010s

