===========
Input Files
===========

Master-Client.go : This file contains the code for both master and slave nodes
slavesfile.txt   : This file contains details about slave nodes. Each line is of form "ip:port"

============
Output Files
============

All the output details are for running 2 slaves and a master node

output.png         : This contains screenshot of terminal which shows the local time each node. This time is printed after every 5s.
log_5555-Log.txt   : This is GoVector log for slave with port 5555
log_5556-Log.txt   : This is GoVector log for slave with port 5556
log_master-Log.txt : This is GoVector log for master node
shiviz.log         : This is the ShiViz log which is formed by combining all the individual logs.
                     To see shiviz output upload this file on "https://bestchai.bitbucket.io/shiviz/"
shiviz.png         : This is a screenshot of shiviz output. To see full result upload shiviz.log on website.

==========
How to Run
==========

You will have to first install GoVector
To run first fill the slavesfile.txt file. This should contain ip and port number of all the slave nodes.
First run each slave using the command
`go run Master-Client.go -s ip:port time logfile`
For each slave in slavesfile you should run this command

Then run master using the command
`go run Master-Client.go -m ip:port time slavesfile logfile`

The time in above commands specifies the initial time that node will start with. This is in seconds.
logfile for each node should be different.
You can see output.png to see example of how to run.

After running all slaves and master, local time to each node will be printed on terminal regularly.
After it is synchronized manually stop all the process
This will generate govec logs for each node, combine that using GoVector and generate a ShiViz compatible log and upload this on website.
